// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
 * Line tables
 */

package gosym

import (
	"bytes"
	"encoding/binary"
	"sync"
)

// version of the pclntab
type version int

const (
	verUnknown version = iota
	ver11
	ver12
	ver116
)

// A LineTable is a data structure mapping program counters to line numbers.
//
// In Go 1.1 and earlier, each function (represented by a Func) had its own LineTable,
// and the line number corresponded to a numbering of all source lines in the
// program, across all files. That absolute line number would then have to be
// converted separately to a file name and line number within the file.
//
// In Go 1.2, the format of the data changed so that there is a single LineTable
// for the entire program, shared by all Funcs, and there are no absolute line
// numbers, just line numbers within specific files.
//
// For the most part, LineTable's methods should be treated as an internal
// detail of the package; callers should use the methods on Table instead.
type LineTable struct {
	Data []byte
	PC   uint64
	Line int

	// This mutex is used to keep parsing of pclntab synchronous.
	mu sync.Mutex

	// Contains the version of the pclntab section.
	version version

	// Go 1.2/1.16 state
	binary      binary.ByteOrder
	quantum     uint32
	ptrsize     uint32
	funcnametab []byte
	cutab       []byte
	funcdata    []byte
	functab     []byte
	nfunctab    uint32
	filetab     []byte
	nfiletab    uint32
	fileMap     map[string]uint32
	funcNames   map[uint32]string // cache the function names
	strings     map[uint32]string // interned substrings of Data, keyed by offset
}

// NOTE(rsc): This is wrong for GOARCH=arm, which uses a quantum of 4,
// but we have no idea whether we're using arm or not. This only
// matters in the old (pre-Go 1.2) symbol table format, so it's not worth
// fixing.
const oldQuantum = 1

func (t *LineTable) parse(targetPC uint64, targetLine int) (b []byte, pc uint64, line int) {
	// The PC/line table can be thought of as a sequence of
	//  <pc update>* <line update>
	// batches. Each update batch results in a (pc, line) pair,
	// where line applies to every PC from pc up to but not
	// including the pc of the next pair.
	//
	// Here we process each update individually, which simplifies
	// the code, but makes the corner cases more confusing.
	b, pc, line = t.Data, t.PC, t.Line
	for pc <= targetPC && line != targetLine && len(b) > 0 {
		code := b[0]
		b = b[1:]
		switch {
		case code == 0:
			if len(b) < 4 {
				b = b[0:0]
				break
			}
			val := binary.BigEndian.Uint32(b)
			b = b[4:]
			line += int(val)
		case code <= 64:
			line += int(code)
		case code <= 128:
			line -= int(code - 64)
		default:
			pc += oldQuantum * uint64(code-128)
			continue
		}
		pc += oldQuantum
	}
	return b, pc, line
}

func (t *LineTable) slice(pc uint64) *LineTable {
	data, pc, line := t.parse(pc, -1)
	return &LineTable{Data: data, PC: pc, Line: line}
}

// PCToLine returns the line number for the given program counter.
//
// Deprecated: Use Table's PCToLine method instead.
func (t *LineTable) PCToLine(pc uint64) int {
	if t.isGo12() {
		return t.go12PCToLine(pc)
	}
	_, _, line := t.parse(pc, -1)
	return line
}

// LineToPC returns the program counter for the given line number,
// considering only program counters before maxpc.
//
// Deprecated: Use Table's LineToPC method instead.
func (t *LineTable) LineToPC(line int, maxpc uint64) uint64 {
	if t.isGo12() {
		return 0
	}
	_, pc, line1 := t.parse(maxpc, line)
	if line1 != line {
		return 0
	}
	// Subtract quantum from PC to account for post-line increment
	return pc - oldQuantum
}

// NewLineTable returns a new PC/line table
// corresponding to the encoded data.
// Text must be the start address of the
// corresponding text segment.
func NewLineTable(data []byte, text uint64) *LineTable {
	return &LineTable{Data: data, PC: text, Line: 0, funcNames: make(map[uint32]string), strings: make(map[uint32]string)}
}

// Go 1.2 symbol table format.
// See golang.org/s/go12symtab.
//
// A general note about the methods here: rather than try to avoid
// index out of bounds errors, we trust Go to detect them, and then
// we recover from the panics and treat them as indicative of a malformed
// or incomplete table.
//
// The methods called by symtab.go, which begin with "go12" prefixes,
// are expected to have that recovery logic.

// isGo12 reports whether this is a Go 1.2 (or later) symbol table.
func (t *LineTable) isGo12() bool {
	t.parsePclnTab()
	return t.version >= ver12
}

const go12magic = 0xfffffffb
const go116magic = 0xfffffffa

// uintptr returns the pointer-sized value encoded at b.
// The pointer size is dictated by the table being read.
func (t *LineTable) uintptr(b []byte) uint64 {
	if t.ptrsize == 4 {
		return uint64(t.binary.Uint32(b))
	}
	return t.binary.Uint64(b)
}

// parsePclnTab parses the pclntab, setting the version.
func (t *LineTable) parsePclnTab() {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.version != verUnknown {
		return
	}

	// Note that during this function, setting the version is the last thing we do.
	// If we set the version too early, and parsing failed (likely as a panic on
	// slice lookups), we'd have a mistaken version.
	//
	// Error paths through this code will default the version to 1.1.
	t.version = ver11

	defer func() {
		// If we panic parsing, assume it's a Go 1.1 pclntab.
		recover()
	}()

	// Check header: 4-byte magic, two zeros, pc quantum, pointer size.
	if len(t.Data) < 16 || t.Data[4] != 0 || t.Data[5] != 0 ||
		(t.Data[6] != 1 && t.Data[6] != 2 && t.Data[6] != 4) || // pc quantum
		(t.Data[7] != 4 && t.Data[7] != 8) { // pointer size
		return
	}

	var possibleVersion version
	leMagic := binary.LittleEndian.Uint32(t.Data)
	beMagic := binary.BigEndian.Uint32(t.Data)
	switch {
	case leMagic == go12magic:
		t.binary, possibleVersion = binary.LittleEndian, ver12
	case beMagic == go12magic:
		t.binary, possibleVersion = binary.BigEndian, ver12
	case leMagic == go116magic:
		t.binary, possibleVersion = binary.LittleEndian, ver116
	case beMagic == go116magic:
		t.binary, possibleVersion = binary.BigEndian, ver116
	default:
		return
	}

	// quantum and ptrSize are the same between 1.2 and 1.16
	t.quantum = uint32(t.Data[6])
	t.ptrsize = uint32(t.Data[7])

	switch possibleVersion {
	case ver116:
		t.nfunctab = uint32(t.uintptr(t.Data[8:]))
		t.nfiletab = uint32(t.uintptr(t.Data[8+t.ptrsize:]))
		offset := t.uintptr(t.Data[8+2*t.ptrsize:])
		t.funcnametab = t.Data[offset:]
		offset = t.uintptr(t.Data[8+3*t.ptrsize:])
		t.cutab = t.Data[offset:]
		offset = t.uintptr(t.Data[8+4*t.ptrsize:])
		t.filetab = t.Data[offset:]
		offset = t.uintptr(t.Data[8+5*t.ptrsize:])
		t.funcdata = t.Data[offset:]
		t.functab = t.Data[offset:]
		functabsize := t.nfunctab*2*t.ptrsize + t.ptrsize
		t.functab = t.functab[:functabsize]
	case ver12:
		t.nfunctab = uint32(t.uintptr(t.Data[8:]))
		t.funcdata = t.Data
		t.funcnametab = t.Data
		t.functab = t.Data[8+t.ptrsize:]
		functabsize := t.nfunctab*2*t.ptrsize + t.ptrsize
		fileoff := t.binary.Uint32(t.functab[functabsize:])
		t.functab = t.functab[:functabsize]
		t.filetab = t.Data[fileoff:]
		t.nfiletab = t.binary.Uint32(t.filetab)
		t.filetab = t.filetab[:t.nfiletab*4]
	default:
		panic("unreachable")
	}
	t.version = possibleVersion
}

// go12Funcs returns a slice of Funcs derived from the Go 1.2 pcln table.
func (t *LineTable) go12Funcs() []Func {
	// Assume it is malformed and return nil on error.
	defer func() {
		recover()
	}()

	n := len(t.functab) / int(t.ptrsize) / 2
	funcs := make([]Func, n)
	for i := range funcs {
		f := &funcs[i]
		f.Entry = t.uintptr(t.functab[2*i*int(t.ptrsize):])
		f.End = t.uintptr(t.functab[(2*i+2)*int(t.ptrsize):])
		info := t.funcdata[t.uintptr(t.functab[(2*i+1)*int(t.ptrsize):]):]
		f.LineTable = t
		f.FrameSize = int(t.binary.Uint32(info[t.ptrsize+2*4:]))
		f.Sym = &Sym{
			Value:  f.Entry,
			Type:   'T',
			Name:   t.funcName(t.binary.Uint32(info[t.ptrsize:])),
			GoType: 0,
			Func:   f,
		}
	}
	return funcs
}

// findFunc returns the func corresponding to the given program counter.
func (t *LineTable) findFunc(pc uint64) []byte {
	if pc < t.uintptr(t.functab) || pc >= t.uintptr(t.functab[len(t.functab)-int(t.ptrsize):]) {
		return nil
	}

	// The function table is a list of 2*nfunctab+1 uintptrs,
	// alternating program counters and offsets to func structures.
	f := t.functab
	nf := t.nfunctab
	for nf > 0 {
		m := nf / 2
		fm := f[2*t.ptrsize*m:]
		if t.uintptr(fm) <= pc && pc < t.uintptr(fm[2*t.ptrsize:]) {
			return t.funcdata[t.uintptr(fm[t.ptrsize:]):]
		} else if pc < t.uintptr(fm) {
			nf = m
		} else {
			f = f[(m+1)*2*t.ptrsize:]
			nf -= m + 1
		}
	}
	return nil
}

// readvarint reads, removes, and returns a varint from *pp.
func (t *LineTable) readvarint(pp *[]byte) uint32 {
	var v, shift uint32
	p := *pp
	for shift = 0; ; shift += 7 {
		b := p[0]
		p = p[1:]
		v |= (uint32(b) & 0x7F) << shift
		if b&0x80 == 0 {
			break
		}
	}
	*pp = p
	return v
}

// funcName returns the name of the function found at off.
func (t *LineTable) funcName(off uint32) string {
	if s, ok := t.funcNames[off]; ok {
		return s
	}
	i := bytes.IndexByte(t.funcnametab[off:], 0)
	s := string(t.funcnametab[off : off+uint32(i)])
	t.funcNames[off] = s
	return s
}

// stringFrom returns a Go string found at off from a position.
func (t *LineTable) stringFrom(arr []byte, off uint32) string {
	if s, ok := t.strings[off]; ok {
		return s
	}
	i := bytes.IndexByte(arr[off:], 0)
	s := string(arr[off : off+uint32(i)])
	t.strings[off] = s
	return s
}

// string returns a Go string found at off.
func (t *LineTable) string(off uint32) string {
	return t.stringFrom(t.funcdata, off)
}

// step advances to the next pc, value pair in the encoded table.
func (t *LineTable) step(p *[]byte, pc *uint64, val *int32, first bool) bool {
	uvdelta := t.readvarint(p)
	if uvdelta == 0 && !first {
		return false
	}
	if uvdelta&1 != 0 {
		uvdelta = ^(uvdelta >> 1)
	} else {
		uvdelta >>= 1
	}
	vdelta := int32(uvdelta)
	pcdelta := t.readvarint(p) * t.quantum
	*pc += uint64(pcdelta)
	*val += vdelta
	return true
}

// pcvalue reports the value associated with the target pc.
// off is the offset to the beginning of the pc-value table,
// and entry is the start PC for the corresponding function.
func (t *LineTable) pcvalue(off uint32, entry, targetpc uint64) int32 {
	p := t.funcdata[off:]

	val := int32(-1)
	pc := entry
	for t.step(&p, &pc, &val, pc == entry) {
		if targetpc < pc {
			return val
		}
	}
	return -1
}

// findFileLine scans one function in the binary looking for a
// program counter in the given file on the given line.
// It does so by running the pc-value tables mapping program counter
// to file number. Since most functions come from a single file, these
// are usually short and quick to scan. If a file match is found, then the
// code goes to the expense of looking for a simultaneous line number match.
func (t *LineTable) findFileLine(entry uint64, filetab, linetab uint32, filenum, line int32) uint64 {
	if filetab == 0 || linetab == 0 {
		return 0
	}

	fp := t.funcdata[filetab:]
	fl := t.funcdata[linetab:]
	fileVal := int32(-1)
	filePC := entry
	lineVal := int32(-1)
	linePC := entry
	fileStartPC := filePC
	for t.step(&fp, &filePC, &fileVal, filePC == entry) {
		if fileVal == filenum && fileStartPC < filePC {
			// fileVal is in effect starting at fileStartPC up to
			// but not including filePC, and it's the file we want.
			// Run the PC table looking for a matching line number
			// or until we reach filePC.
			lineStartPC := linePC
			for linePC < filePC && t.step(&fl, &linePC, &lineVal, linePC == entry) {
				// lineVal is in effect until linePC, and lineStartPC < filePC.
				if lineVal == line {
					if fileStartPC <= lineStartPC {
						return lineStartPC
					}
					if fileStartPC < linePC {
						return fileStartPC
					}
				}
				lineStartPC = linePC
			}
		}
		fileStartPC = filePC
	}
	return 0
}

// go12PCToLine maps program counter to line number for the Go 1.2 pcln table.
func (t *LineTable) go12PCToLine(pc uint64) (line int) {
	defer func() {
		if recover() != nil {
			line = -1
		}
	}()

	f := t.findFunc(pc)
	if f == nil {
		return -1
	}
	entry := t.uintptr(f)
	linetab := t.binary.Uint32(f[t.ptrsize+5*4:])
	return int(t.pcvalue(linetab, entry, pc))
}

// go12PCToFile maps program counter to file name for the Go 1.2 pcln table.
func (t *LineTable) go12PCToFile(pc uint64) (file string) {
	defer func() {
		if recover() != nil {
			file = ""
		}
	}()

	f := t.findFunc(pc)
	if f == nil {
		return ""
	}
	entry := t.uintptr(f)
	filetab := t.binary.Uint32(f[t.ptrsize+4*4:])
	fno := t.pcvalue(filetab, entry, pc)
	if fno <= 0 {
		return ""
	}
	if t.version == ver12 {
		return t.string(t.binary.Uint32(t.filetab[4*fno:]))
	}
	// Go ≥ 1.16
	cuoff := t.binary.Uint32(f[t.ptrsize+7*4:])
	if fnoff := t.binary.Uint32(t.cutab[(cuoff+uint32(fno))*4:]); fnoff != ^uint32(0) {
		return t.stringFrom(t.filetab, fnoff)
	}
	return ""
}

// go12LineToPC maps a (file, line) pair to a program counter for the Go 1.2 pcln table.
func (t *LineTable) go12LineToPC(file string, line int) (pc uint64) {
	defer func() {
		if recover() != nil {
			pc = 0
		}
	}()

	t.initFileMap()
	filenum := t.fileMap[file]
	if filenum == 0 {
		return 0
	}

	// Scan all functions.
	// If this turns out to be a bottleneck, we could build a map[int32][]int32
	// mapping file number to a list of functions with code from that file.
	for i := uint32(0); i < t.nfunctab; i++ {
		f := t.funcdata[t.uintptr(t.functab[2*t.ptrsize*i+t.ptrsize:]):]
		entry := t.uintptr(f)
		filetab := t.binary.Uint32(f[t.ptrsize+4*4:])
		linetab := t.binary.Uint32(f[t.ptrsize+5*4:])
		pc := t.findFileLine(entry, filetab, linetab, int32(filenum), int32(line))
		if pc != 0 {
			return pc
		}
	}
	return 0
}

// initFileMap initializes the map from file name to file number.
func (t *LineTable) initFileMap() {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.fileMap != nil {
		return
	}
	m := make(map[string]uint32)

	if t.version == ver12 {
		for i := uint32(1); i < t.nfiletab; i++ {
			s := t.string(t.binary.Uint32(t.filetab[4*i:]))
			m[s] = i
		}
	} else {
		var pos uint32
		for i := uint32(1); i < t.nfiletab; i++ {
			s := t.stringFrom(t.filetab, pos)
			pos += uint32(len(s) + 1)
			m[s] = i
		}
	}
	t.fileMap = m
}

// go12MapFiles adds to m a key for every file in the Go 1.2 LineTable.
// Every key maps to obj. That's not a very interesting map, but it provides
// a way for callers to obtain the list of files in the program.
func (t *LineTable) go12MapFiles(m map[string]*Obj, obj *Obj) {
	defer func() {
		recover()
	}()

	t.initFileMap()
	for file := range t.fileMap {
		m[file] = obj
	}
}

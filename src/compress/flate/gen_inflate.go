// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package flate implements the DEFLATE compressed data format, described in
// RFC 1951.  The gzip and zlib packages implement access to DEFLATE-based file
// formats.

// +build generate

//go:generate go run $GOFILE && gofmt -w inflate_gen.go

package main

import (
	"os"
	"strings"
)

func main() {
	f, err := os.Create("inflate_gen.go")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	types := []string{"*bytes.Buffer", "*bytes.Reader", "*bufio.Reader", "*strings.Reader"}
	names := []string{"BytesBuffer", "BytesReader", "BufioReader", "StringsReader"}
	imports := []string{"bytes", "bufio", "io", "strings", "math/bits"}
	f.WriteString(`// Code generated by go generate gen_inflate.go. DO NOT EDIT.

package flate

import (
`)

	for _, imp := range imports {
		f.WriteString("\t\"" + imp + "\"\n")
	}
	f.WriteString(")\n\n")

	template := `

// $FUNCNAME$ decodes a single Huffman block from f.
// f.r must be a $TYPE$.
// hl and hd are the Huffman states for the lit/length values
// and the distance values, respectively. If hd == nil, using the
// fixed distance encoding associated with fixed Huffman blocks.
func (f *decompressor) $FUNCNAME$() {
	const (
		stateInit = iota // Zero value must be stateInit
		stateDict
	)
	fr := f.r.($TYPE$)
	moreBits := func() error {
		c, err := fr.ReadByte()
		if err != nil {
			return noEOF(err)
		}
		f.roffset++
		f.b |= uint32(c) << f.nb
		f.nb += 8
		return nil
	}

	switch f.stepState {
	case stateInit:
		goto readLiteral
	case stateDict:
		goto copyHistory
	}

readLiteral:
	// Read literal and/or (length, distance) according to RFC section 3.2.3.
	{
		var v int
		{
			// Inlined v, err := f.huffSym(f.hl)
			// Since a huffmanDecoder can be empty or be composed of a degenerate tree
			// with single element, huffSym must error on these two edge cases. In both
			// cases, the chunks slice will be 0 for the invalid sequence, leading it
			// satisfy the n == 0 check below.
			n := uint(f.hl.maxRead)
			// Optimization. Compiler isn't smart enough to keep f.b,f.nb in registers,
			// but is smart enough to keep local variables in registers, so use nb and b,
			// inline call to moreBits and reassign b,nb back to f on return.
			nb, b := f.nb, f.b
			for {
				for nb < n {
					c, err := fr.ReadByte()
					if err != nil {
						f.b = b
						f.nb = nb
						f.err = noEOF(err)
						return
					}
					f.roffset++
					b |= uint32(c) << (nb & 31)
					nb += 8
				}
				chunk := f.hl.chunks[b&(huffmanNumChunks-1)]
				n = uint(chunk & huffmanCountMask)
				if n > huffmanChunkBits {
					chunk = f.hl.links[chunk>>huffmanValueShift][(b>>huffmanChunkBits)&f.hl.linkMask]
					n = uint(chunk & huffmanCountMask)
				}
				if n <= nb {
					if n == 0 {
						f.b = b
						f.nb = nb
						f.err = CorruptInputError(f.roffset)
						return
					}
					f.b = b >> (n & 31)
					f.nb = nb - n
					v = int(chunk >> huffmanValueShift)
					break
				}
			}
		}

		var n uint // number of bits extra
		var length int
		var err error
		switch {
		case v < 256:
			f.dict.writeByte(byte(v))
			if f.dict.availWrite() == 0 {
				f.toRead = f.dict.readFlush()
				f.step = (*decompressor).$FUNCNAME$
				f.stepState = stateInit
				return
			}
			goto readLiteral
		case v == 256:
			f.finishBlock()
			return
		// otherwise, reference to older data
		case v < 265:
			length = v - (257 - 3)
			n = 0
		case v < 269:
			length = v*2 - (265*2 - 11)
			n = 1
		case v < 273:
			length = v*4 - (269*4 - 19)
			n = 2
		case v < 277:
			length = v*8 - (273*8 - 35)
			n = 3
		case v < 281:
			length = v*16 - (277*16 - 67)
			n = 4
		case v < 285:
			length = v*32 - (281*32 - 131)
			n = 5
		case v < maxNumLit:
			length = 258
			n = 0
		default:
			f.err = CorruptInputError(f.roffset)
			return
		}
		if n > 0 {
			for f.nb < n {
				if err = moreBits(); err != nil {
					f.err = err
					return
				}
			}
			length += int(f.b & uint32(1<<n-1))
			f.b >>= n
			f.nb -= n
		}

		var dist int
		if f.hd == nil {
			for f.nb < 5 {
				if err = moreBits(); err != nil {
					f.err = err
					return
				}
			}
			dist = int(bits.Reverse8(uint8(f.b & 0x1F << 3)))
			f.b >>= 5
			f.nb -= 5
		} else {
			if dist, err = f.huffSym(f.hd); err != nil {
				f.err = err
				return
			}
		}

		switch {
		case dist < 4:
			dist++
		case dist < maxNumDist:
			nb := uint(dist-2) >> 1
			// have 1 bit in bottom of dist, need nb more.
			extra := (dist & 1) << nb
			for f.nb < nb {
				if err = moreBits(); err != nil {
					f.err = err
					return
				}
			}
			extra |= int(f.b & uint32(1<<nb-1))
			f.b >>= nb
			f.nb -= nb
			dist = 1<<(nb+1) + 1 + extra
		default:
			f.err = CorruptInputError(f.roffset)
			return
		}

		// No check on length; encoding can be prescient.
		if dist > f.dict.histSize() {
			f.err = CorruptInputError(f.roffset)
			return
		}

		f.copyLen, f.copyDist = length, dist
		goto copyHistory
	}

copyHistory:
	// Perform a backwards copy according to RFC section 3.2.3.
	{
		cnt := f.dict.tryWriteCopy(f.copyDist, f.copyLen)
		if cnt == 0 {
			cnt = f.dict.writeCopy(f.copyDist, f.copyLen)
		}
		f.copyLen -= cnt

		if f.dict.availWrite() == 0 || f.copyLen > 0 {
			f.toRead = f.dict.readFlush()
			f.step = (*decompressor).$FUNCNAME$ // We need to continue this work
			f.stepState = stateDict
			return
		}
		goto readLiteral
	}
}

`
	for i, t := range types {
		s := strings.Replace(template, "$FUNCNAME$", "huffman"+names[i], -1)
		s = strings.Replace(s, "$TYPE$", t, -1)
		f.WriteString(s)
	}
	f.WriteString("func (f *decompressor) huffmanBlockDecoder() func() {\n")
	f.WriteString("\tswitch f.r.(type) {\n")
	for i, t := range types {
		f.WriteString("\t\tcase " + t + ":\n")
		f.WriteString("\t\t\treturn f.huffman" + names[i] + "\n")
	}
	f.WriteString("\t\tdefault:\n")
	f.WriteString("\t\t\treturn f.huffmanBlockGeneric")
	f.WriteString("\t}\n}\n")
}

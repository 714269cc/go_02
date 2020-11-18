// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ssa

import (
	"cmd/compile/internal/ir"
	"cmd/compile/internal/types"
	"cmd/internal/obj"
	"cmd/internal/obj/arm64"
	"cmd/internal/obj/s390x"
	"cmd/internal/obj/x86"
	"cmd/internal/src"
	"fmt"
	"testing"
)

var CheckFunc = checkFunc
var Opt = opt
var Deadcode = deadcode
var Copyelim = copyelim

var testCtxts = map[string]*obj.Link{
	"amd64": obj.Linknew(&x86.Linkamd64),
	"s390x": obj.Linknew(&s390x.Links390x),
	"arm64": obj.Linknew(&arm64.Linkarm64),
}

func testConfig(tb testing.TB) *Conf      { return testConfigArch(tb, "amd64") }
func testConfigS390X(tb testing.TB) *Conf { return testConfigArch(tb, "s390x") }
func testConfigARM64(tb testing.TB) *Conf { return testConfigArch(tb, "arm64") }

func testConfigArch(tb testing.TB, arch string) *Conf {
	ctxt, ok := testCtxts[arch]
	if !ok {
		tb.Fatalf("unknown arch %s", arch)
	}
	if ctxt.Arch.PtrSize != 8 {
		tb.Fatal("testTypes is 64-bit only")
	}
	c := &Conf{
		config: NewConfig(arch, testTypes, ctxt, true),
		tb:     tb,
	}
	return c
}

type Conf struct {
	config *Config
	tb     testing.TB
	fe     Frontend
}

func (c *Conf) Frontend() Frontend {
	if c.fe == nil {
		c.fe = TestFrontend{t: c.tb, ctxt: c.config.ctxt}
	}
	return c.fe
}

// TestFrontend is a test-only frontend.
// It assumes 64 bit integers and pointers.
type TestFrontend struct {
	t    testing.TB
	ctxt *obj.Link
}

func (TestFrontend) StringData(s string) *obj.LSym {
	return nil
}
func (TestFrontend) Auto(pos src.XPos, t *types.Type) *ir.Node {
	n := ir.NewNameAt(pos, &types.Sym{Name: "aFakeAuto"})
	n.SetClass(ir.PAUTO)
	return n
}
func (d TestFrontend) SplitString(s LocalSlot) (LocalSlot, LocalSlot) {
	return LocalSlot{N: s.N, Type: testTypes.BytePtr, Off: s.Off}, LocalSlot{N: s.N, Type: testTypes.Int, Off: s.Off + 8}
}
func (d TestFrontend) SplitInterface(s LocalSlot) (LocalSlot, LocalSlot) {
	return LocalSlot{N: s.N, Type: testTypes.BytePtr, Off: s.Off}, LocalSlot{N: s.N, Type: testTypes.BytePtr, Off: s.Off + 8}
}
func (d TestFrontend) SplitSlice(s LocalSlot) (LocalSlot, LocalSlot, LocalSlot) {
	return LocalSlot{N: s.N, Type: s.Type.Elem().PtrTo(), Off: s.Off},
		LocalSlot{N: s.N, Type: testTypes.Int, Off: s.Off + 8},
		LocalSlot{N: s.N, Type: testTypes.Int, Off: s.Off + 16}
}
func (d TestFrontend) SplitComplex(s LocalSlot) (LocalSlot, LocalSlot) {
	if s.Type.Size() == 16 {
		return LocalSlot{N: s.N, Type: testTypes.Float64, Off: s.Off}, LocalSlot{N: s.N, Type: testTypes.Float64, Off: s.Off + 8}
	}
	return LocalSlot{N: s.N, Type: testTypes.Float32, Off: s.Off}, LocalSlot{N: s.N, Type: testTypes.Float32, Off: s.Off + 4}
}
func (d TestFrontend) SplitInt64(s LocalSlot) (LocalSlot, LocalSlot) {
	if s.Type.IsSigned() {
		return LocalSlot{N: s.N, Type: testTypes.Int32, Off: s.Off + 4}, LocalSlot{N: s.N, Type: testTypes.UInt32, Off: s.Off}
	}
	return LocalSlot{N: s.N, Type: testTypes.UInt32, Off: s.Off + 4}, LocalSlot{N: s.N, Type: testTypes.UInt32, Off: s.Off}
}
func (d TestFrontend) SplitStruct(s LocalSlot, i int) LocalSlot {
	return LocalSlot{N: s.N, Type: s.Type.FieldType(i), Off: s.Off + s.Type.FieldOff(i)}
}
func (d TestFrontend) SplitArray(s LocalSlot) LocalSlot {
	return LocalSlot{N: s.N, Type: s.Type.Elem(), Off: s.Off}
}

func (d TestFrontend) SplitSlot(parent *LocalSlot, suffix string, offset int64, t *types.Type) LocalSlot {
	return LocalSlot{N: parent.N, Type: t, Off: offset}
}
func (TestFrontend) Line(_ src.XPos) string {
	return "unknown.go:0"
}
func (TestFrontend) AllocFrame(f *Func) {
}
func (d TestFrontend) Syslook(s string) *obj.LSym {
	return d.ctxt.Lookup(s)
}
func (TestFrontend) UseWriteBarrier() bool {
	return true // only writebarrier_test cares
}
func (TestFrontend) SetWBPos(pos src.XPos) {
}

func (d TestFrontend) Logf(msg string, args ...interface{}) { d.t.Logf(msg, args...) }
func (d TestFrontend) Log() bool                            { return true }

func (d TestFrontend) Fatalf(_ src.XPos, msg string, args ...interface{}) { d.t.Fatalf(msg, args...) }
func (d TestFrontend) Warnl(_ src.XPos, msg string, args ...interface{})  { d.t.Logf(msg, args...) }
func (d TestFrontend) Debug_checknil() bool                               { return false }

func (d TestFrontend) MyImportPath() string {
	return "my/import/path"
}

var testTypes Types

func init() {
	// Initialize just enough of the universe and the types package to make our tests function.
	// TODO(josharian): move universe initialization to the types package,
	// so this test setup can share it.

	types.Tconv = func(t *types.Type, flag, mode int) string {
		return t.Etype.String()
	}
	types.Sconv = func(s *types.Sym, flag, mode int) string {
		return "sym"
	}
	types.FormatSym = func(sym *types.Sym, s fmt.State, verb rune, mode int) {
		fmt.Fprintf(s, "sym")
	}
	types.FormatType = func(t *types.Type, s fmt.State, verb rune, mode int) {
		fmt.Fprintf(s, "%v", t.Etype)
	}
	types.Dowidth = func(t *types.Type) {}

	for _, typ := range [...]struct {
		width int64
		et    types.EType
	}{
		{1, types.TINT8},
		{1, types.TUINT8},
		{1, types.TBOOL},
		{2, types.TINT16},
		{2, types.TUINT16},
		{4, types.TINT32},
		{4, types.TUINT32},
		{4, types.TFLOAT32},
		{4, types.TFLOAT64},
		{8, types.TUINT64},
		{8, types.TINT64},
		{8, types.TINT},
		{8, types.TUINTPTR},
	} {
		t := types.New(typ.et)
		t.Width = typ.width
		t.Align = uint8(typ.width)
		types.Types[typ.et] = t
	}
	testTypes.SetTypPtrs()
}

func (d TestFrontend) DerefItab(sym *obj.LSym, off int64) *obj.LSym { return nil }

func (d TestFrontend) CanSSA(t *types.Type) bool {
	// There are no un-SSAable types in test land.
	return true
}

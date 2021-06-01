// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gc

import (
	"cmd/compile/internal/base"
	"cmd/compile/internal/inline"
	"cmd/compile/internal/ir"
	"cmd/compile/internal/typecheck"
	"cmd/compile/internal/types"
	"cmd/internal/bio"
	"fmt"
	"go/constant"
)

func exportf(bout *bio.Writer, format string, args ...interface{}) {
	fmt.Fprintf(bout, format, args...)
	if base.Debug.Export != 0 {
		fmt.Printf(format, args...)
	}
}

func dumpexport(bout *bio.Writer) {
	p := &exporter{marked: make(map[*types.Type]bool)}
	for _, n := range typecheck.Target.Exports {
		p.markObject(n)
	}

	// The linker also looks for the $$ marker - use char after $$ to distinguish format.
	exportf(bout, "\n$$B\n") // indicate binary export format
	off := bout.Offset()
	typecheck.WriteExports(bout.Writer)
	size := bout.Offset() - off
	exportf(bout, "\n$$\n")

	if base.Debug.Export != 0 {
		fmt.Printf("BenchmarkExportSize:%s 1 %d bytes\n", base.Ctxt.Pkgpath, size)
	}
}

func dumpasmhdr() {
	b, err := bio.Create(base.Flag.AsmHdr)
	if err != nil {
		base.Fatalf("%v", err)
	}
	fmt.Fprintf(b, "// generated by compile -asmhdr from package %s\n\n", types.LocalPkg.Name)
	for _, n := range typecheck.Target.Asms {
		if n.Sym().IsBlank() {
			continue
		}
		switch n.Op() {
		case ir.OLITERAL:
			t := n.Val().Kind()
			if t == constant.Float || t == constant.Complex {
				break
			}
			fmt.Fprintf(b, "#define const_%s %#v\n", n.Sym().Name, n.Val())

		case ir.OTYPE:
			t := n.Type()
			if !t.IsStruct() || t.StructType().Map != nil || t.IsFuncArgStruct() {
				break
			}
			fmt.Fprintf(b, "#define %s__size %d\n", n.Sym().Name, int(t.Width))
			for _, f := range t.Fields().Slice() {
				if !f.Sym.IsBlank() {
					fmt.Fprintf(b, "#define %s_%s %d\n", n.Sym().Name, f.Sym.Name, int(f.Offset))
				}
			}
		}
	}

	b.Close()
}

type exporter struct {
	marked map[*types.Type]bool // types already seen by markType
}

// markObject visits a reachable object.
func (p *exporter) markObject(n ir.Node) {
	if n.Op() == ir.ONAME {
		n := n.(*ir.Name)
		if n.Class == ir.PFUNC {
			inline.Inline_Flood(n, typecheck.Export)
		}
	}

	p.markType(n.Type())
}

// markType recursively visits types reachable from t to identify
// functions whose inline bodies may be needed.
func (p *exporter) markType(t *types.Type) {
	if t.IsInstantiatedGeneric() {
		// Re-instantiated types don't add anything new, so don't follow them.
		return
	}
	if p.marked[t] {
		return
	}
	p.marked[t] = true

	// If this is a named type, mark all of its associated
	// methods. Skip interface types because t.Methods contains
	// only their unexpanded method set (i.e., exclusive of
	// interface embeddings), and the switch statement below
	// handles their full method set.
	if t.Sym() != nil && t.Kind() != types.TINTER {
		for _, m := range t.Methods().Slice() {
			if types.IsExported(m.Sym.Name) {
				p.markObject(ir.AsNode(m.Nname))
			}
		}
	}

	// Recursively mark any types that can be produced given a
	// value of type t: dereferencing a pointer; indexing or
	// iterating over an array, slice, or map; receiving from a
	// channel; accessing a struct field or interface method; or
	// calling a function.
	//
	// Notably, we don't mark function parameter types, because
	// the user already needs some way to construct values of
	// those types.
	switch t.Kind() {
	case types.TPTR, types.TARRAY, types.TSLICE:
		p.markType(t.Elem())

	case types.TCHAN:
		if t.ChanDir().CanRecv() {
			p.markType(t.Elem())
		}

	case types.TMAP:
		p.markType(t.Key())
		p.markType(t.Elem())

	case types.TSTRUCT:
		for _, f := range t.FieldSlice() {
			if types.IsExported(f.Sym.Name) || f.Embedded != 0 {
				p.markType(f.Type)
			}
		}

	case types.TFUNC:
		for _, f := range t.Results().FieldSlice() {
			p.markType(f.Type)
		}

	case types.TINTER:
		// TODO(danscales) - will have to deal with the types in interface
		// elements here when implemented in types2 and represented in types1.
		for _, f := range t.AllMethods().Slice() {
			if types.IsExported(f.Sym.Name) {
				p.markType(f.Type)
			}
		}

	case types.TTYPEPARAM:
		// No other type that needs to be followed.
	}
}

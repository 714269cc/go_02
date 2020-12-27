// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"cmd/compile/internal/base"
	"cmd/compile/internal/ir"
	"cmd/compile/internal/typecheck"
	"cmd/compile/internal/types"
)

func walkSelect(sel *ir.SelectStmt) {
	lno := ir.SetPos(sel)
	if len(sel.Compiled) != 0 {
		base.Fatalf("double walkselect")
	}

	init := sel.Init()
	sel.PtrInit().Set(nil)

	init = append(init, walkSelectCases(sel.Cases)...)
	sel.Cases = nil

	sel.Compiled.Set(init)
	walkStmtList(sel.Compiled)

	base.Pos = lno
}

func walkSelectCases(cases []*ir.CommStmt) []ir.Node {
	ncas := len(cases)
	sellineno := base.Pos

	// optimization: zero-case select
	if ncas == 0 {
		return []ir.Node{mkcall("block", nil, nil)}
	}

	// optimization: one-case select: single op.
	if ncas == 1 {
		cas := cases[0]
		ir.SetPos(cas)
		l := cas.Init()
		if cas.Comm != nil { // not default:
			n := cas.Comm
			l = append(l, n.Init()...)
			n.PtrInit().Set(nil)
			switch n.Op() {
			default:
				base.Fatalf("select %v", n.Op())

			case ir.OSEND:
				// already ok

			case ir.OSELRECV2:
				r := n.(*ir.AssignListStmt)
				if ir.IsBlank(r.Lhs[0]) && ir.IsBlank(r.Lhs[1]) {
					n = r.Rhs[0]
					break
				}
				r.SetOp(ir.OAS2RECV)
			}

			l = append(l, n)
		}

		l = append(l, cas.Body...)
		l = append(l, ir.NewBranchStmt(base.Pos, ir.OBREAK, nil))
		return l
	}

	// convert case value arguments to addresses.
	// this rewrite is used by both the general code and the next optimization.
	var dflt *ir.CommStmt
	for _, cas := range cases {
		ir.SetPos(cas)
		n := cas.Comm
		if n == nil {
			dflt = cas
			continue
		}
		switch n.Op() {
		case ir.OSEND:
			n := n.(*ir.SendStmt)
			n.Value = typecheck.NodAddr(n.Value)
			n.Value = typecheck.Expr(n.Value)

		case ir.OSELRECV2:
			n := n.(*ir.AssignListStmt)
			if !ir.IsBlank(n.Lhs[0]) {
				n.Lhs[0] = typecheck.NodAddr(n.Lhs[0])
				n.Lhs[0] = typecheck.Expr(n.Lhs[0])
			}
		}
	}

	// optimization: two-case select but one is default: single non-blocking op.
	if ncas == 2 && dflt != nil {
		cas := cases[0]
		if cas == dflt {
			cas = cases[1]
		}

		n := cas.Comm
		ir.SetPos(n)
		r := ir.NewIfStmt(base.Pos, nil, nil, nil)
		r.PtrInit().Set(cas.Init())
		var call ir.Node
		switch n.Op() {
		default:
			base.Fatalf("select %v", n.Op())

		case ir.OSEND:
			// if selectnbsend(c, v) { body } else { default body }
			n := n.(*ir.SendStmt)
			ch := n.Chan
			call = mkcall1(chanfn("selectnbsend", 2, ch.Type()), types.Types[types.TBOOL], r.PtrInit(), ch, n.Value)

		case ir.OSELRECV2:
			n := n.(*ir.AssignListStmt)
			recv := n.Rhs[0].(*ir.UnaryExpr)
			ch := recv.X
			elem := n.Lhs[0]
			if ir.IsBlank(elem) {
				elem = typecheck.NodNil()
			}
			if ir.IsBlank(n.Lhs[1]) {
				// if selectnbrecv(&v, c) { body } else { default body }
				call = mkcall1(chanfn("selectnbrecv", 2, ch.Type()), types.Types[types.TBOOL], r.PtrInit(), elem, ch)
			} else {
				// TODO(cuonglm): make this use selectnbrecv()
				// if selectnbrecv2(&v, &received, c) { body } else { default body }
				receivedp := typecheck.Expr(typecheck.NodAddr(n.Lhs[1]))
				call = mkcall1(chanfn("selectnbrecv2", 2, ch.Type()), types.Types[types.TBOOL], r.PtrInit(), elem, receivedp, ch)
			}
		}

		r.Cond = typecheck.Expr(call)
		r.Body.Set(cas.Body)
		r.Else.Set(append(dflt.Init(), dflt.Body...))
		return []ir.Node{r, ir.NewBranchStmt(base.Pos, ir.OBREAK, nil)}
	}

	if dflt != nil {
		ncas--
	}
	casorder := make([]*ir.CommStmt, ncas)
	nsends, nrecvs := 0, 0

	var init []ir.Node

	// generate sel-struct
	base.Pos = sellineno
	selv := typecheck.Temp(types.NewArray(scasetype(), int64(ncas)))
	init = append(init, typecheck.Stmt(ir.NewAssignStmt(base.Pos, selv, nil)))

	// No initialization for order; runtime.selectgo is responsible for that.
	order := typecheck.Temp(types.NewArray(types.Types[types.TUINT16], 2*int64(ncas)))

	var pc0, pcs ir.Node
	if base.Flag.Race {
		pcs = typecheck.Temp(types.NewArray(types.Types[types.TUINTPTR], int64(ncas)))
		pc0 = typecheck.Expr(typecheck.NodAddr(ir.NewIndexExpr(base.Pos, pcs, ir.NewInt(0))))
	} else {
		pc0 = typecheck.NodNil()
	}

	// register cases
	for _, cas := range cases {
		ir.SetPos(cas)

		init = append(init, cas.Init()...)
		cas.PtrInit().Set(nil)

		n := cas.Comm
		if n == nil { // default:
			continue
		}

		var i int
		var c, elem ir.Node
		switch n.Op() {
		default:
			base.Fatalf("select %v", n.Op())
		case ir.OSEND:
			n := n.(*ir.SendStmt)
			i = nsends
			nsends++
			c = n.Chan
			elem = n.Value
		case ir.OSELRECV2:
			n := n.(*ir.AssignListStmt)
			nrecvs++
			i = ncas - nrecvs
			recv := n.Rhs[0].(*ir.UnaryExpr)
			c = recv.X
			elem = n.Lhs[0]
		}

		casorder[i] = cas

		setField := func(f string, val ir.Node) {
			r := ir.NewAssignStmt(base.Pos, ir.NewSelectorExpr(base.Pos, ir.ODOT, ir.NewIndexExpr(base.Pos, selv, ir.NewInt(int64(i))), typecheck.Lookup(f)), val)
			init = append(init, typecheck.Stmt(r))
		}

		c = typecheck.ConvNop(c, types.Types[types.TUNSAFEPTR])
		setField("c", c)
		if !ir.IsBlank(elem) {
			elem = typecheck.ConvNop(elem, types.Types[types.TUNSAFEPTR])
			setField("elem", elem)
		}

		// TODO(mdempsky): There should be a cleaner way to
		// handle this.
		if base.Flag.Race {
			r := mkcall("selectsetpc", nil, nil, typecheck.NodAddr(ir.NewIndexExpr(base.Pos, pcs, ir.NewInt(int64(i)))))
			init = append(init, r)
		}
	}
	if nsends+nrecvs != ncas {
		base.Fatalf("walkselectcases: miscount: %v + %v != %v", nsends, nrecvs, ncas)
	}

	// run the select
	base.Pos = sellineno
	chosen := typecheck.Temp(types.Types[types.TINT])
	recvOK := typecheck.Temp(types.Types[types.TBOOL])
	r := ir.NewAssignListStmt(base.Pos, ir.OAS2, nil, nil)
	r.Lhs = []ir.Node{chosen, recvOK}
	fn := typecheck.LookupRuntime("selectgo")
	r.Rhs = []ir.Node{mkcall1(fn, fn.Type().Results(), nil, bytePtrToIndex(selv, 0), bytePtrToIndex(order, 0), pc0, ir.NewInt(int64(nsends)), ir.NewInt(int64(nrecvs)), ir.NewBool(dflt == nil))}
	init = append(init, typecheck.Stmt(r))

	// selv and order are no longer alive after selectgo.
	init = append(init, ir.NewUnaryExpr(base.Pos, ir.OVARKILL, selv))
	init = append(init, ir.NewUnaryExpr(base.Pos, ir.OVARKILL, order))
	if base.Flag.Race {
		init = append(init, ir.NewUnaryExpr(base.Pos, ir.OVARKILL, pcs))
	}

	// dispatch cases
	dispatch := func(cond ir.Node, cas *ir.CommStmt) {
		cond = typecheck.Expr(cond)
		cond = typecheck.DefaultLit(cond, nil)

		r := ir.NewIfStmt(base.Pos, cond, nil, nil)

		if n := cas.Comm; n != nil && n.Op() == ir.OSELRECV2 {
			n := n.(*ir.AssignListStmt)
			if !ir.IsBlank(n.Lhs[1]) {
				x := ir.NewAssignStmt(base.Pos, n.Lhs[1], recvOK)
				r.Body.Append(typecheck.Stmt(x))
			}
		}

		r.Body.Append(cas.Body.Take()...)
		r.Body.Append(ir.NewBranchStmt(base.Pos, ir.OBREAK, nil))
		init = append(init, r)
	}

	if dflt != nil {
		ir.SetPos(dflt)
		dispatch(ir.NewBinaryExpr(base.Pos, ir.OLT, chosen, ir.NewInt(0)), dflt)
	}
	for i, cas := range casorder {
		ir.SetPos(cas)
		dispatch(ir.NewBinaryExpr(base.Pos, ir.OEQ, chosen, ir.NewInt(int64(i))), cas)
	}

	return init
}

// bytePtrToIndex returns a Node representing "(*byte)(&n[i])".
func bytePtrToIndex(n ir.Node, i int64) ir.Node {
	s := typecheck.NodAddr(ir.NewIndexExpr(base.Pos, n, ir.NewInt(i)))
	t := types.NewPtr(types.Types[types.TUINT8])
	return typecheck.ConvNop(s, t)
}

var scase *types.Type

// Keep in sync with src/runtime/select.go.
func scasetype() *types.Type {
	if scase == nil {
		scase = types.NewStruct(types.NoPkg, []*types.Field{
			types.NewField(base.Pos, typecheck.Lookup("c"), types.Types[types.TUNSAFEPTR]),
			types.NewField(base.Pos, typecheck.Lookup("elem"), types.Types[types.TUNSAFEPTR]),
		})
		scase.SetNoalg(true)
	}
	return scase
}

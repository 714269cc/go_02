// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gc

import (
	"log"

	"cmd/internal/obj"
	"cmd/internal/obj/x86" // TODO: remove
	"cmd/internal/ssa"
)

func buildssa(fn *Node) *ssa.Func {
	dumplist("buildssa", Curfn.Nbody)

	var s ssaState

	// TODO(khr): build config just once at the start of the compiler binary
	s.config = ssa.NewConfig(Thearch.Thestring)
	s.f = s.config.NewFunc()
	s.f.Name = fn.Nname.Sym.Name

	// We construct SSA using an algorithm similar to
	// Brau, Buchwald, Hack, Leißa, Mallon, and Zwinkau
	// http://pp.info.uni-karlsruhe.de/uploads/publikationen/braun13cc.pdf
	// TODO: check this comment

	// Allocate starting block
	s.f.Entry = s.f.NewBlock(ssa.BlockPlain)

	// Allocate exit block
	s.exit = s.f.NewBlock(ssa.BlockExit)

	// TODO(khr): all args.  Make a struct containing args/returnvals, declare
	// an FP which contains a pointer to that struct.

	s.vars = map[string]*ssa.Value{}
	s.labels = map[string]*ssa.Block{}
	s.argOffsets = map[string]int64{}

	// Convert the AST-based IR to the SSA-based IR
	s.startBlock(s.f.Entry)
	s.stmtList(fn.Nbody)

	// Finish up exit block
	s.startBlock(s.exit)
	s.exit.Control = s.mem()
	s.endBlock()

	// Link up variable uses to variable definitions
	s.linkForwardReferences()

	// Main call to ssa package to compile function
	ssa.Compile(s.f)

	return s.f
}

type ssaState struct {
	// configuration (arch) information
	config *ssa.Config

	// function we're building
	f *ssa.Func

	// exit block that "return" jumps to (and panics jump to)
	exit *ssa.Block

	// the target block for each label in f
	labels map[string]*ssa.Block

	// current location where we're interpreting the AST
	curBlock *ssa.Block

	// variable assignments in the current block (map from variable name to ssa value)
	vars map[string]*ssa.Value

	// all defined variables at the end of each block.  Indexed by block ID.
	defvars []map[string]*ssa.Value

	// offsets of argument slots
	// unnamed and unused args are not listed.
	argOffsets map[string]int64
}

// startBlock sets the current block we're generating code in to b.
func (s *ssaState) startBlock(b *ssa.Block) {
	s.curBlock = b
	s.vars = map[string]*ssa.Value{}
}

// endBlock marks the end of generating code for the current block.
// Returns the (former) current block.  Returns nil if there is no current
// block, i.e. if no code flows to the current execution point.
func (s *ssaState) endBlock() *ssa.Block {
	b := s.curBlock
	if b == nil {
		return nil
	}
	for len(s.defvars) <= int(b.ID) {
		s.defvars = append(s.defvars, nil)
	}
	s.defvars[b.ID] = s.vars
	s.curBlock = nil
	s.vars = nil
	return b
}

// ssaStmtList converts the statement n to SSA and adds it to s.
func (s *ssaState) stmtList(l *NodeList) {
	for ; l != nil; l = l.Next {
		s.stmt(l.N)
	}
}

// ssaStmt converts the statement n to SSA and adds it to s.
func (s *ssaState) stmt(n *Node) {
	s.stmtList(n.Ninit)
	switch n.Op {

	case OBLOCK:
		s.stmtList(n.List)

	case ODCL:
		// TODO: ???  Assign 0?

	case OLABEL, OGOTO:
		// get block at label, or make one
		t := s.labels[n.Left.Sym.Name]
		if t == nil {
			t = s.f.NewBlock(ssa.BlockPlain)
			s.labels[n.Left.Sym.Name] = t
		}
		// go to that label (we pretend "label:" is preceded by "goto label")
		b := s.endBlock()
		addEdge(b, t)

		if n.Op == OLABEL {
			// next we work on the label's target block
			s.startBlock(t)
		}

	case OAS:
		// TODO(khr): colas?
		val := s.expr(n.Right)
		if n.Left.Op == OINDREG {
			// indirect off a register (TODO: always SP?)
			// used for storing arguments to callees
			addr := s.f.Entry.NewValue(ssa.OpSPAddr, Ptrto(n.Right.Type), n.Left.Xoffset)
			s.vars[".mem"] = s.curBlock.NewValue3(ssa.OpStore, ssa.TypeMem, nil, addr, val, s.mem())
		} else if n.Left.Op != ONAME {
			// some more complicated expression.  Rewrite to a store.  TODO
			addr := s.expr(n.Left) // TODO: wrap in &

			// TODO(khr): nil check
			s.vars[".mem"] = s.curBlock.NewValue3(ssa.OpStore, n.Right.Type, nil, addr, val, s.mem())
		} else if !n.Left.Addable {
			// TODO
			log.Fatalf("assignment to non-addable value")
		} else if n.Left.Class&PHEAP != 0 {
			// TODO
			log.Fatalf("assignment to heap value")
		} else if n.Left.Class == PEXTERN {
			// assign to global variable
			addr := s.f.Entry.NewValue(ssa.OpGlobal, Ptrto(n.Left.Type), n.Left.Sym)
			s.vars[".mem"] = s.curBlock.NewValue3(ssa.OpStore, ssa.TypeMem, nil, addr, val, s.mem())
		} else if n.Left.Class == PPARAMOUT {
			// store to parameter slot
			addr := s.f.Entry.NewValue(ssa.OpFPAddr, Ptrto(n.Right.Type), n.Left.Xoffset)
			s.vars[".mem"] = s.curBlock.NewValue3(ssa.OpStore, ssa.TypeMem, nil, addr, val, s.mem())
		} else {
			// normal variable
			s.vars[n.Left.Sym.Name] = val
		}
	case OIF:
		cond := s.expr(n.Ntest)
		b := s.endBlock()
		b.Kind = ssa.BlockIf
		b.Control = cond
		// TODO(khr): likely direction

		bThen := s.f.NewBlock(ssa.BlockPlain)
		bEnd := s.f.NewBlock(ssa.BlockPlain)
		var bElse *ssa.Block

		if n.Nelse == nil {
			addEdge(b, bThen)
			addEdge(b, bEnd)
		} else {
			bElse = s.f.NewBlock(ssa.BlockPlain)
			addEdge(b, bThen)
			addEdge(b, bElse)
		}

		s.startBlock(bThen)
		s.stmtList(n.Nbody)
		b = s.endBlock()
		if b != nil {
			addEdge(b, bEnd)
		}

		if n.Nelse != nil {
			s.startBlock(bElse)
			s.stmtList(n.Nelse)
			b = s.endBlock()
			if b != nil {
				addEdge(b, bEnd)
			}
		}
		s.startBlock(bEnd)

	case ORETURN:
		s.stmtList(n.List)
		b := s.endBlock()
		addEdge(b, s.exit)

	case OFOR:
		bCond := s.f.NewBlock(ssa.BlockPlain)
		bBody := s.f.NewBlock(ssa.BlockPlain)
		bEnd := s.f.NewBlock(ssa.BlockPlain)

		// first, jump to condition test
		b := s.endBlock()
		addEdge(b, bCond)

		// generate code to test condition
		// TODO(khr): Ntest == nil exception
		s.startBlock(bCond)
		cond := s.expr(n.Ntest)
		b = s.endBlock()
		b.Kind = ssa.BlockIf
		b.Control = cond
		// TODO(khr): likely direction
		addEdge(b, bBody)
		addEdge(b, bEnd)

		// generate body
		s.startBlock(bBody)
		s.stmtList(n.Nbody)
		s.stmt(n.Nincr)
		b = s.endBlock()
		addEdge(b, bCond)

		s.startBlock(bEnd)

	case OVARKILL:
		// TODO(khr): ??? anything to do here?  Only for addrtaken variables?
		// Maybe just link it in the store chain?
	default:
		log.Fatalf("unhandled stmt %s", opnames[n.Op])
	}
}

// expr converts the expression n to ssa, adds it to s and returns the ssa result.
func (s *ssaState) expr(n *Node) *ssa.Value {
	if n == nil {
		// TODO(khr): is this nil???
		return s.f.Entry.NewValue(ssa.OpConst, n.Type, nil)
	}
	switch n.Op {
	case ONAME:
		// TODO: remember offsets for PPARAM names
		if n.Class == PEXTERN {
			// global variable
			addr := s.f.Entry.NewValue(ssa.OpGlobal, Ptrto(n.Type), n.Sym)
			return s.curBlock.NewValue2(ssa.OpLoad, n.Type, nil, addr, s.mem())
		}
		s.argOffsets[n.Sym.Name] = n.Xoffset
		return s.variable(n.Sym.Name, n.Type)
		// binary ops
	case OLITERAL:
		switch n.Val.Ctype {
		case CTINT:
			return s.f.ConstInt(n.Type, Mpgetfix(n.Val.U.Xval))
		default:
			log.Fatalf("unhandled OLITERAL %v", n.Val.Ctype)
			return nil
		}
	case OLT:
		a := s.expr(n.Left)
		b := s.expr(n.Right)
		return s.curBlock.NewValue2(ssa.OpLess, ssa.TypeBool, nil, a, b)
	case OADD:
		a := s.expr(n.Left)
		b := s.expr(n.Right)
		return s.curBlock.NewValue2(ssa.OpAdd, a.Type, nil, a, b)

	case OSUB:
		// TODO:(khr) fold code for all binary ops together somehow
		a := s.expr(n.Left)
		b := s.expr(n.Right)
		return s.curBlock.NewValue2(ssa.OpSub, a.Type, nil, a, b)

	case OIND:
		p := s.expr(n.Left)
		c := s.curBlock.NewValue1(ssa.OpIsNonNil, ssa.TypeBool, nil, p)
		b := s.endBlock()
		b.Kind = ssa.BlockIf
		b.Control = c
		bNext := s.f.NewBlock(ssa.BlockPlain)
		addEdge(b, bNext)
		addEdge(b, s.exit)
		s.startBlock(bNext)
		// TODO(khr): if ptr check fails, don't go directly to exit.
		// Instead, go to a call to panicnil or something.
		// TODO: implicit nil checks somehow?

		return s.curBlock.NewValue2(ssa.OpLoad, n.Type, nil, p, s.mem())
	case ODOTPTR:
		p := s.expr(n.Left)
		// TODO: nilcheck
		p = s.curBlock.NewValue2(ssa.OpAdd, p.Type, nil, p, s.f.ConstInt(s.config.UIntPtr, n.Xoffset))
		return s.curBlock.NewValue2(ssa.OpLoad, n.Type, nil, p, s.mem())

	case OINDEX:
		// TODO: slice vs array?  Map index is already reduced to a function call
		a := s.expr(n.Left)
		i := s.expr(n.Right)
		// convert index to full width
		// TODO: if index is 64-bit and we're compiling to 32-bit, check that high
		// 32 bits are zero (and use a low32 op instead of convnop here).
		i = s.curBlock.NewValue1(ssa.OpConvNop, s.config.UIntPtr, nil, i)

		// bounds check
		len := s.curBlock.NewValue1(ssa.OpSliceLen, s.config.UIntPtr, nil, a)
		cmp := s.curBlock.NewValue2(ssa.OpIsInBounds, ssa.TypeBool, nil, i, len)
		b := s.endBlock()
		b.Kind = ssa.BlockIf
		b.Control = cmp
		bNext := s.f.NewBlock(ssa.BlockPlain)
		addEdge(b, bNext)
		addEdge(b, s.exit)
		s.startBlock(bNext)
		// TODO: don't go directly to s.exit.  Go to a stub that calls panicindex first.

		return s.curBlock.NewValue3(ssa.OpSliceIndex, n.Left.Type.Type, nil, a, i, s.mem())

	case OCALLFUNC:
		// run all argument assignments
		// TODO(khr): do we need to evaluate function first?
		// Or is it already side-effect-free and does not require a call?
		s.stmtList(n.List)

		if n.Left.Op != ONAME {
			// TODO(khr): closure calls?
			log.Fatalf("can't handle CALLFUNC with non-ONAME fn %s", opnames[n.Left.Op])
		}
		bNext := s.f.NewBlock(ssa.BlockPlain)
		call := s.curBlock.NewValue1(ssa.OpStaticCall, ssa.TypeMem, n.Left.Sym, s.mem())
		b := s.endBlock()
		b.Kind = ssa.BlockCall
		b.Control = call
		addEdge(b, bNext)
		addEdge(b, s.exit)

		// read result from stack at the start of the fallthrough block
		s.startBlock(bNext)
		var titer Iter
		fp := Structfirst(&titer, Getoutarg(n.Left.Type))
		a := s.f.Entry.NewValue(ssa.OpSPAddr, Ptrto(fp.Type), fp.Width)
		return s.curBlock.NewValue2(ssa.OpLoad, fp.Type, nil, a, call)
	default:
		log.Fatalf("unhandled expr %s", opnames[n.Op])
		return nil
	}
}

// variable returns the value of a variable at the current location.
func (s *ssaState) variable(name string, t ssa.Type) *ssa.Value {
	if s.curBlock == nil {
		log.Fatalf("nil curblock!")
	}
	v := s.vars[name]
	if v == nil {
		// TODO: get type?  Take Sym as arg?
		v = s.curBlock.NewValue(ssa.OpFwdRef, t, name)
		s.vars[name] = v
	}
	return v
}

func (s *ssaState) mem() *ssa.Value {
	return s.variable(".mem", ssa.TypeMem)
}

func (s *ssaState) linkForwardReferences() {
	// Build ssa graph.  Each variable on its first use in a basic block
	// leaves a FwdRef in that block representing the incoming value
	// of that variable.  This function links that ref up with possible definitions,
	// inserting Phi values as needed.  This is essentially the algorithm
	// described by Brau, Buchwald, Hack, Leißa, Mallon, and Zwinkau:
	// http://pp.info.uni-karlsruhe.de/uploads/publikationen/braun13cc.pdf
	for _, b := range s.f.Blocks {
		for _, v := range b.Values {
			if v.Op != ssa.OpFwdRef {
				continue
			}
			name := v.Aux.(string)
			v.Op = ssa.OpCopy
			v.Aux = nil
			v.SetArgs1(s.lookupVarIncoming(b, v.Type, name))
		}
	}
}

// lookupVarIncoming finds the variable's value at the start of block b.
func (s *ssaState) lookupVarIncoming(b *ssa.Block, t ssa.Type, name string) *ssa.Value {
	// TODO(khr): have lookupVarIncoming overwrite the fwdRef or copy it
	// will be used in, instead of having the result used in a copy value.
	if b == s.f.Entry {
		if name == ".mem" {
			return b.NewValue(ssa.OpArg, t, name)
		}
		// variable is live at the entry block.  Load it.
		a := s.f.Entry.NewValue(ssa.OpFPAddr, Ptrto(t.(*Type)), s.argOffsets[name])
		m := b.NewValue(ssa.OpArg, ssa.TypeMem, ".mem") // TODO: reuse mem starting value
		return b.NewValue2(ssa.OpLoad, t, nil, a, m)
	}
	var vals []*ssa.Value
	for _, p := range b.Preds {
		vals = append(vals, s.lookupVarOutgoing(p, t, name))
	}
	v0 := vals[0]
	for i := 1; i < len(vals); i++ {
		if vals[i] != v0 {
			// need a phi value
			v := b.NewValue(ssa.OpPhi, t, nil)
			v.AddArgs(vals...)
			return v
		}
	}
	return v0
}

// lookupVarOutgoing finds the variable's value at the end of block b.
func (s *ssaState) lookupVarOutgoing(b *ssa.Block, t ssa.Type, name string) *ssa.Value {
	m := s.defvars[b.ID]
	if v, ok := m[name]; ok {
		return v
	}
	// The variable is not defined by b and we haven't
	// looked it up yet.  Generate v, a copy value which
	// will be the outgoing value of the variable.  Then
	// look up w, the incoming value of the variable.
	// Make v = copy(w).  We need the extra copy to
	// prevent infinite recursion when looking up the
	// incoming value of the variable.
	v := b.NewValue(ssa.OpCopy, t, nil)
	m[name] = v
	v.AddArg(s.lookupVarIncoming(b, t, name))
	return v
}

// TODO: the above mutually recursive functions can lead to very deep stacks.  Fix that.

// addEdge adds an edge from b to c.
func addEdge(b, c *ssa.Block) {
	b.Succs = append(b.Succs, c)
	c.Preds = append(c.Preds, b)
}

// an unresolved branch
type branch struct {
	p *obj.Prog  // branch instruction
	b *ssa.Block // target
}

// genssa appends entries to ptxt for each instruction in f.
// gcargs and gclocals are filled in with pointer maps for the frame.
func genssa(f *ssa.Func, ptxt *obj.Prog, gcargs, gclocals *Sym) {
	// TODO: line numbers
	// TODO: layout frame
	stkSize := int64(64)

	if Hasdefer != 0 {
		// deferreturn pretends to have one uintptr argument.
		// Reserve space for it so stack scanner is happy.
		if Maxarg < int64(Widthptr) {
			Maxarg = int64(Widthptr)
		}
	}
	if stkSize+Maxarg > 1<<31 {
		Yyerror("stack frame too large (>2GB)")
		return
	}
	frameSize := stkSize + Maxarg

	ptxt.To.Type = obj.TYPE_TEXTSIZE
	ptxt.To.Val = int32(Rnd(Curfn.Type.Argwid, int64(Widthptr))) // arg size
	ptxt.To.Offset = frameSize - 8                               // TODO: arch-dependent

	// Remember where each block starts.
	bstart := make([]*obj.Prog, f.NumBlocks())

	// Remember all the branch instructions we've seen
	// and where they would like to go
	var branches []branch

	// Emit basic blocks
	for i, b := range f.Blocks {
		bstart[b.ID] = Pc
		// Emit values in block
		for _, v := range b.Values {
			genValue(v, frameSize)
		}
		// Emit control flow instructions for block
		var next *ssa.Block
		if i < len(f.Blocks)-1 {
			next = f.Blocks[i+1]
		}
		branches = genBlock(b, next, branches)
	}

	// Resolve branches
	for _, br := range branches {
		br.p.To.Val = bstart[br.b.ID]
	}

	Pc.As = obj.ARET // overwrite AEND

	// TODO: liveness
	// TODO: gcargs
	// TODO: gclocals

	// TODO: dump frame if -f

	// Emit garbage collection symbols.  TODO: put something in them
	liveness(Curfn, ptxt, gcargs, gclocals)
}

func genValue(v *ssa.Value, frameSize int64) {
	switch v.Op {
	case ssa.OpADDQ:
		// TODO: use addq instead of leaq if target is in the right register.
		p := Prog(x86.ALEAQ)
		p.From.Type = obj.TYPE_MEM
		p.From.Reg = regnum(v.Args[0])
		p.From.Scale = 1
		p.From.Index = regnum(v.Args[1])
		p.To.Type = obj.TYPE_REG
		p.To.Reg = regnum(v)
	case ssa.OpADDCQ:
		// TODO: use addq instead of leaq if target is in the right register.
		p := Prog(x86.ALEAQ)
		p.From.Type = obj.TYPE_MEM
		p.From.Reg = regnum(v.Args[0])
		p.From.Offset = v.Aux.(int64)
		p.To.Type = obj.TYPE_REG
		p.To.Reg = regnum(v)
	case ssa.OpSUBCQ:
		// This code compensates for the fact that the register allocator
		// doesn't understand 2-address instructions yet.  TODO: fix that.
		x := regnum(v.Args[0])
		r := regnum(v)
		if x != r {
			p := Prog(x86.AMOVQ)
			p.From.Type = obj.TYPE_REG
			p.From.Reg = x
			p.To.Type = obj.TYPE_REG
			p.To.Reg = r
			x = r
		}
		p := Prog(x86.ASUBQ)
		p.From.Type = obj.TYPE_CONST
		p.From.Offset = v.Aux.(int64)
		p.To.Type = obj.TYPE_REG
		p.To.Reg = r
	case ssa.OpCMPQ:
		x := regnum(v.Args[0])
		y := regnum(v.Args[1])
		p := Prog(x86.ACMPQ)
		p.From.Type = obj.TYPE_REG
		p.From.Reg = x
		p.To.Type = obj.TYPE_REG
		p.To.Reg = y
	case ssa.OpMOVQconst:
		x := regnum(v)
		p := Prog(x86.AMOVQ)
		p.From.Type = obj.TYPE_CONST
		p.From.Offset = v.Aux.(int64)
		p.To.Type = obj.TYPE_REG
		p.To.Reg = x
	case ssa.OpMOVQloadFP:
		x := regnum(v)
		p := Prog(x86.AMOVQ)
		p.From.Type = obj.TYPE_MEM
		p.From.Reg = x86.REG_SP
		p.From.Offset = v.Aux.(int64) + frameSize
		p.To.Type = obj.TYPE_REG
		p.To.Reg = x
	case ssa.OpMOVQstoreFP:
		x := regnum(v.Args[0])
		p := Prog(x86.AMOVQ)
		p.From.Type = obj.TYPE_REG
		p.From.Reg = x
		p.To.Type = obj.TYPE_MEM
		p.To.Reg = x86.REG_SP
		p.To.Offset = v.Aux.(int64) + frameSize
	case ssa.OpCopy:
		x := regnum(v.Args[0])
		y := regnum(v)
		if x != y {
			p := Prog(x86.AMOVQ)
			p.From.Type = obj.TYPE_REG
			p.From.Reg = x
			p.To.Type = obj.TYPE_REG
			p.To.Reg = y
		}
	case ssa.OpLoadReg8:
		p := Prog(x86.AMOVQ)
		p.From.Type = obj.TYPE_MEM
		p.From.Reg = x86.REG_SP
		p.From.Offset = frameSize - localOffset(v.Args[0])
		p.To.Type = obj.TYPE_REG
		p.To.Reg = regnum(v)
	case ssa.OpStoreReg8:
		p := Prog(x86.AMOVQ)
		p.From.Type = obj.TYPE_REG
		p.From.Reg = regnum(v.Args[0])
		p.To.Type = obj.TYPE_MEM
		p.To.Reg = x86.REG_SP
		p.To.Offset = frameSize - localOffset(v)
	case ssa.OpPhi:
		// just check to make sure regalloc did it right
		f := v.Block.Func
		loc := f.RegAlloc[v.ID]
		for _, a := range v.Args {
			if f.RegAlloc[a.ID] != loc { // TODO: .Equal() instead?
				log.Fatalf("phi arg at different location than phi %v %v %v %v", v, loc, a, f.RegAlloc[a.ID])
			}
		}
	case ssa.OpConst:
		if v.Block.Func.RegAlloc[v.ID] != nil {
			log.Fatalf("const value %v shouldn't have a location", v)
		}
	case ssa.OpArg:
		// memory arg needs no code
		// TODO: only mem arg goes here.
	default:
		log.Fatalf("value %v not implemented yet", v)
	}
}

func genBlock(b, next *ssa.Block, branches []branch) []branch {
	switch b.Kind {
	case ssa.BlockPlain:
		if b.Succs[0] != next {
			p := Prog(obj.AJMP)
			p.To.Type = obj.TYPE_BRANCH
			branches = append(branches, branch{p, b.Succs[0]})
		}
	case ssa.BlockExit:
		Prog(obj.ARET)
	case ssa.BlockLT:
		if b.Succs[0] == next {
			p := Prog(x86.AJGE)
			p.To.Type = obj.TYPE_BRANCH
			branches = append(branches, branch{p, b.Succs[1]})
		} else if b.Succs[1] == next {
			p := Prog(x86.AJLT)
			p.To.Type = obj.TYPE_BRANCH
			branches = append(branches, branch{p, b.Succs[0]})
		} else {
			p := Prog(x86.AJLT)
			p.To.Type = obj.TYPE_BRANCH
			branches = append(branches, branch{p, b.Succs[0]})
			q := Prog(obj.AJMP)
			q.To.Type = obj.TYPE_BRANCH
			branches = append(branches, branch{q, b.Succs[1]})
		}
	default:
		log.Fatalf("branch at %v not implemented yet", b)
	}
	return branches
}

// ssaRegToReg maps ssa register numbers to obj register numbers.
var ssaRegToReg = [...]int16{
	x86.REG_AX,
	x86.REG_CX,
	x86.REG_DX,
	x86.REG_BX,
	x86.REG_SP,
	x86.REG_BP,
	x86.REG_SI,
	x86.REG_DI,
	x86.REG_R8,
	x86.REG_R9,
	x86.REG_R10,
	x86.REG_R11,
	x86.REG_R12,
	x86.REG_R13,
	x86.REG_R14,
	x86.REG_R15,
	// TODO: more
	// TODO: arch-dependent
}

// regnum returns the register (in cmd/internal/obj numbering) to
// which v has been allocated.  Panics if v is not assigned to a
// register.
func regnum(v *ssa.Value) int16 {
	return ssaRegToReg[v.Block.Func.RegAlloc[v.ID].(*ssa.Register).Num]
}

// localOffset returns the offset below the frame pointer where
// a stack-allocated local has been allocated.  Panics if v
// is not assigned to a local slot.
func localOffset(v *ssa.Value) int64 {
	return v.Block.Func.RegAlloc[v.ID].(*ssa.LocalSlot).Idx
}

package fir

type Builder struct {
	ctx     *Context
	fn      *Function
	current *Block
}

func NewBuilder(ctx *Context, fn *Function) *Builder {
	return &Builder{
		ctx:     ctx,
		fn:      fn,
		current: nil,
	}
}

func (b *Builder) SetInsertionPoint(block *Block) {
	b.fn = block.Func
	b.current = block
}

func (b *Builder) NewBlock(name string) *Block {
	block := &Block{
		ID:   BlockID(len(b.fn.Blocks)),
		Name: name,
		Func: b.fn,

		Ins:  nil,
		Term: nil,
	}
	b.fn.Blocks = append(b.fn.Blocks, block)
	b.current = block
	return block
}

func (b *Builder) CurrentBlock() *Block {
	return b.current
}

func (b *Builder) CurrentFunction() *Function {
	return b.fn
}

func (b *Builder) nextValueId() valueId {
	return b.fn.getNextId()
}

func (b *Builder) InsertStackSlot(t Type, name string) Value {
	slot := &StackSlot{
		valueBase: valueBase{
			typ: b.ctx.GetPtrType(),
			id:  b.nextValueId(),
		},
		Name:      name,
		InnerType: t,
	}
	b.fn.Stackslots = append(b.fn.Stackslots, *slot)
	return slot
}

func (b *Builder) InsertLoad(t Type, val Value) Value {
	load := &LoadIns{
		instructionBase: instructionBase{
			parent: b.current,
		},
		valueBase: valueBase{
			typ: t,
			id:  b.nextValueId(),
		},
		T:    t,
		From: val,
	}
	b.current.insertIns(load)
	return load
}

func (b *Builder) InsertStore(val Value, addr Value) *StoreIns {
	store := &StoreIns{
		instructionBase: instructionBase{
			parent: b.current,
		},
		Val:  val,
		Addr: addr,
	}
	b.current.insertIns(store)
	return store
}

func (b *Builder) InsertElemPtr(addr Value, path ...ElemPtrPathElem) Value {
	elemPtr := &ElemPtrIns{
		instructionBase: instructionBase{
			parent: b.current,
		},
		valueBase: valueBase{
			typ: b.ctx.GetPtrType(),
			id:  b.nextValueId(),
		},
		Addr: addr,
		Path: path,
	}
	b.current.insertIns(elemPtr)
	return elemPtr
}

func (b *Builder) InsertBinaryOp(op BinaryOp, lhs Value, rhs Value) Value {
	binaryOp := &BinaryOpIns{
		instructionBase: instructionBase{
			parent: b.current,
		},
		valueBase: valueBase{
			typ: lhs.Type(),
			id:  b.nextValueId(),
		},
		Op:  op,
		Lhs: lhs,
		Rhs: rhs,
	}
	b.current.insertIns(binaryOp)
	return binaryOp
}

func (b *Builder) InsertUnaryOp(op UnaryOp, val Value) Value {
	unaryOp := &UnaryOpIns{
		instructionBase: instructionBase{
			parent: b.current,
		},
		valueBase: valueBase{
			typ: val.Type(),
			id:  b.nextValueId(),
		},
		Op:  op,
		Val: val,
	}
	b.current.insertIns(unaryOp)
	return unaryOp
}

func (b *Builder) InsertCmp(op CmpPredicate, lhs Value, rhs Value) Value {
	cmp := &CmpIns{
		instructionBase: instructionBase{
			parent: b.current,
		},
		valueBase: valueBase{
			typ: b.ctx.GetBoolType(),
			id:  b.nextValueId(),
		},
		Op:  op,
		Lhs: lhs,
		Rhs: rhs,
	}
	b.current.insertIns(cmp)
	return cmp
}

func (b *Builder) InsertCast(val Value, typ Type) Value {
	cast := &CastIns{
		instructionBase: instructionBase{
			parent: b.current,
		},
		valueBase: valueBase{
			typ: typ,
			id:  b.nextValueId(),
		},
		Val: val,
	}
	b.current.insertIns(cast)
	return cast
}

func (b *Builder) InsertCall(fn Value, ret Type, conv CallingConv, args ...Value) Value {
	call := &CallIns{
		instructionBase: instructionBase{
			parent: b.current,
		},
		valueBase: valueBase{
			typ: fn.Type(),
			id:  b.nextValueId(),
		},
		Proc:       fn,
		Conv:       conv,
		Args:       args,
		ReturnType: ret,
	}
	b.current.insertIns(call)
	return call
}

func (b *Builder) InsertBr(target *Block) *BrTerm {
	br := &BrTerm{
		instructionBase: instructionBase{
			parent: b.current,
		},
		Target: target,
	}
	b.current.insertTerm(br)
	return br
}

func (b *Builder) InsertCondBr(cond Value, trueTarget *Block, falseTarget *Block) *CondBrTerm {
	condBr := &CondBrTerm{
		instructionBase: instructionBase{
			parent: b.current,
		},
		Cond:  cond,
		True:  trueTarget,
		False: falseTarget,
	}
	b.current.insertTerm(condBr)
	return condBr
}

func (b *Builder) InsertRet(val Value) *RetTerm {
	ret := &RetTerm{
		instructionBase: instructionBase{
			parent: b.current,
		},
		Value: val,
	}
	b.current.insertTerm(ret)
	return ret
}

func (b *Builder) InsertUnreachable() *UnreachableTerm {
	unreachable := &UnreachableTerm{
		instructionBase: instructionBase{
			parent: b.current,
		},
	}
	b.current.insertTerm(unreachable)
	return unreachable
}

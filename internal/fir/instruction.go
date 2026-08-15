package fir

type instructionBase struct {
	parent *Block
}

func (i *instructionBase) Parent() *Block {
	return i.parent
}

type Instruction interface {
	Parent() *Block
}

type ValueInstruction interface {
	Instruction
	Value
}

// Instruction definitions

// Memory instructions

type LoadIns struct {
	instructionBase
	valueBase

	T    Type
	From Value
}

type StoreIns struct {
	instructionBase

	Addr Value
	Val  Value
}

type EpPathKind uint8

const (
	_ EpPathKind = iota

	EpArray  // Array Element
	EpStruct // Struct Field
	EpUnion  // Union Field
)

type ElemPtrIns struct {
	instructionBase
	valueBase

	Addr Value
	Path []ElemPtrPathElem
}

type ElemPtrPathElem struct {
	Kind  EpPathKind
	Index Value
}

// Binary operations

type BinaryOp uint8

const (
	_ BinaryOp = iota

	BinAdd  // Add
	BinSub  // Subtract
	BinMul  // Multiply
	BinSdiv // Signed Divide
	BinUdiv // Unsigned Divide
	BinSrem // Signed Remainder
	BinUrem // Unsigned Remainder

	BinFadd // Floating Point Add
	BinFsub // Floating Point Subtract
	BinFmul // Floating Point Multiply
	BinFdiv // Floating Point Divide
	BinFrem // Floating Point Remainder

	BinAnd  // Bitwise AND
	BinOr   // Bitwise OR
	BinXor  // Bitwise XOR
	BinShl  // Bitwise Shift Left
	BinLshr // Bitwise Shift Right
	BinAshr // Bitwise Shift Right (Arithmetic)
)

type BinaryOpIns struct {
	instructionBase
	valueBase

	Op  BinaryOp
	Lhs Value
	Rhs Value
}

// Unary operations

type UnaryOp uint8

const (
	_ UnaryOp = iota

	Neg  // Negate
	Fneg // Floating Point Negate
	Not  // Bitwise NOT
)

type UnaryOpIns struct {
	instructionBase
	valueBase

	Op  UnaryOp
	Val Value
}

// Comparison

type CmpPredicate uint8

const (
	_ CmpPredicate = iota

	CmpEq  // Compare Equal
	CmpNe  // Compare Not Equal
	CmpSlt // Compare Signed Less Than
	CmpSle // Compare Signed Less or Equal
	CmpSgt // Compare Signed Greater Than
	CmpSge // Compare Signed Greater or Equal
	CmpUlt // Compare Unsigned Less Than
	CmpUle // Compare Unsigned Less or Equal
	CmpUgt // Compare Unsigned Greater Than
	CmpUge // Compare Unsigned Greater or Equal
	CmpFeq // Compare Float Equal
	CmpFne // Compare Float Not Equal
	CmpFlt // Compare Float Less Than
	CmpFle // Compare Float Less or Equal
	CmpFgt // Compare Float Greater Than
	CmpFge // Compare Float Greater or Equal
)

type CmpIns struct {
	instructionBase
	valueBase

	Op  CmpPredicate
	Lhs Value
	Rhs Value
}

// Conversion

type CastOp uint8

const (
	_ CastOp = iota

	CastIntSext  // Sign Extend
	CastIntZext  // Zero Extend
	CastIntTrunc // Truncate
	CastFpExt    // Float Extend
	CastFpTrunc  // Float Truncate
	CastSitofp   // Signed Integer to Float
	CastUintofp  // Unsigned Integer to Float
	CastFptosi   // Float to Signed Integer
	CastFptoui   // Float to Unsigned Integer
	CastInttoptr // Integer to Pointer
	CastPtrtoint // Pointer to Integer
	CastBitcast  // Bitcast
)

type CastIns struct {
	instructionBase
	valueBase

	Op  CastOp
	Val Value
	T   Type
}

// Call

type CallIns struct {
	instructionBase
	valueBase

	Conv       CallingConv
	Proc       Value
	Args       []Value
	ReturnType Type
}

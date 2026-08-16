package fir

type Type interface {
	Kind() TypeKind
}

type TypeKind uint8

const (
	TypeVoid TypeKind = iota
	TypeBool
	TypeInt
	TypeFloat
	TypePtr
	TypeArray
	TypeStruct
	TypeUnion
	TypeFunction
)

var (
	Void = &VoidType{}
	Bool = &BoolType{}
	I8   = &IntType{Bits: 8}
	I16  = &IntType{Bits: 16}
	I32  = &IntType{Bits: 32}
	I64  = &IntType{Bits: 64}
	F32  = &FloatType{Bits: 32}
	F64  = &FloatType{Bits: 64}
	Ptr  = &PtrType{}
)

type VoidType struct{}

func (t *VoidType) Kind() TypeKind {
	return TypeVoid
}

type BoolType struct{}

func (t *BoolType) Kind() TypeKind {
	return TypeBool
}

type IntType struct {
	Bits uint16
}

func (t *IntType) Kind() TypeKind {
	return TypeInt
}

type FloatType struct {
	Bits uint16
}

func (t *FloatType) Kind() TypeKind {
	return TypeFloat
}

type PtrType struct{}

func (t *PtrType) Kind() TypeKind {
	return TypePtr
}

type ArrayType struct {
	Elem Type
	Len  uint64
}

func (t *ArrayType) Kind() TypeKind {
	return TypeArray
}

type StructType struct {
	Name   string
	Fields []AggregateField
}

type AggregateField struct {
	Type Type
}

func (t *StructType) Kind() TypeKind {
	return TypeStruct
}

type UnionType struct {
	Name   string
	Fields []AggregateField
}

func (t *UnionType) Kind() TypeKind {
	return TypeUnion
}

type FuncType struct {
	Args     []Type
	Ret      Type
	Variadic bool
}

func (t *FuncType) Kind() TypeKind {
	return TypeFunction
}

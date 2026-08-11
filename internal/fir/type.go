package fir

import (
	"fmt"
	"strings"
)

type Type interface {
	Kind() TypeKind
	String() string
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
)

type VoidType struct{}

func (t VoidType) Kind() TypeKind {
	return TypeVoid
}

func (t VoidType) String() string {
	return "void"
}

type BoolType struct{}

func (t BoolType) Kind() TypeKind {
	return TypeBool
}

func (t BoolType) String() string {
	return "bool"
}

type IntType struct {
	Bits uint16
}

func (t IntType) Kind() TypeKind {
	return TypeInt
}

func (t IntType) String() string {
	return fmt.Sprintf("i%d", t.Bits)
}

type FloatType struct {
	Bits uint16
}

func (t FloatType) Kind() TypeKind {
	return TypeFloat
}

func (t FloatType) String() string {
	return fmt.Sprintf("f%d", t.Bits)
}

type PtrType struct{}

func (t PtrType) Kind() TypeKind {
	return TypePtr
}

func (t PtrType) String() string {
	return "ptr"
}

type ArrayType struct {
	Elem Type
	Len  uint64
}

func (t ArrayType) Kind() TypeKind {
	return TypeArray
}

func (t ArrayType) String() string {
	return fmt.Sprintf("[%d]%s", t.Len, t.Elem)
}

type StructType struct {
	Fields []AggregateField
}

type AggregateField struct {
	Type Type
}

func (t AggregateField) String() string {
	return fmt.Sprintf("%s", t.Type)
}

func (t StructType) Kind() TypeKind {
	return TypeStruct
}

func (t StructType) String() string {
	fs := make([]string, len(t.Fields))
	for i, f := range t.Fields {
		fs[i] = f.String()
	}
	return fmt.Sprintf("struct { %s }", strings.Join(fs, ", "))
}

type UnionType struct {
	Fields []AggregateField
}

func (t UnionType) Kind() TypeKind {
	return TypeUnion
}

func (t UnionType) String() string {
	fs := make([]string, len(t.Fields))
	for i, f := range t.Fields {
		fs[i] = f.String()
	}
	return fmt.Sprintf("union { %s }", strings.Join(fs, ", "))
}

type FunctionType struct {
	Args     []Type
	Ret      Type
	Variadic bool
}

func (t FunctionType) Kind() TypeKind {
	return TypeFunction
}

func (t FunctionType) String() string {
	args := make([]string, len(t.Args))
	for i, arg := range t.Args {
		args[i] = arg.String()
	}
	if t.Variadic {
		args = append(args, "...")
	}

	return fmt.Sprintf("func(%s) %s", strings.Join(args, ", "), t.Ret)
}

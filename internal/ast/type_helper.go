package ast

import "reflect"

func CompareTypes(t1, t2 Type) bool {
	if t1 == nil || t2 == nil {
		return false
	}

	at1 := reflect.TypeOf(t1)
	at2 := reflect.TypeOf(t2)

	if at1 != at2 {
		return false
	}

	switch t := t1.(type) {
	case *NamedType:
		t2 := t2.(*NamedType).Name
		return t.Name == t2
	case *BuiltinType:
		t2 := t2.(*BuiltinType).Name
		return t.Name == t2
	default:
		return false
	}
}

func IsNumeric(t Type) bool {
	if t == nil {
		return false
	}

	numericTypes := []BuiltinType{
		{"i8"},
		{"i16"},
		{"i32"},
		{"i64"},
		{"u8"},
		{"u16"},
		{"u32"},
		{"u64"},
		{"f32"},
		{"f64"},
	}

	switch t := t.(type) {
	case *BuiltinType:
		for _, v := range numericTypes {
			if CompareTypes(&v, t) {
				return true
			}
		}

	default:
		break
	}

	return false
}

func FuncDeclToFuncType(fdecl *FunctionDeclaration) *FunctionType {
	argTypes := make([]Type, 0, len(fdecl.Args))

	for _, v := range fdecl.Args {
		argTypes = append(argTypes, v.Type)
	}

	return &FunctionType{
		ReturnType: fdecl.ReturnType,
		ArgTypes:   argTypes,
	}
}

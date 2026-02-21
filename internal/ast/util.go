package ast

import "fracta/internal/token"

var (
	BuiltinTypeNameMap = map[string]*BuiltinType{
		"i8":  {"i8"},
		"i16": {"i16"},
		"i32": {"i32"},
		"i64": {"i64"},

		"u8":  {"u8"},
		"u16": {"u16"},
		"u32": {"u32"},
		"u64": {"u64"},

		"f32": {"f32"},
		"f64": {"f64"},

		"bool": {"bool"},

		"ptr": {"ptr"},
	}

	TokenLiteralMap = map[token.TokenType]Type{
		token.TokI8:  &BuiltinType{"i8"},
		token.TokI16: &BuiltinType{"i16"},
		token.TokI32: &BuiltinType{"i32"},
		token.TokI64: &BuiltinType{"i64"},

		token.TokU8:  &BuiltinType{"u8"},
		token.TokU16: &BuiltinType{"u16"},
		token.TokU32: &BuiltinType{"u32"},
		token.TokU64: &BuiltinType{"u64"},

		token.TokF32: &BuiltinType{"f32"},
		token.TokF64: &BuiltinType{"f64"},
	}
)

type ArgPair struct {
	Type Type
	Name token.Token
}

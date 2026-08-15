package fir

type Context struct {
	ptr     *PtrType
	bool    *BoolType
	ints    map[uint16]*IntType     // key: bits
	floats  map[uint16]*FloatType   // key: bits
	structs map[string]*StructType  // key: name
	unions  map[string]*UnionType   // key: name
	arrays  map[arrayKey]*ArrayType // key: {elem type, len}
	funcs   map[funcKey]*FuncType   // key: name

	types []Type // Canonical list of types
}

type arrayKey struct {
	Elem Type
	Len  uint64
}

type funcKey struct {
	Params   string // hash of []Type
	Return   Type
	Variadic bool
}

func NewContext() *Context {
	ctx := &Context{
		ptr:     Ptr,
		bool:    Bool,
		ints:    make(map[uint16]*IntType),
		floats:  make(map[uint16]*FloatType),
		structs: make(map[string]*StructType),
		unions:  make(map[string]*UnionType),
		arrays:  make(map[arrayKey]*ArrayType),
		funcs:   make(map[funcKey]*FuncType),
		types:   make([]Type, 0, 16),
	}

	ctx.types = append(ctx.types,
		Void,
		Bool,
		I8, I16, I32, I64,
		F32, F64,
		Ptr,
	)

	ctx.ints[8] = I8
	ctx.ints[16] = I16
	ctx.ints[32] = I32
	ctx.ints[64] = I64
	ctx.floats[32] = F32
	ctx.floats[64] = F64

	return ctx
}

func (ctx *Context) GetIntType(bits uint16) *IntType {
	if intType, ok := ctx.ints[bits]; ok {
		return intType
	}
	out := &IntType{Bits: bits}
	ctx.ints[bits] = out
	ctx.types = append(ctx.types, out)
	return out
}

func (ctx *Context) GetFloatType(bits uint16) *FloatType {
	if floatType, ok := ctx.floats[bits]; ok {
		return floatType
	}
	out := &FloatType{Bits: bits}
	ctx.floats[bits] = out
	ctx.types = append(ctx.types, out)
	return out
}

func (ctx *Context) GetStructType(name string) *StructType {
	if structType, ok := ctx.structs[name]; ok {
		return structType
	}
	return nil
}

func (ctx *Context) GetUnionType(name string) *UnionType {
	if unionType, ok := ctx.unions[name]; ok {
		return unionType
	}
	return nil
}

func (ctx *Context) GetArrayType(elem Type, len uint64) *ArrayType {
	if arrayType, ok := ctx.arrays[arrayKey{Elem: elem, Len: len}]; ok {
		return arrayType
	}
	out := &ArrayType{Elem: elem, Len: len}
	ctx.arrays[arrayKey{Elem: elem, Len: len}] = out
	ctx.types = append(ctx.types, out)
	return out
}

func (ctx *Context) GetFuncType(variadic bool, ret Type, params ...Type) *FuncType {
	hash := hashSlice(params)
	key := funcKey{Return: ret, Params: hash, Variadic: variadic}
	if funcType, ok := ctx.funcs[key]; ok {
		return funcType
	}
	out := &FuncType{Ret: ret, Args: params, Variadic: variadic}
	ctx.funcs[key] = out
	ctx.types = append(ctx.types, out)
	return out
}

func (ctx *Context) GetPtrType() *PtrType {
	return ctx.ptr
}

func (ctx *Context) GetBoolType() *BoolType {
	return ctx.bool
}

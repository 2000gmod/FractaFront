package fir

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

func (m *Module) Dump() []byte {
	d := &dumper{m: m}
	d.init()
	d.assignTypeNames()

	var buf bytes.Buffer

	fmt.Fprintf(&buf, "; FIR Module: %s\n", m.Name)

	if len(d.aggregates) > 0 {
		buf.WriteString("\n; types\n")
		for _, t := range d.aggregates {
			buf.WriteString(d.typeDefStr(t))
			buf.WriteString("\n")
		}
	}

	if len(m.Globals) > 0 {
		buf.WriteString("\n; globals\n")
		for _, g := range m.Globals {
			buf.WriteString(d.globalStr(g))
			buf.WriteString("\n")
		}
	}

	if len(d.externs) > 0 {
		buf.WriteString("\n; external functions\n")
		for _, f := range d.externs {
			buf.WriteString(d.externFuncStr(f))
			buf.WriteString("\n")
		}
	}

	if len(d.defs) > 0 {
		buf.WriteString("\n; internal functions\n")
		for i, f := range d.defs {
			if i > 0 {
				buf.WriteString("\n")
			}
			buf.WriteString(d.functionStr(f))
		}
	}

	return buf.Bytes()
}

type dumper struct {
	m *Module

	typeNames  map[Type]string
	aggregates []Type
	externs    []*Function
	defs       []*Function

	valueNames map[Value]string
	usedValues map[Value]bool
	kindCount  map[string]int
	paramIndex map[*Parameter]int
}

func (d *dumper) init() {
	d.typeNames = make(map[Type]string)
	d.valueNames = make(map[Value]string)
	d.usedValues = make(map[Value]bool)

	for _, f := range d.m.Functions {
		if f.External {
			d.externs = append(d.externs, f)
		} else {
			d.defs = append(d.defs, f)
		}
	}

	if d.m.ctx != nil {
		for _, t := range d.m.ctx.types {
			switch t.(type) {
			case *StructType, *UnionType:
				d.aggregates = append(d.aggregates, t)
			}
		}
	} else {
		for _, pt := range d.m.Types {
			switch tt := *pt; tt.(type) {
			case *StructType, *UnionType:
				d.aggregates = append(d.aggregates, tt)
			}
		}
	}

	for _, f := range d.m.Functions {
		for _, blk := range f.Blocks {
			d.collectInsUses(blk)
		}
	}
}

func (d *dumper) collectInsUses(b *Block) {
	for _, ins := range b.Ins {
		switch ii := ins.(type) {
		case *LoadIns:
			d.markUse(ii.From)
		case *StoreIns:
			d.markUse(ii.Addr)
			d.markUse(ii.Val)
		case *ElemPtrIns:
			d.markUse(ii.Addr)
			for _, p := range ii.Path {
				d.markUse(p.Index)
			}
		case *BinaryOpIns:
			d.markUse(ii.Lhs)
			d.markUse(ii.Rhs)
		case *UnaryOpIns:
			d.markUse(ii.Val)
		case *CmpIns:
			d.markUse(ii.Lhs)
			d.markUse(ii.Rhs)
		case *CastIns:
			d.markUse(ii.Val)
		case *CallIns:
			d.markUse(ii.Proc)
			for _, a := range ii.Args {
				d.markUse(a)
			}
		}
	}

	if b.Term != nil {
		switch tt := b.Term.(type) {
		case *BrTerm:
		case *CondBrTerm:
			d.markUse(tt.Cond)
		case *RetTerm:
			d.markUse(tt.Value)
		}
	}
}

func (d *dumper) markUse(v Value) {
	if v == nil {
		return
	}
	d.usedValues[v] = true
}

func (d *dumper) assignTypeNames() {
	gen := 0
	for _, t := range d.aggregates {
		switch tt := t.(type) {
		case *StructType:
			if tt.Name != "" {
				d.typeNames[t] = "$" + tt.Name
			} else {
				d.typeNames[t] = fmt.Sprintf("$t%d", gen)
				gen++
			}
		case *UnionType:
			if tt.Name != "" {
				d.typeNames[t] = "$" + tt.Name
			} else {
				d.typeNames[t] = fmt.Sprintf("$t%d", gen)
				gen++
			}
		}
	}
}

func (d *dumper) typeName(t Type) string {
	if name, ok := d.typeNames[t]; ok {
		return name
	}
	switch tt := t.(type) {
	case *StructType:
		return "$" + tt.Name
	case *UnionType:
		return "$" + tt.Name
	}
	return "?"
}

func (d *dumper) typeDefStr(t Type) string {
	switch tt := t.(type) {
	case *StructType:
		return fmt.Sprintf("%s = type struct {%s}", d.typeName(t), d.fieldsStr(tt.Fields))
	case *UnionType:
		return fmt.Sprintf("%s = type union  {%s}", d.typeName(t), d.fieldsStr(tt.Fields))
	}
	return "?"
}

func (d *dumper) fieldsStr(fields []AggregateField) string {
	parts := make([]string, 0, len(fields))
	for _, f := range fields {
		parts = append(parts, d.typeStr(f.Type))
	}
	return strings.Join(parts, ", ")
}

func (d *dumper) typeStr(t Type) string {
	switch tt := t.(type) {
	case *VoidType:
		return "void"
	case *BoolType:
		return "bool"
	case *IntType:
		return "i" + strconv.Itoa(int(tt.Bits))
	case *FloatType:
		return "f" + strconv.Itoa(int(tt.Bits))
	case *PtrType:
		return "ptr"
	case *ArrayType:
		return "[" + strconv.FormatUint(tt.Len, 10) + "]" + d.typeStr(tt.Elem)
	case *StructType, *UnionType:
		return d.typeName(t)
	case *FuncType:
		var buf bytes.Buffer
		buf.WriteString("func(")
		for i, a := range tt.Args {
			if i > 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(d.typeStr(a))
		}
		if tt.Variadic {
			if len(tt.Args) > 0 {
				buf.WriteString(", ")
			}
			buf.WriteString("...")
		}
		buf.WriteString(") ")
		buf.WriteString(d.typeStr(tt.Ret))
		return buf.String()
	}
	return "?"
}

func (d *dumper) globalStr(g *Global) string {
	val := g.Val
	typ := g.ActualType
	if val != nil {
		typ = val.Type()
	}
	return fmt.Sprintf("@%s = link(%s) %s %s", g.Name, linkageStr(g.Linkage), d.typeStr(typ), d.valueStr(val))
}

func (d *dumper) externFuncStr(f *Function) string {
	return fmt.Sprintf("@%s = link(external) cconv(%s) %s", f.Name, convStr(f.Conv), d.funcSigStr(f))
}

func (d *dumper) functionStr(f *Function) string {
	d.valueNames = make(map[Value]string)
	d.kindCount = make(map[string]int)
	d.paramIndex = make(map[*Parameter]int)
	for i, p := range f.Params {
		d.paramIndex[p] = i
	}

	var buf bytes.Buffer
	fmt.Fprintf(&buf, "@%s = cconv(%s) %s\n", f.Name, convStr(f.Conv), d.funcSigStr(f))

	if len(f.Stackslots) > 0 {
		buf.WriteString("stackslots (\n")
		for _, s := range f.Stackslots {
			fmt.Fprintf(&buf, "    %%%s : %s\n", s.Name, d.typeStr(s.InnerType))
		}
		buf.WriteString(")\n")
	} else {
		buf.WriteString("stackslots()\n")
	}

	buf.WriteString("{\n")
	for _, blk := range f.Blocks {
		d.blockStr(&buf, blk)
	}
	buf.WriteString("}\n")

	return buf.String()
}

func (d *dumper) funcSigStr(f *Function) string {
	var buf bytes.Buffer
	buf.WriteString("func(")

	if f.External {
		for i, a := range f.Sig.Args {
			if i > 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(d.typeStr(a))
		}
	} else {
		for i, p := range f.Params {
			if i > 0 {
				buf.WriteString(", ")
			}
			fmt.Fprintf(&buf, "%%%d %s", i, d.typeStr(p.T))
		}
	}

	if f.Sig.Variadic {
		if len(f.Sig.Args) > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString("...")
	}

	buf.WriteString(") ")
	buf.WriteString(d.typeStr(f.Sig.Ret))
	return buf.String()
}

func (d *dumper) blockStr(buf *bytes.Buffer, blk *Block) {
	fmt.Fprintf(buf, "^%s:\n", blk.Name)

	for _, ins := range blk.Ins {
		buf.WriteString("    ")
		d.insStr(buf, ins)
		buf.WriteString("\n")
	}

	if blk.Term != nil {
		buf.WriteString("    ")
		d.termStr(buf, blk.Term)
		buf.WriteString("\n")
	}
}

func (d *dumper) nameFor(v Value) string {
	if name, ok := d.valueNames[v]; ok {
		return name
	}

	prefix := "v"
	switch v.(type) {
	case *LoadIns:
		prefix = "v"
	case *ElemPtrIns:
		prefix = "f"
	case *BinaryOpIns:
		prefix = "b"
	case *UnaryOpIns:
		prefix = "u"
	case *CmpIns:
		prefix = "c"
	case *CastIns:
		prefix = "x"
	case *CallIns:
		prefix = "r"
	}

	idx := d.kindCount[prefix]
	d.kindCount[prefix] = idx + 1

	name := fmt.Sprintf("%%%s%d", prefix, idx)
	d.valueNames[v] = name
	return name
}

func (d *dumper) insStr(buf *bytes.Buffer, ins Instruction) {
	switch ii := ins.(type) {
	case *LoadIns:
		fmt.Fprintf(buf, "%s = load %s (%s)", d.nameFor(ii), d.typeStr(ii.T), d.valueStr(ii.From))
	case *StoreIns:
		fmt.Fprintf(buf, "store (%s) %s", d.valueStr(ii.Addr), d.valueStr(ii.Val))
	case *ElemPtrIns:
		fmt.Fprintf(buf, "%s = elemptr %s (%s)", d.nameFor(ii), d.typeStr(ii.T), d.valueStr(ii.Addr))
		for _, p := range ii.Path {
			buf.WriteString(" ")
			d.pathStr(buf, p)
		}
	case *BinaryOpIns:
		fmt.Fprintf(buf, "%s = %s %s, %s", d.nameFor(ii), binOpStr(ii.Op), d.valueStr(ii.Lhs), d.valueStr(ii.Rhs))
	case *UnaryOpIns:
		fmt.Fprintf(buf, "%s = %s %s", d.nameFor(ii), unOpStr(ii.Op), d.valueStr(ii.Val))
	case *CmpIns:
		fmt.Fprintf(buf, "%s = cmp %s %s, %s", d.nameFor(ii), cmpStr(ii.Op), d.valueStr(ii.Lhs), d.valueStr(ii.Rhs))
	case *CastIns:
		fmt.Fprintf(buf, "%s = %s %s to %s", d.nameFor(ii), castStr(ii.Op), d.valueStr(ii.Val), d.typeStr(ii.T))
	case *CallIns:
		if d.usedValues[ii] {
			fmt.Fprintf(buf, "%s = ", d.nameFor(ii))
		}
		d.callStr(buf, ii)
	}
}

func (d *dumper) pathStr(buf *bytes.Buffer, p ElemPtrPathElem) {
	switch p.Kind {
	case EpArray:
		fmt.Fprintf(buf, "[elem %s]", d.valueStr(p.Index))
	case EpStruct:
		fmt.Fprintf(buf, "[field %s]", d.valueStr(p.Index))
	case EpUnion:
		fmt.Fprintf(buf, "[member %s]", d.valueStr(p.Index))
	}
}

func (d *dumper) callStr(buf *bytes.Buffer, c *CallIns) {
	fmt.Fprintf(buf, "call cconv(%s) %s ", convStr(c.Conv), d.typeStr(c.ReturnType))

	buf.WriteString("(")
	buf.WriteString(d.valueStr(c.Proc))
	buf.WriteString(") ")

	buf.WriteString("(")
	for i, a := range c.Args {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(d.valueStr(a))
	}
	buf.WriteString(")")
}

func (d *dumper) termStr(buf *bytes.Buffer, t Terminator) {
	switch tt := t.(type) {
	case *BrTerm:
		fmt.Fprintf(buf, "br ^%s", tt.Target.Name)
	case *CondBrTerm:
		fmt.Fprintf(buf, "br %s, true(^%s), false(^%s)", d.valueStr(tt.Cond), tt.True.Name, tt.False.Name)
	case *RetTerm:
		fmt.Fprintf(buf, "ret %s", d.valueStr(tt.Value))
	case *UnreachableTerm:
		buf.WriteString("unreachable")
	}
}

func (d *dumper) valueStr(v Value) string {
	switch vv := v.(type) {
	case nil:
		return "?"
	case Constant:
		return d.constantStr(vv)
	case *Global:
		return "@" + vv.Name
	case *Function:
		return "@" + vv.Name
	case *StackSlot:
		return "%" + vv.Name
	case *Parameter:
		if vv.Name != "" {
			return "%" + vv.Name
		}
		return fmt.Sprintf("%%%d", d.paramIndex[vv])
	case ValueInstruction:
		return d.nameFor(vv)
	}
	return "?"
}

func (d *dumper) constantStr(c Constant) string {
	switch cc := c.(type) {
	case *ConstantInt:
		return strconv.FormatUint(cc.Value, 10)
	case *ConstantFloat:
		return strconv.FormatFloat(cc.Value, 'g', -1, 64)
	case *ConstantBool:
		if cc.Value {
			return "true"
		}
		return "false"
	case *ConstantNull:
		return "null"
	case *ConstantVoid:
		return "void"
	case *ConstantArray:
		if isStringType(cc.T) {
			return escapeString(cc)
		}
		parts := make([]string, 0, len(cc.Elements))
		for _, e := range cc.Elements {
			parts = append(parts, d.constantStr(e))
		}
		return "{" + strings.Join(parts, ", ") + "}"
	}
	return "?"
}

func isStringType(t Type) bool {
	at, ok := t.(*ArrayType)
	if !ok {
		return false
	}
	it, ok := at.Elem.(*IntType)
	return ok && it.Bits == 8
}

func escapeString(arr *ConstantArray) string {
	var buf bytes.Buffer
	buf.WriteString(`"`)
	for _, e := range arr.Elements {
		ci, ok := e.(*ConstantInt)
		if !ok {
			continue
		}
		ch := byte(ci.Value)
		switch {
		case ch == '"' || ch == '\\':
			fmt.Fprintf(&buf, `\%02X`, ch)
		case ch >= 0x20 && ch <= 0x7e:
			buf.WriteByte(ch)
		default:
			fmt.Fprintf(&buf, `\%02X`, ch)
		}
	}
	buf.WriteString(`"`)
	return buf.String()
}

func linkageStr(l LinkageType) string {
	switch l {
	case LinkInternal:
		return "internal"
	case LinkExternal:
		return "external"
	case LinkWeak:
		return "weak"
	}
	return "?"
}

func convStr(c CallingConv) string {
	switch c {
	case ConvC:
		return "c"
	case ConvFracta:
		return "fracta"
	}
	return "?"
}

func binOpStr(op BinaryOp) string {
	switch op {
	case BinAdd:
		return "add"
	case BinSub:
		return "sub"
	case BinMul:
		return "mul"
	case BinSdiv:
		return "sdiv"
	case BinUdiv:
		return "udiv"
	case BinSrem:
		return "srem"
	case BinUrem:
		return "urem"
	case BinFadd:
		return "fadd"
	case BinFsub:
		return "fsub"
	case BinFmul:
		return "fmul"
	case BinFdiv:
		return "fdiv"
	case BinFrem:
		return "frem"
	case BinAnd:
		return "and"
	case BinOr:
		return "or"
	case BinXor:
		return "xor"
	case BinShl:
		return "shl"
	case BinLshr:
		return "lshr"
	case BinAshr:
		return "ashr"
	}
	return "?"
}

func unOpStr(op UnaryOp) string {
	switch op {
	case Neg:
		return "neg"
	case Fneg:
		return "fneg"
	case Not:
		return "not"
	}
	return "?"
}

func cmpStr(op CmpPredicate) string {
	switch op {
	case CmpEq:
		return "eq"
	case CmpNe:
		return "ne"
	case CmpSlt:
		return "slt"
	case CmpSle:
		return "sle"
	case CmpSgt:
		return "sgt"
	case CmpSge:
		return "sge"
	case CmpUlt:
		return "ult"
	case CmpUle:
		return "ule"
	case CmpUgt:
		return "ugt"
	case CmpUge:
		return "uge"
	case CmpFeq:
		return "feq"
	case CmpFne:
		return "fne"
	case CmpFlt:
		return "flt"
	case CmpFle:
		return "fle"
	case CmpFgt:
		return "fgt"
	case CmpFge:
		return "fge"
	}
	return "?"
}

func castStr(op CastOp) string {
	switch op {
	case CastIntSext:
		return "sext"
	case CastIntZext:
		return "zext"
	case CastIntTrunc:
		return "trunc"
	case CastFpExt:
		return "fpext"
	case CastFpTrunc:
		return "fptrunc"
	case CastSitofp:
		return "sitofp"
	case CastUintofp:
		return "uitofp"
	case CastFptosi:
		return "fptosi"
	case CastFptoui:
		return "fptoui"
	case CastInttoptr:
		return "inttoptr"
	case CastPtrtoint:
		return "ptrtoint"
	case CastBitcast:
		return "bitcast"
	}
	return "?"
}
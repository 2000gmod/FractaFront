package lexer_test

import (
	"fmt"
	"fracta/internal/lexer"
	tk "fracta/internal/token"
	"reflect"
	"testing"
)

func TestClassifyNumberLiteral(t *testing.T) {
	type entry struct {
		in    string
		kind  tk.TokenType
		value any
	}

	ok := []entry{
		{"0", tk.TokI64, int64(0)},
		{"1", tk.TokI64, int64(1)},
		{"123", tk.TokI64, int64(123)},

		// base-2/8/16 integer only
		{"0b1010", tk.TokI64, int64(0b1010)},
		{"0o77", tk.TokI64, int64(63)},
		{"0xFF", tk.TokI64, int64(255)},

		// suffix integers
		{"5b", tk.TokI8, int8(5)},
		{"300s", tk.TokI16, int16(300)},
		{"123i", tk.TokI32, int32(123)},
		{"999l", tk.TokI64, int64(999)},

		{"5ub", tk.TokU8, uint8(5)},
		{"300us", tk.TokU16, uint16(300)},
		{"123ui", tk.TokU32, uint32(123)},
		{"999ul", tk.TokU64, uint64(999)},

		// floats
		{"1.0", tk.TokF64, float64(1.0)},
		{"0.5", tk.TokF64, float64(0.5)},
		{".5", tk.TokF64, float64(0.5)},

		{"1.5f", tk.TokF32, float32(1.5)},
		{"2.5d", tk.TokF64, float64(2.5)},

		// scientific notation
		{"1e3", tk.TokF64, float64(1000)},
		{"1.5e2", tk.TokF64, float64(150)},
	}

	for _, v := range ok {
		kind, val, err := lexer.ClassifyNumberLiteral(v.in)
		if err != nil {
			t.Fatalf("unexpected error for %q: %v", v.in, err)
		}
		if kind != v.kind {
			t.Fatalf("wrong kind for %q: got %v want %v", v.in, kind, v.kind)
		}
		if fmt.Sprintf("%T", val) != fmt.Sprintf("%T", v.value) {
			t.Fatalf("wrong type for %q: got %T want %T", v.in, val, v.value)
		}
		if !reflect.DeepEqual(val, v.value) {
			t.Fatalf("wrong value for %q: got %v want %v", v.in, val, v.value)
		}
	}

	bad := []string{
		"1.",      // trailing decimal forbidden
		"0x1.0",   // hex float forbidden
		"0b101e3", // exp on non-decimal forbidden
		"e10",     // no leading digit
		".e10",    // invalid
		"1..2",
	}

	for _, in := range bad {
		_, _, err := lexer.ClassifyNumberLiteral(in)
		if err == nil {
			t.Fatalf("expected error for %q but got nil", in)
		}
	}
}

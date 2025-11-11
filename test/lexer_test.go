package lexer

import (
	"fracta/internal/lexer"
	tk "fracta/internal/token"
	"reflect"
	"testing"
)

func TestNumClassification(t *testing.T) {
	t.Log("Testing number token classification.")
	type testDataEntry struct {
		Input string
		EKind tk.TokenType
		Value any
	}

	testData := []testDataEntry{
		{"1", tk.TokI64, int64(1)},
		{"53us", tk.TokU16, uint16(53)},
	}

	for _, v := range testData {
		t.Logf("Running testcase: %v", v)
		kind, val, err := lexer.ClassifyNumberLiteral(v.Input)

		if err != nil {
			t.Error(err)
			return
		}

		if kind != v.EKind {
			t.Errorf("Wrong kind: %v, %v", kind, v.EKind)
			return
		}
		if reflect.TypeOf(val) != reflect.TypeOf(v.Value) {
			t.Errorf("Value type mismatch")
			return
		}
		if !reflect.DeepEqual(val, v.Value) {
			t.Errorf("Value mismatch: %v, %v", val, v.Value)
			return
		}
	}
}

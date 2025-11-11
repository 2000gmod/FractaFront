package lexer

import (
	"fmt"
	tk "fracta/internal/token"
	"strconv"
	"strings"
	"unicode"
)

func ClassifyNumberLiteral(lit string) (tk.TokenType, any, error) {
	orig := lit
	isFloat := false
	base := 10

	if strings.HasPrefix(lit, "0b") {
		base = 2
		lit = lit[2:]
	} else if strings.HasPrefix(lit, "0o") {
		base = 8
		lit = lit[2:]
	} else if strings.HasPrefix(lit, "0x") {
		base = 16
		lit = lit[2:]
	}

	i := 0
	for i < len(lit) && (unicode.IsDigit(rune(lit[i])) || lit[i] == '.' || (base == 16 && strings.ContainsRune("abcdefABCDEF", rune(lit[i])))) {
		if lit[i] == '.' {
			isFloat = true
		}
		i++
	}

	numPart := lit[:i]
	suffix := lit[i:]

	if strings.ContainsAny(numPart, ".eE") {
		isFloat = true
	}

	var t tk.TokenType
	switch suffix {
	case "b":
		t = tk.TokI8
	case "s":
		t = tk.TokI16
	case "i":
		t = tk.TokI32
	case "l":
		t = tk.TokI64
	case "ub":
		t = tk.TokU8
	case "us":
		t = tk.TokU16
	case "ui":
		t = tk.TokU32
	case "ul":
		t = tk.TokU64
	case "f":
		t = tk.TokF32
		isFloat = true
	case "d":
		t = tk.TokF64
		isFloat = true
	default:
		if isFloat {
			t = tk.TokF64
		} else {
			t = tk.TokI64
		}
	}

	if isFloat {
		val, err := strconv.ParseFloat(numPart, 64)
		if err != nil {
			return tk.TokError, nil, fmt.Errorf("Invalid float literal: %q: %w", orig, err)
		}
		if t == tk.TokF32 {
			return t, float32(val), nil
		}
		return t, val, nil
	}

	val, err := strconv.ParseInt(numPart, base, 64)
	if err == nil {
		switch t {
		case tk.TokI8:
			return t, int8(val), nil
		case tk.TokI16:
			return t, int16(val), nil
		case tk.TokI32:
			return t, int32(val), nil
		case tk.TokI64:
			return t, int64(val), nil
		}
	}

	uval, uerr := strconv.ParseUint(numPart, base, 64)
	if uerr == nil {
		switch t {
		case tk.TokU8:
			return t, uint8(uval), nil
		case tk.TokU16:
			return t, uint16(uval), nil
		case tk.TokU32:
			return t, uint32(uval), nil
		case tk.TokU64:
			return t, uint64(uval), nil
		}
	}

	return tk.TokError, nil, fmt.Errorf("invalid numeric literal %q", orig)
}

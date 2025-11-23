package lexer

import (
	"fmt"
	tk "fracta/internal/token"
	"strconv"
	"strings"
	"unicode"
)

func ClassifyNumberLiteral(orig string) (tk.TokenType, any, error) {
	lit := orig
	i := 0
	n := len(lit)

	if n == 0 {
		return tk.TokError, nil, fmt.Errorf("empty literal")
	}

	// ---------- Detect base prefix ----------
	base := 10
	if strings.HasPrefix(lit, "0b") {
		base = 2
		lit = lit[2:] // REMOVE PREFIX
	} else if strings.HasPrefix(lit, "0o") {
		base = 8
		lit = lit[2:]
	} else if strings.HasPrefix(lit, "0x") {
		base = 16
		lit = lit[2:]
	}
	i = 0
	n = len(lit)

	// ---------- State machine ----------
	// states:
	//   INT        digits
	//   DOT        only for base10 (after seeing '.')
	//   FRAC       digits after dot
	//   EXP_MARK   saw e/E, expecting +,-,digit
	//   EXP_SIGN   saw +/- after E
	//   EXP        digits of exponent
	//
	// Invalid:
	//   - float in non-base10
	//   - exponent in non-base10
	//   - 1.   (requires trailing fractional digit)
	//   - underscore
	state := "INT"
	isFloat := false

	acceptHex := "0123456789abcdefABCDEF"
	//acceptDec := "0123456789"

	digitOK := func(r byte) bool {
		switch base {
		case 2:
			return r == '0' || r == '1'
		case 8:
			return r >= '0' && r <= '7'
		case 10:
			return unicode.IsDigit(rune(r))
		case 16:
			return strings.ContainsRune(acceptHex, rune(r))
		}
		return false
	}

	for i < n {
		c := lit[i]

		if c == '_' {
			return tk.TokError, nil, fmt.Errorf("underscores are not allowed in numeric literals: %q", orig)
		}

		switch state {

		case "INT":
			if digitOK(c) {
				i++
				continue
			}
			if c == '.' {
				if base != 10 {
					return tk.TokError, nil, fmt.Errorf("floating literals are only allowed in decimal: %q", orig)
				}
				state = "DOT"
				isFloat = true
				i++
				continue
			}
			if c == 'e' || c == 'E' {
				if base != 10 {
					return tk.TokError, nil, fmt.Errorf("exponents are only allowed in decimal numbers: %q", orig)
				}
				state = "EXP_MARK"
				isFloat = true
				i++
				continue
			}
			goto endNumber

		case "DOT":
			// must see at least one digit
			if unicode.IsDigit(rune(c)) {
				state = "FRAC"
				i++
				continue
			}
			return tk.TokError, nil, fmt.Errorf("fraction requires at least one digit after '.': %q", orig)

		case "FRAC":
			if unicode.IsDigit(rune(c)) {
				i++
				continue
			}
			if c == 'e' || c == 'E' {
				state = "EXP_MARK"
				i++
				continue
			}
			goto endNumber

		case "EXP_MARK":
			if c == '+' || c == '-' {
				state = "EXP_SIGN"
				i++
				continue
			}
			if unicode.IsDigit(rune(c)) {
				state = "EXP"
				i++
				continue
			}
			return tk.TokError, nil, fmt.Errorf("malformed exponent in literal: %q", orig)

		case "EXP_SIGN":
			if unicode.IsDigit(rune(c)) {
				state = "EXP"
				i++
				continue
			}
			return tk.TokError, nil, fmt.Errorf("malformed exponent in literal: %q", orig)

		case "EXP":
			if unicode.IsDigit(rune(c)) {
				i++
				continue
			}
			goto endNumber
		}
	}

endNumber:
	numEnd := i
	numPart := lit[:numEnd]
	suffix := lit[numEnd:]

	// Reject "1." (float with no digits after '.')
	if isFloat {
		// find the position of the decimal point in the *numPart*
		if dot := strings.IndexRune(numPart, '.'); dot != -1 {
			// check if everything after the dot (until exponent) is empty
			frac := numPart[dot+1:]

			// strip exponent part from frac (if any)
			if ei := strings.IndexAny(frac, "eE"); ei != -1 {
				frac = frac[:ei] // remove exponent from this check
			}

			if len(frac) == 0 {
				return tk.TokError, nil, fmt.Errorf("invalid float literal %q: missing digit after decimal point", orig)
			}
		}
	}

	// ---------- Determine suffix type ----------
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
	case "":
		if isFloat {
			t = tk.TokF64
		} else {
			t = tk.TokI64
		}
	default:
		return tk.TokError, nil, fmt.Errorf("unknown numeric suffix %q in literal %q", suffix, orig)
	}

	// ---------- Parse as float ----------
	if isFloat {
		val, err := strconv.ParseFloat(numPart, 64)
		if err != nil {
			return tk.TokError, nil, fmt.Errorf("invalid float literal %q: %w", orig, err)
		}
		if t == tk.TokF32 {
			return t, float32(val), nil
		}
		return t, val, nil
	}

	// ---------- Parse as integer ----------
	valSigned, errI := strconv.ParseInt(numPart, base, 64)
	if errI == nil {
		switch t {
		case tk.TokI8:
			return t, int8(valSigned), nil
		case tk.TokI16:
			return t, int16(valSigned), nil
		case tk.TokI32:
			return t, int32(valSigned), nil
		case tk.TokI64:
			return t, int64(valSigned), nil
		}
	}

	valUnsigned, errU := strconv.ParseUint(numPart, base, 64)
	if errU == nil {
		switch t {
		case tk.TokU8:
			return t, uint8(valUnsigned), nil
		case tk.TokU16:
			return t, uint16(valUnsigned), nil
		case tk.TokU32:
			return t, uint32(valUnsigned), nil
		case tk.TokU64:
			return t, uint64(valUnsigned), nil
		}
	}

	return tk.TokError, nil, fmt.Errorf("overflow or invalid literal %q", orig)
}

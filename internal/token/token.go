package token

// Represents the type of token
type TokenType byte

const (
	TokNone      TokenType = iota // No token
	TokError                      // Lexing error
	TokEndOfFile                  // End of file

	TokI8  // 8 bit signed literal
	TokI16 // 16 bit signed literal
	TokI32 // 32 bit signed literal
	TokI64 // 64 bit signed literal

	TokU8  // 8 bit unsigned literal
	TokU16 // 16 bit unsigned literal
	TokU32 // 32 bit unsigned literal
	TokU64 // 64 bit unsigned literal

	TokF32 // 32 bit floating point literal
	TokF64 // 64 bit floating point literal

	TokChar   // 8 bit character literal
	TokString // String literal

	TokIdentifier // Any identifier

	TokOpPlus  // Operator '+'
	TokOpMinus // Operator '-'
	TokOpStar  // Operator '*'
	TokOpSlash // Operator '/'
	TokOpMod   // Operator '%'

	TokOpAssign // Operator '='

	TokOpEq           // Operator '=='
	TokOpNotEq        // Operator '!='
	TokOpLessThan     // Operator '<'
	TokOpGreaterThan  // Operator '>'
	TokOpLessEqual    // Operator '<='
	TokOpGreaterEqual // Operator '>='

	TokOpenParen    // Punctuation '('
	TokCloseParen   // Punctuation ')'
	TokOpenSquare   // Punctuation '['
	TokCloseSquare  // Punctuation ']'
	TokOpenBracket  // Punctuation '{'
	TokCloseBracket // Punctuation '}'

	TokOpDot         // Operator '.'
	TokOpColon       // Operator ':'
	TokOpDoubleColon // Operator '::'
	TokOpComma       // Operator ','

	TokSemicolon // Punctuation ';'

	TokKwFunc   // Keyword 'func'
	TokKwReturn // Keyword 'return'
)

// Represents a token from Fracta
type Token struct {
	Kind       TokenType // The kind of token this is
	Lexeme     string    // The source string this token was built out of
	Value      any       // Literal value
	Identifier string    // Identifier name if any
	Line       int       // Position within source file
	File       string    // Source file name of this token
}

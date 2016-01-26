package presta

// Token represents a lexical token.
type Tok int64

const (
	// Special tokens
	ILLEGAL Tok = iota
	EOF

	// Literals
	IDENTIFIER // main
	STRING
	NUMBER

	PAREN_OPEN
	PAREN_CLOSE
	CURLY_OPEN
	CURLY_CLOSE

	// Keywords
	MATCH_ALL
	MATCH_FIRST
	REPEAT
	ASSIGN
	FUNC

	LT
	GT
	LTE
	GTE
	EQ
	NEQ

	ADD
	SUB
	MULT
	DIV
	MOD

	ADD_I
	SUB_I
	MULT_I
	DIV_I
	MOD_I

	INC
	DEC

	AND
	OR

	NOT
	CONCAT
)

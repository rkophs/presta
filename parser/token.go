/*
 * Copyright (c) 2016 Ryan Kophs
 *
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to
 * deal in the Software without restriction, including without limitation the
 * rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
 * sell copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 *
 **/

package parser

// Token represents a lexical token.

type Token struct {
	tok  Tok
	lit  string
	line int64
	pos  int64
}

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

func (t *Token) Type() Tok {
	return t.tok
}

func (t *Token) Lit() string {
	return t.lit
}

func (t *Token) Line() int64 {
	return t.line
}

func (t *Token) Pos() int64 {
	return t.pos
}

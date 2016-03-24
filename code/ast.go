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

package code

import (
	"bytes"
	"github.com/rkophs/presta/err"
	"github.com/rkophs/presta/icg"
	"github.com/rkophs/presta/json"
	"github.com/rkophs/presta/parser"
)

type AstNode interface {
	json.Serializable
	Type() AstNodeType
	GenerateICG(code *icg.Code, s *parser.Semantic) err.Error
}

type AstNodeType int64
type MatchType int64
type BinOpType int64
type DataType int64

const (
	PROG AstNodeType = iota
	CONCAT
	LET
	REPEAT
	ASSIGN
	MATCH
	VAR
	CALL
	DATA
	NOT
	BIN_OP
	FUNC
)

const (
	ALL MatchType = iota
	FIRST
)

const (
	ADD BinOpType = iota
	SUB
	MULT
	DIV
	MOD

	ADD_I
	SUB_I
	MULT_I
	DIV_I
	MOD_I

	LT
	LTE
	GT
	GTE
	EQ
	NEQ
	AND
	OR
)

const (
	STRING DataType = iota
	NUMBER
)

func (a *AstNodeType) String() string {
	switch *a {
	case PROG:
		return "PROG"
	case CONCAT:
		return "CONCAT"
	case LET:
		return "LET"
	case REPEAT:
		return "REPEAT"
	case ASSIGN:
		return "ASSIGN"
	case MATCH:
		return "MATCH"
	case VAR:
		return "VAR"
	case CALL:
		return "CALL"
	case DATA:
		return "DATA"
	case NOT:
		return "NOT"
	case BIN_OP:
		return "BIN_OP"
	default:
		return ""
	}
}

func (a *AstNodeType) Serialize(buffer *bytes.Buffer) {
	json.NewString(a.String()).Serialize(buffer)
}

func (m *MatchType) String() string {
	switch *m {
	case ALL:
		return "MATCH_ALL"
	case FIRST:
		return "MATCH_FIRST"
	default:
		return ""
	}
}

func (m *MatchType) Serialize(buffer *bytes.Buffer) {
	json.NewString(m.String()).Serialize(buffer)
}

func (b *BinOpType) String() string {
	switch *b {
	case ADD:
		return "+"
	case SUB:
		return "-"
	case MULT:
		return "*"
	case DIV:
		return "/"
	case MOD:
		return "%"
	case ADD_I:
		return "+="
	case SUB_I:
		return "-="
	case MULT_I:
		return "*="
	case DIV_I:
		return "/="
	case MOD_I:
		return "%="
	case LT:
		return "<"
	case LTE:
		return "<="
	case GT:
		return ">"
	case GTE:
		return ">="
	case EQ:
		return "=="
	case NEQ:
		return "!"
	case AND:
		return "&&"
	case OR:
		return "||"
	default:
		return ""
	}
}

func (b *BinOpType) Serialize(buffer *bytes.Buffer) {
	json.NewString(b.String()).Serialize(buffer)
}

func (d *DataType) String() string {
	switch *d {
	case STRING:
		return "STRING"
	case NUMBER:
		return "NUMBER"
	default:
		return ""
	}
}

func (d *DataType) Serialize(buffer *bytes.Buffer) {
	json.NewString(d.String()).Serialize(buffer)
}

func parseExit(p *parser.TokenScanner, readCount int) (tree AstNode, e err.Error) {
	p.RollBack(readCount)
	return nil, nil
}

func parseValid(p *parser.TokenScanner, node AstNode) (tree AstNode, e err.Error) {
	return node, nil
}

func parseError(p *parser.TokenScanner, msg string, readCount int) (tree AstNode, e err.Error) {
	p.RollBack(readCount)
	return nil, err.NewSyntaxError(msg)
}

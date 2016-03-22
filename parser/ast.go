package parser

import (
	"bytes"
	"github.com/rkophs/presta/json"
)

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

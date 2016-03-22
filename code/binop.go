package code

import (
	"bytes"
	"github.com/rkophs/presta/err"
	"github.com/rkophs/presta/icg"
	"github.com/rkophs/presta/ir"
	"github.com/rkophs/presta/json"
	"github.com/rkophs/presta/parser"
)

type BinOp struct {
	l  AstNode
	r  AstNode
	op BinOpType
}

func NewBinOp(p *parser.TokenScanner, op BinOpType, readCount int) (tree AstNode, e err.Error) {
	if l, e := NewExpression(p); e != nil {
		return parseError(p, e.Message(), readCount)
	} else if l != nil {
		if r, e := NewExpression(p); e != nil {
			return parseError(p, e.Message(), readCount)
		} else if r != nil {
			node := &BinOp{l: l, r: r, op: op}
			return parseValid(p, node)
		} else {
			return parseError(p, "Binary op needs another expression.", readCount)
		}
	} else {
		return parseError(p, "Binary operation needs 2 expressions.", readCount)
	}
}

func (b *BinOp) Type() AstNodeType {
	return BIN_OP
}

func (b *BinOp) Serialize(buffer *bytes.Buffer) {

	json.BuildMap(buffer,
		&json.KV{K: "opType", V: json.NewString(b.op.String())},
		&json.KV{K: "l", V: b.l},
		&json.KV{K: "r", V: b.r},
		&json.KV{K: "type", V: json.NewString("BINOP")})
}

func (b *BinOp) GenerateICG(code *icg.Code, s *parser.Semantic) err.Error {

	/*Compute left side and push onto stack*/
	if e := b.l.GenerateICG(code, s); e != nil {
		return e
	}
	laccess := ir.NewStackAccess(code.GetFrameOffset())
	code.Append(ir.NewPush(code.Ax))
	code.IncrFrameOffset(1)

	/*Compute right side and push onto stack*/
	if e := b.r.GenerateICG(code, s); e != nil {
		return e
	}
	raccess := ir.NewStackAccess(code.GetFrameOffset())
	code.Append(ir.NewPush(code.Ax))
	code.IncrFrameOffset(1)

	switch b.op {
	case ADD:
		code.Append(ir.NewAdd(laccess, raccess)) //Adds and puts result location
		code.Append(ir.NewMov(code.Ax, laccess))
		break
	default:
		return err.NewSymanticError("Unsupported binary operation")
	}
	return nil
}

package parser

import (
	"bytes"
	"github.com/rkophs/presta/icg"
	"github.com/rkophs/presta/ir"
	"github.com/rkophs/presta/json"
	"github.com/rkophs/presta/semantic"
)

type BinOp struct {
	l  AstNode
	r  AstNode
	op BinOpType
}

func NewBinOp(p *Parser, op BinOpType, readCount int) (tree AstNode, err Error) {
	if l, err := NewExpression(p); err != nil {
		return p.parseError(err.Message(), readCount)
	} else if l != nil {
		if r, err := NewExpression(p); err != nil {
			return p.parseError(err.Message(), readCount)
		} else if r != nil {
			node := &BinOp{l: l, r: r, op: op}
			return p.parseValid(node)
		} else {
			return p.parseError("Binary op needs another expression.", readCount)
		}
	} else {
		return p.parseError("Binary operation needs 2 expressions.", readCount)
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

func (b *BinOp) GenerateICG(code *icg.Code, s *semantic.Semantic) Error {

	/*Compute left side and push onto stack*/
	if err := b.l.GenerateICG(code, s); err != nil {
		return err
	}
	laccess := ir.NewStackAccess(code.GetFrameOffset())
	code.Append(ir.NewPush(code.Ax))
	code.IncrFrameOffset(1)

	/*Compute right side and push onto stack*/
	if err := b.r.GenerateICG(code, s); err != nil {
		return err
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
		return NewSymanticError("Unsupported binary operation")
	}
	return nil
}

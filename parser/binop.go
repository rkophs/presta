package parser

import (
	"bytes"
	"fmt"
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

func (b *BinOp) GenerateICG(doNotUse int64, code *icg.Code, s *semantic.Semantic) (int64, Error) {

	fmt.Println("ICG for binop")

	/*Compute left side and push onto stack*/
	if _, err := b.l.GenerateICG(-1, code, s); err != nil {
		return -1, err
	}
	laccess := ir.NewStackAccess(code.GetFrameOffset())
	code.Append(ir.NewPush(code.Ax))
	code.IncrFrameOffset(1)

	/*Compute right side and push onto stack*/
	if _, err := b.r.GenerateICG(-1, code, s); err != nil {
		return -1, err
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
		return -1, NewSymanticError("Unsupported binary operation")
	}
	return -1, nil
}
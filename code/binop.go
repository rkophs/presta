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

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

type Repeat struct {
	condition AstNode
	exec      AstNode
}

func NewRepeatExpr(p *parser.TokenScanner) (tree AstNode, e err.Error) {
	readCount := 0

	/*Check for ^ */
	readCount++
	if tok, eof := p.Read(); eof {
		return parseError(p, "Premature end.", readCount)
	} else if tok.Type() != parser.REPEAT {
		return parseExit(p, readCount) //Not caller, but data identifier
	}

	/*Get expression*/
	var condition AstNode
	if expr, e := NewExpression(p); e != nil {
		return parseError(p, e.Message(), readCount)
	} else if expr != nil {
		condition = expr
	} else {
		return parseError(p, "Repeat op must have condition", readCount)
	}

	/*Get expression*/
	if expr, e := NewExpression(p); e != nil {
		return parseError(p, e.Message(), readCount)
	} else if expr != nil {
		node := &Repeat{condition: condition, exec: expr}
		return parseValid(p, node)
	} else {
		return parseError(p, "Repeat op must have body", readCount)
	}
}

func (r *Repeat) Type() AstNodeType {
	return REPEAT
}

func (r *Repeat) Serialize(buffer *bytes.Buffer) {

	json.BuildMap(buffer,
		&json.KV{K: "condition", V: r.condition},
		&json.KV{K: "body", V: r.exec},
		&json.KV{K: "type", V: json.NewString("REPEAT")})
}

func (r *Repeat) GenerateICG(code *icg.Code, s *parser.Semantic) err.Error {
	return nil
}

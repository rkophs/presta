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

type Assign struct {
	name  string
	value AstNode
}

func NewAssignExpr(p *parser.TokenScanner) (tree AstNode, e err.Error) {
	readCount := 0

	/*Check for : */
	readCount++
	if tok, eof := p.Read(); eof {
		return parseError(p, "Premature end.", readCount)
	} else if tok.Type() != parser.ASSIGN {
		return parseExit(p, readCount) //Not caller, but data identifier
	}

	/*Get variable name*/
	var name string
	readCount++
	if tok, eof := p.Read(); eof {
		return parseError(p, "Premature end.", readCount)
	} else if tok.Type() != parser.IDENTIFIER {
		return parseError(p, "Assignment operator must precede an identifier.", readCount)
	} else {
		name = tok.Lit()
	}

	/*Get expression*/
	if expr, err := NewExpression(p); err != nil {
		return parseError(p, err.Message(), readCount)
	} else if expr != nil {
		node := &Assign{name: name, value: expr}
		return parseValid(p, node)
	} else {
		return parseError(p, "Assignment operator must have valid assignment expression.", readCount)
	}
}

func (a *Assign) Type() AstNodeType {
	return ASSIGN
}

func (a *Assign) Serialize(buffer *bytes.Buffer) {

	json.BuildMap(buffer,
		&json.KV{K: "name", V: json.NewString(a.name)},
		&json.KV{K: "value", V: a.value},
		&json.KV{K: "type", V: json.NewString("ASSIGN")})
}

func (p *Assign) GenerateICG(code *icg.Code, s *parser.Semantic) err.Error {
	return nil
}

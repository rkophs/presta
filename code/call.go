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

type Call struct {
	name   string
	params []AstNode
}

func NewCallExpr(p *parser.TokenScanner) (tree AstNode, e err.Error) {

	readCount := 0

	/*Get variable name*/
	var name string
	readCount++
	if tok, eof := p.Read(); eof {
		return parseError(p, "Premature end.", readCount)
	} else if tok.Type() != parser.IDENTIFIER {
		return parseExit(p, readCount)
	} else {
		name = tok.Lit()
	}

	/*Check for bracket*/
	readCount++
	if tok, eof := p.Read(); eof {
		return parseError(p, "Premature end.", readCount)
	} else if tok.Type() != parser.CURLY_OPEN {
		return parseExit(p, readCount) //Not caller, but data identifier
	}

	/* Check for arguments */
	args := []AstNode{}
	for {
		if expr, e := NewExpression(p); e != nil {
			return parseError(p, e.Message(), readCount)
		} else if expr != nil {
			args = append(args, expr)
		} else {
			break
		}
	}

	/*Check for bracket*/
	readCount++
	if tok, eof := p.Read(); eof {
		return parseError(p, "Premature end.", readCount)
	} else if tok.Type() != parser.CURLY_CLOSE {
		return parseError(p, "Missing closing bracket.", readCount)
	}

	node := &Call{name: name, params: args}
	return parseValid(p, node)
}

func (c *Call) Type() AstNodeType {
	return CALL
}

func (c *Call) Serialize(buffer *bytes.Buffer) {

	params := []json.Serializable{}
	for _, param := range c.params {
		params = append(params, param)
	}

	json.BuildMap(buffer,
		&json.KV{K: "params", V: json.NewArray(params)},
		&json.KV{K: "name", V: json.NewString(c.name)},
		&json.KV{K: "type", V: json.NewString("CALL")})
}

func (c *Call) GenerateICG(code *icg.Code, s *parser.Semantic) err.Error {

	if !s.FunctionExists(c.name) || s.FunctionArity(c.name) != len(c.params) {
		return err.NewSymanticError("Function not found")
	}

	//Generate params
	offsets := make([]int, len(c.params))
	for i, p := range c.params {
		if e := p.GenerateICG(code, s); e != nil {
			return e
		} else {
			offsets[i] = code.GetFrameOffset()
			code.Append(ir.NewPush(code.Ax))
			code.IncrFrameOffset(1)
		}
	}

	//Push params onto the stack as fn args
	for i, _ := range c.params {
		code.Append(ir.NewPush(ir.NewStackAccess(offsets[i])))
	}
	code.IncrFrameOffset(len(c.params))

	//Call function (which loads AX when finished)
	gotoLoc := code.GetFunctionOffset(c.name)
	code.Append(ir.NewCall(gotoLoc))

	return nil
}

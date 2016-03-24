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

type Let struct {
	params []string
	values []AstNode
	exec   AstNode
}

func NewLetExpr(p *parser.TokenScanner) (tree AstNode, e err.Error) {
	readCount := 0

	/*Check if it starts with ':' */
	readCount++
	if tok, eof := p.Read(); eof {
		return parseError(p, "Premature end.", readCount)
	} else if tok.Type() != parser.ASSIGN {
		return parseExit(p, readCount)
	}

	/*Check for parenthesis*/
	readCount++
	if tok, eof := p.Read(); eof {
		return parseError(p, "Premature end.", readCount)
	} else if tok.Type() != parser.PAREN_OPEN {
		// probably an assignment at this point
		return parseExit(p, readCount)
	}

	/* Check for param names and closing parenthesis*/
	params := []string{}
	for {
		readCount++
		if tok, eof := p.Read(); eof {
			return parseError(p, "Premature end.", readCount)
		} else if tok.Type() == parser.IDENTIFIER {
			params = append(params, tok.Lit())
		} else if tok.Type() == parser.PAREN_CLOSE {
			break
		} else {
			return parseError(p, "Looking for parameter identifiers for function", readCount)
		}
	}

	/*Check for parenthesis*/
	readCount++
	if tok, eof := p.Read(); eof {
		return parseError(p, "Premature end.", readCount)
	} else if tok.Type() != parser.PAREN_OPEN {
		return parseError(p, "Missing opening parenthesis for let assignments", readCount)
	}

	/* Check for assignments */
	values := []AstNode{}
	for {
		if node, e := NewExpression(p); e != nil {
			return parseError(p, e.Message(), readCount)
		} else if node != nil {
			values = append(values, node)
		} else {
			break
		}
	}

	if len(values) != len(params) {
		return parseError(p, "Number of assignments must equal number of variables in let.", readCount)
	}

	/*Check for parenthesis*/
	readCount++
	if tok, eof := p.Read(); eof {
		return parseError(p, "Premature end.", readCount)
	} else if tok.Type() != parser.PAREN_CLOSE {
		return parseError(p, "Missing closing parenthesis for let assignments", readCount)
	}

	body, err := NewExpression(p)
	if err != nil {
		return parseError(p, err.Message(), readCount)
	} else if body == nil {
		return parseError(p, "Missing let statement body", readCount)
	}

	node := &Let{params: params, values: values, exec: body}
	return parseValid(p, node)
}

func (l *Let) Type() AstNodeType {
	return LET
}

func (l *Let) Serialize(buffer *bytes.Buffer) {

	params := []json.Serializable{}
	for _, param := range l.params {
		params = append(params, json.NewString(param))
	}

	values := []json.Serializable{}
	for _, value := range l.values {
		values = append(values, value)
	}

	json.BuildMap(buffer,
		&json.KV{K: "names", V: json.NewArray(params)},
		&json.KV{K: "values", V: json.NewArray(values)},
		&json.KV{K: "body", V: l.exec},
		&json.KV{K: "type", V: json.NewString("LET")})
}

func (l *Let) GenerateICG(code *icg.Code, s *parser.Semantic) err.Error {
	return nil
}

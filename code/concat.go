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

type Concat struct {
	components []AstNode
}

func NewConcatExpr(p *parser.TokenScanner) (tree AstNode, e err.Error) {
	readCount := 0

	/* Get '.' */
	readCount++
	if tok, eof := p.Read(); eof {
		return parseError(p, "Premature end.", readCount)
	} else if tok.Type() != parser.CONCAT {
		return parseExit(p, readCount)
	}

	/*Check for parenthesis*/
	readCount++
	if tok, eof := p.Read(); eof {
		return parseError(p, "Premature end.", readCount)
	} else if tok.Type() != parser.PAREN_OPEN {
		return parseError(p, "Missing opening parenthesis for concat", readCount)
	}

	/*Get List*/
	exprs := []AstNode{}
	for {
		if expr, e := NewExpression(p); e != nil {
			return parseError(p, e.Message(), readCount)
		} else if expr != nil {
			exprs = append(exprs, expr)
		} else {
			break
		}

		/*Exit on closing parenthesis*/
		if tok, eof := p.Peek(); eof {
			return parseError(p, "Premature end.", readCount)
		} else if tok.Type() == parser.PAREN_CLOSE {
			break
		}
	}

	/*Check for parenthesis*/
	readCount++
	if tok, eof := p.Read(); eof {
		return parseError(p, "Premature end.", readCount)
	} else if tok.Type() != parser.PAREN_CLOSE {
		return parseError(p, "Missing closing parenthesis for concat", readCount)
	}

	node := &Concat{components: exprs}
	return parseValid(p, node)
}

func (c *Concat) Type() AstNodeType {
	return CONCAT
}

func (c *Concat) Serialize(buffer *bytes.Buffer) {
	components := []json.Serializable{}
	for _, component := range c.components {
		components = append(components, component)
	}

	json.BuildMap(buffer,
		&json.KV{K: "chunks", V: json.NewArray(components)},
		&json.KV{K: "type", V: json.NewString("CONCAT")})
}

func (c *Concat) GenerateICG(code *icg.Code, s *parser.Semantic) err.Error {
	return nil
}

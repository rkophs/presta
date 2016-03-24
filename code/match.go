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

type Match struct {
	conditions []AstNode
	branches   []AstNode
	matchType  MatchType
}

func NewMatchExpr(p *parser.TokenScanner) (tree AstNode, e err.Error) {
	readCount := 0

	/*Get '@' or '|' */
	var matchType MatchType
	readCount++
	if tok, eof := p.Read(); eof {
		return parseError(p, "Premature end.", readCount)
	} else if tok.Type() == parser.MATCH_ALL {
		matchType = ALL
	} else if tok.Type() == parser.MATCH_FIRST {
		matchType = FIRST
	} else {
		return parseExit(p, readCount)
	}

	/*Check for parenthesis*/
	readCount++
	if tok, eof := p.Read(); eof {
		return parseError(p, "Premature end.", readCount)
	} else if tok.Type() != parser.PAREN_OPEN {
		return parseError(p, "Missing opening parenthesis for match", readCount)
	}

	/*Get branches*/
	conditions, branches, err := branches(p)
	if err != nil {
		return parseError(p, err.Message(), readCount)
	}

	if len(conditions) != len(branches) || len(conditions) == 0 {
		return parseError(p, "Invalid number of conditions and branches", readCount)
	}

	/*Check for parenthesis*/
	readCount++
	if tok, eof := p.Read(); eof {
		return parseError(p, "Premature end.", readCount)
	} else if tok.Type() != parser.PAREN_CLOSE {
		return parseError(p, "Missing closing parenthesis for match", readCount)
	}

	node := &Match{conditions: conditions, branches: branches, matchType: matchType}
	return parseValid(p, node)
}

func branches(p *parser.TokenScanner) (c []AstNode, b []AstNode, e err.Error) {
	conds := []AstNode{}
	branches := []AstNode{}
	for {
		if cond, e := NewExpression(p); e != nil {
			return conds, branches, e
		} else if cond != nil {
			conds = append(conds, cond)
		} else {
			break
		}

		if branch, e := NewExpression(p); e != nil {
			return conds, branches, e
		} else if branch == nil {
			return conds, branches, err.NewSyntaxError("Match expression missing branch")
		} else {
			branches = append(branches, branch)
		}
	}

	return conds, branches, nil
}

func (m *Match) Type() AstNodeType {
	return MATCH
}

func (m *Match) Serialize(buffer *bytes.Buffer) {

	branches := []json.Serializable{}
	for _, branch := range m.branches {
		branches = append(branches, branch)
	}

	conditions := []json.Serializable{}
	for _, condition := range m.conditions {
		conditions = append(conditions, condition)
	}

	json.BuildMap(buffer,
		&json.KV{K: "branches", V: json.NewArray(branches)},
		&json.KV{K: "conditions", V: json.NewArray(conditions)},
		&json.KV{K: "type", V: json.NewString("MATCH")})
}

func (m *Match) GenerateICG(code *icg.Code, s *parser.Semantic) err.Error {
	return nil
}

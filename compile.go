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

package presta

import (
	"bytes"
	"fmt"
	"github.com/rkophs/presta/code"
	"github.com/rkophs/presta/err"
	"github.com/rkophs/presta/icg"
	"github.com/rkophs/presta/ir"
	"github.com/rkophs/presta/parser"
	"io"
)

func Compile(r io.Reader) (i []ir.Instruction, e err.Error) {
	tokens, e := Tokenize(r)
	if e != nil {
		return nil, e
	}

	tree, e := Parse(tokens)
	if e != nil {
		return nil, e
	}

	var buffer1 bytes.Buffer
	tree.Serialize(&buffer1)
	fmt.Println(buffer1.String())

	code, e := Generate(tree)
	if e != nil {
		return nil, e
	}

	var buffer bytes.Buffer
	code.Serialize(&buffer)
	fmt.Println(buffer.String())

	return code.GetInstructions(), nil
}

func Generate(tree code.AstNode) (*icg.Code, err.Error) {
	code := icg.NewCode(icg.NewLinker())
	s := parser.NewSemantic()
	if err := tree.GenerateICG(code, s); err != nil {
		return nil, err
	}
	return code, nil
}

func Parse(tokens []parser.Token) (tree code.AstNode, e err.Error) {
	p := parser.NewTokenScanner(tokens)
	return code.NewProgram(p)
}

func Tokenize(reader io.Reader) (tokens []parser.Token, e err.Error) {
	s := parser.NewLexScanner(reader)
	a := []parser.Token{}
	for {
		tok := s.Scan()
		if tok.Type() == parser.EOF {
			break
		} else if tok.Type() == parser.ILLEGAL {
			return a, err.NewLexicalError(fmt.Sprintf("[%d:%d]\tIllegal token:\t%q\n", tok.Line(), tok.Pos(), tok.Lit()))
		} else {
			a = append(a, *tok)
		}
	}

	return a, nil
}

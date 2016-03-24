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

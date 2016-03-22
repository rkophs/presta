package presta

import (
	"bytes"
	"fmt"
	"github.com/rkophs/presta/err"
	"github.com/rkophs/presta/icg"
	"github.com/rkophs/presta/lexer"
	"github.com/rkophs/presta/parser"
	"github.com/rkophs/presta/semantic"
	"io"
)

func Compile(r io.Reader) err.Error {
	tokens, err := Tokenize(r)
	if err != nil {
		return err
	}

	tree, e := Parse(tokens)
	if e != nil {
		return e
	}

	code, e := Generate(tree)
	if e != nil {
		return e
	}

	var buffer bytes.Buffer
	code.Serialize(&buffer)
	fmt.Println(buffer.String())

	return nil
}

func Generate(tree parser.AstNode) (*icg.Code, err.Error) {
	code := icg.NewCode()
	s := semantic.NewSemantic()
	if err := tree.GenerateICG(code, s); err != nil {
		return nil, err
	}
	return code, nil
}

func Parse(tokens []lexer.Token) (tree parser.AstNode, e err.Error) {
	p := parser.NewParser(tokens)
	return parser.NewProgram(p)
}

func Tokenize(reader io.Reader) (tokens []lexer.Token, e err.Error) {
	s := lexer.NewScanner(reader)
	a := []lexer.Token{}
	for {
		tok := s.Scan()
		if tok.Type() == lexer.EOF {
			break
		} else if tok.Type() == lexer.ILLEGAL {
			return a, err.NewLexicalError(fmt.Sprintf("[%d:%d]\tIllegal token:\t%q\n", tok.Line(), tok.Pos(), tok.Lit()))
		} else {
			a = append(a, *tok)
		}
	}

	return a, nil
}

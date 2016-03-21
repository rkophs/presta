package parser

import (
	"bytes"
	"fmt"
	"github.com/rkophs/presta/icg"
	"github.com/rkophs/presta/ir"
	"github.com/rkophs/presta/json"
	"github.com/rkophs/presta/lexer"
	"github.com/rkophs/presta/semantic"
)

type Function struct {
	json.Serializable
	name   string
	params []string
	exec   AstNode
}

func (p *Parser) function() (tree AstNode, err Error) {
	readCount := 0

	/*Check if it starts with '~' */
	readCount++
	if tok, eof := p.read(); eof {
		return p.parseError("Premature end.", readCount)
	} else if tok.Type() != lexer.FUNC {
		return p.parseExit(readCount)
	}

	/*Check for identifier*/
	readCount++
	tok, eof := p.read()
	if eof {
		return p.parseError("Premature end.", readCount)
	} else if tok.Type() != lexer.IDENTIFIER {
		return p.parseError("Function name must follow ~", readCount)
	}
	funcName := tok.Lit()

	/* Check for parenthesis */
	readCount++
	if tok, eof := p.read(); eof {
		return p.parseError("Premature end.", readCount)
	} else if tok.Type() != lexer.PAREN_OPEN {
		return p.parseError("Parenthesis must follow function name", readCount)
	}

	/* Check for param names */
	params := []string{}
	for {
		readCount++
		if tok, eof := p.read(); eof {
			return p.parseError("Premature end.", readCount)
		} else if tok.Type() == lexer.IDENTIFIER {
			params = append(params, tok.Lit())
		} else if tok.Type() == lexer.PAREN_CLOSE {
			break
		} else {
			return p.parseError("Looking for parameter identifiers for function", readCount)
		}
	}

	/*Check for parenthesis*/
	readCount++
	if tok, eof := p.read(); eof {
		return p.parseError("Premature end.", readCount)
	} else if tok.Type() != lexer.PAREN_OPEN {
		return p.parseError("'(' must prefix function body", readCount)
	}

	/* Check for expression */
	expr, err := p.expression()
	if err != nil {
		return p.parseError(err.Message(), readCount)
	} else if expr == nil {
		return p.parseError("Function body must be an executable expression", readCount)
	}

	/*Check for parenthesis*/
	readCount++
	if tok, eof := p.read(); eof {
		return p.parseError("Premature end.", readCount)
	} else if tok.Type() != lexer.PAREN_CLOSE {
		return p.parseError("Parenthesis must postfix function body", readCount)
	}

	node := &Function{name: funcName, params: params, exec: expr}
	return p.parseValid(node)
}

func (p *Function) Serialize(buffer *bytes.Buffer) {

	params := []json.Serializable{}
	for _, param := range p.params {
		params = append(params, json.NewString(param))
	}

	json.BuildMap(buffer,
		&json.KV{K: "name", V: json.NewString(p.name)},
		&json.KV{K: "params", V: json.NewArray(params)},
		&json.KV{K: "body", V: p.exec},
		&json.KV{K: "type", V: json.NewString("FUNC")})
}

func (f *Function) Type() AstNodeType {
	return FUNC
}

func (f *Function) GenerateICG(doNotUse int64, code *icg.Code, s *semantic.Semantic) (int64, Error) {

	fmt.Println("ICG function")
	//Instantiate stack accessors for each param
	s.PushNewScope(f.params)
	for i, p := range f.params {
		code.SetVariable(s.GetVariableId(p), ir.NewStackAccess(-1*(i+1)))
	}

	//Code generate the body
	if _, err := f.exec.GenerateICG(-1, code, s); err != nil {
		fmt.Println("Error2")
		return -1, err
	}

	//Load result into AX and shrink the stack
	code.Append(ir.NewResult(code.Ax))

	return -1, nil
}

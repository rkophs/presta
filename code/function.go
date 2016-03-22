package code

import (
	"bytes"
	"github.com/rkophs/presta/err"
	"github.com/rkophs/presta/icg"
	"github.com/rkophs/presta/ir"
	"github.com/rkophs/presta/json"
	"github.com/rkophs/presta/parser"
	"github.com/rkophs/presta/semantic"
)

type Function struct {
	json.Serializable
	name   string
	params []string
	exec   AstNode
}

func NewFunction(p *parser.Parser) (tree AstNode, e err.Error) {
	readCount := 0

	/*Check if it starts with '~' */
	readCount++
	if tok, eof := p.Read(); eof {
		return parseError(p, "Premature end.", readCount)
	} else if tok.Type() != parser.FUNC {
		return parseExit(p, readCount)
	}

	/*Check for identifier*/
	readCount++
	tok, eof := p.Read()
	if eof {
		return parseError(p, "Premature end.", readCount)
	} else if tok.Type() != parser.IDENTIFIER {
		return parseError(p, "Function name must follow ~", readCount)
	}
	funcName := tok.Lit()

	/* Check for parenthesis */
	readCount++
	if tok, eof := p.Read(); eof {
		return parseError(p, "Premature end.", readCount)
	} else if tok.Type() != parser.PAREN_OPEN {
		return parseError(p, "Parenthesis must follow function name", readCount)
	}

	/* Check for param names */
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
		return parseError(p, "'(' must prefix function body", readCount)
	}

	/* Check for expression */
	expr, err := NewExpression(p)
	if err != nil {
		return parseError(p, err.Message(), readCount)
	} else if expr == nil {
		return parseError(p, "Function body must be an executable expression", readCount)
	}

	/*Check for parenthesis*/
	readCount++
	if tok, eof := p.Read(); eof {
		return parseError(p, "Premature end.", readCount)
	} else if tok.Type() != parser.PAREN_CLOSE {
		return parseError(p, "Parenthesis must postfix function body", readCount)
	}

	node := &Function{name: funcName, params: params, exec: expr}
	return parseValid(p, node)
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

func (f *Function) GenerateICG(code *icg.Code, s *semantic.Semantic) err.Error {

	//Instantiate stack accessors for each param
	s.PushNewScope(f.params)
	for i, p := range f.params {
		code.SetVariable(s.GetVariableId(p), ir.NewStackAccess(-1*(i+1)))
	}

	//Code generate the body
	if e := f.exec.GenerateICG(code, s); e != nil {
		return e
	}

	//Load result into AX and shrink the stack
	code.Append(ir.NewResult(code.Ax))

	return nil
}

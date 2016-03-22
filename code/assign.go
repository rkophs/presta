package code

import (
	"bytes"
	"github.com/rkophs/presta/err"
	"github.com/rkophs/presta/icg"
	"github.com/rkophs/presta/json"
	"github.com/rkophs/presta/parser"
	"github.com/rkophs/presta/semantic"
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

func (p *Assign) GenerateICG(code *icg.Code, s *semantic.Semantic) err.Error {
	return nil
}

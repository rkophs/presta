package parser

import (
	"bytes"
	"github.com/rkophs/presta/icg"
	"github.com/rkophs/presta/json"
	"github.com/rkophs/presta/lexer"
	"github.com/rkophs/presta/semantic"
)

type Assign struct {
	name  string
	value AstNode
}

func NewAssignExpr(p *Parser) (tree AstNode, err Error) {
	readCount := 0

	/*Check for : */
	readCount++
	if tok, eof := p.read(); eof {
		return p.parseError("Premature end.", readCount)
	} else if tok.Type() != lexer.ASSIGN {
		return p.parseExit(readCount) //Not caller, but data identifier
	}

	/*Get variable name*/
	var name string
	readCount++
	if tok, eof := p.read(); eof {
		return p.parseError("Premature end.", readCount)
	} else if tok.Type() != lexer.IDENTIFIER {
		return p.parseError("Assignment operator must precede an identifier.", readCount)
	} else {
		name = tok.Lit()
	}

	/*Get expression*/
	if expr, err := p.expression(); err != nil {
		return p.parseError(err.Message(), readCount)
	} else if expr != nil {
		node := &Assign{name: name, value: expr}
		return p.parseValid(node)
	} else {
		return p.parseError("Assignment operator must have valid assignment expression.", readCount)
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

func (p *Assign) GenerateICG(offset int64, code *icg.Code, s *semantic.Semantic) (int64, Error) {
	return -1, nil
}

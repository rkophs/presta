package code

import (
	"bytes"
	"github.com/rkophs/presta/err"
	"github.com/rkophs/presta/icg"
	"github.com/rkophs/presta/json"
	"github.com/rkophs/presta/parser"
	"github.com/rkophs/presta/semantic"
)

type Not struct {
	exec AstNode
}

func NewNotExpr(p *parser.Parser) (tree AstNode, e err.Error) {
	readCount := 0
	/*Check for ! */
	readCount++
	if tok, eof := p.Read(); eof {
		return parseError(p, "Premature end.", readCount)
	} else if tok.Type() != parser.NOT {
		return parseExit(p, readCount) //Not caller, but data identifier
	}

	if expr, e := NewExpression(p); e != nil {
		return parseError(p, e.Message(), readCount)
	} else if expr != nil {
		node := &Not{exec: expr}
		return parseValid(p, node)
	} else {
		return parseError(p, "Not operator must precede expression", readCount)
	}

}

func (n *Not) Type() AstNodeType {
	return NOT
}

func (n *Not) Serialize(buffer *bytes.Buffer) {

	json.BuildMap(buffer,
		&json.KV{K: "expression", V: n.exec},
		&json.KV{K: "type", V: json.NewString("NOT")})
}

func (n *Not) GenerateICG(code *icg.Code, s *semantic.Semantic) err.Error {
	return nil
}

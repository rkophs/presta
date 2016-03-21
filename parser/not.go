package parser

import (
	"bytes"
	"github.com/rkophs/presta/icg"
	"github.com/rkophs/presta/json"
	"github.com/rkophs/presta/lexer"
	"github.com/rkophs/presta/semantic"
)

type Not struct {
	exec AstNode
}

func (p *Parser) notExpr() (tree AstNode, yes bool, err Error) {
	readCount := 0
	/*Check for ! */
	readCount++
	if tok, eof := p.read(); eof {
		return p.parseError("Premature end.", readCount)
	} else if tok.Type() != lexer.NOT {
		return p.parseExit(readCount) //Not caller, but data identifier
	}

	if expr, yes, err := p.expression(); err != nil {
		return p.parseError(err.Message(), readCount)
	} else if yes {
		node := &Not{exec: expr}
		return p.parseValid(node)
	} else {
		return p.parseError("Not operator must precede expression", readCount)
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

func (n *Not) GenerateICG(offset int64, code *icg.Code, s *semantic.Semantic) (int64, Error) {
	return -1, nil
}

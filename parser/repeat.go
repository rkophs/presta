package parser

import (
	"bytes"
	"github.com/rkophs/presta/icg"
	"github.com/rkophs/presta/json"
	"github.com/rkophs/presta/lexer"
	"github.com/rkophs/presta/semantic"
)

type Repeat struct {
	condition AstNode
	exec      AstNode
}

func NewRepeatExpr(p *Parser) (tree AstNode, err Error) {
	readCount := 0

	/*Check for ^ */
	readCount++
	if tok, eof := p.read(); eof {
		return p.parseError("Premature end.", readCount)
	} else if tok.Type() != lexer.REPEAT {
		return p.parseExit(readCount) //Not caller, but data identifier
	}

	/*Get expression*/
	var condition AstNode
	if expr, err := p.expression(); err != nil {
		return p.parseError(err.Message(), readCount)
	} else if expr != nil {
		condition = expr
	} else {
		return p.parseError("Repeat op must have condition", readCount)
	}

	/*Get expression*/
	if expr, err := p.expression(); err != nil {
		return p.parseError(err.Message(), readCount)
	} else if expr != nil {
		node := &Repeat{condition: condition, exec: expr}
		return p.parseValid(node)
	} else {
		return p.parseError("Repeat op must have body", readCount)
	}
}

func (r *Repeat) Type() AstNodeType {
	return REPEAT
}

func (r *Repeat) Serialize(buffer *bytes.Buffer) {

	json.BuildMap(buffer,
		&json.KV{K: "condition", V: r.condition},
		&json.KV{K: "body", V: r.exec},
		&json.KV{K: "type", V: json.NewString("REPEAT")})
}

func (r *Repeat) GenerateICG(offset int64, code *icg.Code, s *semantic.Semantic) (int64, Error) {
	return -1, nil
}

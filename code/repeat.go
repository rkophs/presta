package code

import (
	"bytes"
	"github.com/rkophs/presta/err"
	"github.com/rkophs/presta/icg"
	"github.com/rkophs/presta/json"
	"github.com/rkophs/presta/lexer"
	"github.com/rkophs/presta/parser"
	"github.com/rkophs/presta/semantic"
)

type Repeat struct {
	condition AstNode
	exec      AstNode
}

func NewRepeatExpr(p *parser.Parser) (tree AstNode, e err.Error) {
	readCount := 0

	/*Check for ^ */
	readCount++
	if tok, eof := p.Read(); eof {
		return parseError(p, "Premature end.", readCount)
	} else if tok.Type() != lexer.REPEAT {
		return parseExit(p, readCount) //Not caller, but data identifier
	}

	/*Get expression*/
	var condition AstNode
	if expr, e := NewExpression(p); e != nil {
		return parseError(p, e.Message(), readCount)
	} else if expr != nil {
		condition = expr
	} else {
		return parseError(p, "Repeat op must have condition", readCount)
	}

	/*Get expression*/
	if expr, e := NewExpression(p); e != nil {
		return parseError(p, e.Message(), readCount)
	} else if expr != nil {
		node := &Repeat{condition: condition, exec: expr}
		return parseValid(p, node)
	} else {
		return parseError(p, "Repeat op must have body", readCount)
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

func (r *Repeat) GenerateICG(code *icg.Code, s *semantic.Semantic) err.Error {
	return nil
}

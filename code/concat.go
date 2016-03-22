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

type Concat struct {
	components []AstNode
}

func NewConcatExpr(p *parser.Parser) (tree AstNode, e err.Error) {
	readCount := 0

	/* Get '.' */
	readCount++
	if tok, eof := p.Read(); eof {
		return parseError(p, "Premature end.", readCount)
	} else if tok.Type() != lexer.CONCAT {
		return parseExit(p, readCount)
	}

	/*Check for parenthesis*/
	readCount++
	if tok, eof := p.Read(); eof {
		return parseError(p, "Premature end.", readCount)
	} else if tok.Type() != lexer.PAREN_OPEN {
		return parseError(p, "Missing opening parenthesis for concat", readCount)
	}

	/*Get List*/
	exprs := []AstNode{}
	for {
		if expr, e := NewExpression(p); e != nil {
			return parseError(p, e.Message(), readCount)
		} else if expr != nil {
			exprs = append(exprs, expr)
		} else {
			break
		}

		/*Exit on closing parenthesis*/
		if tok, eof := p.Peek(); eof {
			return parseError(p, "Premature end.", readCount)
		} else if tok.Type() == lexer.PAREN_CLOSE {
			break
		}
	}

	/*Check for parenthesis*/
	readCount++
	if tok, eof := p.Read(); eof {
		return parseError(p, "Premature end.", readCount)
	} else if tok.Type() != lexer.PAREN_CLOSE {
		return parseError(p, "Missing closing parenthesis for concat", readCount)
	}

	node := &Concat{components: exprs}
	return parseValid(p, node)
}

func (c *Concat) Type() parser.AstNodeType {
	return parser.CONCAT
}

func (c *Concat) Serialize(buffer *bytes.Buffer) {
	components := []json.Serializable{}
	for _, component := range c.components {
		components = append(components, component)
	}

	json.BuildMap(buffer,
		&json.KV{K: "chunks", V: json.NewArray(components)},
		&json.KV{K: "type", V: json.NewString("CONCAT")})
}

func (c *Concat) GenerateICG(code *icg.Code, s *semantic.Semantic) err.Error {
	return nil
}

package parser

import (
	"bytes"
	"github.com/rkophs/presta/icg"
	"github.com/rkophs/presta/json"
	"github.com/rkophs/presta/lexer"
	"github.com/rkophs/presta/semantic"
)

type Concat struct {
	components []AstNode
}

func (p *Parser) concatExpr() (tree AstNode, yes bool, err Error) {
	readCount := 0

	/* Get '.' */
	readCount++
	if tok, eof := p.read(); eof {
		return p.parseError("Premature end.", readCount)
	} else if tok.Type() != lexer.CONCAT {
		return p.parseExit(readCount)
	}

	/*Check for parenthesis*/
	readCount++
	if tok, eof := p.read(); eof {
		return p.parseError("Premature end.", readCount)
	} else if tok.Type() != lexer.PAREN_OPEN {
		return p.parseError("Missing opening parenthesis for concat", readCount)
	}

	/*Get List*/
	exprs := []AstNode{}
	for {
		if expr, yes, err := p.expression(); err != nil {
			return p.parseError(err.Message(), readCount)
		} else if yes {
			exprs = append(exprs, expr)
		} else {
			break
		}

		/*Exit on closing parenthesis*/
		if tok, eof := p.peek(); eof {
			return p.parseError("Premature end.", readCount)
		} else if tok.Type() == lexer.PAREN_CLOSE {
			break
		}
	}

	/*Check for parenthesis*/
	readCount++
	if tok, eof := p.read(); eof {
		return p.parseError("Premature end.", readCount)
	} else if tok.Type() != lexer.PAREN_CLOSE {
		return p.parseError("Missing closing parenthesis for concat", readCount)
	}

	node := &Concat{components: exprs}
	return p.parseValid(node)
}

func (c *Concat) Type() AstNodeType {
	return CONCAT
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

func (c *Concat) GenerateICG(offset int64, code *icg.Code, s *semantic.Semantic) (int64, Error) {
	return -1, nil
}

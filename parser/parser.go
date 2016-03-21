package parser

import (
	"github.com/rkophs/presta/lexer"
	// "strconv"
)

type Parser struct {
	tokens []lexer.Token
	at     int
}

func NewParser(tokens []lexer.Token) *Parser {
	return &Parser{tokens: tokens, at: 0}
}

func (p *Parser) Scan() (tree AstNode, err Error) {
	if tree, err := NewProgram(p); err != nil {
		return nil, err
	} else if tree == nil {
		return p.parseError("No program available.", 0)
	} else {
		return p.parseValid(tree)
	}
}

func (p *Parser) parseExit(readCount int) (tree AstNode, err Error) {
	p.rollBack(readCount)
	return nil, nil
}

func (p *Parser) parseValid(node AstNode) (tree AstNode, err Error) {
	return node, nil
}

func (p *Parser) parseError(msg string, readCount int) (tree AstNode, err Error) {
	p.rollBack(readCount)
	return nil, NewSyntaxError(msg)
}

func (p *Parser) rollBack(amount int) {
	for i := 0; i < amount; i++ {
		p.unread()
	}
}

func (p *Parser) read() (tok lexer.Token, eof bool) {
	tok, eof = p.peek()
	if !eof {
		p.at++
	}
	return tok, eof
}

func (p *Parser) peek() (tok lexer.Token, eof bool) {
	if p.at >= len(p.tokens) {
		var ret lexer.Token
		return ret, true
	}
	tok = p.tokens[p.at]
	return tok, false
}

func (p *Parser) unread() {
	p.at--
}

package parser

import (
	"github.com/rkophs/presta/err"
	"github.com/rkophs/presta/lexer"
)

type Parser struct {
	tokens []lexer.Token
	at     int
}

func NewParser(tokens []lexer.Token) *Parser {
	return &Parser{tokens: tokens, at: 0}
}

func (p *Parser) Scan() (tree AstNode, e err.Error) {
	if tree, e := NewProgram(p); e != nil {
		return nil, e
	} else {
		return p.parseValid(tree)
	}
}

func (p *Parser) parseExit(readCount int) (tree AstNode, e err.Error) {
	p.rollBack(readCount)
	return nil, nil
}

func (p *Parser) parseValid(node AstNode) (tree AstNode, e err.Error) {
	return node, nil
}

func (p *Parser) parseError(msg string, readCount int) (tree AstNode, e err.Error) {
	p.rollBack(readCount)
	return nil, err.NewSyntaxError(msg)
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

package parser

import (
	"github.com/rkophs/presta/lexer"
)

type Parser struct {
	tokens []lexer.Token
	at     int
}

func NewParser(tokens []lexer.Token) *Parser {
	return &Parser{tokens: tokens, at: 0}
}

func (p *Parser) RollBack(amount int) {
	for i := 0; i < amount; i++ {
		p.Unread()
	}
}

func (p *Parser) Read() (tok lexer.Token, eof bool) {
	tok, eof = p.Peek()
	if !eof {
		p.at++
	}
	return tok, eof
}

func (p *Parser) Peek() (tok lexer.Token, eof bool) {
	if p.at >= len(p.tokens) {
		var ret lexer.Token
		return ret, true
	}
	tok = p.tokens[p.at]
	return tok, false
}

func (p *Parser) Unread() {
	p.at--
}

package parser

type Parser struct {
	tokens []Token
	at     int
}

func NewParser(tokens []Token) *Parser {
	return &Parser{tokens: tokens, at: 0}
}

func (p *Parser) RollBack(amount int) {
	for i := 0; i < amount; i++ {
		p.Unread()
	}
}

func (p *Parser) Read() (tok Token, eof bool) {
	tok, eof = p.Peek()
	if !eof {
		p.at++
	}
	return tok, eof
}

func (p *Parser) Peek() (tok Token, eof bool) {
	if p.at >= len(p.tokens) {
		var ret Token
		return ret, true
	}
	tok = p.tokens[p.at]
	return tok, false
}

func (p *Parser) Unread() {
	p.at--
}

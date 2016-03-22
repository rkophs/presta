package parser

type TokenScanner struct {
	tokens []Token
	at     int
}

func NewTokenScanner(tokens []Token) *TokenScanner {
	return &TokenScanner{tokens: tokens, at: 0}
}

func (p *TokenScanner) RollBack(amount int) {
	for i := 0; i < amount; i++ {
		p.Unread()
	}
}

func (p *TokenScanner) Read() (tok Token, eof bool) {
	tok, eof = p.Peek()
	if !eof {
		p.at++
	}
	return tok, eof
}

func (p *TokenScanner) Peek() (tok Token, eof bool) {
	if p.at >= len(p.tokens) {
		var ret Token
		return ret, true
	}
	tok = p.tokens[p.at]
	return tok, false
}

func (p *TokenScanner) Unread() {
	p.at--
}

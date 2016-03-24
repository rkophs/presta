/*
 * Copyright (c) 2016 Ryan Kophs
 *
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to
 * deal in the Software without restriction, including without limitation the
 * rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
 * sell copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 *
 **/

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

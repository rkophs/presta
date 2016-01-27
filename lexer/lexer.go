package lexer

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
)

type Scanner struct {
	r    *bufio.Reader
	line int64
	pos  int64
}

func Tokenize(reader io.Reader) (tokens []Token, err bool, msg string) {
	s := NewScanner(reader)
	a := make([]Token, 1)
	for {
		tok := s.Scan()
		if tok.tok == EOF {
			break
		} else if tok.tok == ILLEGAL {
			return nil, true, fmt.Sprintf("[%d:%d]\tIllegal token:\t%q\n", tok.line, tok.pos, tok.lit)
		} else {
			a = append(a, *tok)
		}
	}

	return a, false, ""
}

func NewScanner(r io.Reader) *Scanner {
	return &Scanner{r: bufio.NewReader(r), line: 0, pos: 0}
}

func (s *Scanner) Scan() *Token {

	if ch, _, _ := s.peek(); isWhitespace(ch) {
		s.discardWhitespace()
	}

	if ch, _, _ := s.peek(); isLetter(ch) {
		return s.scanIdent()
	} else if isQuote(ch) {
		return s.scanStringLiteral()
	} else if isSymbol(ch) {
		return s.scanOperation()
	} else if isDigit(ch) {
		return s.scanNumber()
	}

	ch, l, p := s.read()
	var tok Tok
	var lit string
	switch ch {
	case eof:
		tok, lit = EOF, ""
	case '(':
		tok, lit = PAREN_OPEN, string(ch)
	case ')':
		tok, lit = PAREN_CLOSE, string(ch)
	case '{':
		tok, lit = CURLY_OPEN, string(ch)
	case '}':
		tok, lit = CURLY_CLOSE, string(ch)
	default:
		tok, lit = ILLEGAL, string(ch)
	}

	return &Token{tok: tok, lit: lit, line: l, pos: p}
}

func (s *Scanner) scanOperation() *Token {
	var buf bytes.Buffer
	ch, l, p := s.read()
	buf.WriteRune(ch)

	var lit string
	var tok Tok

	switch ch {
	case '>':
		tok, lit = s.handleTwoOptions('=', GTE, GT, buf)
	case '<':
		tok, lit = s.handleTwoOptions('=', LTE, LT, buf)
	case '=':
		tok, lit = s.handleTwoOptions('=', EQ, ILLEGAL, buf)
	case '&':
		tok, lit = s.handleTwoOptions('&', AND, ILLEGAL, buf)
	case '|':
		tok, lit = s.handleTwoOptions('|', OR, MATCH_FIRST, buf)
	case '@':
		tok, lit = MATCH_ALL, buf.String()
	case '^':
		tok, lit = REPEAT, buf.String()
	case '!':
		tok, lit = s.handleTwoOptions('=', NEQ, NOT, buf)
	case '+':
		tok, lit = s.handleThreeOptions('=', '+', ADD_I, INC, ADD, buf)
	case '-':
		tok, lit = s.handleThreeOptions('=', '-', SUB_I, DEC, SUB, buf)
	case '/':
		tok, lit = s.handleTwoOptions('=', DIV_I, DIV, buf)
	case '*':
		tok, lit = s.handleTwoOptions('=', MULT_I, MULT, buf)
	case '%':
		tok, lit = s.handleTwoOptions('=', MOD_I, MOD, buf)
	case ':':
		tok, lit = ASSIGN, buf.String()
	case '~':
		tok, lit = FUNC, buf.String()
	case '.':
		tok, lit = CONCAT, buf.String()
	default:
		tok, lit = ILLEGAL, buf.String()
	}

	return &Token{tok: tok, lit: lit, line: l, pos: p}
}

func (s *Scanner) handleTwoOptions(cmp rune, yes Tok, no Tok, buf bytes.Buffer) (tok Tok, lit string) {
	if ch, _, _ := s.peek(); ch == cmp {
		s.read()
		buf.WriteRune(ch)
		return yes, buf.String()
	} else {
		return no, buf.String()
	}
}

func (s *Scanner) handleThreeOptions(ifCmp rune, elifCmp rune, a Tok, b Tok, c Tok, buf bytes.Buffer) (tok Tok, lit string) {
	if ch, _, _ := s.peek(); ch == ifCmp {
		s.read()
		buf.WriteRune(ch)
		return a, buf.String()
	} else if ch == elifCmp {
		s.read()
		buf.WriteRune(ch)
		return b, buf.String()
	} else {
		return c, buf.String()
	}
}

func (s *Scanner) scanNumber() *Token {

	var buf bytes.Buffer
	ch, l, p := s.read()
	buf.WriteRune(ch)
	dotUsed := false

	for {
		if ch, _, _ := s.peek(); ch == '.' {
			if dotUsed {
				break
			} else {
				dotUsed = true
				s.read()
				buf.WriteRune(ch)
				if next, _, _ := s.peek(); !isDigit(next) {
					return &Token{tok: ILLEGAL, lit: buf.String(), line: l, pos: p}
				}
			}
		} else if isDigit(ch) {
			s.read()
			buf.WriteRune(ch)
		} else {
			break
		}
	}

	return &Token{tok: NUMBER, lit: buf.String(), line: l, pos: p}
}

func (s *Scanner) scanStringLiteral() *Token {
	var buf bytes.Buffer
	_, l, p := s.read() //throw away string opener

	for {
		if ch, _, _ := s.read(); ch == eof || isQuote(ch) { //throw away string closer
			return &Token{tok: STRING, lit: buf.String(), line: l, pos: p}
		} else {
			buf.WriteRune(ch)
		}
	}
}

func (s *Scanner) scanIdent() *Token {
	var buf bytes.Buffer
	ch, l, p := s.read()
	buf.WriteRune(ch)

	for {
		if ch, _, _ := s.peek(); ch == eof || (!isLetter(ch) && !isDigit(ch) && ch != '_') {
			break
		} else {
			ch, _, _ = s.read()
			_, _ = buf.WriteRune(ch)
		}
	}

	return &Token{tok: IDENTIFIER, lit: buf.String(), line: l, pos: p}
}

func (s *Scanner) discardWhitespace() {
	for {
		if ch, _, _ := s.peek(); ch == eof || !isWhitespace(ch) {
			break
		}
		s.read()
	}
}

func (s *Scanner) peek() (rune, int64, int64) {
	ch, _, err := s.r.ReadRune()
	_ = s.r.UnreadRune()
	if err != nil {
		return eof, -1, -1
	}
	if ch == '\n' {
		return ch, s.line + 1, 0
	} else {
		return ch, s.line, s.pos + 1
	}
}

func (s *Scanner) read() (rune, int64, int64) {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eof, -1, -1
	}
	if ch == '\n' {
		s.line++
		s.pos = 0
	} else {
		s.pos++
	}
	return ch, s.line, s.pos
}

func isSymbol(ch rune) bool {
	symbols := []rune{'>', '<', '=', '&', '|', '@', '^', '!', '+', '-', '/', '*', '%', '.', ':', '~'}
	return runeInSlice(ch, symbols)
}

func isQuote(ch rune) bool { return ch == '\'' }

func isWhitespace(ch rune) bool { return ch == ' ' || ch == '\t' || ch == '\n' }

func isLetter(ch rune) bool { return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') }

func isDigit(ch rune) bool { return (ch >= '0' && ch <= '9') }

func runeInSlice(a rune, list []rune) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

var eof = rune(0)

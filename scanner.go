package presta

import (
	"bufio"
	"bytes"
	"io"
)

type Scanner struct {
	r *bufio.Reader
}

func NewScanner(r io.Reader) *Scanner {
	return &Scanner{r: bufio.NewReader(r)}
}

func (s *Scanner) Scan() (tok Tok, lit string) {

	if ch := s.peek(); isWhitespace(ch) {
		s.discardWhitespace()
	}

	if ch := s.peek(); isLetter(ch) {
		return s.scanIdent()
	} else if isQuote(ch) {
		return s.scanStringLiteral()
	} else if isSymbol(ch) {
		return s.scanOperation()
	} else if isDigit(ch) {
		return s.scanNumber()
	}

	switch ch := s.read(); ch {
	case eof:
		return EOF, ""
	case '(':
		return PAREN_OPEN, string(ch)
	case ')':
		return PAREN_CLOSE, string(ch)
	case '{':
		return CURLY_OPEN, string(ch)
	case '}':
		return CURLY_CLOSE, string(ch)
	default:
		return ILLEGAL, string(ch)
	}
}

func (s *Scanner) scanOperation() (tok Tok, lit string) {
	var buf bytes.Buffer
	ch := s.read()
	buf.WriteRune(ch)

	switch ch {
	case '>':
		return s.handleTwoOptions('=', GTE, GT, buf)
	case '<':
		return s.handleTwoOptions('=', LTE, LT, buf)
	case '=':
		return s.handleTwoOptions('=', EQ, ILLEGAL, buf)
	case '&':
		return s.handleTwoOptions('&', AND, ILLEGAL, buf)
	case '|':
		return s.handleTwoOptions('|', OR, MATCH_FIRST, buf)
	case '@':
		return MATCH_ALL, buf.String()
	case '^':
		return REPEAT, buf.String()
	case '!':
		return s.handleTwoOptions('=', NEQ, NOT, buf)
	case '+':
		return s.handleThreeOptions('=', '+', ADD_I, INC, ADD, buf)
	case '-':
		return s.handleThreeOptions('=', '-', SUB_I, DEC, SUB, buf)
	case '/':
		return s.handleTwoOptions('=', DIV_I, DIV, buf)
	case '*':
		return s.handleTwoOptions('=', MULT_I, MULT, buf)
	case '%':
		return s.handleTwoOptions('=', MOD_I, MOD, buf)
	case ':':
		return ASSIGN, buf.String()
	case '~':
		return FUNC, buf.String()
	case '.':
		return CONCAT, buf.String()
	default:
		return ILLEGAL, buf.String()
	}
}

func (s *Scanner) handleTwoOptions(cmp rune, yes Tok, no Tok, buf bytes.Buffer) (tok Tok, lit string) {
	if ch := s.peek(); ch == cmp {
		buf.WriteRune(s.read())
		return yes, buf.String()
	} else {
		return no, buf.String()
	}
}

func (s *Scanner) handleThreeOptions(ifCmp rune, elifCmp rune, a Tok, b Tok, c Tok, buf bytes.Buffer) (tok Tok, lit string) {
	if ch := s.peek(); ch == ifCmp {
		buf.WriteRune(s.read())
		return a, buf.String()
	} else if ch == elifCmp {
		buf.WriteRune(s.read())
		return b, buf.String()
	} else {
		return c, buf.String()
	}
}

func (s *Scanner) scanNumber() (tok Tok, lit string) {

	var buf bytes.Buffer
	buf.WriteRune(s.read())
	dotUsed := false

	for {
		if ch := s.peek(); ch == '.' {
			if dotUsed {
				break
			} else {
				dotUsed = true
				buf.WriteRune(s.read())
				if next := s.peek(); !isDigit(next) {
					return ILLEGAL, buf.String()
				}
			}
		} else if isDigit(ch) {
			buf.WriteRune(s.read())
		} else {
			break
		}
	}

	return NUMBER, buf.String()
}

func (s *Scanner) scanStringLiteral() (tok Tok, lit string) {
	var buf bytes.Buffer
	s.read() //throw away string opener

	for {
		if ch := s.read(); ch == eof || isQuote(ch) { //throw away string closer
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return STRING, buf.String()
}

func (s *Scanner) discardWhitespace() {
	for {
		if ch := s.peek(); ch == eof || !isWhitespace(ch) {
			break
		}
		s.read()
	}
}

func (s *Scanner) scanIdent() (tok Tok, lit string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	for {
		if ch := s.peek(); ch == eof || (!isLetter(ch) && !isDigit(ch) && ch != '_') {
			break
		} else {
			_, _ = buf.WriteRune(s.read())
		}
	}

	return IDENTIFIER, buf.String()
}

func (s *Scanner) peek() rune {
	ch, _, err := s.r.ReadRune()
	_ = s.r.UnreadRune()
	if err != nil {
		return eof
	}
	return ch
}

func (s *Scanner) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}
	return ch
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

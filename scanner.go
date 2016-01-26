package presta

import (
	"bufio"
	"bytes"
	"io"
)

// Scanner represents a lexical scanner.
type Scanner struct {
	r *bufio.Reader
}

// NewScanner returns a new instance of Scanner.
func NewScanner(r io.Reader) *Scanner {
	return &Scanner{r: bufio.NewReader(r)}
}

// Scan returns the next token and literal value.
func (s *Scanner) Scan() (tok Tok, lit string) {
	// Read the next rune.
	ch := s.read()

	//Skip whitespace
	if isWhitespace(ch) {
		s.discardWhitespace()
		ch = s.read()
	}

	if isLetter(ch) {
		s.unread()
		return s.scanIdent()
	} else if isQuote(ch) {
		s.unread()
		return s.scanStringLiteral()
	} else if isSymbol(ch) {
		s.unread()
		return s.scanOperation()
	} else if isDigit(ch) {
		s.unread()
		return s.scanNumber()
	}

	// Otherwise read the individual character.
	switch ch {
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
		// case ',':
		// 	return COMMA, string(ch)
	}

	return ILLEGAL, string(ch)
}

func (s *Scanner) scanOperation() (tok Tok, lit string) {
	var buf bytes.Buffer
	ch := s.read()
	buf.WriteRune(ch)

	//{'>', '<', '=', '&', '|', '@', '^', '!', '+', '-', '/', '*', '%', '.'}

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
	ch := s.read()
	if ch == cmp {
		buf.WriteRune(ch)
		return yes, buf.String()
	} else {
		s.unread()
		return no, buf.String()
	}
}

func (s *Scanner) handleThreeOptions(ifCmp rune, elifCmp rune, a Tok, b Tok, c Tok, buf bytes.Buffer) (tok Tok, lit string) {
	ch := s.read()
	if ch == ifCmp {
		buf.WriteRune(ch)
		return a, buf.String()
	} else if ch == elifCmp {
		buf.WriteRune(ch)
		return b, buf.String()
	} else {
		s.unread()
		return c, buf.String()
	}
}

func (s *Scanner) scanNumber() (tok Tok, lit string) {

	var buf bytes.Buffer
	buf.WriteRune(s.read())
	dotUsed := false

	for {
		if ch := s.read(); ch == eof {
			break
		} else if ch == '.' {
			if dotUsed {
				s.unread()
				break
			} else {
				dotUsed = true
				buf.WriteRune(ch)
				if next := s.peek(); !isDigit(next) {
					return ILLEGAL, buf.String()
				}
			}
		} else if isDigit(ch) {
			buf.WriteRune(ch)
		} else {
			s.unread()
			break
		}
	}

	return NUMBER, buf.String()
}

func (s *Scanner) scanStringLiteral() (tok Tok, lit string) {
	// Create a buffer and throw away current character.
	var buf bytes.Buffer
	s.read()

	// Read every subsequent whitespace character into the buffer.
	// Non-whitespace characters and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof {
			break
		} else if isQuote(ch) {
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return STRING, buf.String()
}

// scanWhitespace consumes the current rune and all contiguous whitespace.
func (s *Scanner) discardWhitespace() {

	// Read every subsequent whitespace character into the buffer.
	// Non-whitespace characters and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isWhitespace(ch) {
			s.unread()
			break
		}
	}
}

// scanIdent consumes the current rune and all contiguous ident runes.
func (s *Scanner) scanIdent() (tok Tok, lit string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isLetter(ch) && !isDigit(ch) && ch != '_' {
			s.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
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

func (s *Scanner) unread() { _ = s.r.UnreadRune() }

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

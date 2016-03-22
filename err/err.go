package err

type ErrorCode int64

const (
	LEXICAL_ERROR ErrorCode = iota
	SYNTAX_ERROR
	SEMANTIC_ERROR
)

type Error interface {
	Message() string
	Code() ErrorCode
}

type SyntaxError struct {
	msg string
}

func NewSyntaxError(msg string) *SyntaxError {
	return &SyntaxError{msg: msg}
}

func (s *SyntaxError) Message() string {
	return s.msg
}

func (s *SyntaxError) Code() ErrorCode {
	return SYNTAX_ERROR
}

type SymanticError struct {
	msg string
}

func NewSymanticError(msg string) *SymanticError {
	return &SymanticError{msg: msg}
}

func (s *SymanticError) Message() string {
	return s.msg
}

func (s *SymanticError) Code() ErrorCode {
	return SEMANTIC_ERROR
}

type LexicalError struct {
	msg string
}

func NewLexicalError(msg string) *LexicalError {
	return &LexicalError{msg: msg}
}

func (l *LexicalError) Message() string {
	return l.msg
}

func (l *LexicalError) Code() ErrorCode {
	return LEXICAL_ERROR
}

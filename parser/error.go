package parser

type ErrorCode int64

const (
	SYNTAX_ERROR ErrorCode = iota
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

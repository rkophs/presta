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

package err

type ErrorCode int64

const (
	LEXICAL_ERROR ErrorCode = iota
	SYNTAX_ERROR
	SEMANTIC_ERROR
	RUNTIME_ERROR
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

type RuntimeError struct {
	msg string
}

func NewRuntimeError(msg string) *RuntimeError {
	return &RuntimeError{msg: msg}
}

func (r *RuntimeError) Message() string {
	return r.msg
}

func (r *RuntimeError) Code() ErrorCode {
	return RUNTIME_ERROR
}

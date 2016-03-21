package ir

/*This file will go into a separate package than ir */
import (
	"encoding/binary"
	"encoding/hex"
	"math"
	"strconv"
)

type Error struct {
	message string
	code    int
}

type Stack struct {
	stack []StackEntry
}

func (s *Stack) Push(a StackEntry) *Error {
	return nil
}

func (s *Stack) Pop() (StackEntry, *Error) {
	return nil, nil
}

func (s *Stack) FetchS(offset int) (StackEntry, *Error) {
	return nil, nil
}

func (s *Stack) FetchM(memAddr int) (StackEntry, *Error) {
	return nil, nil
}

func (s *Stack) FetchR(id int) (StackEntry, *Error) {
	return nil, nil
}

func (s *Stack) SetS(offset int, entry StackEntry) *Error {
	return nil
}

func (s *Stack) SetM(memAddr int, entry StackEntry) *Error {
	return nil
}

func (s *Stack) SetR(id int, entry StackEntry) *Error {
	return nil
}

func (s *Stack) New(memAddr int) *Error {
	return nil
}

func (s *Stack) Release(addr int) *Error {
	return nil
}

func (s *Stack) Shrink(result StackEntry) *Error {
	return nil //Save the result into %AX
}

func (s *Stack) Goto(offset int) *Error {
	return nil
}

/* Garbage to be removed */
func (s *Stack) GetNumber(addr int) (float64, *Error) {
	return (s.stack[addr]).ToNumber()
}

func (s *Stack) GetString(addr int) (string, *Error) {
	return (s.stack[addr]).ToString()
}

func (s *Stack) GetEntry(addr int) (StackEntry, *Error) {
	return s.stack[addr], nil
}

func (s *Stack) SetEntry(addr int, entry StackEntry) {
	s.stack[addr] = entry
}

type StackEntry interface {
	ToNumber() (float64, *Error)
	ToString() (string, *Error)
	ToHex() (string, *Error)
	Clone() StackEntry
	//ToArray() []StackEntry
}

type Number struct {
	number float64
}

func NewNumber(number float64) *Number {
	return &Number{number: number}
}

func (n *Number) SetNumber(number float64) {
	n.number = number
}

func (n *Number) ToNumber() (float64, *Error) {
	return n.number, nil
}

func (n *Number) ToString() (string, *Error) {
	return strconv.FormatFloat(n.number, 'E', -1, 64), nil
}

func (n *Number) ToHex() (string, *Error) {
	bits := math.Float64bits(n.number)
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, bits)
	return hex.EncodeToString(bytes), nil
}

func (n *Number) Clone() StackEntry {
	return NewNumber(n.number)
}

type String struct {
	str string
}

func NewString(str string) *String {
	return &String{str: str}
}

func (s *String) ToNumber() (float64, *Error) {
	return -1, &Error{message: "", code: -1}
}

func (s *String) ToString() (string, *Error) {
	return s.str, nil
}

func (s *String) ToHex() (string, *Error) {
	str := []byte(s.str)
	return hex.EncodeToString(str), nil
}

func (s *String) Clone() StackEntry {
	return NewString(s.str)
}

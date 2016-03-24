package system

import (
	"encoding/binary"
	"encoding/hex"
	"github.com/rkophs/presta/err"
	"math"
	"strconv"
)

type StackEntry interface {
	ToNumber() (float64, err.Error)
	ToString() (string, err.Error)
	ToHex() (string, err.Error)
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

func (n *Number) ToNumber() (float64, err.Error) {
	return n.number, nil
}

func (n *Number) ToString() (string, err.Error) {
	return strconv.FormatFloat(n.number, 'f', -1, 64), nil
}

func (n *Number) ToHex() (string, err.Error) {
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

func (s *String) ToNumber() (float64, err.Error) {
	return -1, err.NewRuntimeError("string type not convertable to number.")
}

func (s *String) ToString() (string, err.Error) {
	return s.str, nil
}

func (s *String) ToHex() (string, err.Error) {
	str := []byte(s.str)
	return hex.EncodeToString(str), nil
}

func (s *String) Clone() StackEntry {
	return NewString(s.str)
}

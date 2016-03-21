package json

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"math"
)

type Serializable interface {
	Serialize(*bytes.Buffer)
}

type KV struct {
	K string
	V Serializable
}

type Array struct {
	l []Serializable
}

type Number struct {
	n float64
}

type String struct {
	Serializable
	v string
}

func NewString(input string) *String {
	return &String{v: input}
}

func NewNumber(input float64) *Number {
	return &Number{n: input}
}

func NewArray(elems []Serializable) *Array {
	return &Array{l: elems}
}

func (s *String) Serialize(buffer *bytes.Buffer) {
	buffer.WriteRune('"')
	buffer.WriteString(s.v)
	buffer.WriteRune('"')
}

func (n *Number) Serialize(buffer *bytes.Buffer) {
	bits := math.Float64bits(n.n)
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, bits)

	buffer.WriteRune('"')
	buffer.WriteString(hex.EncodeToString(bytes))
	buffer.WriteRune('"')
}

func (a *Array) Serialize(buffer *bytes.Buffer) {
	buffer.WriteRune('[')
	last := len(a.l) - 1
	for it, elem := range a.l {
		elem.Serialize(buffer)
		if it != last {
			buffer.WriteRune(',')
		}
	}
	//buffer.UnreadRune()
	buffer.WriteRune(']')
}

func BuildMap(buffer *bytes.Buffer, tuples ...*KV) {
	buffer.WriteString("{")
	last := len(tuples) - 1
	for it, tuple := range tuples {
		buffer.WriteString("\"")
		buffer.WriteString(tuple.K)
		buffer.WriteString("\":")
		tuple.V.Serialize(buffer)
		if it != last {
			buffer.WriteRune(',')
		}
	}
	//buffer.UnreadRune()
	buffer.WriteString("}")
}

func BuildArray(buffer *bytes.Buffer, values []Serializable) {
	buffer.WriteRune('[')
	last := len(values) - 1
	for it, value := range values {
		value.Serialize(buffer)
		if it != last {
			buffer.WriteRune(',')
		}
	}
	//buffer.UnreadRune()
	buffer.WriteRune(']')
}

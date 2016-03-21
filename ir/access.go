package ir

import (
	"bytes"
	"fmt"
	"strconv"
)

type Accessor interface {
	Assign(s *Stack, entry StackEntry) *Error
	ToValue(s *Stack) (StackEntry, *Error)
	Serialize(buffer *bytes.Buffer)
	//ToArray()
}

/*=================================================================================*/
type StackAccess struct {
	offset int
}

func NewStackAccess(offset int) *StackAccess {
	fmt.Println(offset)
	return &StackAccess{offset: offset}
}

func (o *StackAccess) ToValue(s *Stack) (StackEntry, *Error) {
	return s.FetchS(o.offset)
}

func (o *StackAccess) Assign(s *Stack, entry StackEntry) *Error {
	return s.SetS(o.offset, entry)
}

func (s *StackAccess) Serialize(buffer *bytes.Buffer) {
	buffer.WriteString("BP(")
	num := s.offset
	if num < 0 {
		buffer.WriteString("-0x")
		num *= -1
	} else {
		buffer.WriteString("+0x")
	}
	buffer.WriteString(strconv.FormatInt(int64(num), 16))
	buffer.WriteRune(')')
}

/*=================================================================================*/
type MemoryAccess struct {
	addr int
}

func (m *MemoryAccess) ToValue(s *Stack) (StackEntry, *Error) {
	return s.FetchM(m.addr)
}

func (m *MemoryAccess) Assign(s *Stack, entry StackEntry) *Error {
	return s.SetM(m.addr, entry)
}

func (m *MemoryAccess) Release(s *Stack) *Error {
	return s.Release(m.addr)
}

func (m *MemoryAccess) New(s *Stack) *Error {
	return s.New(m.addr)
}

func (m *MemoryAccess) Serialize(buffer *bytes.Buffer) {
	buffer.WriteString("M(0x")
	buffer.WriteString(strconv.FormatInt(int64(m.addr), 16))
	buffer.WriteRune(')')
}

/*=================================================================================*/
type RegisterAccess struct {
	id int
}

func NewRegisterAccess(id int) *RegisterAccess {
	return &RegisterAccess{id: id}
}

func (r *RegisterAccess) ToValue(s *Stack) (StackEntry, *Error) {
	return s.FetchR(r.id)
}

func (r *RegisterAccess) Assign(s *Stack, entry StackEntry) *Error {
	return s.SetR(r.id, entry)
}

func (r *RegisterAccess) Serialize(buffer *bytes.Buffer) {
	buffer.WriteString("%")
	buffer.WriteString(strconv.FormatInt(int64(r.id), 16))
}

/*=================================================================================*/
type ConstantAccess struct {
	entry StackEntry
}

func NewConstantAccess(entry StackEntry) *ConstantAccess {
	return &ConstantAccess{entry: entry}
}

func (c *ConstantAccess) ToValue(s *Stack) (StackEntry, *Error) {
	return c.entry, nil
}

func (c *ConstantAccess) Assign(s *Stack, entry StackEntry) *Error {
	return &Error{message: "Cannot reassign a constant."}
}

func (c *ConstantAccess) Serialize(buffer *bytes.Buffer) {
	hex, _ := c.entry.ToHex()
	buffer.WriteString(hex)
}

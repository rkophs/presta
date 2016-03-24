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

package ir

import (
	"bytes"
	"github.com/rkophs/presta/system"
	"strconv"
)

type Accessor interface {
	Assign(s system.System, entry system.StackEntry)
	ToValue(s system.System) system.StackEntry
	Serialize(buffer *bytes.Buffer)
	//ToArray()
}

/*=================================================================================*/
type StackAccess struct {
	offset int
}

func NewStackAccess(offset int) *StackAccess {
	return &StackAccess{offset: offset}
}

func (o *StackAccess) ToValue(s system.System) system.StackEntry {
	return s.FetchS(o.offset)
}

func (o *StackAccess) Assign(s system.System, entry system.StackEntry) {
	s.SetS(o.offset, entry)
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

func (m *MemoryAccess) ToValue(s system.System) system.StackEntry {
	return s.FetchM(m.addr)
}

func (m *MemoryAccess) Assign(s system.System, entry system.StackEntry) {
	s.SetM(m.addr, entry)
}

func (m *MemoryAccess) Release(s system.System) {
	s.Release(m.addr)
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

func (r *RegisterAccess) ToValue(s system.System) system.StackEntry {
	return s.FetchR(r.id)
}

func (r *RegisterAccess) Assign(s system.System, entry system.StackEntry) {
	s.SetR(r.id, entry)
}

func (r *RegisterAccess) Serialize(buffer *bytes.Buffer) {
	buffer.WriteString("%")
	buffer.WriteString(strconv.FormatInt(int64(r.id), 16))
}

/*=================================================================================*/
type ConstantAccess struct {
	entry system.StackEntry
}

func NewConstantAccess(entry system.StackEntry) *ConstantAccess {
	return &ConstantAccess{entry: entry}
}

func (c *ConstantAccess) ToValue(s system.System) system.StackEntry {
	return c.entry
}

func (c *ConstantAccess) Assign(s system.System, entry system.StackEntry) {
	s.SetError("Cannot reasign a constant accessor.")
}

func (c *ConstantAccess) Serialize(buffer *bytes.Buffer) {
	hex, _ := c.entry.ToHex()
	buffer.WriteString(hex)
}

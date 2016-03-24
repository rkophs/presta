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
	"github.com/rkophs/presta/err"
	"github.com/rkophs/presta/system"
	"strconv"
)

type InstructionLocation struct {
	location int
}

func NewInstructionLocation(input int) *InstructionLocation {
	return &InstructionLocation{location: input}
}

func (i *InstructionLocation) SetLocation(location int) {
	i.location = location
}

func (i *InstructionLocation) GetLocation() int {
	return i.location
}

func (i *InstructionLocation) Serialize(buffer *bytes.Buffer) {
	buffer.WriteString("0x")
	buffer.WriteString(strconv.FormatInt(int64(i.location), 16))
}

func mergeErrors(errs ...err.Error) err.Error {
	for _, e := range errs {
		if e != nil {
			return e
		}
	}
	return nil
}

func writeInstr(buffer *bytes.Buffer, instr string, params ...string) {
	buffer.WriteString(instr)
	buffer.WriteRune('\t')

	for _, param := range params {
		buffer.WriteString(param)
		buffer.WriteRune(',')
	}
	buffer.WriteRune('\n')
}

/*====================================================================================*/

type Instruction interface {
	Serialize(buffer *bytes.Buffer)
	Execute(system.System)
}

type InstructionType byte

const (
	ADD InstructionType = iota
	PUSH
	RELEASE
	POP
	NEW
	MOV
	CALL
	RESULT
)

type Add struct {
	l Accessor
	r Accessor
}

func NewAdd(l, r Accessor) *Add {
	return &Add{l: l, r: r}
}

func (a *Add) Execute(s system.System) {
	l := a.l.ToValue(s)
	r := a.r.ToValue(s)

	lv, e := l.ToNumber()
	if e != nil {
		s.SetError("Addition requires 2 numbers")
		return
	}
	rv, e := r.ToNumber()
	if e != nil {
		s.SetError("Addition requires 2 numbers")
	}
	a.l.Assign(s, system.NewNumber(lv+rv))
}

func (a *Add) Serialize(buffer *bytes.Buffer) {
	buffer.WriteString("add\t")
	a.l.Serialize(buffer)
	buffer.WriteRune(',')
	a.r.Serialize(buffer)
	buffer.WriteRune('\n')
}

type Push struct {
	v Accessor
}

func NewPush(v Accessor) *Push {
	return &Push{v: v}
}

func (p *Push) Execute(s system.System) {
	s.Push(p.v.ToValue(s))
}

func (p *Push) Serialize(buffer *bytes.Buffer) {
	buffer.WriteString("push\t")
	p.v.Serialize(buffer)
	buffer.WriteRune('\n')
}

type Release struct {
	v *MemoryAccess
}

func (r *Release) Execute(s system.System) {
	r.v.Release(s)
}

func (r *Release) Serialize(buffer *bytes.Buffer) {
	buffer.WriteString("push\t")
	r.v.Serialize(buffer)
	buffer.WriteRune('\n')
}

type Mov struct {
	l Accessor
	r Accessor
}

func NewMov(l Accessor, r Accessor) *Mov {
	return &Mov{l: l, r: r}
}

func (m *Mov) Execute(s system.System) {
	m.l.Assign(s, m.r.ToValue(s))
}

func (m *Mov) Serialize(buffer *bytes.Buffer) {
	buffer.WriteString("mov\t")
	m.l.Serialize(buffer)
	buffer.WriteRune(',')
	m.r.Serialize(buffer)
	buffer.WriteRune('\n')
}

type Call struct {
	location *InstructionLocation
}

func NewCall(location *InstructionLocation) *Call {
	return &Call{location: location}
}

func (c *Call) Execute(s system.System) {
	s.Call(c.location.GetLocation())
}

func (c *Call) Serialize(buffer *bytes.Buffer) {
	buffer.WriteString("call\t")
	c.location.Serialize(buffer)
	buffer.WriteRune('\n')
}

type Result struct {
	from Accessor
}

func NewResult(from Accessor) *Result {
	return &Result{from: from}
}

func (r *Result) Execute(s system.System) {
	s.Return(r.from.ToValue(s))
}

func (r *Result) Serialize(buffer *bytes.Buffer) {
	buffer.WriteString("ret\t")
	r.from.Serialize(buffer)
	buffer.WriteRune('\n')
}

type Exit struct {
	from Accessor
}

func NewExit(from Accessor) *Exit {
	return &Exit{from: from}
}

func (r *Exit) Execute(s system.System) {
	s.Exit(r.from.ToValue(s))
}

func (r *Exit) Serialize(buffer *bytes.Buffer) {
	buffer.WriteString("exit\t")
	r.from.Serialize(buffer)
	buffer.WriteRune('\n')
}

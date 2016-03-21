package ir

import (
	"bytes"
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

func mergeErrors(errs ...*Error) *Error {
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
	Execute(*Stack) *Error
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

func (a *Add) Execute(s *Stack) *Error {
	l, el := a.l.ToValue(s)
	lv, elv := l.ToNumber()
	r, er := a.r.ToValue(s)
	rv, erv := r.ToNumber()
	if err := mergeErrors(el, elv, er, erv); err != nil {
		return err
	}
	return a.l.Assign(s, NewNumber(lv+rv))
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

func (p *Push) Execute(s *Stack) *Error {
	if v, e := p.v.ToValue(s); e != nil {
		return e
	} else {
		return s.Push(v)
	}
}

func (p *Push) Serialize(buffer *bytes.Buffer) {
	buffer.WriteString("push\t")
	p.v.Serialize(buffer)
	buffer.WriteRune('\n')
}

type Release struct {
	v *MemoryAccess
}

func (r *Release) Execute(s *Stack) *Error {
	return r.v.Release(s)
}

func (r *Release) Serialize(buffer *bytes.Buffer) {
	buffer.WriteString("push\t")
	r.v.Serialize(buffer)
	buffer.WriteRune('\n')
}

type New struct {
	v *MemoryAccess
}

func (n *New) Execute(s *Stack) *Error {
	return n.v.New(s)
}

func (n *New) Serialize(buffer *bytes.Buffer) {
	buffer.WriteString("new\t")
	n.v.Serialize(buffer)
	buffer.WriteRune('\n')
}

type Pop struct {
	v Accessor
}

func (p *Pop) Execute(s *Stack) *Error {
	if entry, e := s.Pop(); e != nil {
		return e
	} else {
		return p.v.Assign(s, entry)
	}
}

func (p *Pop) Serialize(buffer *bytes.Buffer) {
	buffer.WriteString("pop\t")
	p.v.Serialize(buffer)
	buffer.WriteRune('\n')
}

type Mov struct {
	l Accessor
	r Accessor
}

func NewMov(l Accessor, r Accessor) *Mov {
	return &Mov{l: l, r: r}
}

func (m *Mov) Execute(s *Stack) *Error {
	if v, e := m.r.ToValue(s); e != nil {
		return e
	} else {
		return m.l.Assign(s, v)
	}
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

func (c *Call) Execute(s *Stack) *Error {
	return s.Goto(c.location.GetLocation())
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

func (r *Result) Execute(s *Stack) *Error {
	if from, e := r.from.ToValue(s); e != nil {
		return e
	} else {
		return s.Shrink(from)
	}
}

func (r *Result) Serialize(buffer *bytes.Buffer) {
	buffer.WriteString("ret\t")
	r.from.Serialize(buffer)
	buffer.WriteRune('\n')
}

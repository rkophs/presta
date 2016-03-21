package ir

import (
	"bytes"
)

type Instruction2 interface {
	//Serializable
	GetNunber() (float64, *Error)
	GetString() (string, *Error)
	GetValue() (StackEntry, *Error)
	// GetBoolean() bool
}

/*=================================================================================*/
type Concat struct {
	elements []Instruction2
}

func NewConcat(elems []Instruction2) *Concat {
	return &Concat{elements: elems}
}

func (c *Concat) execute() (string, *Error) {
	var buffer bytes.Buffer
	for _, e := range c.elements {
		if str, err := e.GetString(); err != nil {
			return "", &Error{message: "", code: -1}
		} else {
			buffer.WriteString(str)
		}
	}
	return buffer.String(), nil
}

func (c *Concat) GetNumber() (float64, *Error) {
	return -1, &Error{message: "", code: -1}
}

func (c *Concat) GetString() (string, *Error) {
	return c.execute()
}

func (c *Concat) GetValue() (StackEntry, *Error) {
	if str, err := c.execute(); err != nil {
		return nil, err
	} else {
		return NewString(str), nil
	}
}

/*=================================================================================*/
type Variable struct {
	memAddr int
	stack   *Stack
}

func NewVariable(memAddr int, stack *Stack) *Variable {
	return &Variable{memAddr: memAddr, stack: stack}
}

func (v *Variable) GetNumber() (float64, *Error) {
	return v.stack.GetNumber(v.memAddr)
}

func (v *Variable) GetString() (string, *Error) {
	return v.stack.GetString(v.memAddr)
}

func (v *Variable) GetValue() (StackEntry, *Error) {
	return v.stack.GetEntry(v.memAddr)
}

/*=================================================================================*/
type Let struct {
	stack  *Stack
	values []Instruction2
	toArrs []int
	body   Instruction2
}

func NewLet(values []Instruction2, toArrs []int, body Instruction2, stack *Stack) *Let {
	return &Let{values: values, toArrs: toArrs, body: body, stack: stack}
}

func (l *Let) assignScope() *Error {
	for i, v := range l.values {
		if val, err := v.GetValue(); err != nil {
			return err
		} else {
			l.stack.SetEntry(l.toArrs[i], val.Clone())
		}
	}
	return nil
}

func (l *Let) GetNumber() (float64, *Error) {
	if err := l.assignScope(); err != nil {
		return -1, err
	}
	return l.body.GetNunber()
}

func (l *Let) GetString() (string, *Error) {
	if err := l.assignScope(); err != nil {
		return "", err
	}
	return l.body.GetString()
}

func (l *Let) GetValue() (StackEntry, *Error) {
	if err := l.assignScope(); err != nil {
		return nil, err
	}
	return l.body.GetValue()
}

/*=================================================================================*/
type Assign struct {
	value  Instruction2
	toAddr int
	stack  *Stack
}

func NewAssign(value Instruction2, toAddr int, stack *Stack) *Assign {
	return &Assign{value: value, toAddr: toAddr, stack: stack}
}

func (a *Assign) execute() *Error {
	if val, err := a.value.GetValue(); err != nil {
		return err
	} else {
		a.stack.SetEntry(a.toAddr, val.Clone())
	}
	return nil
}

func (a *Assign) GetNumber() (float64, *Error) {
	if err := a.execute(); err != nil {
		return -1, err
	} else {
		return a.stack.GetNumber(a.toAddr)
	}
}

func (a *Assign) GetString() (string, *Error) {
	if err := a.execute(); err != nil {
		return "", err
	} else {
		return a.stack.GetString(a.toAddr)
	}
}

func (a *Assign) GetValue() (StackEntry, *Error) {
	if err := a.execute(); err != nil {
		return nil, err
	} else {
		return a.stack.GetEntry(a.toAddr)
	}
}

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

package vm

import (
	"fmt"
	"github.com/rkophs/presta/err"
	"github.com/rkophs/presta/ir"
	"github.com/rkophs/presta/system"
)

type VM struct {
	stack     *Stack
	heap      *Heap
	flow      *Flow
	registers []system.StackEntry
	exited    bool
	err       err.Error
	interrupt bool
}

func NewVM(instructions []ir.Instruction) *VM {
	return &VM{
		flow:      NewFlow(instructions),
		stack:     NewStack(),
		heap:      NewHeap(),
		interrupt: false,
		err:       nil,
		registers: make([]system.StackEntry, 1),
		exited:    false,
	}
}

func (v *VM) Run() err.Error {
	for !v.exited && !v.interrupt {
		v.Print()
		v.flow.Execute(v)
	}
	fmt.Println("Complete:")
	v.Print()
	return v.err
}

func (v *VM) Push(a system.StackEntry) {
	v.stack.Push(a)
}

func (v *VM) Print() {
	fmt.Println("============")
	var s string
	if len(v.registers) > 0 && v.registers[0] != nil {
		if k, e := v.registers[0].ToString(); e != nil {
			s = ""
		} else {
			s = k
		}
	} else {
		s = ""
	}
	fmt.Println("PC: ", v.flow.pc, " BP: ", v.stack.bp, " SP: ", v.stack.sp, " AX: ", s)
	fmt.Println("Stack:")
	for i, v := range v.stack.stack {
		s, _ := v.ToString()
		fmt.Println(i, " ", s)
	}
	fmt.Println("Heap:")
	for k, v := range v.heap.heap {
		s, _ := v.ToString()
		fmt.Println(k, " ", s)
	}
	fmt.Println("============")

}

func (v *VM) FetchS(offset int) system.StackEntry {
	return v.stack.Fetch(offset)
}

func (v *VM) FetchM(memAddr int) system.StackEntry {
	return v.heap.Fetch(memAddr)
}

func (v *VM) FetchR(id int) system.StackEntry {
	return v.registers[id]
}

func (v *VM) SetS(offset int, entry system.StackEntry) {
	v.stack.Set(offset, entry)
}

func (v *VM) SetM(memAddr int, entry system.StackEntry) {
	v.heap.Set(memAddr, entry)
}

func (v *VM) SetR(id int, entry system.StackEntry) {
	v.registers[id] = entry
}

func (v *VM) Release(addr int) {
	v.heap.Release(addr)
}

func (v *VM) Return(result system.StackEntry) {
	v.stack.PopFrame()
	v.flow.Return()
	v.registers[0] = result
}

func (v *VM) Exit(result system.StackEntry) {
	v.exited = true
	v.stack.PopFrame()
	v.flow.Return()
	v.registers[0] = result
}

func (v *VM) SetError(e string) {
	v.interrupt = true
	v.err = err.NewRuntimeError(e)
}

func (v *VM) Call(offset int) {
	v.flow.Call(offset)
	v.stack.PushFrame()
}

func (v *VM) Goto(offset int) {
	v.flow.GoTo(offset)
}

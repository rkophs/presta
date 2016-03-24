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

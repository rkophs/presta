package vm

import (
	"fmt"
	"github.com/rkophs/presta/err"
	"github.com/rkophs/presta/ir"
	"github.com/rkophs/presta/system"
)

type VM struct {
	instr     []ir.Instruction
	stack     []system.StackEntry
	heap      map[int]system.StackEntry
	pc        int
	funcs     []int
	frames    []int
	bp        int
	sp        int
	registers []system.StackEntry
	exited    bool
}

func NewVM(instructions []ir.Instruction) *VM {
	return &VM{
		instr:     instructions,
		stack:     []system.StackEntry{},
		heap:      make(map[int]system.StackEntry),
		pc:        0,
		funcs:     []int{},
		frames:    []int{},
		bp:        0,
		sp:        0,
		registers: make([]system.StackEntry, 1),
		exited:    false,
	}
}

func (v *VM) Run() err.Error {
	for !v.exited {
		v.Print()
		if e := v.instr[v.pc].Execute(v); e != nil {
			return e
		}
		v.pc++
	}
	v.Print()
	return nil
}

func (v *VM) Push(a system.StackEntry) err.Error {
	v.stack = append(v.stack, a)
	v.sp++
	return nil
}

func (v *VM) Pop() (system.StackEntry, err.Error) {
	v.sp--
	ret := v.stack[v.sp]
	v.stack = v.stack[:v.sp]
	return ret, nil
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
	fmt.Println("PC: ", v.pc, " BP: ", v.bp, " SP: ", v.sp, " AX: ", s)
	fmt.Println("Stack:")
	for i, v := range v.stack {
		s, _ := v.ToString()
		fmt.Println(i, " ", s)
	}
	fmt.Println("Heap:")
	for k, v := range v.heap {
		s, _ := v.ToString()
		fmt.Println(k, " ", s)
	}
	fmt.Println("============")

}

func (v *VM) FetchS(offset int) (system.StackEntry, err.Error) {
	return v.stack[v.bp+offset], nil
}

func (v *VM) FetchM(memAddr int) (system.StackEntry, err.Error) {
	return v.heap[memAddr], nil
}

func (v *VM) FetchR(id int) (system.StackEntry, err.Error) {
	return v.registers[id], nil
}

func (v *VM) SetS(offset int, entry system.StackEntry) err.Error {
	v.stack[v.bp+offset] = entry
	return nil
}

func (v *VM) SetM(memAddr int, entry system.StackEntry) err.Error {
	v.heap[memAddr] = entry
	return nil
}

func (v *VM) SetR(id int, entry system.StackEntry) err.Error {
	v.registers[id] = entry
	return nil
}

func (v *VM) New(memAddr int) err.Error {
	v.heap[memAddr] = nil
	return nil
}

func (v *VM) Release(addr int) err.Error {
	delete(v.heap, addr)
	return nil
}

func (v *VM) Shrink(result system.StackEntry) err.Error {
	v.stack = v.stack[:v.bp]
	v.sp = v.bp
	v.bp = v.frames[len(v.frames)-1]
	v.registers[0] = result

	if pc_len := (len(v.funcs) - 1); pc_len >= 0 {
		v.pc = v.funcs[pc_len]
		v.funcs = v.funcs[:pc_len]
	}
	return nil
}

func (v *VM) Exit() err.Error {
	v.exited = true
	return nil
}

func (v *VM) Expand() err.Error {
	v.frames = append(v.frames, v.bp)
	v.bp = v.sp
	return nil
}

func (v *VM) Goto(offset int) err.Error {
	v.funcs = append(v.funcs, v.pc)
	fmt.Println("Going to: ", offset-1)
	v.pc = offset - 1
	return nil
}

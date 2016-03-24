package vm

import (
	"github.com/rkophs/presta/ir"
	"github.com/rkophs/presta/system"
)

type Flow struct {
	instr []ir.Instruction
	pc    int
	funcs []int
}

func NewFlow(instr []ir.Instruction) *Flow {
	return &Flow{instr: instr, pc: 0, funcs: []int{-1}}
}

func (f *Flow) Execute(v system.System) {
	f.instr[f.pc].Execute(v)
	f.pc++
}

func (f *Flow) Return() {
	pc_len := (len(f.funcs) - 1)
	f.pc = f.funcs[pc_len]
	f.funcs = f.funcs[:pc_len]
}

func (f *Flow) Call(offset int) {
	f.funcs = append(f.funcs, f.pc)
	f.pc = offset - 1
}

func (f *Flow) GoTo(offset int) {
	f.pc = offset - 1
}

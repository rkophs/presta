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

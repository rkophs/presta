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

package icg

import (
	"bytes"
	"fmt"
	"github.com/rkophs/presta/ir"
)

type Error struct {
	message string
	code    int
}

type Code struct {
	instructions []ir.Instruction
	linker       *Linker             //FunctionId -> instruction offset
	vars         map[int]ir.Accessor //varId -> access location
	Ax           *ir.RegisterAccess
	count        int
	frameOffset  int
}

func NewCode(linker *Linker) *Code {
	return &Code{
		instructions: make([]ir.Instruction, 0),
		linker:       linker,
		vars:         make(map[int]ir.Accessor),
		count:        0,
		frameOffset:  0,
		Ax:           ir.NewRegisterAccess(0),
	}
}

func (c *Code) Append(elem ir.Instruction) {
	c.instructions = append(c.instructions, elem)
	c.count++
}

func (c *Code) GetCount() int {
	return c.count
}

func (c *Code) GetFrameOffset() int {
	return c.frameOffset
}

func (c *Code) IncrFrameOffset(amount int) {
	c.frameOffset += amount
}

func (c *Code) ResetFrameOffset(amount int) {
	c.frameOffset = 0
}

func (c *Code) AppendBlock(block *Code) {
	c.instructions = append(c.instructions, block.instructions...)
	c.count += block.count
}

func (c *Code) SetFunctionOffset(id string, offset *ir.InstructionLocation) {
	c.linker.SetFunctionOffset(id, offset)
}

func (c *Code) GetFunctionOffset(id string) *ir.InstructionLocation {
	return c.linker.GetFunctionOffset(id)
}

func (c *Code) GetLinker() *Linker {
	return c.linker
}

func (c *Code) SetVariable(id int, location ir.Accessor) {
	c.vars[id] = location
}

func (c *Code) GetVariableLocation(id int) ir.Accessor {
	return c.vars[id]
}

func (c *Code) GetInstructions() []ir.Instruction {
	return c.instructions
}

func (c *Code) Serialize(buffer *bytes.Buffer) {
	for i, instr := range c.instructions {
		fmt.Println(i)
		instr.Serialize(buffer)
	}
}

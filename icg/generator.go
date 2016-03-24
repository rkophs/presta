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

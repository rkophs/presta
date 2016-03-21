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
	linker       map[string]*ir.InstructionLocation //FunctionId -> instruction offset
	vars         map[int]ir.Accessor                //varId -> access location
	Ax           *ir.RegisterAccess
	count        int
	frameOffset  int
}

func NewCode() *Code {
	return &Code{
		instructions: make([]ir.Instruction, 0),
		linker:       make(map[string]*ir.InstructionLocation),
		vars:         make(map[int]ir.Accessor),
		count:        0,
		frameOffset:  0,
		Ax:           ir.NewRegisterAccess(0),
	}
}

func (c *Code) Append(elem ir.Instruction) {
	fmt.Println("Calling append")
	if elem == nil {
		fmt.Println("Nil input")
	}
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
	fmt.Println("Appending block")
	fmt.Println(block.GetCount())
	c.instructions = append(c.instructions, block.instructions...)
	c.count += block.count
}

func (c *Code) SetFunctionOffset(id string, offset *ir.InstructionLocation) {
	c.linker[id] = offset
}

func (c *Code) GetFunctionOffset(id string) *ir.InstructionLocation {
	return c.linker[id]
}

func (c *Code) SetVariable(id int, location ir.Accessor) {
	c.vars[id] = location
}

func (c *Code) GetVariableLocation(id int) ir.Accessor {
	return c.vars[id]
}

func (c *Code) Serialize(buffer *bytes.Buffer) {
	for _, instr := range c.instructions {
		if instr == nil {
			fmt.Println("nil instr")
		}
		instr.Serialize(buffer)
	}
}

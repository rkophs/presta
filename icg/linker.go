package icg

import (
	"github.com/rkophs/presta/ir"
)

type Linker struct {
	linker map[string]*ir.InstructionLocation
}

func NewLinker() *Linker {
	return &Linker{linker: make(map[string]*ir.InstructionLocation)}
}

func (c *Linker) SetFunctionOffset(id string, offset *ir.InstructionLocation) {
	c.linker[id] = offset
}

func (c *Linker) GetFunctionOffset(id string) *ir.InstructionLocation {
	return c.linker[id]
}

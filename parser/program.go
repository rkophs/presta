package parser

import (
	"bytes"
	"github.com/rkophs/presta/icg"
	"github.com/rkophs/presta/ir"
	"github.com/rkophs/presta/json"
	"github.com/rkophs/presta/semantic"
)

type Program struct {
	funcs []*Function
	exec  AstNode
}

func (p *Program) Type() AstNodeType {
	return PROG
}

func NewProgram(p *Parser) (tree AstNode, err Error) {
	readCount := 0

	/*Check for function declarations*/
	functions := []*Function{}
	for {
		if function, err := NewFunction(p); err != nil {
			return p.parseError(err.Message(), readCount)
		} else if function != nil {
			functions = append(functions, function.(*Function))
		} else {
			break
		}
	}

	/*Check for exec*/
	expr, err := NewExpression(p)
	if err != nil {
		return p.parseError(err.Message(), readCount)
	} else if expr == nil {
		return p.parseError("Program must contain an executable expression", readCount)
	}

	program := &Program{funcs: functions, exec: expr}
	return p.parseValid(program)
}

func (p *Program) Serialize(buffer *bytes.Buffer) {

	fns := []json.Serializable{}
	for _, fn := range p.funcs {
		fns = append(fns, fn)
	}

	json.BuildMap(buffer,
		&json.KV{K: "functions", V: json.NewArray(fns)},
		&json.KV{K: "body", V: p.exec},
		&json.KV{K: "type", V: json.NewString("PROG")})
}

func (p *Program) GenerateICG(code *icg.Code, s *semantic.Semantic) Error {

	/* Add function linker symbols */
	for _, f := range p.funcs {
		s.AddFunction(f.name, len(f.params))
		code.SetFunctionOffset(f.name, ir.NewInstructionLocation(-1))
	}

	if err := p.exec.GenerateICG(code, s); err != nil {
		return err
	}

	//Return result & shrink stack
	code.Append(ir.NewResult(code.Ax))

	//Concatenate instruction lists and set correct function offsets
	offset := code.GetCount()
	for _, f := range p.funcs {
		code.GetFunctionOffset(f.name).SetLocation(offset)
		fnBlock := icg.NewCode()
		if err := f.GenerateICG(fnBlock, s); err != nil {
			return err
		}
		code.AppendBlock(fnBlock)
		offset += fnBlock.GetCount()
	}

	return nil
}

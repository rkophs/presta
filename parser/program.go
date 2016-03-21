package parser

import (
	"bytes"
	"fmt"
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

func (p *Parser) program() (tree *Program, yes bool, err Error) {
	readCount := 0
	var invalid Program

	/*Check for function declarations*/
	functions := []*Function{}
	for {
		if function, yes, err := p.function(); err != nil {
			_, yes, err := p.parseError(err.Message(), readCount)
			return &invalid, yes, err
		} else if yes {
			functions = append(functions, function)
		} else {
			break
		}
	}

	/*Check for exec*/
	expr, yes, err := p.expression()
	if err != nil {
		_, yes, err := p.parseError(err.Message(), readCount)
		return &invalid, yes, err
	} else if !yes {
		_, yes, err := p.parseError("Program must contain an executable expression", readCount)
		return &invalid, yes, err
	}

	program := &Program{funcs: functions, exec: expr}
	_, yes, err = p.parseValid(program)
	return program, yes, err
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

func (p *Program) GenerateICG(doNotUse int64, code *icg.Code, s *semantic.Semantic) (int64, Error) {

	fmt.Println("ICG for program")

	/* Add function linker symbols */
	for _, f := range p.funcs {
		s.AddFunction(f.name, len(f.params))
		code.SetFunctionOffset(f.name, ir.NewInstructionLocation(-1))
	}

	if _, err := p.exec.GenerateICG(-1, code, s); err != nil {
		return -1, err
	}

	//Return result & shrink stack
	code.Append(ir.NewResult(code.Ax))

	//Concatenate instruction lists and set correct function offsets
	offset := code.GetCount()
	for _, f := range p.funcs {
		code.GetFunctionOffset(f.name).SetLocation(offset)
		fnBlock := icg.NewCode()
		if _, err := f.GenerateICG(-1, fnBlock, s); err != nil {
			return -1, err
		}
		code.AppendBlock(fnBlock)
		offset += fnBlock.GetCount()
	}

	return -1, nil
}

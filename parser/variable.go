package parser

import (
	"bytes"
	"github.com/rkophs/presta/err"
	"github.com/rkophs/presta/icg"
	"github.com/rkophs/presta/ir"
	"github.com/rkophs/presta/json"
	"github.com/rkophs/presta/semantic"
)

type Variable struct {
	name string
}

func (v *Variable) Type() AstNodeType {
	return VAR
}

func (v *Variable) Serialize(buffer *bytes.Buffer) {

	json.BuildMap(buffer,
		&json.KV{K: "name", V: json.NewString(v.name)},
		&json.KV{K: "type", V: json.NewString("VAR")})
}

func (v *Variable) GenerateICG(code *icg.Code, s *semantic.Semantic) err.Error {
	if !s.VariableExists(v.name) {
		return err.NewSymanticError("Undefined variable.")
	}

	code.Append(ir.NewMov(code.Ax, code.GetVariableLocation(s.GetVariableId(v.name))))
	return nil
}

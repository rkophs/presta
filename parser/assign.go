package parser

import (
	"bytes"
	"github.com/rkophs/presta/icg"
	"github.com/rkophs/presta/json"
	"github.com/rkophs/presta/semantic"
)

type Assign struct {
	name  string
	value AstNode
}

func (a *Assign) Type() AstNodeType {
	return ASSIGN
}

func (a *Assign) Serialize(buffer *bytes.Buffer) {

	json.BuildMap(buffer,
		&json.KV{K: "name", V: json.NewString(a.name)},
		&json.KV{K: "value", V: a.value},
		&json.KV{K: "type", V: json.NewString("ASSIGN")})
}

func (p *Assign) GenerateICG(offset int64, code *icg.Code, s *semantic.Semantic) (int64, Error) {
	return -1, nil
}

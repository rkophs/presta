package parser

import (
	"bytes"
	"github.com/rkophs/presta/icg"
	"github.com/rkophs/presta/json"
	"github.com/rkophs/presta/semantic"
)

type Not struct {
	exec AstNode
}

func (n *Not) Type() AstNodeType {
	return NOT
}

func (n *Not) Serialize(buffer *bytes.Buffer) {

	json.BuildMap(buffer,
		&json.KV{K: "expression", V: n.exec},
		&json.KV{K: "type", V: json.NewString("NOT")})
}

func (n *Not) GenerateICG(offset int64, code *icg.Code, s *semantic.Semantic) (int64, Error) {
	return -1, nil
}

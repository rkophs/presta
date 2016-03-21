package parser

import (
	"bytes"
	"github.com/rkophs/presta/icg"
	"github.com/rkophs/presta/json"
	"github.com/rkophs/presta/semantic"
)

type Repeat struct {
	condition AstNode
	exec      AstNode
}

func (r *Repeat) Type() AstNodeType {
	return REPEAT
}

func (r *Repeat) Serialize(buffer *bytes.Buffer) {

	json.BuildMap(buffer,
		&json.KV{K: "condition", V: r.condition},
		&json.KV{K: "body", V: r.exec},
		&json.KV{K: "type", V: json.NewString("REPEAT")})
}

func (r *Repeat) GenerateICG(offset int64, code *icg.Code, s *semantic.Semantic) (int64, Error) {
	return -1, nil
}

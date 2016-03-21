package parser

import (
	"bytes"
	"github.com/rkophs/presta/icg"
	"github.com/rkophs/presta/json"
	"github.com/rkophs/presta/semantic"
)

type Let struct {
	params []string
	values []AstNode
	exec   AstNode
}

func (l *Let) Type() AstNodeType {
	return LET
}

func (l *Let) Serialize(buffer *bytes.Buffer) {

	params := []json.Serializable{}
	for _, param := range l.params {
		params = append(params, json.NewString(param))
	}

	values := []json.Serializable{}
	for _, value := range l.values {
		values = append(values, value)
	}

	json.BuildMap(buffer,
		&json.KV{K: "names", V: json.NewArray(params)},
		&json.KV{K: "values", V: json.NewArray(values)},
		&json.KV{K: "body", V: l.exec},
		&json.KV{K: "type", V: json.NewString("LET")})
}

func (l *Let) GenerateICG(offset int64, code *icg.Code, s *semantic.Semantic) (int64, Error) {
	return -1, nil
}

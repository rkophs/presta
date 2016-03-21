package parser

import (
	"bytes"
	"github.com/rkophs/presta/icg"
	"github.com/rkophs/presta/json"
	"github.com/rkophs/presta/semantic"
)

type Concat struct {
	components []AstNode
}

func (c *Concat) Type() AstNodeType {
	return CONCAT
}

func (c *Concat) Serialize(buffer *bytes.Buffer) {
	components := []json.Serializable{}
	for _, component := range c.components {
		components = append(components, component)
	}

	json.BuildMap(buffer,
		&json.KV{K: "chunks", V: json.NewArray(components)},
		&json.KV{K: "type", V: json.NewString("CONCAT")})
}

func (c *Concat) GenerateICG(offset int64, code *icg.Code, s *semantic.Semantic) (int64, Error) {
	return -1, nil
}

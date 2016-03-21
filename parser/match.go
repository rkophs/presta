package parser

import (
	"bytes"
	"github.com/rkophs/presta/icg"
	"github.com/rkophs/presta/json"
	"github.com/rkophs/presta/semantic"
)

type Match struct {
	conditions []AstNode
	branches   []AstNode
	matchType  MatchType
}

func (m *Match) Type() AstNodeType {
	return MATCH
}

func (m *Match) Serialize(buffer *bytes.Buffer) {

	branches := []json.Serializable{}
	for _, branch := range m.branches {
		branches = append(branches, branch)
	}

	conditions := []json.Serializable{}
	for _, condition := range m.conditions {
		conditions = append(conditions, condition)
	}

	json.BuildMap(buffer,
		&json.KV{K: "branches", V: json.NewArray(branches)},
		&json.KV{K: "conditions", V: json.NewArray(conditions)},
		&json.KV{K: "type", V: json.NewString("MATCH")})
}

func (m *Match) GenerateICG(offset int64, code *icg.Code, s *semantic.Semantic) (int64, Error) {
	return -1, nil
}

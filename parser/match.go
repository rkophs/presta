package parser

import (
	"bytes"
	"github.com/rkophs/presta/icg"
	"github.com/rkophs/presta/json"
	"github.com/rkophs/presta/lexer"
	"github.com/rkophs/presta/semantic"
)

type Match struct {
	conditions []AstNode
	branches   []AstNode
	matchType  MatchType
}

func (p *Parser) matchExpr() (tree AstNode, yes bool, err Error) {
	readCount := 0

	/*Get '@' or '|' */
	var matchType MatchType
	readCount++
	if tok, eof := p.read(); eof {
		return p.parseError("Premature end.", readCount)
	} else if tok.Type() == lexer.MATCH_ALL {
		matchType = ALL
	} else if tok.Type() == lexer.MATCH_FIRST {
		matchType = FIRST
	} else {
		return p.parseExit(readCount)
	}

	/*Check for parenthesis*/
	readCount++
	if tok, eof := p.read(); eof {
		return p.parseError("Premature end.", readCount)
	} else if tok.Type() != lexer.PAREN_OPEN {
		return p.parseError("Missing opening parenthesis for match", readCount)
	}

	/*Get branches*/
	conditions, branches, err := p.branches()
	if err != nil {
		return p.parseError(err.Message(), readCount)
	}

	if len(conditions) != len(branches) || len(conditions) == 0 {
		return p.parseError("Invalid number of conditions and branches", readCount)
	}

	/*Check for parenthesis*/
	readCount++
	if tok, eof := p.read(); eof {
		return p.parseError("Premature end.", readCount)
	} else if tok.Type() != lexer.PAREN_CLOSE {
		return p.parseError("Missing closing parenthesis for match", readCount)
	}

	node := &Match{conditions: conditions, branches: branches, matchType: matchType}
	return p.parseValid(node)
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

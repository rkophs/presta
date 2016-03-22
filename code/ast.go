package code

import (
	"github.com/rkophs/presta/err"
	"github.com/rkophs/presta/icg"
	"github.com/rkophs/presta/json"
	"github.com/rkophs/presta/parser"
	"github.com/rkophs/presta/semantic"
)

type AstNode interface {
	json.Serializable
	Type() parser.AstNodeType
	GenerateICG(code *icg.Code, s *semantic.Semantic) err.Error
}

func parseExit(p *parser.Parser, readCount int) (tree AstNode, e err.Error) {
	p.RollBack(readCount)
	return nil, nil
}

func parseValid(p *parser.Parser, node AstNode) (tree AstNode, e err.Error) {
	return node, nil
}

func parseError(p *parser.Parser, msg string, readCount int) (tree AstNode, e err.Error) {
	p.RollBack(readCount)
	return nil, err.NewSyntaxError(msg)
}

package code

import (
	"github.com/rkophs/presta/err"
	"github.com/rkophs/presta/lexer"
	"github.com/rkophs/presta/parser"
)

func NewExpression(p *parser.Parser) (tree AstNode, e err.Error) {

	readCount := 0
	parens := false

	/*Check for parenthesis*/
	if tok, eof := p.Peek(); eof {
		return parseError(p, "Premature end.", readCount)
	} else if tok.Type() == lexer.PAREN_OPEN {
		readCount++
		p.Read()
		parens = true
	}

	if node, e := NewLetExpr(p); e != nil {
		return parseError(p, e.Message(), readCount)
	} else if node != nil {
		return validExprEnding(p, node, parens, readCount)
	}

	if node, e := parseUnaryExpression(p); e != nil {
		return parseError(p, e.Message(), readCount)
	} else if node != nil {
		return validExprEnding(p, node, parens, readCount)
	}

	if node, e := parseBinaryExpression(p); e != nil {
		return parseError(p, e.Message(), readCount)
	} else if node != nil {
		return validExprEnding(p, node, parens, readCount)
	}

	if node, e := NewData(p); e != nil {
		return parseError(p, e.Message(), readCount)
	} else if node != nil {
		return validExprEnding(p, node, parens, readCount)
	}

	return parseExit(p, readCount)
}

func validExprEnding(p *parser.Parser, node AstNode, hasOpening bool, readCount int) (tree AstNode, e err.Error) {
	if !hasOpening {
		return parseValid(p, node)
	}

	if yes, e := closeParen(p); e != nil {
		return parseError(p, e.Message(), readCount)
	} else if !yes {
		return parseError(p, "Missing closing parenthesis for expression", readCount)
	} else {
		return node, e
	}
}

func closeParen(p *parser.Parser) (yes bool, e err.Error) {
	if tok, eof := p.Peek(); eof {
		return false, err.NewSyntaxError("Premature end.")
	} else if tok.Type() != lexer.PAREN_CLOSE {
		return false, nil
	} else {
		p.Read()
		return true, nil
	}
}

func parseUnaryExpression(p *parser.Parser) (tree AstNode, e err.Error) {

	readCount := 0

	if node, e := NewMatchExpr(p); e != nil {
		return parseError(p, e.Message(), readCount)
	} else if node != nil {
		return parseValid(p, node)
	}

	if node, e := NewConcatExpr(p); e != nil {
		return parseError(p, e.Message(), readCount)
	} else if node != nil {
		return parseValid(p, node)
	}

	if node, e := NewCallExpr(p); e != nil {
		return parseError(p, e.Message(), readCount)
	} else if node != nil {
		return parseValid(p, node)
	}

	if node, e := NewNotExpr(p); e != nil {
		return parseError(p, e.Message(), readCount)
	} else if node != nil {
		return parseValid(p, node)
	}

	if node, e := parseIncrExpression(p); e != nil {
		return parseError(p, e.Message(), readCount)
	} else if node != nil {
		return parseValid(p, node)
	}

	return parseExit(p, readCount)
}

func parseIncrExpression(p *parser.Parser) (tree AstNode, e err.Error) {
	readCount := 0

	/* Get op type */
	var opType parser.BinOpType
	readCount++
	tok, eof := p.Read()
	if eof {
		return parseError(p, "Premature end.", readCount)
	} else if tok.Type() == lexer.INC {
		opType = parser.ADD_I
	} else if tok.Type() == lexer.DEC {
		opType = parser.SUB_I
	} else {
		return parseExit(p, readCount)
	}

	/*Get variable name*/
	var variable AstNode
	readCount++
	if tok, eof = p.Read(); eof {
		return parseError(p, "Premature end.", readCount)
	} else if tok.Type() != lexer.IDENTIFIER {
		return parseError(p, "Inc/Dec operator must precede an identifier", readCount)
	} else {
		name := tok.Lit()
		variable = &Variable{name: name}
	}

	one := &Data{dataType: parser.NUMBER, num: 1}
	node := &BinOp{l: variable, r: one, op: opType}
	return parseValid(p, node)
}

func parseBinaryExpression(p *parser.Parser) (tree AstNode, e err.Error) {
	readCount := 0

	if node, e := NewAssignExpr(p); e != nil {
		return parseError(p, e.Message(), readCount)
	} else if node != nil {
		return parseValid(p, node)
	}

	if node, e := NewRepeatExpr(p); e != nil {
		return parseError(p, e.Message(), readCount)
	} else if node != nil {
		return parseValid(p, node)
	}

	readCount++
	tok, eof := p.Read()
	if eof {
		return parseError(p, "Premature end.", readCount)
	}
	switch tok.Type() {
	case lexer.GT:
		return NewBinOp(p, parser.GT, readCount)
	case lexer.LT:
		return NewBinOp(p, parser.LT, readCount)
	case lexer.GTE:
		return NewBinOp(p, parser.GTE, readCount)
	case lexer.LTE:
		return NewBinOp(p, parser.LTE, readCount)
	case lexer.EQ:
		return NewBinOp(p, parser.EQ, readCount)
	case lexer.NEQ:
		return NewBinOp(p, parser.NEQ, readCount)
	case lexer.OR:
		return NewBinOp(p, parser.OR, readCount)
	case lexer.AND:
		return NewBinOp(p, parser.AND, readCount)
	case lexer.ADD:
		return NewBinOp(p, parser.ADD, readCount)
	case lexer.SUB:
		return NewBinOp(p, parser.SUB, readCount)
	case lexer.MULT:
		return NewBinOp(p, parser.MULT, readCount)
	case lexer.DIV:
		return NewBinOp(p, parser.DIV, readCount)
	case lexer.MOD:
		return NewBinOp(p, parser.MOD, readCount)
	case lexer.ADD_I:
		return NewBinOp(p, parser.ADD_I, readCount)
	case lexer.SUB_I:
		return NewBinOp(p, parser.SUB_I, readCount)
	case lexer.MULT_I:
		return NewBinOp(p, parser.MULT_I, readCount)
	case lexer.DIV_I:
		return NewBinOp(p, parser.DIV_I, readCount)
	case lexer.MOD_I:
		return NewBinOp(p, parser.MOD_I, readCount)
	default:
		return parseExit(p, readCount)
	}
}

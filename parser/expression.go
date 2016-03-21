package parser

import (
	"github.com/rkophs/presta/lexer"
)

func NewExpression(p *Parser) (tree AstNode, err Error) {

	readCount := 0
	parens := false

	/*Check for parenthesis*/
	if tok, eof := p.peek(); eof {
		return p.parseError("Premature end.", readCount)
	} else if tok.Type() == lexer.PAREN_OPEN {
		readCount++
		p.read()
		parens = true
	}

	if node, err := NewLetExpr(p); err != nil {
		return p.parseError(err.Message(), readCount)
	} else if node != nil {
		return validExprEnding(p, node, parens, readCount)
	}

	if node, err := parseUnaryExpression(p); err != nil {
		return p.parseError(err.Message(), readCount)
	} else if node != nil {
		return validExprEnding(p, node, parens, readCount)
	}

	if node, err := parseBinaryExpression(p); err != nil {
		return p.parseError(err.Message(), readCount)
	} else if node != nil {
		return validExprEnding(p, node, parens, readCount)
	}

	if node, err := NewData(p); err != nil {
		return p.parseError(err.Message(), readCount)
	} else if node != nil {
		return validExprEnding(p, node, parens, readCount)
	}

	return p.parseExit(readCount)
}

func validExprEnding(p *Parser, node AstNode, hasOpening bool, readCount int) (tree AstNode, err Error) {
	if !hasOpening {
		return p.parseValid(node)
	}

	if yes, err := closeParen(p); err != nil {
		return p.parseError(err.Message(), readCount)
	} else if !yes {
		return p.parseError("Missing closing parenthesis for expression", readCount)
	} else {
		return node, err
	}
}

func closeParen(p *Parser) (yes bool, err Error) {
	if tok, eof := p.peek(); eof {
		return false, NewSyntaxError("Premature end.")
	} else if tok.Type() != lexer.PAREN_CLOSE {
		return false, nil
	} else {
		p.read()
		return true, nil
	}
}

func parseUnaryExpression(p *Parser) (tree AstNode, err Error) {

	readCount := 0

	if node, err := NewMatchExpr(p); err != nil {
		return p.parseError(err.Message(), readCount)
	} else if node != nil {
		return p.parseValid(node)
	}

	if node, err := NewConcatExpr(p); err != nil {
		return p.parseError(err.Message(), readCount)
	} else if node != nil {
		return p.parseValid(node)
	}

	if node, err := NewCallExpr(p); err != nil {
		return p.parseError(err.Message(), readCount)
	} else if node != nil {
		return p.parseValid(node)
	}

	if node, err := NewNotExpr(p); err != nil {
		return p.parseError(err.Message(), readCount)
	} else if node != nil {
		return p.parseValid(node)
	}

	if node, err := parseIncrExpression(p); err != nil {
		return p.parseError(err.Message(), readCount)
	} else if node != nil {
		return p.parseValid(node)
	}

	return p.parseExit(readCount)
}

func parseIncrExpression(p *Parser) (tree AstNode, err Error) {
	readCount := 0

	/* Get op type */
	var opType BinOpType
	readCount++
	tok, eof := p.read()
	if eof {
		return p.parseError("Premature end.", readCount)
	} else if tok.Type() == lexer.INC {
		opType = ADD_I
	} else if tok.Type() == lexer.DEC {
		opType = SUB_I
	} else {
		return p.parseExit(readCount)
	}

	/*Get variable name*/
	var variable AstNode
	readCount++
	if tok, eof = p.read(); eof {
		return p.parseError("Premature end.", readCount)
	} else if tok.Type() != lexer.IDENTIFIER {
		return p.parseError("Inc/Dec operator must precede an identifier", readCount)
	} else {
		name := tok.Lit()
		variable = &Variable{name: name}
	}

	one := &Data{dataType: NUMBER, num: 1}
	node := &BinOp{l: variable, r: one, op: opType}
	return p.parseValid(node)
}

func parseBinaryExpression(p *Parser) (tree AstNode, err Error) {
	readCount := 0

	if node, err := NewAssignExpr(p); err != nil {
		return p.parseError(err.Message(), readCount)
	} else if node != nil {
		return p.parseValid(node)
	}

	if node, err := NewRepeatExpr(p); err != nil {
		return p.parseError(err.Message(), readCount)
	} else if node != nil {
		return p.parseValid(node)
	}

	readCount++
	tok, eof := p.read()
	if eof {
		return p.parseError("Premature end.", readCount)
	}
	switch tok.Type() {
	case lexer.GT:
		return NewBinOp(p, GT, readCount)
	case lexer.LT:
		return NewBinOp(p, LT, readCount)
	case lexer.GTE:
		return NewBinOp(p, GTE, readCount)
	case lexer.LTE:
		return NewBinOp(p, LTE, readCount)
	case lexer.EQ:
		return NewBinOp(p, EQ, readCount)
	case lexer.NEQ:
		return NewBinOp(p, NEQ, readCount)
	case lexer.OR:
		return NewBinOp(p, OR, readCount)
	case lexer.AND:
		return NewBinOp(p, AND, readCount)
	case lexer.ADD:
		return NewBinOp(p, ADD, readCount)
	case lexer.SUB:
		return NewBinOp(p, SUB, readCount)
	case lexer.MULT:
		return NewBinOp(p, MULT, readCount)
	case lexer.DIV:
		return NewBinOp(p, DIV, readCount)
	case lexer.MOD:
		return NewBinOp(p, MOD, readCount)
	case lexer.ADD_I:
		return NewBinOp(p, ADD_I, readCount)
	case lexer.SUB_I:
		return NewBinOp(p, SUB_I, readCount)
	case lexer.MULT_I:
		return NewBinOp(p, MULT_I, readCount)
	case lexer.DIV_I:
		return NewBinOp(p, DIV_I, readCount)
	case lexer.MOD_I:
		return NewBinOp(p, MOD_I, readCount)
	default:
		return p.parseExit(readCount)
	}
}
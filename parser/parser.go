package parser

import (
	"github.com/rkophs/presta/lexer"
	// "strconv"
)

type Parser struct {
	tokens []lexer.Token
	at     int
}

func NewParser(tokens []lexer.Token) *Parser {
	return &Parser{tokens: tokens, at: 0}
}

func (p *Parser) Scan() (tree AstNode, err Error) {
	if tree, yes, err := p.program(); !yes {
		if err != nil {
			return tree, err
		} else {
			return tree, NewSyntaxError("No program available.")
		}
	} else {
		return tree, nil
	}
}

func (p *Parser) expression() (tree AstNode, yes bool, err Error) {

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

	if node, yes, err := p.letExpr(); err != nil {
		return p.parseError(err.Message(), readCount)
	} else if yes {
		return p.validExprEnding(node, parens, readCount)
	}

	if node, yes, err := p.unExpr(); err != nil {
		return p.parseError(err.Message(), readCount)
	} else if yes {
		return p.validExprEnding(node, parens, readCount)
	}

	if node, yes, err := p.binExpr(); err != nil {
		return p.parseError(err.Message(), readCount)
	} else if yes {
		return p.validExprEnding(node, parens, readCount)
	}

	if node, yes, err := p.data(); err != nil {
		return p.parseError(err.Message(), readCount)
	} else if yes {
		return p.validExprEnding(node, parens, readCount)
	}

	return p.parseExit(readCount)
}

func (p *Parser) validExprEnding(node AstNode, hasOpening bool, readCount int) (tree AstNode, yes bool, err Error) {
	if !hasOpening {
		return p.parseValid(node)
	}

	if yes, err := p.closeParen(); err != nil {
		return p.parseError(err.Message(), readCount)
	} else if !yes {
		return p.parseError("Missing closing parenthesis for expression", readCount)
	} else {
		return node, true, err
	}
}

func (p *Parser) closeParen() (yes bool, err Error) {
	if tok, eof := p.peek(); eof {
		return false, NewSyntaxError("Premature end.")
	} else if tok.Type() != lexer.PAREN_CLOSE {
		return false, nil
	} else {
		p.read()
		return true, nil
	}
}

func (p *Parser) unExpr() (tree AstNode, yes bool, err Error) {

	readCount := 0

	if node, yes, err := p.matchExpr(); err != nil {
		return p.parseError(err.Message(), readCount)
	} else if yes {
		return p.parseValid(node)
	}

	if node, yes, err := p.concatExpr(); err != nil {
		return p.parseError(err.Message(), readCount)
	} else if yes {
		return p.parseValid(node)
	}

	if node, yes, err := p.callExpr(); err != nil {
		return p.parseError(err.Message(), readCount)
	} else if yes {
		return p.parseValid(node)
	}

	if node, yes, err := p.notExpr(); err != nil {
		return p.parseError(err.Message(), readCount)
	} else if yes {
		return p.parseValid(node)
	}

	if node, yes, err := p.incExpr(); err != nil {
		return p.parseError(err.Message(), readCount)
	} else if yes {
		return p.parseValid(node)
	}

	return p.parseExit(readCount)
}

func (p *Parser) branches() (c []AstNode, b []AstNode, err Error) {
	conds := []AstNode{}
	branches := []AstNode{}
	for {
		if cond, yes, err := p.expression(); err != nil {
			return conds, branches, err
		} else if yes {
			conds = append(conds, cond)
		} else {
			break
		}

		if branch, yes, err := p.expression(); err != nil {
			return conds, branches, err
		} else if !yes {
			return conds, branches, NewSyntaxError("Match expression missing branch")
		} else {
			branches = append(branches, branch)
		}
	}

	return conds, branches, nil
}

func (p *Parser) incExpr() (tree AstNode, yes bool, err Error) {
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

	one := &Data{dataType: NUMBER, num: 0}
	node := &BinOp{l: variable, r: one, op: opType}
	return p.parseValid(node)
}

func (p *Parser) binExpr() (tree AstNode, yes bool, err Error) {
	readCount := 0

	if node, yes, err := p.assignExpr(); err != nil {
		return p.parseError(err.Message(), readCount)
	} else if yes {
		return p.parseValid(node)
	}

	if node, yes, err := p.repeatExpr(); err != nil {
		return p.parseError(err.Message(), readCount)
	} else if yes {
		return p.parseValid(node)
	}

	readCount++
	tok, eof := p.read()
	if eof {
		return p.parseError("Premature end.", readCount)
	}
	switch tok.Type() {
	case lexer.GT:
		return p.parseBinaryOp(GT, readCount)
	case lexer.LT:
		return p.parseBinaryOp(LT, readCount)
	case lexer.GTE:
		return p.parseBinaryOp(GTE, readCount)
	case lexer.LTE:
		return p.parseBinaryOp(LTE, readCount)
	case lexer.EQ:
		return p.parseBinaryOp(EQ, readCount)
	case lexer.NEQ:
		return p.parseBinaryOp(NEQ, readCount)
	case lexer.OR:
		return p.parseBinaryOp(OR, readCount)
	case lexer.AND:
		return p.parseBinaryOp(AND, readCount)
	case lexer.ADD:
		return p.parseBinaryOp(ADD, readCount)
	case lexer.SUB:
		return p.parseBinaryOp(SUB, readCount)
	case lexer.MULT:
		return p.parseBinaryOp(MULT, readCount)
	case lexer.DIV:
		return p.parseBinaryOp(DIV, readCount)
	case lexer.MOD:
		return p.parseBinaryOp(MOD, readCount)
	case lexer.ADD_I:
		return p.parseBinaryOp(ADD_I, readCount)
	case lexer.SUB_I:
		return p.parseBinaryOp(SUB_I, readCount)
	case lexer.MULT_I:
		return p.parseBinaryOp(MULT_I, readCount)
	case lexer.DIV_I:
		return p.parseBinaryOp(DIV_I, readCount)
	case lexer.MOD_I:
		return p.parseBinaryOp(MOD_I, readCount)
	default:
		return p.parseExit(readCount)
	}
}

func (p *Parser) parseBinaryOp(op BinOpType, readCount int) (tree AstNode, yes bool, err Error) {
	if l, yes, err := p.expression(); err != nil {
		return p.parseError(err.Message(), readCount)
	} else if yes {
		if r, yes, err := p.expression(); err != nil {
			return p.parseError(err.Message(), readCount)
		} else if yes {
			node := &BinOp{l: l, r: r, op: op}
			return p.parseValid(node)
		} else {
			return p.parseError("Binary op needs another expression.", readCount)
		}
	} else {
		return p.parseError("Binary operation needs 2 expressions.", readCount)
	}
}

func (p *Parser) parseExit(readCount int) (tree AstNode, yes bool, err Error) {
	var node AstNode
	p.rollBack(readCount)
	return node, false, nil
}

func (p *Parser) parseValid(node AstNode) (tree AstNode, yes bool, err Error) {
	return node, true, nil
}

func (p *Parser) parseError(msg string, readCount int) (tree AstNode, yes bool, err Error) {
	var node AstNode
	p.rollBack(readCount)
	return node, false, NewSyntaxError(msg)
}

func (p *Parser) rollBack(amount int) {
	for i := 0; i < amount; i++ {
		p.unread()
	}
}

func (p *Parser) read() (tok lexer.Token, eof bool) {
	tok, eof = p.peek()
	if !eof {
		p.at++
	}
	return tok, eof
}

func (p *Parser) peek() (tok lexer.Token, eof bool) {
	if p.at >= len(p.tokens) {
		var ret lexer.Token
		return ret, true
	}
	tok = p.tokens[p.at]
	return tok, false
}

func (p *Parser) unread() {
	p.at--
}

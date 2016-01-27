package parser

import (
	"github.com/rkophs/presta/lexer"
	"strconv"
)

type Error struct {
	err bool
	msg string
}

type Parser struct {
	tokens []lexer.Token
	at     int
}

func NewParser(tokens []lexer.Token) *Parser {
	return &Parser{tokens: tokens, at: 0}
}

func (p *Parser) Scan() (tree AstNode, err *Error) {
	if tree, yes, err := p.program(); !yes {
		if err.err {
			return tree, err
		} else {
			return tree, &Error{err: true, msg: "No program available."}
		}
	} else {
		return tree, &Error{err: false}
	}
}

func (p *Parser) program() (tree AstNode, yes bool, err *Error) {
	readCount := 0

	/*Check for function declarations*/
	functions := []*Function{}
	for {
		if function, yes, err := p.function(); err.err {
			return p.parseError(err.msg, readCount)
		} else if yes {
			functions = append(functions, function)
		} else {
			break
		}
	}

	/*Check for exec*/
	expr, yes, err := p.expression()
	if err.err {
		return p.parseError(err.msg, readCount)
	} else if !yes {
		return p.parseError("Program must contain an executable expression", readCount)
	}

	program := &Program{funcs: functions, exec: expr}
	return p.parseValid(program)
}

func (p *Parser) function() (tree *Function, yes bool, err *Error) {
	readCount := 0

	var badFunc *Function

	/*Check if it starts with '~' */
	readCount++
	if tok, eof := p.read(); eof {
		_, yes, err := p.parseError("Premature end.", readCount)
		return badFunc, yes, err
	} else if tok.Type() != lexer.FUNC {
		_, yes, err := p.parseExit(readCount)
		return badFunc, yes, err
	}

	/*Check for identifier*/
	readCount++
	tok, eof := p.read()
	if eof {
		_, yes, err := p.parseError("Premature end.", readCount)
		return badFunc, yes, err
	} else if tok.Type() != lexer.IDENTIFIER {
		_, yes, err := p.parseError("Function name must follow ~", readCount)
		return badFunc, yes, err
	}
	funcName := tok.Lit()

	/* Check for parenthesis */
	readCount++
	if tok, eof := p.read(); eof {
		_, yes, err := p.parseError("Premature end.", readCount)
		return badFunc, yes, err
	} else if tok.Type() != lexer.PAREN_OPEN {
		_, yes, err := p.parseError("Parenthesis must follow function name", readCount)
		return badFunc, yes, err
	}

	/* Check for param names */
	params := []string{}
	for {
		readCount++
		if tok, eof := p.read(); eof {
			_, yes, err := p.parseError("Premature end.", readCount)
			return badFunc, yes, err
		} else if tok.Type() == lexer.IDENTIFIER {
			params = append(params, tok.Lit())
		} else if tok.Type() == lexer.PAREN_CLOSE {
			break
		} else {
			_, yes, err := p.parseError("Looking for parameter identifiers for function", readCount)
			return badFunc, yes, err
		}
	}

	/*Check for parenthesis*/
	readCount++
	if tok, eof := p.read(); eof {
		_, yes, err := p.parseError("Premature end.", readCount)
		return badFunc, yes, err
	} else if tok.Type() != lexer.PAREN_OPEN {
		_, yes, err := p.parseError("'(' must prefix function body", readCount)
		return badFunc, yes, err
	}

	/* Check for expression */
	expr, yes, err := p.expression()
	if err.err {
		_, yes, err := p.parseError(err.msg, readCount)
		return badFunc, yes, err
	} else if !yes {
		_, yes, err := p.parseError("Function body must be an executable expression", readCount)
		return badFunc, yes, err
	}

	/*Check for parenthesis*/
	readCount++
	if tok, eof := p.read(); eof {
		_, yes, err := p.parseError("Premature end.", readCount)
		return badFunc, yes, err
	} else if tok.Type() != lexer.PAREN_CLOSE {
		_, yes, err := p.parseError("Parenthesis must postfix function body", readCount)
		return badFunc, yes, err
	}

	node := &Function{name: funcName, params: params, exec: expr}
	return node, yes, &Error{err: false}
}

func (p *Parser) expression() (tree AstNode, yes bool, err *Error) {

	readCount := 0
	parens := false

	/*Check for parenthesis*/
	if tok, eof := p.peek(); eof {
		return p.parseError("Premature end.", readCount)
	} else if tok.Type() != lexer.PAREN_OPEN {
		readCount++
		p.read()
		parens = true
	}

	if node, yes, err := p.letExpr(); err.err {
		return p.parseError(err.msg, readCount)
	} else if yes {
		return p.validExprEnding(node, parens, readCount)
	}

	if node, yes, err := p.unExpr(); err.err {
		return p.parseError(err.msg, readCount)
	} else if yes {
		return p.validExprEnding(node, parens, readCount)
	}

	if node, yes, err := p.binExpr(); err.err {
		return p.parseError(err.msg, readCount)
	} else if yes {
		return p.validExprEnding(node, parens, readCount)
	}

	if node, yes, err := p.data(); err.err {
		return p.parseError(err.msg, readCount)
	} else if yes {
		return p.validExprEnding(node, parens, readCount)
	}

	return p.parseExit(readCount)
}

func (p *Parser) validExprEnding(node AstNode, hasOpening bool, readCount int) (tree AstNode, yes bool, err *Error) {
	if !hasOpening {
		return p.parseValid(node)
	}

	if yes, err := p.closeParen(); err.err {
		return p.parseError(err.msg, readCount)
	} else if !yes {
		return p.parseError("Missing closing parenthesis.", readCount)
	} else {
		return node, true, err
	}
}

func (p *Parser) closeParen() (yes bool, err *Error) {
	if tok, eof := p.peek(); eof {
		return false, &Error{err: true, msg: "Premature end."}
	} else if tok.Type() != lexer.PAREN_CLOSE {
		return false, &Error{err: false}
	} else {
		p.read()
		return true, &Error{err: false}
	}
}

func (p *Parser) letExpr() (tree AstNode, yes bool, err *Error) {
	readCount := 0

	/*Check if it starts with ':' */
	readCount++
	if tok, eof := p.read(); eof {
		return p.parseError("Premature end.", readCount)
	} else if tok.Type() != lexer.ASSIGN {
		return p.parseExit(readCount)
	}

	/*Check for parenthesis*/
	readCount++
	if tok, eof := p.read(); eof {
		return p.parseError("Premature end.", readCount)
	} else if tok.Type() != lexer.PAREN_OPEN {
		return p.parseError("Missing opening parenthesis", readCount)
	}

	/* Check for param names */
	params := []string{}
	for {
		readCount++
		if tok, eof := p.read(); eof {
			return p.parseError("Premature end.", readCount)
		} else if tok.Type() == lexer.IDENTIFIER {
			params = append(params, tok.Lit())
		} else if tok.Type() == lexer.PAREN_CLOSE {
			break
		} else {
			return p.parseError("Looking for parameter identifiers for function", readCount)
		}
	}

	/*Check for parenthesis*/
	readCount++
	if tok, eof := p.read(); eof {
		return p.parseError("Premature end.", readCount)
	} else if tok.Type() != lexer.PAREN_CLOSE {
		return p.parseError("Missing closing parenthesis", readCount)
	}

	/*Check for parenthesis*/
	readCount++
	if tok, eof := p.read(); eof {
		return p.parseError("Premature end.", readCount)
	} else if tok.Type() != lexer.PAREN_OPEN {
		return p.parseError("Missing opening parenthesis", readCount)
	}

	/* Check for assignments */
	values := []AstNode{}
	for {
		if node, yes, err := p.expression(); err.err {
			return p.parseError(err.msg, readCount)
		} else if yes {
			values = append(values, node)
		} else {
			break
		}
	}

	/*Check for parenthesis*/
	readCount++
	if tok, eof := p.read(); eof {
		return p.parseError("Premature end.", readCount)
	} else if tok.Type() != lexer.PAREN_CLOSE {
		return p.parseError("Missing closing parenthesis", readCount)
	}

	body, yes, err := p.expression()
	if err.err {
		return p.parseError(err.msg, readCount)
	} else if !yes {
		return p.parseError("Missing let statement body", readCount)
	}

	node := &Let{params: params, values: values, exec: body}
	return p.parseValid(node)
}

func (p *Parser) unExpr() (tree AstNode, yes bool, err *Error) {

	readCount := 0

	if node, yes, err := p.matchExpr(); err.err {
		return p.parseError(err.msg, readCount)
	} else if yes {
		return p.parseValid(node)
	}

	if node, yes, err := p.concatExpr(); err.err {
		return p.parseError(err.msg, readCount)
	} else if yes {
		return p.parseValid(node)
	}

	if node, yes, err := p.callExpr(); err.err {
		return p.parseError(err.msg, readCount)
	} else if yes {
		return p.parseValid(node)
	}

	if node, yes, err := p.notExpr(); err.err {
		return p.parseError(err.msg, readCount)
	} else if yes {
		return p.parseValid(node)
	}

	if node, yes, err := p.incExpr(); err.err {
		return p.parseError(err.msg, readCount)
	} else if yes {
		return p.parseValid(node)
	}

	return p.parseExit(readCount)
}

func (p *Parser) callExpr() (tree AstNode, yes bool, err *Error) {

	readCount := 0

	/*Get variable name*/
	var name string
	readCount++
	if tok, eof := p.read(); eof {
		return p.parseError("Premature end.", readCount)
	} else if tok.Type() != lexer.IDENTIFIER {
		return p.parseExit(readCount)
	} else {
		name = tok.Lit()
	}

	/*Check for bracket*/
	readCount++
	if tok, eof := p.read(); eof {
		return p.parseError("Premature end.", readCount)
	} else if tok.Type() != lexer.CURLY_OPEN {
		return p.parseExit(readCount) //Not caller, but data identifier
	}

	/* Check for arguments */
	args := []AstNode{}
	for {
		if expr, yes, err := p.expression(); err.err {
			return p.parseError(err.msg, readCount)
		} else if yes {
			args = append(args, expr)
		} else {
			break
		}
	}

	/*Check for bracket*/
	readCount++
	if tok, eof := p.read(); eof {
		return p.parseError("Premature end.", readCount)
	} else if tok.Type() != lexer.CURLY_CLOSE {
		return p.parseError("Missing closing bracket.", readCount)
	}

	node := &Call{name: name, params: args}
	return p.parseValid(node)
}

func (p *Parser) notExpr() (tree AstNode, yes bool, err *Error) {
	readCount := 0
	/*Check for ! */
	readCount++
	if tok, eof := p.read(); eof {
		return p.parseError("Premature end.", readCount)
	} else if tok.Type() != lexer.NOT {
		return p.parseExit(readCount) //Not caller, but data identifier
	}

	if expr, yes, err := p.expression(); err.err {
		return p.parseError(err.msg, readCount)
	} else if yes {
		node := &Not{exec: expr}
		return p.parseValid(node)
	} else {
		return p.parseError("Not operator must precede expression", readCount)
	}

}

func (p *Parser) concatExpr() (tree AstNode, yes bool, err *Error) {
	readCount := 0

	/*Get '.'*/
	readCount++
	if tok, eof := p.read(); eof {
		return p.parseError("Premature end.", readCount)
	} else if tok.Type() != lexer.CONCAT {
		return p.parseExit(readCount)
	}

	/*Check for parenthesis*/
	readCount++
	if tok, eof := p.read(); eof {
		return p.parseError("Premature end.", readCount)
	} else if tok.Type() != lexer.PAREN_OPEN {
		return p.parseError("Missing opening parenthesis", readCount)
	}

	/*Get List*/
	exprs := []AstNode{}
	for {
		if expr, yes, err := p.expression(); err.err {
			return p.parseError(err.msg, readCount)
		} else if yes {
			exprs = append(exprs, expr)
		} else {
			break
		}
	}

	/*Check for parenthesis*/
	readCount++
	if tok, eof := p.read(); eof {
		return p.parseError("Premature end.", readCount)
	} else if tok.Type() != lexer.PAREN_CLOSE {
		return p.parseError("Missing opening parenthesis", readCount)
	}

	node := &Concat{components: exprs}
	return p.parseValid(node)
}

func (p *Parser) matchExpr() (tree AstNode, yes bool, err *Error) {
	readCount := 0

	/*Get '@' or '|' */
	var matchType MatchType
	readCount++
	if tok, eof := p.read(); eof {
		return p.parseError("Premature end.", readCount)
	} else if tok.Type() != lexer.MATCH_ALL {
		matchType = ALL
	} else if tok.Type() != lexer.MATCH_FIRST {
		matchType = FIRST
	} else {
		return p.parseExit(readCount)
	}

	/*Check for parenthesis*/
	readCount++
	if tok, eof := p.read(); eof {
		return p.parseError("Premature end.", readCount)
	} else if tok.Type() != lexer.PAREN_OPEN {
		return p.parseError("Missing opening parenthesis", readCount)
	}

	/*Get branches*/
	conditions, branches, err := p.branches()
	if err.err {
		return p.parseError(err.msg, readCount)
	}

	/*Check for parenthesis*/
	readCount++
	if tok, eof := p.read(); eof {
		return p.parseError("Premature end.", readCount)
	} else if tok.Type() != lexer.PAREN_CLOSE {
		return p.parseError("Missing opening parenthesis", readCount)
	}

	node := &Match{conditions: conditions, branches: branches, matchType: matchType}
	return p.parseValid(node)
}

func (p *Parser) branches() (c []AstNode, b []AstNode, err *Error) {
	conds := []AstNode{}
	branches := []AstNode{}
	for {
		if cond, yes, err := p.expression(); err.err {
			return conds, branches, err
		} else if yes {
			conds = append(conds, cond)
		} else {
			break
		}

		if branch, yes, err := p.expression(); err.err {
			return conds, branches, err
		} else if !yes {
			return conds, branches, &Error{err: true, msg: "Match expression missing branch"}
		} else {
			branches = append(branches, branch)
		}
	}

	return conds, branches, &Error{err: false}
}

func (p *Parser) incExpr() (tree AstNode, yes bool, err *Error) {

	readCount := 0

	/* Get op type */
	var opType BinOpType
	readCount++
	tok, eof := p.read()
	if eof {
		return p.parseError("Premature end.", readCount)
	} else if tok.Type() == lexer.INC {
		opType = ADD_I
	} else if tok.Type() != lexer.DEC {
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
	node := &BinOp{a: variable, b: one, op: opType}
	return p.parseValid(node)
}

func (p *Parser) binExpr() (tree AstNode, yes bool, err *Error) {
	readCount := 0

	if node, yes, err := p.assignExpr(); err.err {
		return p.parseError(err.msg, readCount)
	} else if yes {
		return p.parseValid(node)
	}

	if node, yes, err := p.repeatExpr(); err.err {
		return p.parseError(err.msg, readCount)
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

func (p *Parser) repeatExpr() (tree AstNode, yes bool, err *Error) {
	readCount := 0

	/*Check for ^ */
	readCount++
	if tok, eof := p.read(); eof {
		return p.parseError("Premature end.", readCount)
	} else if tok.Type() != lexer.REPEAT {
		return p.parseExit(readCount) //Not caller, but data identifier
	}

	/*Get expression*/
	var condition AstNode
	if expr, yes, err := p.expression(); err.err {
		return p.parseError(err.msg, readCount)
	} else if yes {
		condition = expr
	} else {
		return p.parseError("Repeat op must have condition", readCount)
	}

	/*Get expression*/
	if expr, yes, err := p.expression(); err.err {
		return p.parseError(err.msg, readCount)
	} else if yes {
		node := &Repeat{condition: condition, exec: expr}
		return p.parseValid(node)
	} else {
		return p.parseError("Repeat op must have body", readCount)
	}
}

func (p *Parser) assignExpr() (tree AstNode, yes bool, err *Error) {
	readCount := 0

	/*Check for : */
	readCount++
	if tok, eof := p.read(); eof {
		return p.parseError("Premature end.", readCount)
	} else if tok.Type() != lexer.ASSIGN {
		return p.parseExit(readCount) //Not caller, but data identifier
	}

	/*Get variable name*/
	var name string
	readCount++
	if tok, eof := p.read(); eof {
		return p.parseError("Premature end.", readCount)
	} else if tok.Type() != lexer.IDENTIFIER {
		return p.parseError("Assignment operator must precede an identifier.", readCount)
	} else {
		name = tok.Lit()
	}

	/*Get expression*/
	if expr, yes, err := p.expression(); err.err {
		return p.parseError(err.msg, readCount)
	} else if yes {
		node := &Assign{name: name, value: expr}
		return p.parseValid(node)
	} else {
		return p.parseError("Assignment operator must have valid assignment expression.", readCount)
	}
}

func (p *Parser) parseBinaryOp(op BinOpType, readCount int) (tree AstNode, yes bool, err *Error) {
	if a, yes, err := p.expression(); err.err {
		return p.parseError(err.msg, readCount)
	} else if yes {
		if b, yes, err := p.expression(); err.err {
			return p.parseError(err.msg, readCount)
		} else if yes {
			node := &BinOp{a: a, b: b, op: op}
			return p.parseValid(node)
		} else {
			return p.parseError("Binary operation must be followed by 2 expressions.", readCount)
		}
	} else {
		return p.parseError("Binary operation must be followed by 2 expressions.", readCount)
	}
}

func (p *Parser) data() (tree AstNode, yes bool, err *Error) {
	readCount := 1
	if tok, err := p.read(); err {
		return p.parseError("Premature end.", readCount)
	} else if tok.Type() == lexer.STRING {
		node := &Data{str: tok.Lit(), dataType: STRING}
		return p.parseValid(node)
	} else if tok.Type() == lexer.NUMBER {
		if num, err := strconv.ParseFloat(tok.Lit(), 64); err != nil {
			return p.parseError("Error parsing numeric.", readCount)
		} else {
			node := &Data{num: num, dataType: NUMBER}
			return p.parseValid(node)
		}
	} else if tok.Type() == lexer.IDENTIFIER {
		if next, err := p.peek(); err {
			node := &Variable{name: tok.Lit()}
			return p.parseValid(node)
		} else if next.Type() == lexer.CURLY_OPEN { //Not identifer - but caller
			p.parseExit(readCount)
		}
	}

	return p.parseExit(readCount)
}

func (p *Parser) parseExit(readCount int) (tree AstNode, yes bool, err *Error) {
	var node AstNode
	p.rollBack(readCount)
	return node, false, &Error{err: false}
}

func (p *Parser) parseValid(node AstNode) (tree AstNode, yes bool, err *Error) {
	return node, true, &Error{err: false}
}

func (p *Parser) parseError(msg string, readCount int) (tree AstNode, yes bool, err *Error) {
	var node AstNode
	p.rollBack(readCount)
	return node, false, &Error{err: true, msg: msg}
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

package parser

import (
	"github.com/rkophs/presta/lexer"
	"strconv"
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

func (p *Parser) program() (tree *Program, yes bool, err Error) {
	readCount := 0
	var invalid Program

	/*Check for function declarations*/
	functions := []*Function{}
	for {
		if function, yes, err := p.function(); err != nil {
			_, yes, err := p.parseError(err.Message(), readCount)
			return &invalid, yes, err
		} else if yes {
			functions = append(functions, function)
		} else {
			break
		}
	}

	/*Check for exec*/
	expr, yes, err := p.expression()
	if err != nil {
		_, yes, err := p.parseError(err.Message(), readCount)
		return &invalid, yes, err
	} else if !yes {
		_, yes, err := p.parseError("Program must contain an executable expression", readCount)
		return &invalid, yes, err
	}

	program := &Program{funcs: functions, exec: expr}
	_, yes, err = p.parseValid(program)
	return program, yes, err
}

func (p *Parser) function() (tree *Function, yes bool, err Error) {
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
	if err != nil {
		_, yes, err := p.parseError(err.Message(), readCount)
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
	return node, yes, nil
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

func (p *Parser) letExpr() (tree AstNode, yes bool, err Error) {
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
		// probably an assignment at this point
		return p.parseExit(readCount)
	}

	/* Check for param names and closing parenthesis*/
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
	} else if tok.Type() != lexer.PAREN_OPEN {
		return p.parseError("Missing opening parenthesis for let assignments", readCount)
	}

	/* Check for assignments */
	values := []AstNode{}
	for {
		if node, yes, err := p.expression(); err != nil {
			return p.parseError(err.Message(), readCount)
		} else if yes {
			values = append(values, node)
		} else {
			break
		}
	}

	if len(values) != len(params) {
		return p.parseError("Number of assignments must equal number of variables in let.", readCount)
	}

	/*Check for parenthesis*/
	readCount++
	if tok, eof := p.read(); eof {
		return p.parseError("Premature end.", readCount)
	} else if tok.Type() != lexer.PAREN_CLOSE {
		return p.parseError("Missing closing parenthesis for let assignments", readCount)
	}

	body, yes, err := p.expression()
	if err != nil {
		return p.parseError(err.Message(), readCount)
	} else if !yes {
		return p.parseError("Missing let statement body", readCount)
	}

	node := &Let{params: params, values: values, exec: body}
	return p.parseValid(node)
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

func (p *Parser) callExpr() (tree AstNode, yes bool, err Error) {

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
		if expr, yes, err := p.expression(); err != nil {
			return p.parseError(err.Message(), readCount)
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

func (p *Parser) notExpr() (tree AstNode, yes bool, err Error) {
	readCount := 0
	/*Check for ! */
	readCount++
	if tok, eof := p.read(); eof {
		return p.parseError("Premature end.", readCount)
	} else if tok.Type() != lexer.NOT {
		return p.parseExit(readCount) //Not caller, but data identifier
	}

	if expr, yes, err := p.expression(); err != nil {
		return p.parseError(err.Message(), readCount)
	} else if yes {
		node := &Not{exec: expr}
		return p.parseValid(node)
	} else {
		return p.parseError("Not operator must precede expression", readCount)
	}

}

func (p *Parser) concatExpr() (tree AstNode, yes bool, err Error) {
	readCount := 0

	/* Get '.' */
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
		return p.parseError("Missing opening parenthesis for concat", readCount)
	}

	/*Get List*/
	exprs := []AstNode{}
	for {
		if expr, yes, err := p.expression(); err != nil {
			return p.parseError(err.Message(), readCount)
		} else if yes {
			exprs = append(exprs, expr)
		} else {
			break
		}

		/*Exit on closing parenthesis*/
		if tok, eof := p.peek(); eof {
			return p.parseError("Premature end.", readCount)
		} else if tok.Type() == lexer.PAREN_CLOSE {
			break
		}
	}

	/*Check for parenthesis*/
	readCount++
	if tok, eof := p.read(); eof {
		return p.parseError("Premature end.", readCount)
	} else if tok.Type() != lexer.PAREN_CLOSE {
		return p.parseError("Missing closing parenthesis for concat", readCount)
	}

	node := &Concat{components: exprs}
	return p.parseValid(node)
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

func (p *Parser) repeatExpr() (tree AstNode, yes bool, err Error) {
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
	if expr, yes, err := p.expression(); err != nil {
		return p.parseError(err.Message(), readCount)
	} else if yes {
		condition = expr
	} else {
		return p.parseError("Repeat op must have condition", readCount)
	}

	/*Get expression*/
	if expr, yes, err := p.expression(); err != nil {
		return p.parseError(err.Message(), readCount)
	} else if yes {
		node := &Repeat{condition: condition, exec: expr}
		return p.parseValid(node)
	} else {
		return p.parseError("Repeat op must have body", readCount)
	}
}

func (p *Parser) assignExpr() (tree AstNode, yes bool, err Error) {
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
	if expr, yes, err := p.expression(); err != nil {
		return p.parseError(err.Message(), readCount)
	} else if yes {
		node := &Assign{name: name, value: expr}
		return p.parseValid(node)
	} else {
		return p.parseError("Assignment operator must have valid assignment expression.", readCount)
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

func (p *Parser) data() (tree AstNode, yes bool, err Error) {
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
			return p.parseError("Premature end.", readCount)
		} else if next.Type() == lexer.CURLY_OPEN { //Not identifer - but caller
			p.parseExit(readCount)
		} else {
			node := &Variable{name: tok.Lit()}
			return p.parseValid(node)
		}
	}

	return p.parseExit(readCount)
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

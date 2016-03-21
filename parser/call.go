package parser

import (
	"bytes"
	"github.com/rkophs/presta/icg"
	"github.com/rkophs/presta/ir"
	"github.com/rkophs/presta/json"
	"github.com/rkophs/presta/lexer"
	"github.com/rkophs/presta/semantic"
)

type Call struct {
	name   string
	params []AstNode
}

func (c *Call) Type() AstNodeType {
	return CALL
}

func (c *Call) Serialize(buffer *bytes.Buffer) {

	params := []json.Serializable{}
	for _, param := range c.params {
		params = append(params, param)
	}

	json.BuildMap(buffer,
		&json.KV{K: "params", V: json.NewArray(params)},
		&json.KV{K: "name", V: json.NewString(c.name)},
		&json.KV{K: "type", V: json.NewString("CALL")})
}

func NewCallExpr(p *Parser) (tree AstNode, err Error) {

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
		if expr, err := p.expression(); err != nil {
			return p.parseError(err.Message(), readCount)
		} else if expr != nil {
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

func (c *Call) GenerateICG(offset int64, code *icg.Code, s *semantic.Semantic) (int64, Error) {

	if !s.FunctionExists(c.name) || s.FunctionArity(c.name) != len(c.params) {
		return -1, NewSymanticError("Function not found")
	}

	//Generate params
	offsets := make([]int, len(c.params))
	for i, p := range c.params {
		if _, err := p.GenerateICG(-1, code, s); err != nil {
			return -1, err
		} else {
			offsets[i] = code.GetFrameOffset()
			code.Append(ir.NewPush(code.Ax))
			code.IncrFrameOffset(1)
		}
	}

	//Push params onto the stack as fn args
	for i, _ := range c.params {
		code.Append(ir.NewPush(ir.NewStackAccess(offsets[i])))
	}
	code.IncrFrameOffset(len(c.params))

	//Call function (which loads AX when finished)
	gotoLoc := code.GetFunctionOffset(c.name)
	code.Append(ir.NewCall(gotoLoc))

	return -1, nil
}

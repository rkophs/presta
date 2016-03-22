package code

import (
	"bytes"
	"encoding/hex"
	"github.com/rkophs/presta/err"
	"github.com/rkophs/presta/icg"
	"github.com/rkophs/presta/ir"
	"github.com/rkophs/presta/json"
	"github.com/rkophs/presta/lexer"
	"github.com/rkophs/presta/parser"
	"github.com/rkophs/presta/semantic"
	"strconv"
)

type Data struct {
	str      string
	num      float64
	dataType DataType
}

func NewData(p *parser.Parser) (tree AstNode, e err.Error) {
	readCount := 1
	if tok, e := p.Read(); e {
		return parseError(p, "Premature end.", readCount)
	} else if tok.Type() == lexer.STRING {
		node := &Data{str: tok.Lit(), dataType: STRING}
		return parseValid(p, node)
	} else if tok.Type() == lexer.NUMBER {
		if num, e := strconv.ParseFloat(tok.Lit(), 64); e != nil {
			return parseError(p, "Error parsing numeric.", readCount)
		} else {
			node := &Data{num: num, dataType: NUMBER}
			return parseValid(p, node)
		}
	} else if tok.Type() == lexer.IDENTIFIER {
		if next, e := p.Peek(); e {
			return parseError(p, "Premature end.", readCount)
		} else if next.Type() == lexer.CURLY_OPEN { //Not identifer - but caller
			parseExit(p, readCount)
		} else {
			node := &Variable{name: tok.Lit()}
			return parseValid(p, node)
		}
	}

	return parseExit(p, readCount)
}

func (d *Data) Type() AstNodeType {
	return DATA
}

func (d *Data) Serialize(buffer *bytes.Buffer) {
	if d.dataType == NUMBER {
		json.BuildMap(buffer,
			&json.KV{K: "dataType", V: json.NewString(d.dataType.String())},
			&json.KV{K: "value", V: json.NewNumber(d.num)},
			&json.KV{K: "type", V: json.NewString("DATA")})
	} else {
		str := []byte(d.str)
		hexStr := hex.EncodeToString(str)
		json.BuildMap(buffer,
			&json.KV{K: "dataType", V: json.NewString("string")},
			&json.KV{K: "value", V: json.NewString(hexStr)},
			&json.KV{K: "type", V: json.NewString("DATA")})
	}
}

func (d *Data) GenerateICG(code *icg.Code, s *semantic.Semantic) err.Error {

	var entry ir.StackEntry
	switch d.dataType {
	case STRING:
		entry = ir.NewString(d.str)
		break
	case NUMBER:
		entry = ir.NewNumber(d.num)
		break
	}

	code.Append(ir.NewMov(code.Ax, ir.NewConstantAccess(entry)))
	return nil
}

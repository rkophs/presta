/*
 * Copyright (c) 2016 Ryan Kophs
 *
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to
 * deal in the Software without restriction, including without limitation the
 * rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
 * sell copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 *
 **/

package code

import (
	"bytes"
	"encoding/hex"
	"github.com/rkophs/presta/err"
	"github.com/rkophs/presta/icg"
	"github.com/rkophs/presta/ir"
	"github.com/rkophs/presta/json"
	"github.com/rkophs/presta/parser"
	"github.com/rkophs/presta/system"
	"strconv"
)

type Data struct {
	str      string
	num      float64
	dataType DataType
}

func NewData(p *parser.TokenScanner) (tree AstNode, e err.Error) {
	readCount := 1
	if tok, e := p.Read(); e {
		return parseError(p, "Premature end.", readCount)
	} else if tok.Type() == parser.STRING {
		node := &Data{str: tok.Lit(), dataType: STRING}
		return parseValid(p, node)
	} else if tok.Type() == parser.NUMBER {
		if num, e := strconv.ParseFloat(tok.Lit(), 64); e != nil {
			return parseError(p, "Error parsing numeric.", readCount)
		} else {
			node := &Data{num: num, dataType: NUMBER}
			return parseValid(p, node)
		}
	} else if tok.Type() == parser.IDENTIFIER {
		if next, e := p.Peek(); e {
			return parseError(p, "Premature end.", readCount)
		} else if next.Type() == parser.CURLY_OPEN { //Not identifer - but caller
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

func (d *Data) GenerateICG(code *icg.Code, s *parser.Semantic) err.Error {

	var entry system.StackEntry
	switch d.dataType {
	case STRING:
		entry = system.NewString(d.str)
		break
	case NUMBER:
		entry = system.NewNumber(d.num)
		break
	}

	code.Append(ir.NewMov(code.Ax, ir.NewConstantAccess(entry)))
	return nil
}

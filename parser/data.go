package parser

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/rkophs/presta/icg"
	"github.com/rkophs/presta/ir"
	"github.com/rkophs/presta/json"
	"github.com/rkophs/presta/semantic"
)

type Data struct {
	str      string
	num      float64
	dataType DataType
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

func (d *Data) GenerateICG(doNothing int64, code *icg.Code, s *semantic.Semantic) (int64, Error) {

	fmt.Println("ICG for Data")

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
	return -1, nil
}

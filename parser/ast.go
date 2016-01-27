package parser

type AstNode interface {
	Type() AstNodeType
}

type Function struct {
	name   string
	params []string
	exec   AstNode
}

type AstNodeType int64
type MatchType int64
type BinOpType int64
type DataType int64

const (
	PROG AstNodeType = iota
	CONCAT
	LET
	REPEAT
	ASSIGN
	MATCH
	VAR
	CALL
	DATA
	NOT
	BIN_OP
)

const (
	ALL MatchType = iota
	FIRST
)

const (
	ADD BinOpType = iota
	SUB
	MULT
	DIV
	MOD

	ADD_I
	SUB_I
	MULT_I
	DIV_I
	MOD_I

	LT
	LTE
	GT
	GTE
	EQ
	NEQ
	AND
	OR
)

const (
	STRING DataType = iota
	NUMBER
)

type Program struct {
	funcs []*Function
	exec  AstNode
}

type Concat struct {
	components []AstNode
}

type Let struct {
	params []string
	values []AstNode
	exec   AstNode
}

type Repeat struct {
	condition AstNode
	exec      AstNode
}

type Assign struct {
	name  string
	value AstNode
}

type Match struct {
	conditions []AstNode
	branches   []AstNode
	matchType  MatchType
}

type Variable struct {
	name string
}

type Call struct {
	name   string
	params []AstNode
}

type Data struct {
	str      string
	num      float64
	dataType DataType
}

type Not struct {
	exec AstNode
}

type BinOp struct {
	a  AstNode
	b  AstNode
	op BinOpType
}

func (p *Program) Type() AstNodeType {
	return PROG
}

func (c *Concat) Type() AstNodeType {
	return CONCAT
}

func (l *Let) Type() AstNodeType {
	return LET
}

func (r *Repeat) Type() AstNodeType {
	return REPEAT
}

func (a *Assign) Type() AstNodeType {
	return ASSIGN
}

func (m *Match) Type() AstNodeType {
	return MATCH
}

func (v *Variable) Type() AstNodeType {
	return VAR
}

func (c *Call) Type() AstNodeType {
	return CALL
}

func (n *Data) Type() AstNodeType {
	return DATA
}

func (n *Not) Type() AstNodeType {
	return NOT
}

func (b *BinOp) Type() AstNodeType {
	return BIN_OP
}

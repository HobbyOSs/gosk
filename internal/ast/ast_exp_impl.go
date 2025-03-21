package ast

import "strings"

// TODO: go generateで作成できないか
type DataType string

const (
	Byte  DataType = "BYTE"
	Word  DataType = "WORD"
	Dword DataType = "DWORD"
	None  DataType = ""
)

var stringToDataType = map[string]DataType{
	"BYTE":  Byte,
	"WORD":  Word,
	"DWORD": Dword,
	"":      None,
}

func NewDataType(s string) (DataType, bool) {
	c, ok := stringToDataType[s]
	return c, ok
}

// TODO: go generateで作成できないか
type JumpType string

const (
	Short JumpType = "SHORT"
	Near  JumpType = "NEAR"
	Far   JumpType = "FAR"
	Empty JumpType = ""
)

var stringToJumpType = map[string]JumpType{
	"SHORT": Short,
	"NEAR":  Near,
	"FAR":   Far,
	"":      Empty,
}

func NewJumpType(s string) (JumpType, bool) {
	c, ok := stringToJumpType[s]
	return c, ok
}

//go:generate newc
type SegmentExp struct {
	BaseExp
	DataType DataType
	Left     *AddExp
	Right    *AddExp // nullable
}

func (s SegmentExp) expressionNode() {}
func (s SegmentExp) TokenLiteral() string {
	leftStr := ExpToString(s.Left)
	rightStr := ""
	if s.Right != nil {
		rightStr = ExpToString(s.Right)
	}
	dataTypeStr := ""
	if s.DataType != None {
		dataTypeStr = string(s.DataType) + " "
	}
	if rightStr == "" {
		return dataTypeStr + leftStr
	} else {
		return dataTypeStr + leftStr + ":" + rightStr
	}
}

//go:generate newc
type MemoryAddrExp struct {
	BaseExp
	DataType DataType
	JumpType JumpType
	Left     *AddExp
	Right    *AddExp // nullable
}

func (m MemoryAddrExp) expressionNode() {}
func (m MemoryAddrExp) TokenLiteral() string {
	var str = ""
	if m.DataType != None {
		str += string(m.DataType)
		str += " "
	}
	str += "[ "
	str += m.Left.TokenLiteral()
	if m.Right != nil {
		str += " : "
		str += m.Right.TokenLiteral()
	}
	str += " ]"
	return str
}

//go:generate newc
type AddExp struct {
	BaseExp
	HeadExp   *MultExp
	Operators []string
	TailExps  []*MultExp
}

func (a AddExp) expressionNode() {}
func (a AddExp) TokenLiteral() string {
	// 頭の項
	head := ExpToString(a.HeadExp)
	// 後続の Operators & TailExps をまとめて文字列に
	var buf strings.Builder
	buf.WriteString(head)
	for i, op := range a.Operators {
		buf.WriteByte(' ')
		buf.WriteString(op)
		buf.WriteByte(' ')
		tailStr := ExpToString(a.TailExps[i])
		buf.WriteString(tailStr)
	}
	return buf.String()
}

//go:generate newc
type MultExp struct {
	BaseExp
	HeadExp   *ImmExp
	Operators []string
	TailExps  []*ImmExp
}

func (m MultExp) expressionNode() {}
func (m MultExp) TokenLiteral() string {
	head := ExpToString(m.HeadExp)
	// 例： "4 * ESI" など
	var buf strings.Builder
	buf.WriteString(head)
	for i, op := range m.Operators {
		buf.WriteByte(' ')
		buf.WriteString(op)
		buf.WriteByte(' ')
		tailStr := ExpToString(m.TailExps[i])
		buf.WriteString(tailStr)
	}
	return buf.String()
}

//go:generate newc
type ImmExp struct {
	BaseExp
	Factor Factor
}

func (imm ImmExp) expressionNode() {}
func (imm ImmExp) TokenLiteral() string {
	return imm.Factor.TokenLiteral()
}

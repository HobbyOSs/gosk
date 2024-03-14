package ast

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

//go:generate newc
type SegmentExp struct {
	BaseExp
	DataType DataType
	Left     *AddExp
	Right    *AddExp // nullable
}

func (s SegmentExp) String() string {
	str := s.Left.String()
	if s.Right != nil {
		str += " : "
		str += s.Right.String()
	}
	return str
}

//go:generate newc
type MemoryAddrExp struct {
	BaseExp
	DataType DataType
	Left     *AddExp
	Right    *AddExp // nullable
}

func (m MemoryAddrExp) String() string {
	str := "[ "
	str += m.Left.String()
	if m.Right != nil {
		str += " : "
		str += m.Right.String()
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

func (a AddExp) String() string {
	return a.HeadExp.String()
}

//go:generate newc
type MultExp struct {
	BaseExp
	HeadExp   *ImmExp
	Operators []string
	TailExps  []*ImmExp
}

func (m MultExp) String() string {
	return m.HeadExp.String()
}

//go:generate newc
type ImmExp struct {
	BaseExp
	Factor Factor
}

func (imm ImmExp) String() string {
	return imm.Factor.String()
}

package pass1

import (
	"github.com/HobbyOSs/gosk/ast"
	"github.com/HobbyOSs/gosk/token"
	"github.com/zeroflucs-given/generics/collections/stack"
)

type BitMode int

const (
	ID_16BIT_MODE BitMode = 16
	ID_32BIT_MODE BitMode = 32
	ID_64BIT_MODE BitMode = 64
)

var intToBitMode = map[int]BitMode{
	16: ID_16BIT_MODE,
	32: ID_32BIT_MODE,
	64: ID_64BIT_MODE,
}

func NewBitMode(i int) (BitMode, bool) {
	b, ok := intToBitMode[i]
	return b, ok
}

type Pass1 struct {
	// LOC(location of counter)
	LOC     int32
	BitMode BitMode
	EquMap  map[string]*token.ParseToken
	// Pass1のシンボルテーブル
	SymTable         map[string]uint32
	GlobalSymbolList []string
	ExternSymbolList []string
	Ctx              *stack.Stack[*token.ParseToken]
}

func (p *Pass1) Eval(program ast.Prog) {
	v := NewVisitor(p)
	v.Visit(program)
}

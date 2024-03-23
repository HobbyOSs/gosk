package pass1

import (
	"github.com/HobbyOSs/gosk/ast"
	"github.com/HobbyOSs/gosk/token"
	"github.com/zeroflucs-given/generics/collections/stack"
)

//go:generate newc
type Pass1 struct {
	// LOC(location of counter)
	LOC int32
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

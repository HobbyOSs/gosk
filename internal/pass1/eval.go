package pass1

import (
	"github.com/HobbyOSs/gosk/internal/ast"
	client "github.com/HobbyOSs/gosk/internal/ocode_client"
	"github.com/HobbyOSs/gosk/internal/token"
	"github.com/zeroflucs-given/generics/collections/stack"
)

type Pass1 struct {
	// LOC(location of counter)
	LOC     int32
	BitMode ast.BitMode
	EquMap  map[string]*token.ParseToken
	// Pass1のシンボルテーブル
	SymTable         map[string]int32
	GlobalSymbolList []string
	ExternSymbolList []string
	Ctx              *stack.Stack[*token.ParseToken]
	// 中間言語
	Client client.CodegenClient
}

func (p *Pass1) Eval(program ast.Prog) {
	TraverseAST(program, p)
}

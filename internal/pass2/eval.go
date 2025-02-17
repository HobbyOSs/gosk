package pass2

import (
	"github.com/HobbyOSs/gosk/internal/ast"
	client "github.com/HobbyOSs/gosk/internal/ocode_client"
	"github.com/HobbyOSs/gosk/internal/token"
	"github.com/zeroflucs-given/generics/collections/stack"
)

type Pass2 struct {
	BitMode          ast.BitMode
	EquMap           map[string]*token.ParseToken
	SymTable         map[string]int32
	GlobalSymbolList []string
	ExternSymbolList []string
	Ctx              *stack.Stack[*token.ParseToken]
	DollarPos        uint32 // $ の位置
	// 中間言語
	Client client.CodegenClient
}

func (p *Pass2) Eval(program ast.Prog) ([]byte, error) {
	//TraverseAST(program, p)
	return p.Client.Exec()
}

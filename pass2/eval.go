package pass2

import (
	"github.com/HobbyOSs/gosk/ast"
	"github.com/HobbyOSs/gosk/junkjit"
	"github.com/HobbyOSs/gosk/token"
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
	Asm              junkjit.Assembler
}

func (p *Pass2) Eval(program ast.Prog) {
	TraverseAST(program, p)
}

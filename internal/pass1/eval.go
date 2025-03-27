package pass1

import (
	"github.com/HobbyOSs/gosk/internal/ast" // Restored ast import
	"github.com/HobbyOSs/gosk/internal/client"
	"github.com/HobbyOSs/gosk/internal/token"
	"github.com/HobbyOSs/gosk/pkg/asmdb"
	"github.com/HobbyOSs/gosk/pkg/cpu" // Keep one cpu import
	"github.com/zeroflucs-given/generics/collections/stack"
)

type Pass1 struct {
	LOC              int32 // LOC(location of counter)
	BitMode          cpu.BitMode // Keep cpu.BitMode
	EquMap           map[string]*token.ParseToken
	SymTable         map[string]int32 // Pass1のシンボルテーブル
	NextImmJumpID    int              // 即値用のカウンタ
	DollarPosition   uint32           // ORG命令で設定されるエントリーポイントのアドレス
	GlobalSymbolList []string
	ExternSymbolList []string
	Ctx              *stack.Stack[*token.ParseToken]
	Client           client.CodegenClient // 中間言語
	AsmDB            *asmdb.InstructionDB
}

func (p *Pass1) Eval(program ast.Prog) { // Restored ast.Prog
	TraverseAST(program, p)
}

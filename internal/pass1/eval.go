package pass1

import (
	"github.com/HobbyOSs/gosk/internal/ast" // Restored ast import
	"github.com/HobbyOSs/gosk/internal/client"
	"github.com/HobbyOSs/gosk/pkg/asmdb"
	"github.com/HobbyOSs/gosk/pkg/cpu" // Keep one cpu import
)

type Pass1 struct {
	LOC              int32              // LOC(location of counter)
	BitMode          cpu.BitMode        // Keep cpu.BitMode
	MacroMap         map[string]ast.Exp // New map to store evaluated macro expressions
	SymTable         map[string]int32   // Pass1のシンボルテーブル
	NextImmJumpID    int                // 即値用のカウンタ
	DollarPosition   uint32             // ORG命令で設定されるエントリーポイントのアドレス
	GlobalSymbolList []string
	ExternSymbolList []string
	Client           client.CodegenClient // 中間言語
	AsmDB            *asmdb.InstructionDB
}

func (p *Pass1) Eval(program ast.Prog) { // Restored ast.Prog
	TraverseAST(program, p)
}

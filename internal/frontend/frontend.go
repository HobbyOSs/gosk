package frontend

import (
	"fmt"
	"os"

	"github.com/HobbyOSs/gosk/internal/ast"
	"github.com/HobbyOSs/gosk/internal/codegen"
	ocode_client "github.com/HobbyOSs/gosk/internal/ocode_client"
	"github.com/HobbyOSs/gosk/internal/pass1"
	"github.com/HobbyOSs/gosk/internal/pass2"
	"github.com/HobbyOSs/gosk/internal/token"
	"github.com/HobbyOSs/gosk/pkg/asmdb"
	"github.com/HobbyOSs/gosk/pkg/operand" // Add operand import
	"github.com/zeroflucs-given/generics/collections/stack"
)

// pass1, pass2を操作するモジュール
func Exec(parseTree any, assemblyDst string) (*pass1.Pass1, *pass2.Pass2) {

	// 読み書き可能, 新規作成, ファイル内容あっても切り詰め
	dstFile, err := os.OpenFile(assemblyDst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Printf("GOSK : can't open %s", assemblyDst)
		os.Exit(17)
	}
	defer dstFile.Close()

	prog, ok := (parseTree).(ast.Prog)
	if !ok {
		fmt.Printf("GOSK : failed to parse")
		os.Exit(-1)
	}

	// pass1のEvalを実行
	ctx := &codegen.CodeGenContext{BitMode: operand.MODE_16BIT} // Change ast.MODE_16BIT to operand.MODE_16BIT
	client, _ := ocode_client.NewCodegenClient(ctx)

	pass1 := &pass1.Pass1{
		LOC:              0,
		BitMode:          operand.MODE_16BIT, // Change ast.MODE_16BIT to operand.MODE_16BIT
		EquMap:           make(map[string]*token.ParseToken, 0),
		SymTable:         make(map[string]int32, 0),
		GlobalSymbolList: []string{},
		ExternSymbolList: []string{},
		Ctx:              stack.NewStack[*token.ParseToken](100),
		Client:           client,
		AsmDB:            asmdb.NewInstructionDB(),
	}
	pass1.Eval(prog)

	pass2 := &pass2.Pass2{
		BitMode:          pass1.BitMode,
		EquMap:           pass1.EquMap,
		SymTable:         pass1.SymTable,
		GlobalSymbolList: pass1.GlobalSymbolList,
		ExternSymbolList: pass1.ExternSymbolList,
		Ctx:              stack.NewStack[*token.ParseToken](100),
		Client:           pass1.Client,
		DollarPos:        pass1.DollarPosition,
	}
	code, err := pass2.Eval(prog)
	if err != nil {
		fmt.Printf("GOSK : failed to generate %s", err)
		os.Exit(-1)
	}

	_, err = dstFile.Write(code)
	if err != nil {
		fmt.Printf("GOSK : can't write %s", assemblyDst)
		os.Exit(-1)
	}

	return pass1, pass2
}

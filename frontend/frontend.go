package frontend

import (
	"fmt"
	"os"

	"github.com/HobbyOSs/gosk/ast"
	"github.com/HobbyOSs/gosk/junkjit"
	"github.com/HobbyOSs/gosk/pass1"
	"github.com/HobbyOSs/gosk/pass2"
	"github.com/HobbyOSs/gosk/token"
	"github.com/zeroflucs-given/generics/collections/stack"
)

// pass1, pass2を操作するモジュール
func Exec(parseTree any, assemblyDst string) {

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
	pass1 := &pass1.Pass1{
		LOC:              0,
		BitMode:          ast.ID_16BIT_MODE,
		EquMap:           make(map[string]*token.ParseToken, 0),
		SymTable:         make(map[string]int32, 0),
		GlobalSymbolList: []string{},
		ExternSymbolList: []string{},
		Ctx:              stack.NewStack[*token.ParseToken](100),
	}
	pass1.Eval(prog)

	code := &junkjit.CodeHolder{}
	asm := junkjit.NewAssembler(code)

	pass2 := &pass2.Pass2{
		BitMode:          pass1.BitMode,
		EquMap:           pass1.EquMap,
		SymTable:         pass1.SymTable,
		GlobalSymbolList: pass1.GlobalSymbolList,
		ExternSymbolList: pass1.ExternSymbolList,
		Ctx:              stack.NewStack[*token.ParseToken](100),
		Asm:              asm,
	}
	pass2.Eval(prog)

	_, err = dstFile.Write(code.Buffer())
	if err != nil {
		fmt.Printf("GOSK : can't write %s", assemblyDst)
		os.Exit(-1)
	}

	return
}

package frontend

import (
	"fmt"
	"os"

	"github.com/HobbyOSs/gosk/ast"
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

	// TODO: pass1のEvalを実行
	pass1 := pass1.NewPass1(
		0,
		make(map[string]uint32),
		[]string{},
		[]string{},
		stack.NewStack[*token.ParseToken](0),
	)
	pass1.Eval(prog)

	// TODO: pass2のEvalを実行
	pass2.Eval(prog)
}

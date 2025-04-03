package frontend

import (
	"fmt"
	"log" // Added missing import
	"os"
	"path/filepath" // Added missing import

	"github.com/HobbyOSs/gosk/internal/ast" // Restored ast import
	"github.com/HobbyOSs/gosk/internal/codegen"
	"github.com/HobbyOSs/gosk/internal/filefmt" // filefmt をインポート
	"github.com/HobbyOSs/gosk/internal/gen"     // Added missing import
	ocode_client "github.com/HobbyOSs/gosk/internal/ocode_client"
	"github.com/HobbyOSs/gosk/internal/pass1"
	"github.com/HobbyOSs/gosk/internal/pass2"
	"github.com/HobbyOSs/gosk/pkg/asmdb"
	"github.com/HobbyOSs/gosk/pkg/cpu" // cpu インポートを1つ保持
)

// pass1, pass2 を操作するモジュール
func Exec(parseTree any, assemblyDst string) (*pass1.Pass1, *pass2.Pass2) {

	// 読み書き可能、新規作成、ファイル内容があっても切り詰め
	dstFile, err := os.OpenFile(assemblyDst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Printf("GOSK : can't open %s", assemblyDst)
		os.Exit(17)
	}
	defer dstFile.Close()

	prog, ok := (parseTree).(ast.Prog) // ast.Prog を復元
	if !ok {
		fmt.Printf("GOSK : failed to parse")
		os.Exit(-1)
	}

	// pass1 の Eval を実行
	// CodeGenContext の初期化に GlobalSymbolList を追加
	ctx := &codegen.CodeGenContext{
		BitMode:          cpu.MODE_16BIT, // cpu.MODE_16BIT を保持
		SymTable:         make(map[string]int32),
		GlobalSymbolList: []string{},
		MachineCode:      []byte{},
	}
	client, _ := ocode_client.NewCodegenClient(ctx)

	pass1 := &pass1.Pass1{
		LOC:              0,
		BitMode:          cpu.MODE_16BIT,       // cpu.MODE_16BIT を保持
		SymTable:         ctx.SymTable,         // CodeGenContext の SymTable を共有
		GlobalSymbolList: ctx.GlobalSymbolList, // CodeGenContext の GlobalSymbolList を共有
		ExternSymbolList: []string{},
		Client:           client,
		AsmDB:            asmdb.NewInstructionDB(),
		MacroMap:         make(map[string]ast.Exp), // activeContext.md に基づいて MacroMap の初期化を追加
	}
	// pass1.Eval を呼び出し、ctx を直接更新
	pass1.Eval(prog, ctx) // Pass ctx, no return value assignment needed

	// pass2 の初期化に GlobalSymbolList を追加 (更新された ctx のリストを使用)
	pass2 := &pass2.Pass2{
		BitMode:          pass1.BitMode,
		OutputFormat:     pass1.OutputFormat,
		SourceFileName:   pass1.SourceFileName,
		CurrentSection:   pass1.CurrentSection,
		SymTable:         pass1.SymTable,       // pass1 と共有 (SymTable はマップなので参照が渡る)
		GlobalSymbolList: ctx.GlobalSymbolList, // 更新された ctx のリストを使用
		ExternSymbolList: pass1.ExternSymbolList,
		Client:           pass1.Client,
		DollarPos:        pass1.DollarPosition,
	}
	// pass2.Eval はエラーのみを返すように変更 (機械語は Client/CodeGenContext に格納)
	err = pass2.Eval(prog)
	if err != nil {
		fmt.Printf("GOSK : failed in pass2 %s", err)
		os.Exit(-1)
	}

	// ファイルフォーマットを選択して書き出し
	var format filefmt.FileFormat
	switch pass2.OutputFormat {
	case "WCOFF":
		format = &filefmt.CoffFormat{}
	// case "ELF": // 将来的に ELF もサポートする場合
	// 	format = &filefmt.ElfFormat{}
	default:
		// デフォルトはバイナリ直接書き出し (既存の動作)
		// ただし、pass2.Eval が []byte を返さなくなったため、CodeGenContext から取得
		_, err = dstFile.Write(ctx.MachineCode) // ctx.MachineCode を使用
		if err != nil {
			fmt.Printf("GOSK : can't write raw binary %s", assemblyDst)
			os.Exit(-1)
		}
		log.Printf("info: Output format '%s' not explicitly handled, writing raw binary.", pass2.OutputFormat)
		return pass1, pass2 // Raw バイナリ書き出しの場合はここで終了
	}

	// 選択されたフォーマットでファイルに書き出す
	err = format.Write(ctx, assemblyDst) // CodeGenContext を渡す
	if err != nil {
		fmt.Printf("GOSK : failed to write %s format file: %s", pass2.OutputFormat, err)
		os.Exit(-1)
	}

	return pass1, pass2
}

// ParseFile は指定されたファイルを解析します。
func ParseFile(filename string) (any, error) { // main.go での使用に基づいて ParseFile 関数を追加
	content, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}
	// パーサーインスタンスを作成する代わりに gen.Parse を直接使用します
	return gen.Parse(filepath.Base(filename), content)
}

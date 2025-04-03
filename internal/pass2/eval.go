package pass2

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/HobbyOSs/gosk/internal/ast" // Restored ast import
	"github.com/HobbyOSs/gosk/internal/client"
	"github.com/HobbyOSs/gosk/internal/token"
	"github.com/HobbyOSs/gosk/pkg/cpu" // Keep one cpu import
	"github.com/zeroflucs-given/generics/collections/stack"
)

type Pass2 struct {
	BitMode          cpu.BitMode // Keep cpu.BitMode
	OutputFormat     string      // [FORMAT "WCOFF"] の値を保持
	SourceFileName   string      // [FILE "naskfunc.nas"] の値を保持
	CurrentSection   string      // [SECTION .text] の値を保持
	EquMap           map[string]*token.ParseToken
	SymTable         map[string]int32
	GlobalSymbolList []string
	ExternSymbolList []string
	Ctx              *stack.Stack[*token.ParseToken]
	DollarPos        uint32               // $ の位置
	Client           client.CodegenClient // 中間言語
}

// Eval メソッドの戻り値を error のみに変更
func (p *Pass2) Eval(program ast.Prog) error { // Restored ast.Prog
	// TODO: このあたりの受け渡しおかしい
	ocodes := p.Client.GetOcodes()
	p.Client.SetDollarPosition(p.DollarPos)
	p.Client.SetSymbolTable(p.SymTable)
	// GlobalSymbolList も Client/CodeGenContext に渡す必要があるかもしれない
	// p.Client.SetGlobalSymbolList(p.GlobalSymbolList) // Client にメソッドを追加する必要がある

	for i, ocode := range ocodes {
		for j, operand := range ocode.Operands {
			if strings.Contains(operand, "{{.") {
				tmpl, err := template.New("").Parse(operand)
				if err != nil {
					return fmt.Errorf("failed to parse template: %v", err) // エラーのみを返す
				}

				var buf bytes.Buffer
				err = tmpl.Execute(&buf, p.SymTable)
				if err != nil {
					return fmt.Errorf("failed to execute template: %v", err) // エラーのみを返す
				}

				ocodes[i].Operands[j] = buf.String()
			}
		}
	}
	p.Client.SetOcodes(ocodes)

	// Client.Exec() を呼び出し、結果を CodeGenContext にセットする
	// (Client.Exec() が []byte, error を返すと仮定)
	// machineCode 変数は使わないので _ で受ける
	_, err := p.Client.Exec()
	if err != nil {
		return fmt.Errorf("codegen client execution failed: %w", err)
	}
	// CodeGenContext に MachineCode をセットする (Client 経由でアクセスする必要があるかもしれない)
	// 現状の Client インターフェースでは直接 CodeGenContext を取得できないため、
	// Client に SetMachineCode のようなメソッドを追加するか、
	// Exec の戻り値を使わずに Client 内部でセットする設計にする必要がある。
	// ここでは、Client が内部で CodeGenContext の MachineCode を更新すると仮定する。
	// もし Client.Exec が []byte を返すなら、それを捨てるか、あるいは
	// frontend 側で ctx.MachineCode = machineCode のようにセットする。
	// frontend.go の修正で ctx.MachineCode を使うようにしたので、ここでは何もしない。

	return nil // エラーがない場合は nil を返す
}

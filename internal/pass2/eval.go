package pass2

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/HobbyOSs/gosk/internal/ast"
	"github.com/HobbyOSs/gosk/internal/client"
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
	DollarPos        uint32               // $ の位置
	Client           client.CodegenClient // 中間言語
}

func (p *Pass2) Eval(program ast.Prog) ([]byte, error) {
	// TODO: このあたりの受け渡しおかしい
	ocodes := p.Client.GetOcodes()
	p.Client.SetDollarPosition(p.DollarPos)
	p.Client.SetSymbolTable(p.SymTable)

	for i, ocode := range ocodes {
		for j, operand := range ocode.Operands {
			if strings.Contains(operand, "{{.") {
				tmpl, err := template.New("").Parse(operand)
				if err != nil {
					return nil, fmt.Errorf("failed to parse template: %v", err)
				}

				var buf bytes.Buffer
				err = tmpl.Execute(&buf, p.SymTable)
				if err != nil {
					return nil, fmt.Errorf("failed to execute template: %v", err)
				}

				ocodes[i].Operands[j] = buf.String()
			}
		}
	}
	p.Client.SetOcodes(ocodes)
	return p.Client.Exec()
}

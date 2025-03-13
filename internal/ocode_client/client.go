package client

import (
	"fmt"
	"log"
	"strings"

	"github.com/HobbyOSs/gosk/internal/ast"
	"github.com/HobbyOSs/gosk/internal/codegen"
	"github.com/HobbyOSs/gosk/pkg/ocode"
)

// CodegenClient インターフェースの定義
type CodegenClient interface {
	Emit(line string) error
	EmitAll(text string) error
	Exec() ([]byte, error)
}

// ocodeClient 構造体の定義
type ocodeClient struct {
	ocodes []ocode.Ocode
}

// NewCodegenClient は新しい CodegenClient を返す
func NewCodegenClient() CodegenClient {
	return &ocodeClient{}
}

// Emit メソッドの実装
func (c *ocodeClient) Emit(line string) error {
	log.Printf("debug: emit %s\n", line)
	ocode, err := parseLineToOcode(line)
	if err != nil {
		return err
	}
	c.ocodes = append(c.ocodes, ocode)
	return nil
}

func parseLineToOcode(line string) (ocode.Ocode, error) {
	// スペースで分割
	parts := strings.Fields(line)
	if len(parts) == 0 {
		return ocode.Ocode{}, fmt.Errorf("empty line")
	}

	// Ocodeのバリデーション
	kind, err := ocode.OcodeKindString("Op" + parts[0])
	if err != nil {
		return ocode.Ocode{}, fmt.Errorf("invalid OcodeKind: %s", parts[0])
	}

	// オペランドを,区切りで取得
	var operands []string
	if len(parts) > 1 {
		operands_str := strings.Join(parts[1:], "")
		operands = strings.Split(operands_str, ",")
	}

	return ocode.Ocode{
		Kind:     kind,
		Operands: operands,
	}, nil
}

// EmitAll メソッドの実装
func (c *ocodeClient) EmitAll(text string) error {
	lines := strings.Split(strings.TrimSpace(text), "\n")
	for _, line := range lines {
		if err := c.Emit(strings.TrimSpace(line)); err != nil {
			return err
		}
	}
	return nil
}

// Exec メソッドの実装
func (c *ocodeClient) Exec() ([]byte, error) {
	machineCode := codegen.GenerateX86(c.ocodes, ast.MODE_32BIT)
	return machineCode, nil
}

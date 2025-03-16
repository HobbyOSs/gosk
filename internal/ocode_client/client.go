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
	GetOcodes() []ocode.Ocode
	SetOcodes(ocodes []ocode.Ocode)
	SetDollarPosition(pos uint32)
	SetLOC(loc int32) // SetLOCメソッドを追加
}

// ocodeClient 構造体の定義
type ocodeClient struct {
	Ocodes         []ocode.Ocode
	bitMode        ast.BitMode
	DollarPosition uint32      // エントリーポイントのアドレス
	LOC            int32       // Location Counter
	ctx            *codegen.CodeGenContext // CodeGenContextを保持
}

// NewCodegenClient は新しい CodegenClient を返す
func NewCodegenClient(bitMode ast.BitMode, ctx *codegen.CodeGenContext) CodegenClient {
	if ctx == nil {
		// デフォルトのContextを作成
		ctx = &codegen.CodeGenContext{
			MachineCode:    make([]byte, 0),
			VS:             nil,
			BitMode:        bitMode,
			DollarPosition: 0x7c00, // デフォルト値
			LOC:            0,
		}
	}

	return &ocodeClient{
		bitMode:        bitMode,
		DollarPosition: ctx.DollarPosition,
		LOC:            ctx.LOC,
		ctx:            ctx,
		Ocodes:         make([]ocode.Ocode, 0),
	}
}

// Emit メソッドの実装
func (c *ocodeClient) Emit(line string) error {
	log.Printf("debug: emit %s\n", line)
	ocode, err := parseLineToOcode(line)
	if err != nil {
		return err
	}
	c.Ocodes = append(c.Ocodes, ocode)
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

func (c *ocodeClient) GetOcodes() []ocode.Ocode {
	return c.Ocodes
}

func (c *ocodeClient) SetOcodes(ocodes []ocode.Ocode) {
	c.Ocodes = ocodes
}

func (c *ocodeClient) SetDollarPosition(pos uint32) {
	c.DollarPosition = pos
}

// SetLOC メソッドの実装
func (c *ocodeClient) SetLOC(loc int32) {
	c.LOC = loc
}

// Exec メソッドの実装
func (c *ocodeClient) Exec() ([]byte, error) {
	// 保持しているContextを使用
	c.ctx.DollarPosition = c.DollarPosition
	c.ctx.LOC = c.LOC
	return codegen.GenerateX86(c.Ocodes, c.bitMode, c.ctx), nil
}

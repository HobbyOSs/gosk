package client

import (
	"fmt"
	"log"
	"strings"

	"github.com/HobbyOSs/gosk/internal/ast"
	"github.com/HobbyOSs/gosk/internal/client"
	"github.com/HobbyOSs/gosk/internal/codegen"
	"github.com/HobbyOSs/gosk/internal/pass1"
	"github.com/HobbyOSs/gosk/pkg/ocode"
)

// ocodeClient 構造体の定義
type ocodeClient struct {
	Ocodes  []ocode.Ocode
	bitMode ast.BitMode
	ctx     *codegen.CodeGenContext // CodeGenContextを保持
	pass1   *pass1.Pass1            // pass1の結果を保持
}

// NewCodegenClient は新しい CodegenClient を返す
func NewCodegenClient(ctx *codegen.CodeGenContext, pass1 *pass1.Pass1) (client.CodegenClient, error) {
	if ctx == nil {
		return nil, fmt.Errorf("CodeGenContext must not be nil")
	}

	return &ocodeClient{
		Ocodes:  make([]ocode.Ocode, 0),
		bitMode: ctx.BitMode, // ctxから取得
		ctx:     ctx,
		pass1:   pass1,
	}, nil
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
	c.ctx.DollarPosition = uint64(pos)
}

func (c *ocodeClient) SetLOC(loc int32) {
	// Do nothing
}

// Exec メソッドの実装
func (c *ocodeClient) Exec() ([]byte, error) {
	// 保持しているContextを使用
	return codegen.GenerateX86(c.Ocodes, c.ctx.BitMode, c.ctx), nil
}

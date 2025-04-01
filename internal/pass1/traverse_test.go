package pass1

import (
	"strconv" // Add strconv import
	"testing"

	"github.com/HobbyOSs/gosk/internal/ast" // Import ast package
	"github.com/HobbyOSs/gosk/pkg/cpu"      // Add cpu import
	"github.com/comail/colog"
	"github.com/stretchr/testify/suite"

	"github.com/HobbyOSs/gosk/internal/client" // Import client for mock
)

// --- Mock CodegenClient for testing ---
type mockCodegenClient struct {
	client.CodegenClient // Embed the interface to avoid implementing all methods
	emittedLines         []string
}

func (m *mockCodegenClient) Emit(line string) error {
	m.emittedLines = append(m.emittedLines, line)
	return nil
}
func (m *mockCodegenClient) SetBitMode(mode cpu.BitMode) {} // No-op for this test
// Implement other methods as needed, or leave them unimplemented if embedding works

// --- End Mock CodegenClient ---

// ヘルパー関数: 単純な NumberExp (完全に評価済み) を作成します
func newExpectedNumberExp(val int64) *ast.NumberExp {
	// BaseExp と Factor は NumberExp の構造に必要ですが、値の比較には重要ではありません
	baseFactor := ast.NewNumberFactor(ast.BaseFactor{}, int(val))
	baseImmExp := ast.NewImmExp(ast.BaseExp{}, baseFactor)
	return ast.NewNumberExp(*baseImmExp, val)
}

// ヘルパー関数: 単純な IdentExp (IdentFactor を持つ ImmExp) を作成します
func newExpectedIdentExp(name string) *ast.ImmExp {
	baseFactor := ast.NewIdentFactor(ast.BaseFactor{}, name)
	return ast.NewImmExp(ast.BaseExp{}, baseFactor)
}

// ヘルパー関数: テストの期待値用に単純な MultExp を作成します
// 注意: 単純化のため、未解決の部分のファクターは ImmExp(IdentFactor) であると仮定します
func newExpectedMultExp(head ast.Exp, ops []string, tails []ast.Exp) *ast.MultExp { // head と tails を ast.Exp に変更
	return ast.NewMultExp(ast.BaseExp{}, head, ops, tails)
}

// ヘルパー関数: テストの期待値用に単純な AddExp を作成します
// 注意: 単純化のため、未解決の部分のファクターは ImmExp(IdentFactor) であると仮定します
func newExpectedAddExp(head *ast.MultExp, ops []string, tails []*ast.MultExp) *ast.AddExp {
	return ast.NewAddExp(ast.BaseExp{}, head, ops, tails)
}

// ヘルパー: 単一の ImmExp から MultExp を作成します (AddExp 構造用)
func multExpFromImm(imm *ast.ImmExp) *ast.MultExp {
	return ast.NewMultExp(ast.BaseExp{}, imm, nil, nil)
}

// ヘルパー: 数値文字列から ImmExp を作成します (AddExp の期待値で使用)
func immExpFromNumStr(numStr string) *ast.ImmExp {
	// これは、数値文字列が NumberFactor に正しく解析されることを前提としています
	// より堅牢なヘルパーは、潜在的な解析エラーを処理する可能性があります
	numVal, _ := strconv.Atoi(numStr)
	return ast.NewImmExp(ast.BaseExp{}, ast.NewNumberFactor(ast.BaseFactor{}, numVal))
}

type Pass1TraverseSuite struct { // 構造体の名前を変更
	suite.Suite
}

func TestPass1TraverseSuite(t *testing.T) { // テスト関数の名前を変更
	suite.Run(t, new(Pass1TraverseSuite)) // 名前変更された構造体を使用
}

func (s *Pass1TraverseSuite) SetupSuite() { // 名前変更された構造体を使用
	setUpColog(colog.LDebug)
}

// これらは test_helper.go に存在するはずです

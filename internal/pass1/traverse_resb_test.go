package pass1

import (
	"testing"

	"github.com/HobbyOSs/gosk/internal/ast" // Import ast package
	"github.com/HobbyOSs/gosk/internal/gen"
	"github.com/HobbyOSs/gosk/pkg/cpu" // Add cpu import
	"github.com/stretchr/testify/assert"
)

// TestResbExpression は、'$' を含む RESB 式が正しく評価され、LOC が更新されることをテストします。
func (s *Pass1TraverseSuite) TestResbExpression() {
	tests := []struct {
		name        string
		text        string      // RESB を含むアセンブリコードスニペット
		initialLOC  int32       // テスト開始時の LOC
		expectedLOC int32       // 処理後の期待される LOC
		bitMode     cpu.BitMode // Pass1 コンテキストのビットモード
	}{
		{
			name:        "RESB with $",
			text:        "RESB 0x7dfe - $",
			initialLOC:  0x7000, // $ がこの値に評価されると仮定
			expectedLOC: 0x7dfe, // 0x7000 + (0x7dfe - 0x7000) = 0x7dfe
			bitMode:     cpu.MODE_16BIT,
		},
		{
			name:        "RESB with simple number", // 既存の RESB の動作も確認
			text:        "RESB 100",
			initialLOC:  0x100,
			expectedLOC: 0x100 + 100,
			bitMode:     cpu.MODE_16BIT,
		},
		{
			name:        "RESB with label + $", // ラベルと $ を含むケース
			text:        "MY_LABEL EQU 0x8000\nRESB MY_LABEL - $",
			initialLOC:  0x7F00, // $ がこの値に評価されると仮定
			expectedLOC: 0x8000, // 0x7F00 + (0x8000 - 0x7F00) = 0x8000
			bitMode:     cpu.MODE_16BIT,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// 1. スニペット全体を Program として解析します
			parseTree, err := gen.Parse("", []byte(tt.text), gen.Entrypoint("Program"))
			s.Require().NoError(err, "Parsing program snippet failed")
			program, ok := parseTree.(*ast.Program)
			s.Require().True(ok, "Parsed result is not *ast.Program")

			// 2. Pass1 環境をセットアップします
			p := &Pass1{
				LOC:      tt.initialLOC, // 初期 LOC を設定
				BitMode:  tt.bitMode,
				SymTable: make(map[string]int32),
				MacroMap: make(map[string]ast.Exp),
				Client:   &mockCodegenClient{}, // Initialize Client with the mock
				// AsmDB は LOC 計算には不要
			}

			// 3. プログラム全体を TraverseAST で処理します
			TraverseAST(program, p)

			// 4. 最終的な LOC が期待値と一致するかアサートします
			assert.Equal(t, tt.expectedLOC, p.LOC, "Final LOC mismatch after processing RESB")
		})
	}
}

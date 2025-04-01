package pass1

import (
	"testing"

	"github.com/HobbyOSs/gosk/internal/ast" // Import ast package
	"github.com/HobbyOSs/gosk/internal/gen"
	"github.com/stretchr/testify/assert"
)

func (s *Pass1TraverseSuite) TestMultExp() { // 名前変更された構造体を使用
	tests := []struct {
		name         string
		text         string
		expectedNode ast.Exp // 期待される評価済みノード
	}{
		// --- NumberExp に評価されるべきケース ---
		{
			name:         "simple math 1",
			text:         "1005*8",
			expectedNode: newExpectedNumberExp(8040),
		},
		{
			name:         "simple math 2",
			text:         "512/4",
			expectedNode: newExpectedNumberExp(128),
		},
		{
			name:         "simple math 3",
			text:         "512*1024/4",
			expectedNode: newExpectedNumberExp(131072),
		},
		// --- Cases that should not be fully evaluated ---
		{
			name: "scale 1 (not evaluable)",
			text: "EDX*4",
			// Expected: MultExp{ ImmExp{EDX}, "*", ImmExp{4} }
			expectedNode: newExpectedMultExp(
				newExpectedIdentExp("EDX"), // ast.Exp compatible
				[]string{"*"},
				[]ast.Exp{immExpFromNumStr("4")}, // Pass as []ast.Exp
			),
		},
		{
			name: "scale 2 (not evaluable)",
			text: "ESI*8",
			// Expected: MultExp{ ImmExp{ESI}, "*", ImmExp{8} }
			expectedNode: newExpectedMultExp(
				newExpectedIdentExp("ESI"), // ast.Exp compatible
				[]string{"*"},
				[]ast.Exp{immExpFromNumStr("8")}, // Pass as []ast.Exp
			),
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// 入力テキストを MultExp として解析します
			got, err := gen.Parse("", []byte(tt.text), gen.Entrypoint("MultExp"))
			if !assert.NoError(t, err, "Parsing failed for input: %s", tt.text) {
				t.FailNow()
			}

			// 解析されたノードが実際に Exp であることを確認します
			startNode, ok := got.(ast.Exp)
			if !ok {
				t.Fatalf("Parsed node is not an ast.Exp, but %T", got)
			}

			// Pass1 環境をセットアップします
			p := &Pass1{
				SymTable: make(map[string]int32),
				MacroMap: make(map[string]ast.Exp),
			}

			// ノードを評価します
			evaluatedNode := TraverseAST(startNode, p) // Pass1 は ast.Env を実装します

			// 評価されたノードを期待されるノードと比較します
			switch expected := tt.expectedNode.(type) {
			case *ast.NumberExp:
				actual, ok := evaluatedNode.(*ast.NumberExp)
				assert.True(t, ok, "Expected *ast.NumberExp, got %T", evaluatedNode)
				if ok {
					assert.Equal(t, expected.Value, actual.Value, "Evaluated number value mismatch")
				}
			case *ast.MultExp:
				actual, ok := evaluatedNode.(*ast.MultExp)
				assert.True(t, ok, "Expected *ast.MultExp, got %T", evaluatedNode)
				if ok {
					// 構造チェックのための TokenLiteral を使用した基本的な比較
					assert.Equal(t, expected.TokenLiteral(), actual.TokenLiteral(), "Evaluated MultExp structure mismatch")
				}
			default:
				t.Fatalf("Unhandled expected node type: %T", tt.expectedNode)
			}
		})
	}
}

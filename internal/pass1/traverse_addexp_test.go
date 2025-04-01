package pass1

import (
	"strings" // Add strings import for MYLABEL check
	"testing"

	"github.com/HobbyOSs/gosk/internal/ast" // Import ast package
	"github.com/HobbyOSs/gosk/internal/gen"
	"github.com/stretchr/testify/assert"
)

func (s *Pass1TraverseSuite) TestAddExp() { // 名前変更された構造体を使用
	tests := []struct {
		name         string
		text         string
		expectedNode ast.Exp // 期待される評価済みノード
	}{
		// --- NumberExp に評価されるべきケース ---
		{
			name:         "+int",
			text:         "30",
			expectedNode: newExpectedNumberExp(30),
		},
		{
			name:         "-int",
			text:         "-30",
			expectedNode: newExpectedNumberExp(-30),
		},
		{
			name:         "hex",
			text:         "0x0ff0",
			expectedNode: newExpectedNumberExp(0x0ff0),
		},
		// CharFactor の評価は、必要に応じて ImmExp.Eval で調整が必要になる場合があります
		// {
		// 	name:         "char",
		// 	text:         "'A'", // 単純な文字
		// 	expectedNode: newExpectedNumberExp(65),
		// },
		{
			name:         "simple math 1",
			text:         "1 + 1",
			expectedNode: newExpectedNumberExp(2),
		},
		{
			name:         "simple math 2",
			text:         "4 - 2",
			expectedNode: newExpectedNumberExp(2),
		},
		{
			name:         "simple math 3",
			text:         "1 + 3 - 2 + 4",
			expectedNode: newExpectedNumberExp(6),
		},
		{
			name:         "complex math 1 (fully evaluated)",
			text:         "8 * 3 - 1", // Parser creates AddExp{ MultExp{8*3}, "-", ImmExp{1} }
			expectedNode: newExpectedNumberExp(23),
		},
		{
			name:         "label + constant (fully evaluated)",
			text:         "MYLABEL + 512", // MYLABEL = 0x8000 (defined below)
			expectedNode: newExpectedNumberExp(0x8000 + 512),
		},
		{
			name:         "label - constant (fully evaluated)",
			text:         "MYLABEL - 10", // MYLABEL = 0x8000
			expectedNode: newExpectedNumberExp(0x8000 - 10),
		},
		// --- Cases that should not be fully evaluated (contain unresolved identifiers) ---
		{
			name: "ident (not evaluable)",
			text: `_testZ009$`,
			// Expected: AddExp -> MultExp -> ImmExp -> IdentFactor
			expectedNode: newExpectedAddExp(
				multExpFromImm(newExpectedIdentExp(`_testZ009$`)),
				nil, nil,
			),
		},
		{
			name: "displacement 1 (not evaluable, constant folding)",
			text: "ESP+4",
			// Expected: AddExp{ MultExp{ImmExp{ESP}}, "+", MultExp{ImmExp{4}} }
			expectedNode: newExpectedAddExp(
				multExpFromImm(newExpectedIdentExp("ESP")), // Head is non-constant term ESP
				[]string{"+"}, // Operator
				[]*ast.MultExp{
					multExpFromImm(immExpFromNumStr("4")), // Tail is constant 4
				},
			),
		},
		{
			name: "displacement 2 (not evaluable, constant folding)",
			text: "ESP+12+8",
			// Expected: AddExp{ MultExp{ImmExp{ESP}}, "+", MultExp{ImmExp{20}} }
			expectedNode: newExpectedAddExp(
				multExpFromImm(newExpectedIdentExp("ESP")), // Head is non-constant term ESP
				[]string{"+"}, // Operator
				[]*ast.MultExp{
					multExpFromImm(immExpFromNumStr("20")), // Tail is sum of constants 20
				},
			),
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// テストケースに基づいてエントリポイントを調整します
			entrypoint := "AddExp"
			if tt.name == "ident (評価不可)" {
				entrypoint = "Exp" // 単一の識別子を Exp (ImmExp) として解析します
			}

			// 入力テキストを解析します
			got, err := gen.Parse("", []byte(tt.text), gen.Entrypoint(entrypoint))
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
			// 関連するテストのために MYLABEL マクロを定義します
			if strings.Contains(tt.text, "MYLABEL") {
				// ast.Exp としての格納を処理する DefineMacro メソッドを使用します
				p.DefineMacro("MYLABEL", newExpectedNumberExp(0x8000))
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
			case *ast.AddExp:
				actual, ok := evaluatedNode.(*ast.AddExp)
				assert.True(t, ok, "Expected *ast.AddExp, got %T", evaluatedNode)
				if ok {
					// 構造チェックのための TokenLiteral を使用した基本的な比較
					// 複雑なケースでは、より詳細な比較が必要になる場合があります
					assert.Equal(t, expected.TokenLiteral(), actual.TokenLiteral(), "Evaluated AddExp structure mismatch")
				}
			case *ast.ImmExp: // 未解決の識別子のようなケースの場合
				actual, ok := evaluatedNode.(*ast.ImmExp)
				assert.True(t, ok, "Expected *ast.ImmExp, got %T", evaluatedNode)
				if ok {
					assert.Equal(t, expected.TokenLiteral(), actual.TokenLiteral(), "Evaluated ImmExp mismatch")
				}
			default:
				t.Fatalf("Unhandled expected node type: %T", tt.expectedNode)
			}
		})
	}
}

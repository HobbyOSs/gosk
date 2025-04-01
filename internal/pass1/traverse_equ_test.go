package pass1

import (
	"testing"

	"github.com/HobbyOSs/gosk/internal/ast" // Import ast package
	"github.com/HobbyOSs/gosk/internal/gen"
	"github.com/HobbyOSs/gosk/pkg/cpu" // Add cpu import
)

// TestEQUExpansionInExpression は、TraverseAST が MacroMap を使用して、
// EQU ステートメントで定義された定数を含む式を正しく評価することを確認します。
func (s *Pass1TraverseSuite) TestEQUExpansionInExpression() {
	tests := []struct {
		name           string
		text           string      // EQU を含むアセンブリコードスニペット
		expressionText string      // EQU 処理後に評価する特定の式部分
		expectedValue  int64       // 評価後の期待される数値
		bitMode        cpu.BitMode // Pass1 コンテキストのビットモード
	}{
		{
			name: "Simple EQU constant evaluation",
			text: `
				MY_EQU_CONST EQU 500
				MOV AX, MY_EQU_CONST ; MY_EQU_CONST を評価
			`,
			expressionText: "MY_EQU_CONST",
			expectedValue:  500,
			bitMode:        cpu.MODE_16BIT,
		},
		{
			name: "EQU constant + number evaluation",
			text: `
				MY_EQU_CONST2 EQU 100
				ADD BX, MY_EQU_CONST2 + 20 ; MY_EQU_CONST2 + 20 を評価
			`,
			expressionText: "MY_EQU_CONST2 + 20",
			expectedValue:  120,
			bitMode:        cpu.MODE_16BIT,
		},
		{
			name: "EQU constant used in another EQU",
			text: `
				BASE_VAL EQU 1000
				OFFSET_VAL EQU BASE_VAL + 50
				MOV CX, OFFSET_VAL ; OFFSET_VAL を評価
			`,
			expressionText: "OFFSET_VAL",
			expectedValue:  1050,
			bitMode:        cpu.MODE_16BIT,
		},
		{
			name: "EQU constant in multiplication",
			text: `
				FACTOR EQU 8
				IMUL DX, FACTOR * 2 ; FACTOR * 2 を評価
			`,
			expressionText: "FACTOR * 2",
			expectedValue:  16,
			bitMode:        cpu.MODE_16BIT,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// 1. EQU を処理するために、スニペット全体を Program として解析します
			parseTree, err := gen.Parse("", []byte(tt.text), gen.Entrypoint("Program"))
			s.Require().NoError(err, "Parsing program snippet failed")
			// Statements フィールドにアクセスするために、具象型 *ast.Program にアサートします
			program, ok := parseTree.(*ast.Program)
			s.Require().True(ok, "Parsed result is not *ast.Program")

			// 2. Pass1 環境をセットアップします
			p := &Pass1{
				LOC:      0,
				BitMode:  tt.bitMode,
				SymTable: make(map[string]int32),
				MacroMap: make(map[string]ast.Exp),
				// TraverseAST による純粋な式評価には Client と AsmDB は不要です
			}

			// 3. EQU ステートメントを処理して MacroMap を設定します
			// これは、EQU を処理する Pass1.Eval の部分をシミュレートします。
			// まず EQU 式自体を評価する必要があります。
			// これで program.Statements に直接アクセスできます
			for _, stmt := range program.Statements {
				// ステートメントが EQU ステートメント (DeclareStmt で表される) かどうかを確認します
				if declareStmt, ok := stmt.(*ast.DeclareStmt); ok {
					// 現在の環境 (p) を使用して EQU で割り当てられた式を評価します
					// これは、EQU が以前の EQU に依存する場合を処理します。
					evaluatedEquNode := TraverseAST(declareStmt.Value, p)  // declareStmt.Value を使用
					evaluatedEquExpr, okExpr := evaluatedEquNode.(ast.Exp) // ast.Exp にアサート
					if !okExpr {
						t.Fatalf("EQU expression evaluation did not return an ast.Exp: %T", evaluatedEquNode)
					}
					_, isEvaluable := evaluatedEquExpr.(*ast.NumberExp) // 数値に評価されたかどうかを確認します
					if !isEvaluable {
						// EQU 式自体が完全に評価できなかった場合 (例: ラベルを含む)、
						// 部分的に評価された式を格納します。これらのテストでは、EQU は数値に解決されると仮定します。
						t.Logf("Warning: EQU expression for %s did not evaluate to a number: %T", declareStmt.Id.Value, evaluatedEquExpr) // declareStmt.Id.Value を使用
					}
					p.DefineMacro(declareStmt.Id.Value, evaluatedEquExpr) // declareStmt.Id.Value を使用し、アサートされた ast.Exp を渡します
				}
			}

			// 4. ターゲットの式文字列を個別に解析します
			// 評価をテストしたい特定の式を解析する必要があります。
			exprTree, err := gen.Parse("", []byte(tt.expressionText), gen.Entrypoint("Exp")) // "Exp" エントリポイントを使用
			s.Require().NoError(err, "Parsing target expression failed: %s", tt.expressionText)
			exprToEval, ok := exprTree.(ast.Exp)
			s.Require().True(ok, "Parsed expression is not ast.Exp")

			// 5. TraverseAST と設定された Pass1 環境を使用してターゲットの式を評価します
			evaluatedNode := TraverseAST(exprToEval, p)      // ast.Node を返します
			evaluatedExpr, okExpr := evaluatedNode.(ast.Exp) // ast.Exp にアサート
			if !okExpr {
				t.Fatalf("Target expression evaluation did not return an ast.Exp: %T", evaluatedNode)
			}

			// 6. 結果が期待される値を持つ NumberExp であることをアサートします
			expectedNode := newExpectedNumberExp(tt.expectedValue)
			actualNode, ok := evaluatedExpr.(*ast.NumberExp) // アサートされた evaluatedExpr を使用
			s.True(ok, "Expected evaluated node to be *ast.NumberExp, got %T for expression '%s'", evaluatedExpr, tt.expressionText)
			if ok {
				s.Equal(expectedNode.Value, actualNode.Value, "Evaluated value mismatch for expression '%s'", tt.expressionText)
			}
		})
	}
}

package pass1

import (
	"strconv" // Add strconv import
	"strings" // Add strings import for MYLABEL check
	"testing"

	"github.com/HobbyOSs/gosk/internal/ast" // Import ast package
	"github.com/HobbyOSs/gosk/internal/gen"
	"github.com/HobbyOSs/gosk/pkg/cpu" // Add cpu import
	"github.com/comail/colog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

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
			name:         "complex math 1 (完全に評価)",
			text:         "8 * 3 - 1", // パーサーは AddExp{ MultExp{8*3}, "-", ImmExp{1} } を作成します
			expectedNode: newExpectedNumberExp(23),
		},
		{
			name:         "label + constant (完全に評価)",
			text:         "MYLABEL + 512", // MYLABEL = 0x8000 (以下で定義)
			expectedNode: newExpectedNumberExp(0x8000 + 512),
		},
		{
			name:         "label - constant (完全に評価)",
			text:         "MYLABEL - 10", // MYLABEL = 0x8000
			expectedNode: newExpectedNumberExp(0x8000 - 10),
		},
		// --- 完全に評価されるべきではないケース (未解決の識別子を含む) ---
		{
			name: "ident (評価不可)",
			text: `_testZ009$`,
			// 期待値: AddExp -> MultExp -> ImmExp -> IdentFactor
			expectedNode: newExpectedAddExp(
				multExpFromImm(newExpectedIdentExp(`_testZ009$`)),
				nil, nil,
			),
		},
		{
			name: "displacement 1 (評価不可、定数畳み込み)",
			text: "ESP+4",
			// 期待値: AddExp{ MultExp{ImmExp{4}}, "+", MultExp{ImmExp{ESP}} }
			// 定数畳み込みにより、定数 4 が最初に配置されるはずです。
			expectedNode: newExpectedAddExp(
				multExpFromImm(immExpFromNumStr("4")), // Head は定数 4
				[]string{"+"},                         // 演算子
				[]*ast.MultExp{
					multExpFromImm(newExpectedIdentExp("ESP")), // Tail は非定数項 ESP
				},
			),
		},
		{
			name: "displacement 2 (評価不可、定数畳み込み)",
			text: "ESP+12+8",
			// 期待値: AddExp{ MultExp{ImmExp{20}}, "+", MultExp{ImmExp{ESP}} }
			// 定数畳み込みにより、12 + 8 = 20 が結合され、定数が最初に配置されるはずです。
			expectedNode: newExpectedAddExp(
				multExpFromImm(immExpFromNumStr("20")), // Head は定数の合計 20
				[]string{"+"},                          // 演算子
				[]*ast.MultExp{
					multExpFromImm(newExpectedIdentExp("ESP")), // Tail は非定数項 ESP
				},
			),
		},
		// 複数の文字を持つ StringFactor と CharFactor は、通常、算術的に評価できません
		// {
		// 	name: "string (評価不可)",
		// 	text: `"0x0ff0"`,
		// 	expectedNode: ast.NewImmExp(ast.BaseExp{}, ast.NewStringFactor(ast.BaseFactor{}, "0x0ff0")),
		// },
		// {
		// 	name: "char multi (評価不可)",
		// 	text: "'AB'",
		// 	expectedNode: ast.NewImmExp(ast.BaseExp{}, ast.NewCharFactor(ast.BaseFactor{}, "AB")),
		// },
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
		// --- 完全に評価されるべきではないケース ---
		{
			name: "scale 1 (評価不可)",
			text: "EDX*4",
			// 期待値: MultExp{ ImmExp{EDX}, "*", ImmExp{4} }
			expectedNode: newExpectedMultExp(
				newExpectedIdentExp("EDX"), // ast.Exp 互換
				[]string{"*"},
				[]ast.Exp{immExpFromNumStr("4")}, // []ast.Exp として渡す
			),
		},
		{
			name: "scale 2 (評価不可)",
			text: "ESI*8",
			// 期待値: MultExp{ ImmExp{ESI}, "*", ImmExp{8} }
			expectedNode: newExpectedMultExp(
				newExpectedIdentExp("ESI"), // ast.Exp 互換
				[]string{"*"},
				[]ast.Exp{immExpFromNumStr("8")}, // []ast.Exp として渡す
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

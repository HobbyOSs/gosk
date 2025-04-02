package pass1

import (
	"log"
	"strings"

	"github.com/HobbyOSs/gosk/internal/ast"
	"github.com/HobbyOSs/gosk/pkg/cpu"
)

// TraverseAST は ast.Node と *Pass1 (Env として) を受け取り、変換される可能性のある ast.Node を返します。
// スタック (env.Ctx) は使用しません。
func TraverseAST(node ast.Node, env *Pass1) ast.Node {
	if node == nil {
		return nil
	}

	// Pass1 の Env インターフェースを実装します。
	// これにより、'env' を Eval メソッドに直接渡すことができます。
	var evalEnv ast.Env = env // Pass1 は DefineMacro/LookupMacro メソッドを通じて ast.Env を実装します。

	switch n := node.(type) {
	case *ast.Program:
		newStatements := make([]ast.Statement, 0, len(n.Statements))
		for _, stmt := range n.Statements {
			// 各ステートメントを走査します。TraverseAST は Node を返すようになりました。
			processedStmt := TraverseAST(stmt, env)
			if processedStmt != nil {
				// 返されたノードが実際に Statement であることを確認します。
				if statement, ok := processedStmt.(ast.Statement); ok {
					newStatements = append(newStatements, statement)
				} else {
					// TraverseAST が評価された式 (EQU からの NumberExp など) を返した場合、ステートメントレベルで破棄します。
					log.Printf("info: TraverseAST returned a non-Statement node (%T) for a statement, discarding.", processedStmt)
				}
			}
		}
		// 変更されたステートメントを持つ新しい Program ノードを返します。
		return ast.NewProgram(newStatements) // NewProgram の呼び出しを修正

	case *ast.DeclareStmt: // EQU ステートメント
		// まず値の式を評価します。
		evalValueNode := TraverseAST(n.Value, env)
		evalValueExp, ok := evalValueNode.(ast.Exp)
		if !ok {
			log.Printf("error: EQU value expression %s evaluated to non-expression type %T", n.Value.TokenLiteral(), evalValueNode)
			return nil // またはエラーを適切に処理します
		}

		// Pass1 のメソッドを使用して環境にマクロを定義します。
		env.DefineMacro(n.Id.Value, evalValueExp)
		log.Printf("debug: Defined macro '%s' = %s", n.Id.Value, evalValueExp.TokenLiteral())
		// EQU ステートメント自体は出力を生成しないため、nil を返します。
		return nil

	case *ast.LabelStmt:
		label := strings.TrimSuffix(n.Label.Value, ":")
		log.Printf("debug: [LOC Before Label] LOC: 0x%x (%d) before processing label '%s'", env.LOC, env.LOC, label) // ラベル処理前のLOC (traceレベルに変更)
		env.SymTable[label] = env.LOC
		log.Printf("debug: [LabelStmt] Defined label '%s' at LOC 0x%x (%d)", label, env.LOC, env.LOC) // ラベル定義時のログ追加 (traceレベルに変更)
		// ラベルステートメント自体は処理後に出力ノードを生成しません。
		return nil // または、pass2 でラベルを AST に残す場合は n を返します

	case *ast.MnemonicStmt:
		opcode := n.Opcode.Value

		// TraverseAST -> Eval を使用して最初にオペランドを評価します
		evalOperands := make([]ast.Exp, len(n.Operands))
		canProcess := true
		for i, operand := range n.Operands {
			evalOperandNode := TraverseAST(operand, env)
			if expOperand, ok := evalOperandNode.(ast.Exp); ok {
				evalOperands[i] = expOperand
			} else {
				log.Printf("error: Operand %d for %s evaluated to non-expression type %T", i, opcode, evalOperandNode)
				canProcess = false
				break // オペランドが無効な場合は処理を停止します
			}
		}

		if !canProcess {
			// オペランドが無効な場合、この命令を処理できません
			return n // オペランドが無効な場合は元のノードを返します
		}

		// 適切なハンドラ関数を見つけて呼び出します
		if handler, ok := opcodeEvalFns[opcode]; ok {
			handler(env, evalOperands) // 評価されたオペランドでハンドラを呼び出します
		} else {
			log.Printf("error: No handler found for opcode %s", opcode)
			// 不明なオペコードの処理方法を決定します (例: スキップ、エラー、デフォルトサイズ?)
			// 現時点ではログのみ。LOC は更新されません。
		}

		// LOC と Emit は特定のハンドラ関数内で処理される
		log.Printf("debug: [LOC After Mnemonic] LOC: 0x%x (%d) after processing '%s'", env.LOC, env.LOC, opcode) // 命令処理後のLOC (traceレベルに変更)
		return nil                                                                                               // Mnemonic ステートメントが処理されたため、nil を返します

	case *ast.OpcodeStmt: // オペランドのない命令
		opcode := n.Opcode.Value

		// 適切なハンドラ関数を見つけて呼び出します (空のオペランドを渡します)
		if handler, ok := opcodeEvalFns[opcode]; ok {
			handler(env, []ast.Exp{}) // 空のオペランドでハンドラを呼び出します
		} else {
			log.Printf("error: No handler found for opcode %s", opcode)
		}

		// LOC と Emit は特定のハンドラ関数内で処理されるようになりました。
		log.Printf("debug: [LOC After Opcode] LOC: 0x%x (%d) after processing '%s'", env.LOC, env.LOC, opcode) // 命令処理後のLOC (traceレベルに変更)
		return nil                                                                                             // Opcode ステートメントが処理されたため、nil を返します

	// --- 式の評価 ---
	// AddExp と MultExp には ast_exp_impl.go で実装された特定の Eval ロジックがあります
	// ImmExp、SegmentExp、MemoryAddrExp にも Eval メソッドがあります
	case ast.Exp: // すべての式タイプをキャッチします
		evalExp, _ := n.Eval(evalEnv) // ast.Env 型である evalEnv を使用します
		return evalExp                // 評価された式ノードを返します

	// --- ファクターの処理 ---
	case *ast.NumberFactor, *ast.StringFactor, *ast.HexFactor, *ast.IdentFactor, *ast.CharFactor:
		log.Printf("warning: TraverseAST encountered a Factor type (%T) directly. Wrapping in ImmExp.", n)
		// ファクターは式の一部である必要があるため、返す前に ImmExp でラップします。
		return ast.NewImmExp(ast.BaseExp{}, n.(ast.Factor))

	// --- その他のステートメントタイプ ---
	case *ast.ExportSymStmt:
		// TODO: 必要に応じてロジックを実装します (例: env.GlobalSymbolList に追加)
		return nil // または n を返します
	case *ast.ExternSymStmt:
		// TODO: 必要に応じてロジックを実装します (例: env.ExternSymbolList に追加)
		return nil // または n を返します
	case *ast.ConfigStmt:
		if n.ConfigType == ast.Bits {
			// ファクターを評価してビットモード値を取得します
			factorNode := TraverseAST(n.Factor, env)
			// ファクターは上記の Factor ケースで ImmExp でラップされている必要があります
			immExp, ok := factorNode.(*ast.ImmExp)
			if !ok {
				log.Printf("error: BITS directive requires a constant value, got %T", factorNode)
				return nil
			}
			evalExp, _ := immExp.Eval(evalEnv) // evalEnv を使用します
			numExp, ok := evalExp.(*ast.NumberExp)
			if !ok {
				log.Printf("error: BITS directive value did not evaluate to a number: %s", evalExp.TokenLiteral())
				return nil
			}

			bitModeVal := int(numExp.Value)
			bitMode, ok := cpu.NewBitMode(bitModeVal)
			if !ok {
				log.Printf("error: Invalid bit mode value %d for BITS directive", bitModeVal)
				return nil
			}
			env.BitMode = bitMode
			env.Client.SetBitMode(bitMode)
			log.Printf("debug: Set bit mode to %d", bitModeVal)
		}
		return nil // Config ステートメントは出力ノードを生成しません

	default:
		log.Printf("Unknown AST node type in TraverseAST: %T\n", node)
		return node // 不明なノードは変更せずに返します
	}
}

// DefineMacro は、メソッドとして定義することにより、Pass1 の ast.Env インターフェースを実装します。
func (p *Pass1) DefineMacro(name string, exp ast.Exp) {
	// 新しいマップが nil の場合は初期化します
	if p.MacroMap == nil {
		p.MacroMap = make(map[string]ast.Exp)
	}
	p.MacroMap[name] = exp
	log.Printf("debug: Defined macro '%s' = %s (stored as ast.Exp)", name, exp.TokenLiteral())

}

// LookupMacro は、メソッドとして定義することにより、Pass1 の ast.Env インターフェースを実装します。
func (p *Pass1) LookupMacro(name string) (ast.Exp, bool) {
	// 新しい MacroMap を使用します
	if p.MacroMap == nil {
		return nil, false // マップが初期化されていません
	}
	exp, ok := p.MacroMap[name]
	// Eval ロジックでは、古い EquMap へのフォールバックは不要です
	return exp, ok
}

// GetLOC は Pass1 の ast.Env インターフェースを実装します。
// 現在のロケーションカウンタを返します。
func (p *Pass1) GetLOC() int32 {
	return p.LOC
}

// GetConstValue は Pass1 の ast.Env インターフェースを実装します。
// ローカルの getConstValue ヘルパー関数をラップします。
func (p *Pass1) GetConstValue(exp ast.Exp) (int, bool) {
	return getConstValue(exp)
}

// getConstValue は、式が定数である場合に整数値を抽出します。
func getConstValue(exp ast.Exp) (int, bool) {
	// まず、すでに NumberExp (以前の評価の結果) であるかどうかを確認します
	if numExp, ok := exp.(*ast.NumberExp); ok {
		return int(numExp.Value), true
	}
	// そうでない場合は、NumberFactor を含む ImmExp であるかどうかを確認します
	if imm, ok := exp.(*ast.ImmExp); ok {
		if num, ok := imm.Factor.(*ast.NumberFactor); ok {
			return num.Value, true
		}
		// 必要に応じてここで HexFactor を処理します (数値に評価されると仮定)
		// ImmExp.Eval が最初に HexFactor の評価を処理する必要があるため、現時点では空白識別子を使用します。
		if _, ok := imm.Factor.(*ast.HexFactor); ok {
			// 単純化のため、ImmExp.Eval が最初に HexFactor の評価を処理するようにします
			// ImmExp.Eval が NumberExp を返す場合、最初のチェックでキャッチされます。
		}
	}
	return 0, false
}

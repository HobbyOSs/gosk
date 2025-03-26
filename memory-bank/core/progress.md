# Progress

## 実装済み機能 (要約)
- x86アセンブラ解析基盤 (pass1, pass2)
- Ocode・PEGパーサ等の基礎部分
- 主要命令実装 (システム命令, 算術命令, CMP, JE, MOV, ADD, JMP)
- CodegenClient 関連機能
- EQU命令の展開
- `getModRMFromOperands`の返り値の型変更 (`uint32` -> `[]byte`)
- InstructionForm の Encoding を ModRM の要否で振り分け
- OUT命令の実装
- CALL命令の実装
- 論理命令の実装 (AND, OR, XOR, NOT)

## まだ必要な実装
- JMP系命令 (Jcc命令) のrel32オフセット対応
- RESBの計算処理の実装
- メモリアドレッシング
- `internal/codegen` パッケージのリファクタリング (CodeGenContext 導入)
- `internal/codegen` パッケージの不要パラメータ削除

## 実装済み機能 (詳細)
- ASTノードの文字列化ヘルパー関数 `ExpToString` を `internal/ast` パッケージに実装
- `FactorToString` 関数を `internal/ast/ast_factor_impl.go` に実装
- `SegmentExp`, `AddExp`, `MultExp` の `TokenLiteral()` メソッドを `internal/ast/ast_exp_impl.go` で修正
- `ExpToString` 関数のテストコードを `internal/ast/ast_string_test.go` に実装し、テストをパス
- `pkg/asmdb/instruction_search.go` の `filterForms` 関数を修正し、ModRM 要否によるフィルタリングロジックを実装
- `pkg/asmdb/instruction_search_test.go` にテストケースを追加
- CALL命令の実装
- 論理命令の実装 (AND, OR, XOR, NOT) (pass1, codegen, test)

## 関連情報
[technical_notes.md](../details/technical_notes.md)

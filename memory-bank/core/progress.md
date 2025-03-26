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
- 論理シフト/算術シフト命令の実装 (SHR, SHL, SAR) (一部テストはコメントアウト)
- IN命令の実装 (pass1, codegen, fallback table)
- RET命令の実装 (pass1, codegen, test)

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
- 論理シフト/算術シフト命令の実装 (SHR, SHL, SAR) (pass1, ocode, codegen, test) (一部テストはコメントアウト)
- IN命令の実装 (pass1, codegen, fallback table)
- RET命令の実装 (pass1, ocode, codegen, test)
- `internal/codegen/x86gen.go`: `processOcode` 関数を修正し、オペランドなし命令 (`CLI` など) を `opcodeMap` を使って処理するように変更。
- `internal/codegen/x86gen_lgdt.go`: `handleLGDT` 関数を修正し、`LGDT [label]` 形式を正しく処理するように変更。不要なインポートを削除。
- `internal/codegen/x86gen_utils.go`: `ResolveOpcode` 関数を修正し、複数バイトのオペコード文字列 (`0F20` など) を処理できるように変更。戻り値を `[]byte` に変更。
- `internal/codegen/x86gen_logical.go`, `x86gen_arithmetic.go`, `x86gen_mov.go`: `ResolveOpcode` の変更に合わせて `append` を修正 (`opcode...`)。

## 関連情報
[technical_notes.md](../details/technical_notes.md)

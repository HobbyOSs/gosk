# Progress

## 実装済み機能 (要約)
- x86アセンブラ解析基盤 (pass1, pass2)
- Ocode・PEGパーサ等の基礎部分
- システム命令 (INT, HLT)
  - grammar_test.go の HLT 命令関連テストを修正
- 算術命令(一部)
- `internal/ocode_client/client.go` の `Exec()` メソッドで `BitMode` を呼び出し元で渡せるように修正
- CMP命令の呼び出し修正
- CMP命令のテストケース追加
- `pkg/operand/operand.go`に`InternalStrings() []string`メソッドを追加
- `pkg/operand/operand_impl.go`に`InternalStrings()`メソッドを実装し、キャッシュ処理を追加
- `pkg/asmdb/instruction_search.go`の`matchOperands`関数を修正してアキュムレータレジスタの優先検索を実装
- JE命令の追加
- MOV命令 (レジスタ間, 即値)
- ADD命令 (フラグ更新)
- JMP命令のラベル解決
  - pass1でラベルをテンプレート文字列としてEmit
  - pass2でテンプレート文字列をアドレスに置換
- CodegenClientインターフェースの拡張
  - GetOcodes/SetOcodesメソッドを追加
- CodeGenContextへのBitModeの移動
  - `internal/codegen/typedef.go` に `BitMode` を追加
  - `internal/ocode_client/client.go` で `NewCodeGenContext` を呼び出す際に `bitMode` を渡すように変更
  - `internal/ocode_client/client.go` の `Exec()` メソッドで `CodeGenContext` から `bitMode` を取得するように変更
  - 関連するテストファイル(`internal/ocode_client/client_test.go`, `test/pass1_test.go`, `internal/frontend/frontend.go`)の`NewCodegenClient`呼び出しを修正 (test/pass1_test.go, internal/frontend/frontend.go はユーザーが修正)
- `internal/ocode_client/client.go` の `NewCodegenClient` 関数を修正し、`ctx == nil` の場合にエラーを返すように変更
- `internal/ocode_client/client_test.go` を修正し、上記変更に対応

## 実装内容 (day01, day02)
- day01
  - x86アセンブラ解析基盤 (pass1, pass2)
  - Ocode・PEGパーサ等の基礎部分
  - システム命令 (INT, HLT)
    - grammar_test.go の HLT 命令関連テストを修正
  - 算術命令(一部)
- day02
  - `internal/ocode_client/client.go` の `Exec()` メソッドで `BitMode` を呼び出し元で渡せるように修正
  - CMP命令の呼び出し修正
  - CMP命令のテストケース追加
  - `pkg/operand/operand.go`に`InternalStrings() []string`メソッドを追加
  - `pkg/operand/operand_impl.go`に`InternalStrings()`メソッドを実装し、キャッシュ処理を追加
  - `pkg/asmdb/instruction_search.go`の`matchOperands`関数を修正してアキュムレータレジスタの優先検索を実装
  - JE命令の追加
  - MOV命令 (レジスタ間, 即値)
  - ADD命令 (フラグ更新)
  - JMP命令のラベル解決
    - pass1でラベルをテンプレート文字列としてEmit
    - pass2でテンプレート文字列をアドレスに置換
  - CodegenClientインターフェースの拡張
    - GetOcodes/SetOcodesメソッドを追加
  - CodeGenContextへのBitModeの移動
    - `internal/codegen/typedef.go` に `BitMode` を追加
    - `internal/ocode_client/client.go` で `NewCodeGenContext` を呼び出す際に `bitMode` を渡すように変更
    - `internal/ocode_client/client.go` の `Exec()` メソッドで `CodeGenContext` から `bitMode` を取得するように変更
    - 関連するテストファイル(`internal/ocode_client/client_test.go`, `test/pass1_test.go`, `internal/frontend/frontend.go`)の`NewCodegenClient`呼び出しを修正 (test/pass1_test.go, internal/frontend/frontend.go はユーザーが修正)
  - `internal/ocode_client/client.go` の `NewCodegenClient` 関数を修正し、`ctx == nil` の場合にエラーを返すように変更
  - `internal/ocode_client/client_test.go` を修正し、上記変更に対応

## まだ必要な実装
- メモリアドレッシング
- `internal/codegen` パッケージ内の関数で、`CodeGenContext` をパラメータオブジェクトとして使用するようにリファクタリング
- `internal/codegen` パッケージ内の不要になったパラメータを削除

(細かな実装ステップや過去履歴は [implementation_details.md](../details/implementation_details.md) に記載)

## 関連情報
[technical_notes.md](../details/technical_notes.md)

---
## 実装済み機能 (詳細) - 2025/03/27 アーカイブ
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
- `internal/codegen/x86gen_utils.go`:
    - `ResolveOpcode` 関数を修正し、複数バイトのオペコード文字列 (`0F20` など) を処理できるように変更。戻り値を `[]byte` に変更。
    - `GetRegisterNumber` 関数を修正し、制御レジスタ (CR0, CR2, CR3, CR4) に対応。
    - `ModRMByOperand` 関数を修正し、`bitMode` に基づいて 16bit/32bit メモリオペランド処理を分岐。16bit モードの処理を改善。
    - 未使用の `regexp` インポートを削除。
    - ローカルヘルパー関数 `parseNumeric` を追加。
- `internal/codegen/x86gen_utils.go` のリファクタリング:
    - `modeStr` の switch 文を共通関数 `parseMode` として切り出し。
    - `ModRMByOperand` および `ModRMByValue` がメモリオペランド解析に `pkg/operand.ParseMemoryOperand` を使用するように修正。
    - 冗長な16bitモードの手動解析ロジック、`parseNumeric` 関数、`encoding/binary` インポートを削除。
    - 英語コメントを日本語に翻訳。
- `internal/codegen/x86gen_logical.go`, `x86gen_arithmetic.go`, `x86gen_mov.go`: `ResolveOpcode` の変更に合わせて `append` を修正 (`opcode...`)。

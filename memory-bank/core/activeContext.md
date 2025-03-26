# Active Context

## 現在の作業の焦点
- `test/day03_harib00i_test.go` のテスト失敗に対するデバッグ作業中。

## 直近の変更点
- `internal/codegen/x86gen.go`: `processOcode` 関数を修正し、オペランドなし命令 (`CLI` など) を `opcodeMap` を使って処理するように変更。
- `internal/codegen/x86gen_lgdt.go`: `handleLGDT` 関数を修正し、`LGDT [label]` 形式を正しく処理するように変更。不要なインポートを削除。
- `internal/codegen/x86gen_utils.go`: `ResolveOpcode` 関数を修正し、複数バイトのオペコード文字列 (`0F20` など) を処理できるように変更。戻り値を `[]byte` に変更。
- `internal/codegen/x86gen_logical.go`: `ResolveOpcode` の変更に合わせて `append` を修正 (`opcode...`)。
- `internal/codegen/x86gen_arithmetic.go`: `ResolveOpcode` の変更に合わせて `append` を修正 (`opcode...`)。
- `internal/codegen/x86gen_mov.go`: `ResolveOpcode` の変更に合わせて `append` を修正 (`opcode...`)。

## 次のステップ
- `MOV r32, CR0` および `MOV CR0, r32` 命令の処理に関するエラー (`Failed to generate ModR/M: failed to get register number for CR0`) の修正。
    - `pkg/asmdb/instruction_table_fallback.go` のエンコーディング定義修正。
    - `internal/codegen/x86gen_mov.go` の `handleMOV` 関数修正。

## 関連情報
[technical_notes.md](../details/technical_notes.md)
[implementation_details.md](../details/implementation_details.md)

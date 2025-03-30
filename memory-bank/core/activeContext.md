# 現在の状況 (Active Context)

## 完了した作業 (2025/03/30)
- `internal/codegen/x86gen_test.go` のテスト構造を改善し、`BitMode` をテストケースごとに指定可能にした。
- `internal/codegen/x86gen_utils.go` の `calculateModRM` を修正し、16bitモードでの32bitアドレッシングモード（`67h` プレフィックスが必要なケース）に対応。
- 上記修正に伴い、`TestGenerateX86` に関連テストケースを追加し、既存テストの期待値を修正。これにより `TestGenerateX86` が成功するようになった。
- `TestDay03Suite/TestHarib00i` を再実行し、ModR/M関連のエラーは解消されたことを確認。
- `pkg/asmdb/instruction_table_fallback.go` の `OUT imm8, AL` 等の fallback 定義のオペランド順序を修正。テスト実行結果から、この修正によりOUT命令関連のエラーは解消されたと判断。
- `pkg/ng_operand/operand_grammar.peg` を修正し、`MemoryAddress` ルールのアクションで `MemoryBody` の結果を正しく処理するように変更。これにより `LGDT [ GDTR0 ]` のパースエラーを解消。
- `pkg/asmdb/instruction_table_fallback.go` の `MOV creg, r32` / `MOV r32, creg` のオペランドタイプを `"cr"` から `"creg"` に修正。
- `internal/codegen/x86gen_lgdt.go` を修正し、ビットモードに応じて `LGDT` 命令の ModR/M とディスプレースメントを正しく生成するように変更。
- `internal/pass1/pass1_inst_lgdt.go` を修正し、`LGDT` 命令のサイズをビットモードに応じて正しく計算するように変更 (asmdb 呼び出しを削除)。これにより `LGDT` 命令のエンコーディングエラーは解消。
- **`pass1` JMP/CALL LOC 計算修正**:
    - `internal/pass1/pass1_inst_call.go`: ラベル参照時に `CALL rel16` (3バイト) を仮定するように修正。
    - `internal/pass1/pass1_inst_jmp.go`: ラベル参照時に `JMP rel16` (3バイト) または `Jcc rel16` (4バイト) を仮定するように修正。
- **`codegen` JMP/CALL エンコーディング修正**:
    - `internal/codegen/x86gen_call.go`: オフセットに応じて `CALL rel16`/`rel32` を生成するように修正。
    - `internal/codegen/x86gen_jmp.go`: `Jcc rel16` のオフセット計算を修正。
- **`codegen` IMUL 命令処理修正**:
    - `internal/codegen/x86gen_arithmetic.go`: `handleIMUL`, `handleSUB` を追加。
    - `internal/codegen/x86gen.go`: `OpIMUL`, `OpSUB` を `switch` 文に追加。
    - `internal/codegen/x86gen_no_param.go`: `opcodeMap` から `OpIMUL` を削除。
    - `pkg/ocode/ocode.go`: `OpSUB` 定数を追加し、`go generate` を実行。
    - `pkg/asmdb/instruction_table_fallback.go`: `IMUL r/m, imm` (オペコード `6B`) の `ModRM` 定義を修正。

## 次のステップ
- `TestDay03Suite/TestHarib00i` のテスト失敗を調査・修正。
    - 現在の問題: バイナリ差分と長さの不一致 (expected 304, actual 300)。
    - 原因: `pass1` での LOC 計算のずれが依然として残っている可能性が高い。特に 32bit アドレッシング (`67h`) や `MOV r32, CR0` (`66h`) 関連。
    - 調査対象:
        - `pkg/ng_operand/requires.go` (`Require66h`, `Require67h`)
        - `pkg/ng_operand/operand_impl.go` (`CalcOffsetByteSize`)
        - `pkg/asmdb` (`FindMinOutputSize`)

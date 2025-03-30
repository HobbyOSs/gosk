# Active Context

## 現在の作業の焦点
- `internal/codegen` パッケージ内の関数で、`CodeGenContext` をパラメータオブジェクトとして使用するようにリファクタリング

## 直近の変更点
- day02までの実装完了
  - システム命令 (INT, HLT)
  - 算術命令(一部)
  - CMP命令の呼び出し修正とテストケース追加
  - `pkg/operand/operand.go`に`InternalStrings() []string`メソッドを追加
  - `pkg/asmdb/instruction_search.go`の`matchOperands`関数を修正
  - JE命令、MOV命令 (レジスタ間, 即値)、ADD命令 (フラグ更新)、JMP命令のラベル解決を追加
  - CodegenClientインターフェースの拡張 (GetOcodes/SetOcodesメソッドを追加)
  - CodeGenContextへのBitModeの移動
  - `internal/ocode_client/client.go` の `NewCodegenClient` 関数を修正し、エラーハンドリングを追加

## 次のステップ
- `internal/codegen` パッケージのリファクタリング完了 (CodeGenContextパラメータオブジェクト化)
- `internal/codegen` パッケージ内の不要になったパラメータを削除
- メモリアドレッシングの実装

## 関連情報
[technical_notes.md](../details/technical_notes.md)

---
## 直近の変更点 - 2025/03/27 アーカイブ
- `internal/codegen/x86gen_utils.go`:
    - `GetRegisterNumber` 関数を修正し、制御レジスタ (CR0, CR2, CR3, CR4) に対応。
    - `ModRMByOperand` 関数を修正し、`bitMode` に基づいて 16bit/32bit メモリオペランド処理を分岐。
        - 16bit モード処理を改善し、単純なレジスタ間接参照 (`[SI]`, `[DI]`) および直接アドレス (`[imm16]`) の ModR/M とディスプレースメントを生成するように修正。これにより `TestGenerateX86` スイートのデグレを解消。
    - 未使用の `regexp` インポートを削除。
    - `operand.ParseNumeric` の代わりにローカルヘルパー関数 `parseNumeric` を追加・使用。
- `internal/codegen/x86gen_utils.go` のリファクタリング:
    - `modeStr` の switch 文を共通関数 `parseMode` として切り出し。
    - `ModRMByOperand` および `ModRMByValue` がメモリオペランド解析に `pkg/operand.ParseMemoryOperand` を使用するように修正。
    - 冗長な16bitモードの手動解析ロジック、`parseNumeric` 関数、`encoding/binary` インポートを削除。
    - 英語コメントを日本語に翻訳。

---
## 2025/03/29 アーカイブ (pkg/operand パーサー基本修正後)

## 現在の作業の焦点
- `pkg/operand` のパーサー (`participle` ベース) の基本的な問題を修正。
- 引き続き `test/day03_harib00i_test.go` のエラー (エンコーディングエラー、バイナリ長不一致) 対応。

## 直近の変更点
- **`pkg/operand` パーサー修正 (2025/03/29):**
    - `pkg/operand` のテスト (`TestBaseOperand_OperandType`) で失敗していたオペランド解析の問題を特定。
    - TDD アプローチで `pkg/operand/operand_impl_test.go` にテストケースを追加。
    - `pkg/operand/operand_impl.go` のレキサールール (`Reg`, `Seg` の順序) を修正。
    - `pkg/operand/operand.go` のパーサー定義 (`Instruction`, `CommaOperand` 構造体) を修正。
    - `pkg/operand/operand_impl.go` の `OperandTypes` ロジック (ラベルの扱い) を修正。
    - `pkg/operand/operand_impl_test.go` のテストケース (`MOV r32, label`) で `ForceRelAsImm=true` を設定。
    - 上記修正により `pkg/operand` のテスト (`TestBaseOperand_OperandType`) が成功。

## 次のステップ
- `test/day03_harib00i_test.go` を再実行し、エラー内容を確認する。
- エラー内容に基づき、エンコーディングエラーやバイナリ長不一致の原因を調査・修正する。
    - `asmdb` (JSON, fallback) のエンコーディング定義確認
    - `codegen` (MOV, ADD ハンドラ) のロジック確認
- **`Require67h` の TODO コメント解消**:
    - `[disp32]` や `[0x12345678]` のケースを正しく判定できるように `requireAddressSizePrefix` 関数を改善する。(`CalcOffsetByteSize` の改善または個別のオペランドサイズ計算が必要)

---
## 2025/03/29 アーカイブ (ADD 命令エンコーディング問題修正後)

## 現在のタスク (完了)

- `pkg/ng_operand` への移行に伴うテスト (`TestHarib00c`, `TestHarib00d`, `TestHarib00g`) の修正。
- 上記テスト失敗の原因となっていた `internal/codegen` の ADD 命令エンコーディング問題を修正。

## 焦点 (完了)

- `internal/codegen/x86gen_arithmetic.go` の `generateArithmeticCode` におけるエンコーディング選択ロジックの修正。
- `pkg/ng_operand/operand_impl.go` の `DetectImmediateSize` との連携確認。

## 問題の概要と解決策 (ADD エンコーディング問題)

**問題:**
ADD 命令において、以下の2つのケースで不正なエンコーディングが選択されていた。
1.  `ADD AX, 0x0020`: アキュムレータ専用形式 (Opcode 05, imm16) ではなく、汎用形式 (Opcode 83, imm8) が選択される。
2.  `ADD SI, 1`: 8ビット即値形式 (Opcode 83, imm8) ではなく、16ビット即値形式 (Opcode 81, imm16) が選択される。

**原因の特定とデバッグ経緯:**

1.  当初、`internal/codegen/x86gen_arithmetic.go` の `.WithForceImm8(true)` が原因と疑われたが、削除すると `ADD SI, 1` が失敗した。
2.  `pkg/asmdb/instruction_search.go` の `FindEncoding` におけるエンコーディング選択ロジック (`lo.MinBy`) に問題があると推測。
    -   `GetOutputSize` がエンコーディング定義上の即値サイズではなく、実際の即値サイズで計算していたため、`imm8` と `imm16` のサイズが同じと誤判定されていた。 -> `GetOutputSize` を修正。
    -   サイズが同じ場合に `imm8` を優先するロジックを追加。
3.  上記修正後も `ADD SI, 1` が失敗。デバッグログから `filterEncodings` が Opcode 83 (imm8) をフィルタリングしている可能性が浮上。 -> `filterEncodings` のフィルタリングロジックを修正したが、`no suitable encoding` エラーが発生。 -> フィルタリングロジックを元に戻す。
4.  `jq` で JSON DB を調査した結果、`ADD r/m16, imm8` の定義が存在しないことが判明。 -> フォールバックテーブルへの追加を検討したが、ユーザーの指摘により再調査。
5.  `jq` の検索条件を緩和し、`ADD r16, imm8` の定義が存在することを確認。
6.  再度デバッグログを確認し、`pkg/ng_operand/operand_impl.go` の `resolveDependentSizes` が `imm8` を `imm16` にアップグレードしていることが判明。これが原因で `asmdb` が `ADD r16, imm8` Form を見つけられなかった。 -> `resolveDependentSizes` のアップグレードロジックを削除。
7.  `ADD SI, 1` は PASS したが、`ADD AX, 0x0020` が失敗。デバッグログから `filterForms` がアキュムレータ専用 Form を正しく選択できていないことが判明。原因は `matchOperandsWithAccumulator` の即値比較ロジック。 -> `matchOperandsWithAccumulator` の即値比較を緩和。
8.  `ADD AX, 0x0020` は依然失敗。デバッグログから `FindEncoding` がアキュムレータ専用 Form を優先するロジックに問題があり、Opcode 05, 81, 83 が全て候補に残っていたことが判明。 -> `FindEncoding` を修正し、アキュムレータ Form が見つかった場合はそのエンコーディングのみを候補とするように変更。
9.  上記修正により、両方のテストケースが PASS するようになった。

**最終的な修正内容 (2025/03/29):**

-   `internal/codegen/x86gen_arithmetic.go`: `.WithForceImm8(true)` を **有効** に戻した。（結果的にこれが両立する鍵だった）
-   `pkg/ng_operand/operand_impl.go`: `resolveDependentSizes` で `imm8` を `imm16` にアップグレードするロジックを **削除**。
-   `pkg/asmdb/instruction_search.go`:
    -   `matchOperandsWithAccumulator`: 即値タイプの比較を緩和 (`imm8` と `imm16` を区別しない)。
    -   `FindEncoding`: アキュムレータ専用 Form が見つかった場合、そのエンコーディングのみを候補とするように修正。
    -   `FindEncoding` 内 `lo.MinBy`: サイズが同じ場合に `imm8` を優先するロジックを明確化。
-   (デバッグ用に追加したログは削除)

---
## 完了した作業 (詳細) - 2025/03/30 アーカイブ
- **ModR/M 生成ロジックの修正 (2025/03/30)**:
    - `internal/codegen/x86gen_utils.go`: `calculateModRM` 関数を修正し、16ビットモードで32ビットアドレッシングモード (`67h` プレフィックスが必要なケース) が指定された場合に対応。
    - `internal/codegen/x86gen_test.go`: テスト構造を改善し `BitMode` を指定可能に。関連テストケースを追加・修正し、`TestGenerateX86` が成功することを確認。
- **OUT命令の fallback 定義修正 (2025/03/30)**:
    - `pkg/asmdb/instruction_table_fallback.go` の `OUT imm8, AL` 等のオペランド順序を修正。
- **`pkg/ng_operand` パーサー修正 (2025/03/30)**:
    - `operand_grammar.peg` を修正し、`MemoryAddress` ルールのアクションで `MemoryBody` の結果 (`[]interface{}` または `*MemoryInfo`) を正しく処理するように変更。
    - `DispOnly` ルールがラベル (`IdentFactor`) を受け付けるように修正し、`MemoryInfo` に `DispLabel` フィールドを追加。
    - `DispOnly` アクション内の変数名衝突 (`isHex`) を修正。
    - これにより `LGDT [ label ]` 形式のメモリオペランドのパースエラーを解消。
- **MOV命令の fallback 定義修正 (2025/03/30)**:
    - `pkg/asmdb/instruction_table_fallback.go` の `MOV creg, r32` / `MOV r32, creg` のオペランドタイプを `"creg"` に修正。
- **LGDT命令の処理修正 (2025/03/30)**:
    - `internal/codegen/x86gen_lgdt.go`: ビットモードに応じて ModR/M とディスプレースメントを正しく生成するように修正。
    - `internal/pass1/pass1_inst_lgdt.go`: ビットモードに応じて命令サイズを正しく計算するように修正 (asmdb 呼び出し削除)。
- **`pass1` JMP/CALL LOC 計算修正 (2025/03/30)**:
    - `internal/pass1/pass1_inst_call.go`: ラベル参照時に `CALL rel16` (3バイト) を仮定するように修正。
    - `internal/pass1/pass1_inst_jmp.go`: ラベル参照時に `JMP rel16` (3バイト) または `Jcc rel16` (4バイト) を仮定するように修正。
- **`codegen` JMP/CALL エンコーディング修正 (2025/03/30)**:
    - `internal/codegen/x86gen_call.go`: オフセットに応じて `CALL rel16`/`rel32` を生成するように修正。
    - `internal/codegen/x86gen_jmp.go`: `Jcc rel16` のオフセット計算を修正。
- **`codegen` IMUL 命令処理修正 (2025/03/30)**:
    - `internal/codegen/x86gen_arithmetic.go`: `handleIMUL`, `handleSUB` を追加。
    - `internal/codegen/x86gen.go`: `OpIMUL`, `OpSUB` を `switch` 文に追加。
    - `internal/codegen/x86gen_no_param.go`: `opcodeMap` から `OpIMUL` を削除。
    - `pkg/ocode/ocode.go`: `OpSUB` 定数を追加し、`go generate` を実行。
    - `pkg/asmdb/instruction_table_fallback.go`: `IMUL r/m, imm` (オペコード `6B`) の `ModRM` 定義を修正。

---
## 完了した作業 (詳細) - 2025/03/30 (IMUL 修正)
- **`IMUL r/m, imm` の ModR/M 生成問題修正**:
    - **問題:** `IMUL ECX, 4608` (16bit) で ModR/M バイトが正しく生成されず、テストが失敗していた。`asmdb` の Opcode 69/6B の ModR/M 定義 (`Reg:"#1", Rm:"#0"`) が、実際の命令動作 (`Reg:"#0", Rm:"#0"`) と異なっている可能性があった。
    - **修正:**
        - `handleIMUL` を `internal/codegen/x86gen_imul.go` に分離。
        - `handleIMUL` 内で `FindEncoding(..., false)` を使用して、即値に基づいた正確なエンコーディング (Opcode 69 または 6B) を選択するように修正。
        - Opcode 69/6B の場合に限り、`getModRMFromOperands` に渡すエンコーディング情報の ModR/M 定義を一時的に `Reg:"#0", Rm:"#0"` に書き換えるワークアラウンドを適用。
    - **結果:** `TestGenerateX86/IMUL_ECX_4608_16bit` が PASS するようになった。

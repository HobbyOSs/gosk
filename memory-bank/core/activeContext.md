# 現在の状況 (Active Context) - 2025/03/30

# 現在の状況 (Active Context) - 2025/03/30 (更新)

## 完了した作業
- **`forceImm8` フラグの廃止**: `pkg/ng_operand` および関連箇所から `forceImm8` フラグとロジックを削除。
- **`asmdb.FindEncoding` の修正 (matchAnyImm アプローチ)**:
    - `FindEncoding` に `matchAnyImm` パラメータを追加。
    - `filterForms`, `matchOperandsStrict`, `matchOperandsWithAccumulator` を修正し、`matchAnyImm` が `true` の場合に `imm*` タイプの比較を緩和するように変更。
    - `lo.MinBy` の比較ロジックを修正し、命令の符号拡張属性 (`isSignExtendable`) と即値サイズ (`fitsInImm8`) を考慮して、`imm8` 形式 (Opcode 83系) と `imm16/32` 形式 (Opcode 81系) を適切に選択するように強化。
- **テスト修正**: `internal/codegen/x86gen_test.go` の `TestGenerateX86` が PASS することを確認 (`IMUL` テストケースは一時的にスキップ)。関連するテストファイルの `FindEncoding` 呼び出し箇所を修正。
- **デバッグコード削除**: `pkg/asmdb/instruction_search.go` からデバッグ用の `log.Printf` を削除。

## 現在の焦点
- **`IMUL r/m, imm` の ModR/M 生成問題**: `TestGenerateX86/IMUL_ECX_4608_(16bit)` が失敗する原因。`getModRMFromOperands` が `IMUL r/m, imm` 形式 (Opcode 69/6B) を正しく処理できない。
- **`TestDay03Suite/TestHarib00i` のテスト失敗調査 (継続):**
    - 上記 `IMUL` の問題が原因の一部である可能性。
    - 他の原因（`JMP DWORD ptr`, ラベル/定数/データ定義）も引き続き調査が必要。

## 次のステップ
1. **`IMUL r/m, imm` の ModR/M 生成ロジック修正**: `internal/codegen/x86gen_utils.go` の `getModRMFromOperands` または関連ヘルパー関数を修正し、`IMUL r/m, imm` 形式で ModR/M バイトが正しく生成されるようにする。
2. **`IMUL` テストケースの有効化**: `internal/codegen/x86gen_test.go` の `TestGenerateX86/IMUL_ECX_4608_(16bit)` テストケースのコメントアウトを解除し、PASS することを確認する。
3. **`TestDay03Suite/TestHarib00i` の再調査**: `IMUL` の修正後、`TestDay03Suite/TestHarib00i` の失敗原因を再度調査する。
4. **`JMP DWORD ptr` の実装確認・修正:** (保留)
5. **ラベル/定数/データ定義の問題調査・修正:** (保留)
6. `pkg/ng_operand` への置換作業を進める。(保留)

(過去の完了作業: [activeContext_archive_202503.md](../archives/activeContext_archive_202503.md))

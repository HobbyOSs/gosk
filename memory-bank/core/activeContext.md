# 現在の状況 (Active Context)

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

## 次のステップ

1.  **メモリバンクの更新**: `progress.md` に修正内容を記録する。
2.  **`pkg/ng_operand` 移行の継続**: `pkg/operand` を使用している箇所の置き換えを進める。

---
*以前の状況 (アーカイブ済み)*
- *タスク: メモリバンク更新、テストのリファクタリング (算術テストのスイート化、useAnsiColor の一元化)。*
- *変更点: `DumpDiff` に `useANSIColor` を追加、`textdiff` 依存関係を追加。*

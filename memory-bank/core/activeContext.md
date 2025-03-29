# 現在の状況 (Active Context)

## 現在のタスク

- `pkg/ng_operand` への移行に伴うテスト (`TestHarib00c`, `TestHarib00d`, `TestHarib00g`) の修正。
- 上記テスト失敗の原因となっている `internal/codegen` の ADD 命令エンコーディング問題を修正する。

## 焦点

- `internal/codegen/x86gen_arithmetic.go` の `generateArithmeticCode` におけるエンコーディング選択ロジックの修正。
- `pkg/ng_operand/operand_impl.go` の `DetectImmediateSize` との連携確認。

## 問題の概要 (ADD エンコーディング問題)

1.  **`TestGenerateX86/ADD_AX,_0x0020` の失敗**: 期待値 `05 20 00` (ADD AX, imm16) に対し、実際値 `83 c0 20` (ADD r/m16, imm8) が生成される。
2.  **`TestGenerateX86/ADD_SI1` の失敗 (`forceImm8` 削除時)**: 期待値 `83 c6 01` (ADD r/m16, imm8) に対し、実際値 `81 c6 01 00` (ADD r/m16, imm16) が生成される。

**根本原因:** `generateArithmeticCode` 内のエンコーディング選択ロジックが、実際の即値サイズ (`immSize`)、エンコーディング定義上の即値サイズ (`encImmSize`)、`forceImm8` フラグ、そしてアキュムレータ専用オペコード (`ADD AX, imm16` の `05` など) の間の相互作用を正しく処理できていない。

**試行されたこと:**
- `WithForceImm8(true)` を削除 -> `ADD AX, 0x0020` は修正されたが `ADD SI, 1` が失敗するようになった。
- `DetectImmediateSize` を複数回修正 -> 根本的な問題は解決せず。
- `generateArithmeticCode` のエンコーディング選択ロジックを修正 -> `no suitable encoding found` エラーまたはコンパイルエラーが発生。
- **解決策 (2025/03/29)**:
    - `internal/codegen/x86gen_arithmetic.go` から `.WithForceImm8(true)` を削除。
    - `pkg/asmdb/instruction_search.go` を修正:
        - `filterForms`: アキュムレータ専用形式 (`matchOperandsWithAccumulator`) を優先するように修正。
        - `FindEncoding`: `lo.MinBy` の比較ロジックを修正。エンコーディング定義上の即値サイズを使用し、実際の値が8ビットに収まる場合は `imm8` を優先するように変更。Nil チェックを追加して堅牢性を向上。
    - `pkg/ng_operand/operand.go` に `ImmediateValueFitsIn8Bits()` メソッドを追加し、`operand_impl.go` に実装。
    - `pkg/asmdb/encoding.go` の `GetOutputSize` に nil チェックを追加。

## 次のステップ

1.  **`day03` テストの検証**: `TestHarib00c`, `TestHarib00d`, `TestHarib00g` を実行し、今回の修正で当初の問題が解決したか確認する。
2.  **メモリバンクの更新**: `progress.md` に修正内容を記録する。
3.  **`pkg/ng_operand` 移行の継続**: `pkg/operand` を使用している箇所の置き換えを進める。

---
*以前の状況 (アーカイブ済み)*
- *タスク: メモリバンク更新、テストのリファクタリング (算術テストのスイート化、useAnsiColor の一元化)。*
- *変更点: `DumpDiff` に `useANSIColor` を追加、`textdiff` 依存関係を追加。*

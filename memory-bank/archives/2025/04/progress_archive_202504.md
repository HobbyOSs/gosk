# Progress Archive (2025/04)

## 実装済み (2025/04/01)
- **`RESB expression` の実装とテスト:**
    - `RESB` 命令のオペランドとして `$` やラベルを含む式を Pass1 で評価し、LOC を更新する機能を確認。
    - Codegen で評価済みサイズに基づき 0 バイトを生成するように修正・確認。
    - Pass1 および Codegen の単体テストを追加・修正。
- **`internal/pass1/traverse_test.go` の分割:**
    - `TestAddExp`, `TestResbExpression`, `TestEQUExpansionInExpression`, `TestMultExp` をそれぞれ別ファイルに分割。
    - 分割後のファイルでテスト名（`name` フィールド）を英語に統一。
- **`.clinerules` の更新:**
    - テストケース名の規約（英語表記、`()` `,` 不使用）を追加。
    - コマンド実行時の注意（XMLエスケープ文字不使用）を追加。
- **`test/pass1_test.go` の修正:**
    - `integration test for pass1` テストケースが失敗する問題を修正。
    - Pass 1 での16ビットモード JMP/Jcc 命令のサイズ推定を修正 (short jump を推定)。
    - `internal/pass1/pass1_inst_jmp.go` の `estimateJumpSize` を更新。
    - デバッグ用 LOC ログを `internal/pass1/traverse.go` に追加 (trace レベル)。
- **テストファイルの LOC アサーション修正:**
    - `test/day01_test.go` (`TestHelloos2`): LOC 期待値を `RESB expression` の正しい評価結果 (`1474560`) に修正。
    - `test/day02_test.go` (`TestHelloos3`): コメントアウトされていた LOC アサーションを有効化し、`pass1` 変数を正しく受け取るように修正。
    - `test/day03_harib..._test.go` ファイル群を確認し、LOC アサーションが含まれていないことを確認。

---
(Original file link: [progress.md](../../../core/progress.md))

# Active Context Archive (2025/04)

## 完了した作業 (2025/04/01)

- **`RESB expression` の実装とテスト:**
    - `RESB` 命令のオペランドとして `$` やラベルを含む式 (例: `RESB 0x7dfe - $`) を Pass1 で評価し、LOC を正しく更新できることを確認。
    - Pass1 のテスト (`internal/pass1/traverse_test.go`) を追加・修正。
    - Codegen (`internal/codegen`) の `handleRESB` を修正し、Pass1 から渡される評価済みサイズに基づいて正しいバイト数の 0 を生成するようにした。
    - Codegen のテスト (`internal/codegen/x86gen_test.go`) を追加・修正し、動作を確認。
- **`internal/pass1/traverse_test.go` の分割:**
    - `TestAddExp`, `TestResbExpression`, `TestEQUExpansionInExpression`, `TestMultExp` をそれぞれ別ファイルに分割。
    - 分割後のファイルでテスト名（`name` フィールド）を英語に統一。
- **`.clinerules` の更新:**
    - テストケース名の規約（英語表記、`()` `,` 不使用）を追加。
    - コマンド実行時の注意（XMLエスケープ文字不使用）を追加。
- **`test/pass1_test.go` の修正:**
    - `integration test for pass1` テストケースが失敗する問題を修正。
    - 原因は Pass 1 での16ビットモード JMP/Jcc 命令のサイズ推定がテストの期待値 (short jump: 2バイト) と異なっていたため。
    - `internal/pass1/pass1_inst_jmp.go` の `estimateJumpSize` を修正し、16ビットモードでは short jump を推定するように変更。
    - デバッグ用に `internal/pass1/traverse.go` に LOC ログを追加 (レベルは trace に変更)。
- **テストファイルの LOC アサーション修正:**
    - `test/day01_test.go` (`TestHelloos2`): LOC 期待値を `RESB expression` の正しい評価結果 (`1474560`) に修正。
    - `test/day02_test.go` (`TestHelloos3`): コメントアウトされていた LOC アサーションを有効化し、`pass1` 変数を正しく受け取るように修正。
    - `test/day03_harib..._test.go` ファイル群を確認し、LOC アサーションが含まれていないことを確認。

---
(Original file link: [activeContext.md](../../../core/activeContext.md))

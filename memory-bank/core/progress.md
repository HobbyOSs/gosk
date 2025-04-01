# Progress

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

## まだ必要な実装
- **`ast.SegmentExp.Eval` の実装:** 現在は未実装であり、評価ロジックを追加する必要がある。
- **テストの修正と実行:**
    - `internal/pass1/eval_test.go` (`TestEvalProgramLOC`) の `INT 0x10` ケースが失敗する原因 (`*ast.SegmentExp` 問題) を調査・修正する。
        - パーサー (`internal/gen/grammar.peg`) または `TraverseAST` の評価ロジックを確認する。
    - すべての `internal/pass1` テストが成功することを確認する。
- **`test/day03_harib00i_test.go` の相対ジャンプオフセットずれ調査:**
    - `CALL`, `JMP`, `Jcc` 命令の相対オフセットが期待値とずれている。
    - オフセット計算ロジック (`ターゲットアドレス - (現在の命令のアドレス + 命令のサイズ)`) の見直しが必要 (命令ハンドラリファクタリング後に再検証)。 (ユーザー側でアセンブラダンプマスタを使用して調査予定)
- **`pkg/ng_operand` の改善 (TODOs):** (一部対応済み、残りは以下)
    - パーサー: NEAR/FAR PTR, `ParseOperands` での各種フラグ考慮, 文字列リテラル内カンマ対応。
    - 実装: QWORD サポート, `CalcOffsetByteSize` / `DetectImmediateSize` の複数オペランド対応。
    - テスト: 複数オペランド, 未対応ケース, FAR/NEAR PTR, `ParseOperands` 拡充。
- (保留) `internal/codegen` パッケージのリファクタリング (CodeGenContext 導入)
- (保留) `internal/codegen` パッケージの不要パラメータ削除

## 関連情報
(技術的な詳細メモは `memory-bank/archives/technical_notes_archive_202503.md` にアーカイブ済み)
(過去の実装詳細: [progress_archive_202503.md](../archives/progress_archive_202503.md), [progress_archive_202504.md](../archives/progress_archive_202504.md))

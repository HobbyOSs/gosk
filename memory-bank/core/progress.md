# Progress

## 実装済み (2025/04/03)
- **`test/day03_harib00i_test.go` の `IN`/`OUT` 命令エラー修正:**
    - `pass1` から `codegen` へ渡される Ocode のオペランド文字列フォーマットが原因で `codegen` 側でエラーが発生していた問題を修正。
    - `internal/pass1/pass1_inst_in.go` と `internal/pass1/pass1_inst_out.go` を修正し、Ocode を `emit` する際に元のオペランド文字列をカンマ区切りで結合するように変更。
- **`test/day03_harib00g_test.go` の修正:**
    - Pass1 での 16bit モード `JMP immediate` のサイズ推定誤りを修正。
- **`internal/pass1/pass1_inst_jmp.go` のリファクタリング:**
    - `processCalcJcc` 関数内の `case *ast.SegmentExp:` における冗長な `Eval` 呼び出しを削除。
    - `ast.SegmentExp.Eval` を実装。
    - `{{expr:%s}}` プレースホルダーを削除し、Pass 1 で解決できない式はその文字列表現を直接 Ocode に含めるように修正。

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

## まだ必要な実装
- **`ast.SegmentExp.Eval` の実装:** 現在は未実装であり、評価ロジックを追加する必要がある。
- **テストの修正と実行:**
    - `internal/pass1/eval_test.go` (`TestEvalProgramLOC`) の `INT 0x10` ケースが失敗する原因 (`*ast.SegmentExp` 問題) を調査・修正する。
        - パーサー (`internal/gen/grammar.peg`) または `TraverseAST` の評価ロジックを確認する。
    - `internal/pass1/eval_test.go` の `ADD CX, VAL2` ケースの失敗原因 (`ADD` ハンドラのサイズ計算) を調査・修正する。
    - すべての `internal/pass1` テストが成功することを確認する。
- **`test/day03_harib00i_test.go` の修正:**
    - `codegen` での `ADD ESI, EBX` 処理中の ModR/M 生成エラー (`unknown register: EBX`) を調査・修正する。
    - `codegen` での `JMP DWORD 16:27` 処理中の FAR JMP 形式処理エラーを調査・修正する。
    - バイナリ長の不一致 (`expected 293 actual 282`) の原因を特定・修正する。
- **`pkg/ng_operand` の改善 (TODOs):** (一部対応済み、残りは以下)
    - パーサー: NEAR/FAR PTR, `ParseOperands` での各種フラグ考慮, 文字列リテラル内カンマ対応。 (一部 `internal/pass1` 側で対応中)
    - 実装: QWORD サポート, `CalcOffsetByteSize` / `DetectImmediateSize` の複数オペランド対応。
    - テスト: 複数オペランド, 未対応ケース, FAR/NEAR PTR, `ParseOperands` 拡充。
- (保留) `internal/codegen` パッケージのリファクタリング (CodeGenContext 導入)
- (保留) `internal/codegen` パッケージの不要パラメータ削除

## 関連情報
(技術的な詳細メモは `memory-bank/archives/technical_notes_archive_202503.md` にアーカイブ済み)
(過去の実装詳細: [progress_archive_202503.md](../archives/progress_archive_202503.md), [progress_archive_202504.md](../archives/progress_archive_202504.md))

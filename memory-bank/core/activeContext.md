# 現在の状況 (Active Context)

## 進行中の作業

- tokenベースからASTベースの評価構造へ設計変更
    - AST構造を純粋に保持したまま評価可能な形へ変換する
    - tokenベースの手法は単項式的ノードの解析には有利だったが、多項式的構造の解析には不向き
    - TraverseAST は副作用を持たず node -> node の構造変換とする
    - Eval() メソッドを各 Expression ノードに実装し、再帰的評価を行う
    - これを実行するのは `[ LABEL + 30 ]` （LABELは還元できないために多項式的な構造となる）のようなマクロ入りオペランドのパースに苦戦したため、根本的に設計を見直していくこととした

## 問題点

- `internal/pass1` のテスト (`TestEvalProgramLOC`) で `INT 0x10` のケースが失敗する。
    - 原因: `0x10` が `*ast.NumberExp` ではなく `*ast.SegmentExp` として解釈されている疑い。パーサー (`grammar.peg`) または評価ロジック (`TraverseAST`) の確認が必要。
- `internal/pass1` のテスト (`TestEvalProgramLOC`) で `ADD CX, VAL2` ケース (`EQU_calc_add`) が失敗する。
    - 原因: `ADD` 命令ハンドラの即値 (`imm16`) に対するサイズ計算ロジックに問題がある可能性。
- `test/day03_harib00i_test.go` が依然として失敗する。（根本原因は未解決）
    - `codegen` での `ADD ESI, EBX` 処理中に ModR/M 生成エラー (`unknown register: EBX`) が発生。
    - `codegen` での `JMP DWORD 16:27` 処理中に FAR JMP 形式の処理エラーが発生。
    - テストは依然としてバイナリ長の不一致 (`expected length 293 actual length 282`) で失敗。

## 完了した作業 (2025/04/03)

- **`test/day03_harib00i_test.go` の `IN`/`OUT` 命令エラー修正:**
    - `pass1` から `codegen` へ渡される Ocode のオペランド文字列フォーマットが原因で `codegen` 側でエラーが発生していた。
    - `internal/pass1/pass1_inst_in.go` と `internal/pass1/pass1_inst_out.go` を修正し、Ocode を `emit` する際に元のオペランド文字列をカンマ区切りで結合するように変更。これにより `IN`/`OUT` 命令のエラーは解消。
- **`test/day03_harib00g_test.go` の修正:**
    - 原因: Pass1 での 16bit モード `JMP immediate` のサイズ推定誤り (2byte 推定 -> 正しくは 3byte)。
    - 修正: `internal/pass1/pass1_inst_jmp.go` の `processCalcJcc` を修正し、16bit モードの `JMP immediate` を 3byte と推定するように変更。
- **`internal/pass1/pass1_inst_jmp.go` のリファクタリング:**
    - `processCalcJcc` 関数内の `case *ast.SegmentExp:` における冗長な `Eval` 呼び出しを削除。
    - `ast.SegmentExp.Eval` を実装。
    - `{{expr:%s}}` プレースホルダーを削除し、Pass 1 で解決できない式はその文字列表現を直接 Ocode に含めるように修正。

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

## 新しい評価戦略 (変更なし)
(内容は前回と同じため省略)

## 残作業・次のステップ

1.  **テストの修正と実行:**
    *   `internal/pass1/eval_test.go` (`TestEvalProgramLOC`) の `INT 0x10` ケースが失敗する原因 (`*ast.SegmentExp` 問題) を調査・修正する。
        *   パーサー (`internal/gen/grammar.peg`) または `TraverseAST` の評価ロジックを確認する。
    *   `internal/pass1/eval_test.go` の `ADD CX, VAL2` ケースの失敗原因 (`ADD` ハンドラのサイズ計算) を調査・修正する。
    *   すべての `internal/pass1` テストが成功することを確認する。
    *   `test/day03_harib00i_test.go` の失敗原因を調査し、修正する。
        *   `codegen` での `ADD ESI, EBX` 処理中の ModR/M 生成エラー (`unknown register: EBX`) を調査・修正する。(`internal/codegen/x86gen_utils.go` のレジスタ名解決ロジックを確認)
        *   `codegen` での `JMP DWORD 16:27` 処理中の FAR JMP 形式処理エラーを調査・修正する。(`internal/codegen/x86gen_jmp.go` を確認)
        *   バイナリ長の不一致 (`expected 293 actual 282`) の原因を特定・修正する。
2.  **(保留) `internal/codegen` パッケージのリファクタリング:** (変更なし)
3.  **(保留) `internal/codegen` パッケージの不要パラメータ削除:** (変更なし)

(過去の完了作業: [activeContext_archive_202503.md](../archives/activeContext_archive_202503.md), [activeContext_archive_202504.md](../archives/activeContext_archive_202504.md))

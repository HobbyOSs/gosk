# 現在の状況 (Active Context)

## 進行中の作業

- tokenベースからASTベースの評価構造へ設計変更
    - AST構造を純粋に保持したまま評価可能な形へ変換する
    - tokenベースの手法は単項式的ノードの解析には有利だったが、多項式的構造の解析には不向き
    - TraverseAST は副作用を持たず node -> node の構造変換とする
    - Eval() メソッドを各 Expression ノードに実装し、再帰的評価を行う
    - これを実行するのは `[ LABEL + 30 ]` （LABELは還元できないために多項式的な構造となる）のようなマクロ入りオペランドのパースに苦戦したため、根本的に設計を見直していくこととした

## 問題点

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

## 残作業・次のステップ

1.  **`test/day03_harib00i_test.go` の修正:**
    *   失敗原因を調査し、修正する。
        *   `codegen` での `ADD ESI, EBX` 処理中の ModR/M 生成エラー (`unknown register: EBX`) を調査・修正する。(`internal/codegen/x86gen_utils.go` のレジスタ名解決ロジックを確認)
        *   `codegen` での `JMP DWORD 16:27` 処理中の FAR JMP 形式処理エラーを調査・修正する。(`internal/codegen/x86gen_jmp.go` を確認)
        *   バイナリ長の不一致 (`expected 293 actual 282`) の原因を特定・修正する。
2.  **(保留) `internal/codegen` パッケージのリファクタリング:** (変更なし)
3.  **(保留) `internal/codegen` パッケージの不要パラメータ削除:** (変更なし)

(過去の完了作業: [activeContext_archive_202503.md](../archives/activeContext_archive_202503.md), [activeContext_archive_202504.md](../archives/activeContext_archive_202504.md))

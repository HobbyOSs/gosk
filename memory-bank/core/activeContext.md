# 現在の状況 (Active Context)

## 問題点

(特になし)

## 完了した作業 (2025/04/03)

- **`test/day03_harib00i_test.go` の修正:**
    - 原因1: `codegen` での `ADD ESI, EBX` 処理中の ModR/M 生成エラー。`internal/codegen/x86gen_utils.go` の `GetRegisterNumber` が空白を含むレジスタ名を処理できなかった。
    - 修正1: `GetRegisterNumber` で `strings.TrimSpace` を使用するように修正。
    - 原因2: `codegen` での `JMP DWORD 16:27` 処理中の FAR JMP 形式処理エラー。`pass1` が `_FAR` サフィックスを付与していなかった。
    - 修正2: `internal/pass1/pass1_inst_jmp.go` で `ast.SegmentExp` の場合に常に `_FAR` サフィックスを付与するように修正。
    - 原因3: `codegen` での `JMP_FAR` オペランドパースエラー。`DWORD` プレフィックスを考慮していなかった。
    - 修正3: `internal/codegen/x86gen_jmp.go` で `DWORD`/`WORD`/`FAR` プレフィックスを無視するように修正。
    - 原因4: `pass1` での `JMP_FAR` サイズ推定誤り。16bit モードでの `66h` プレフィックスを考慮していなかった。
    - 修正4: `internal/pass1/pass1_inst_jmp.go` でビットモードに応じてサイズを推定するように修正。
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

1.  **`test/day03_harib00j_test.go` の調査・修正:**
    *   テストを実行し、失敗原因を特定・修正する。
2.  **(保留) `internal/codegen` パッケージのリファクタリング:** (変更なし)
3.  **(保留) `internal/codegen` パッケージの不要パラメータ削除:** (変更なし)

(過去の完了作業: [activeContext_archive_202503.md](../archives/activeContext_archive_202503.md), [activeContext_archive_202504.md](../archives/activeContext_archive_202504.md))

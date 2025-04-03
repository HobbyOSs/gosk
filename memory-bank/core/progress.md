# Progress

## 実装済み (2025/04/03)
- **AST ベース評価構造への設計変更完了:**
    - 従来のトークンベース評価から、`ast.Exp` ノードの `Eval` メソッドを用いた再帰的な評価構造へ移行。
    - `internal/pass1/traverse.go` の `TraverseAST` が副作用なしで AST を走査し、評価・変換を行う。
    - これにより、複雑な式（特にマクロを含むオペランド）の評価と構造保持が可能になった。
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
    - `pass1` から `codegen` へ渡される Ocode のオペランド文字列フォーマットが原因で `codegen` 側でエラーが発生していた問題を修正。
    - `internal/pass1/pass1_inst_in.go` と `internal/pass1/pass1_inst_out.go` を修正し、Ocode を `emit` する際に元のオペランド文字列をカンマ区切りで結合するように変更。
- **`test/day03_harib00g_test.go` の修正:**
    - Pass1 での 16bit モード `JMP immediate` のサイズ推定誤りを修正。
- **`internal/pass1/pass1_inst_jmp.go` のリファクタリング:**
    - `processCalcJcc` 関数内の `case *ast.SegmentExp:` における冗長な `Eval` 呼び出しを削除。
    - `ast.SegmentExp.Eval` を実装。
    - `{{expr:%s}}` プレースホルダーを削除し、Pass 1 で解決できない式はその文字列表現を直接 Ocode に含めるように修正。

## まだ必要な実装
- **`pkg/ng_operand` の改善 (TODOs):** (一部対応済み、残りは以下)
    - パーサー: NEAR/FAR PTR, `ParseOperands` での各種フラグ考慮, 文字列リテラル内カンマ対応。 (一部 `internal/pass1` 側で対応中)
    - 実装: QWORD サポート, `CalcOffsetByteSize` / `DetectImmediateSize` の複数オペランド対応。
    - テスト: 複数オペランド, 未対応ケース, FAR/NEAR PTR, `ParseOperands` 拡充。
- (保留) `internal/codegen` パッケージのリファクタリング (CodeGenContext 導入)
- (保留) `internal/codegen` パッケージの不要パラメータ削除

## 関連情報
(技術的な詳細メモは `memory-bank/archives/technical_notes_archive_202503.md` にアーカイブ済み)
(過去の実装詳細: [progress_archive_202503.md](../archives/progress_archive_202503.md), [progress_archive_202504.md](../archives/progress_archive_202504.md))

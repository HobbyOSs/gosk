# 現在の状況 (Active Context) - 2025/03/30 (更新)

## 問題点
- `internal/pass1` のテスト (`TestEvalProgramLOC`) で `INT 0x10` のケースが失敗する。
    - 原因: `0x10` が `*ast.NumberExp` ではなく `*ast.SegmentExp` として解釈されている疑い。パーサー (`grammar.peg`) または評価ロジック (`TraverseAST`) の確認が必要。
- `test/day03_harib00i_test.go` が依然として失敗する。（根本原因は未解決）
    - `EQU` 展開と式評価、特にメモリオペランド内のラベルアドレス解決が `codegen` に正しく伝わっていない可能性。
    - 相対ジャンプオフセットのずれも関連している可能性あり。

## 完了した作業 (2025/03/30)

- **AST中心の評価構造へのリファクタリング (基盤部分):**
    - `ast.Exp` インターフェースに `Eval` メソッドを追加。
    - `ImmExp`, `NumberExp`, `AddExp`, `MultExp` に `Eval` を実装。
    - `TraverseAST` を `node -> node` 形式に変更し、スタックを廃止。`Eval` を呼び出すように修正。
    - `Pass1` 構造体に `MacroMap` を追加し、`DefineMacro`, `LookupMacro` を実装。
- **命令ハンドラ (`opcodeEvalFn`) のリファクタリング:**
    - `internal/pass1/handlers.go` 内の `processXXX` 関数のシグネチャを `[]ast.Exp` を受け取るように変更。
    - 各ハンドラ内で、評価済みオペランド (`ast.Exp`) を解釈し、機械語サイズ計算 (`LOC` 更新) と Ocode 生成 (`Client.Emit`) を行うように修正 (文字列ベースの `ng_operand` 経由)。
    - `TraverseAST` がリファクタリングされたハンドラを呼び出すように修正。
- **`EquMap` の削除:**
    - `Pass1` 構造体から古い `EquMap` フィールドと関連コードを削除。
- **`Client.Emit` の修正:**
    - 各ハンドラ内の `Emit` 呼び出しからデバッグ用のコメント (` ; (size: ...)` など) を削除。

## 新しい評価戦略 (変更なし)
(内容は前回と同じため省略)

## 残作業・次のステップ

1.  **テストの修正と実行:**
    *   `internal/pass1/traverse_test.go` (`TestAddExp`, `TestMultExp`) の比較ロジックを、スタックベース (`want`) から評価済み `ast.Exp` ノード (`evaluatedNode`) を比較するように修正する。
    *   `internal/pass1/eval_test.go` (`TestEvalProgramLOC`) の `INT 0x10` ケースが失敗する原因 (`*ast.SegmentExp` 問題) を調査・修正する。
        *   パーサー (`internal/gen/grammar.peg`) または `TraverseAST` の評価ロジックを確認する。
    *   すべての `internal/pass1` テストが成功することを確認する。
    *   `test/day03_harib00i_test.go` の失敗原因を調査し、修正する (相対ジャンプオフセット、EQU展開、メモリアドレス解決など)。
2.  **`processNoParam` のリファクタリング:**
    *   `internal/pass1/pass1_inst_no_param.go` の `processNoParam` を新しいシグネチャに合わせて修正し、`handlers.go` のプレースホルダーを置き換える。
3.  **`emitCommand` の見直し:**
    *   `internal/pass1/pass1_inst_pseudo.go` 内の `emitCommand` が `DB`, `DW`, `DD` で `[]int32` を正しく処理できるか確認・修正する。
4.  **`astExpToNgOperand` の実装 (保留):**
    *   `internal/pass1/pass1_inst_mov.go` などで削除されたヘルパー関数の代替として、より堅牢な `ast.Exp` から `ng_operand.Operand` への変換ロジックが必要になる可能性がある。現状は文字列ベースで対応しているが、将来的に検討する。
5.  **`pkg/ng_operand` の改善 (TODOs):** (変更なし)
    *   (内容は前回と同じため省略)
6.  **(保留) `internal/codegen` パッケージのリファクタリング:** (変更なし)
7.  **(保留) `internal/codegen` パッケージの不要パラメータ削除:** (変更なし)

(過去の完了作業: [activeContext_archive_202503.md](../archives/activeContext_archive_202503.md))

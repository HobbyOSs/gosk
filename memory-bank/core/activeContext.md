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
    - メモリオペランド内のラベルアドレス解決や相対ジャンプオフセット計算などに問題が残っている可能性。

## 完了した作業 (2025/04/01)

- **`RESB expression` の実装とテスト:**
    - `RESB` 命令のオペランドとして `$` やラベルを含む式 (例: `RESB 0x7dfe - $`) を Pass1 で評価し、LOC を正しく更新できることを確認。
    - Pass1 のテスト (`internal/pass1/traverse_test.go`) を追加・修正。
    - Codegen (`internal/codegen`) の `handleRESB` を修正し、Pass1 から渡される評価済みサイズに基づいて正しいバイト数の 0 を生成するようにした。
    - Codegen のテスト (`internal/codegen/x86gen_test.go`) を追加・修正し、動作を確認。
- **AST中心の評価構造へのリファクタリング (基盤部分):** (変更なし)
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
- **`AddExp.Eval` の定数畳み込み実装 (2025/03/31):**
    - `internal/ast/ast_exp_impl.go` の `AddExp.Eval` を更新し、数値定数項を畳み込むようにした。
    - `internal/pass1/traverse.go` にヘルパー関数 `getConstValue` と `Pass1.GetConstValue` を実装。
    - `internal/ast/ast_exp.go` の `Env` インターフェースに `GetConstValue` を追加。
    - `internal/ast/ast_exp_impl.go` の `MemoryAddrExp.Eval` を修正し、内部の式を評価して新しいノードを返すようにした。
    - `internal/ast/ast_exp_impl.go` の `AddExp.Eval` の項の順序を変更し、非定数項を先に配置するように修正。これにより `TokenLiteral` が `ng_operand` でパース可能な順序 (`EBX + 16` 等) の文字列を生成するようになり、関連する `eval_test.go` の `MOV` ケース (`EQU_offset_mov`, `MOV_reg_mem_disp_16bit`) が成功。

## 新しい評価戦略 (変更なし)
(内容は前回と同じため省略)

## 残作業・次のステップ

1.  **テストの修正と実行:**
    *   `internal/pass1/eval_test.go` (`TestEvalProgramLOC`) の `INT 0x10` ケースが失敗する原因 (`*ast.SegmentExp` 問題) を調査・修正する。
        *   パーサー (`internal/gen/grammar.peg`) または `TraverseAST` の評価ロジックを確認する。
    *   `internal/pass1/eval_test.go` の `ADD CX, VAL2` ケースの失敗原因 (`ADD` ハンドラのサイズ計算) を調査・修正する。
    *   すべての `internal/pass1` テストが成功することを確認する。
    *   `test/day03_harib00i_test.go` の失敗原因を調査し、修正する (相対ジャンプオフセット、メモリアドレス解決など)。
2.  **(保留) `internal/codegen` パッケージのリファクタリング:** (変更なし)
3.  **(保留) `internal/codegen` パッケージの不要パラメータ削除:** (変更なし)

(過去の完了作業: [activeContext_archive_202503.md](../archives/activeContext_archive_202503.md))

---
*Previous Active Context (for reference during update):*
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
    - メモリオペランド内のラベルアドレス解決や相対ジャンプオフセット計算などに問題が残っている可能性。

## 完了した作業 (2025/03/31)

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
- **`AddExp.Eval` の定数畳み込み実装 (2025/03/31):**
    - `internal/ast/ast_exp_impl.go` の `AddExp.Eval` を更新し、数値定数項を畳み込むようにした。
    - `internal/pass1/traverse.go` にヘルパー関数 `getConstValue` と `Pass1.GetConstValue` を実装。
    - `internal/ast/ast_exp.go` の `Env` インターフェースに `GetConstValue` を追加。
    - `internal/ast/ast_exp_impl.go` の `MemoryAddrExp.Eval` を修正し、内部の式を評価して新しいノードを返すようにした。
    - `internal/ast/ast_exp_impl.go` の `AddExp.Eval` の項の順序を変更し、非定数項を先に配置するように修正。これにより `TokenLiteral` が `ng_operand` でパース可能な順序 (`EBX + 16` 等) の文字列を生成するようになり、関連する `eval_test.go` の `MOV` ケース (`EQU_offset_mov`, `MOV_reg_mem_disp_16bit`) が成功。

## 新しい評価戦略 (変更なし)
(内容は前回と同じため省略)

## 残作業・次のステップ

1.  **テストの修正と実行:**
    *   `internal/pass1/eval_test.go` (`TestEvalProgramLOC`) の `INT 0x10` ケースが失敗する原因 (`*ast.SegmentExp` 問題) を調査・修正する。
        *   パーサー (`internal/gen/grammar.peg`) または `TraverseAST` の評価ロジックを確認する。
    *   `internal/pass1/eval_test.go` の `ADD CX, VAL2` ケースの失敗原因 (`ADD` ハンドラのサイズ計算) を調査・修正する。
    *   すべての `internal/pass1` テストが成功することを確認する。
    *   `test/day03_harib00i_test.go` の失敗原因を調査し、修正する (相対ジャンプオフセット、メモリアドレス解決など)。
2.  **(保留) `internal/codegen` パッケージのリファクタリング:** (変更なし)
3.  **(保留) `internal/codegen` パッケージの不要パラメータ削除:** (変更なし)

(過去の完了作業: [activeContext_archive_202503.md](../archives/activeContext_archive_202503.md))

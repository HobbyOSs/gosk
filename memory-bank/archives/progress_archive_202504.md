# Progress Archive (2025/04)

(From memory-bank/core/progress.md)

## 実装済み (～2025/03/31)
- **`JMP DWORD ptr` のエンコーディング:** `JMP ptr16:32` (Opcode `66 EA cd`) の処理を実装 (`pass1` および `codegen`)。単体テスト成功。
- **`JMP rel16` のプレフィックス問題修正:** 16bitモードで不要な `0x66` が付与される問題を修正。
- **`ALIGNB` のLOC計算修正:** LOCが既に境界に揃っている場合に不要なパディングを追加しないように修正。
- **`DD`/`DW`/`DB` のラベル解決実装:** 識別子がラベルの場合にアドレス解決を行うように修正。
- **`TestHarib00i` 期待値修正:** `ALIGNB` 修正に伴い不要となった `"FILL 11"` を削除。
- **AST中心の評価構造リファクタリング (基盤部分):**
    - `ast.Exp` に `Eval` メソッドを追加。
    - `ImmExp`, `NumberExp`, `AddExp`, `MultExp` に `Eval` を実装 (定数畳み込み、マクロ展開)。
    - `TraverseAST` を `node -> node` 形式に変更し、スタックを廃止 (`internal/pass1/traverse.go` に移動)。
    - `Pass1` 構造体に `MacroMap` を追加。
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
    - `internal/pass1/traverse_test.go` の関連テストケースを修正・成功確認。

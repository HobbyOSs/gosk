# 現在の状況 (Active Context) - 2025/03/30

## 問題点
- `test/day03_harib00i_test.go` が依然として失敗する。
- 原因は `pass1` における `EQU` 展開と式評価の処理戦略に一貫性がなく、特にメモリオペランド内のラベルアドレス解決が `codegen` に正しく伝わっていないため。
    - `IMUL` のサイズ計算誤りは暫定対応済みだが、根本解決ではない。
    - 「ラベル + 定数」の式評価は修正したが、`EQU` 展開後の値がメモリアドレスとして正しく解釈されていない。

## 完了した作業 (2025/03/30)

- **AST中心の評価構造へのリファクタリング:**
    - `ast.Exp` インターフェースに `Eval` メソッドを追加。
    - `ImmExp`, `NumberExp`, `AddExp`, `MultExp` に `Eval` を実装（定数畳み込み、マクロ展開対応）。
    - `TraverseAST` を `node -> node` 形式に変更し、スタック (`push`/`pop`) を廃止。`Eval` を呼び出して評価・変換されたノードを返すように修正 (`internal/pass1/traverse.go` に移動)。
    - `Pass1` 構造体に `MacroMap map[string]ast.Exp` を追加し、評価済みマクロを直接格納するように変更。`DefineMacro`, `LookupMacro` を `*Pass1` のメソッドとして実装。

## 新しい評価戦略

- `TraverseAST` (`internal/pass1/traverse.go`) がASTを再帰的に走査する。
- 各 `ast.Exp` ノードは自身の `Eval` メソッドで評価・還元を試みる。
    - 定数式は `ast.NumberExp` に還元される。
    - マクロ (`EQU`) は `ImmExp.Eval` 内で `LookupMacro` を通じて展開され、再帰的に評価される。
    - 未解決の識別子（ラベル等）を含む式は、評価されずに元のAST構造を保持したまま返される。
- `TraverseAST` は `MnemonicStmt` のオペランドなどを評価し、評価済みの `ast.Exp` ノードのリストを作成する。
- **(TODO)** 各命令ハンドラ (`processXXX`) は、評価済みの `ast.Exp` オペランドを受け取り、それに基づいて機械語サイズの計算とOcode生成を行う必要がある。

## 残作業・次のステップ

1.  **命令ハンドラ (`opcodeEvalFn`) のリファクタリング:**
    *   `internal/pass1/handlers.go` 内の `processXXX` 関数のシグネチャを、`[]*token.ParseToken` ではなく `[]ast.Exp` を受け取るように変更する。
    *   各ハンドラ内で、渡された評価済み `ast.Exp` オペランド（`NumberExp`, `ImmExp` (未解決ラベル含む), `MemoryAddrExp` など）を解釈し、適切な処理を行うようにロジックを修正する。
2.  **機械語サイズ計算ロジックの修正:**
    *   各命令ハンドラ内で、評価済みオペランドに基づいて正確な機械語サイズを計算し、`env.LOC` を更新するように修正する。
3.  **`EquMap` の削除:**
    *   命令ハンドラのリファクタリング完了後、`Pass1` 構造体の古い `EquMap` フィールドと、`DefineMacro` 内の互換性維持コードを削除する。
4.  **テストの修正と実行:**
    *   既存のテスト (`TestAddExp`, `TestEvalProgramLOC`, `day03_harib00i_test.go` など) を新しい評価構造に合わせて修正し、実行して動作を確認する。

(過去の完了作業: [activeContext_archive_202503.md](../archives/activeContext_archive_202503.md))

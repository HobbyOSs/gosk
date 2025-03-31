# Progress

## 実装済み
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
- (保留) RESBの計算処理の実装
- (保留) `internal/codegen` パッケージのリファクタリング (CodeGenContext 導入)
- (保留) `internal/codegen` パッケージの不要パラメータ削除

## 関連情報
(技術的な詳細メモは `memory-bank/archives/technical_notes_archive_202503.md` にアーカイブ済み)
(過去の実装詳細: [progress_archive_202503.md](../archives/progress_archive_202503.md))

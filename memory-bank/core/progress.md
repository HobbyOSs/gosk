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

## まだ必要な実装
- **命令ハンドラ (`opcodeEvalFn`) のリファクタリング:**
    - `internal/pass1/handlers.go` 内の `processXXX` 関数を、古いスタックベース (`[]*token.ParseToken`) から新しい評価済みオペランド (`[]ast.Exp`) を受け取るように修正する。
    - 上記に伴い、各ハンドラ内の機械語サイズ計算ロジックも修正する。
- **`test/day03_harib00i_test.go` の相対ジャンプオフセットずれ調査:**
    - `CALL`, `JMP`, `Jcc` 命令の相対オフセットが期待値とずれている。
    - オフセット計算ロジック (`ターゲットアドレス - (現在の命令のアドレス + 命令のサイズ)`) の見直しが必要 (命令ハンドラリファクタリング後に再検証)。 (ユーザー側でアセンブラダンプマスタを使用して調査予定)
- **テストの修正と実行:**
    - 命令ハンドラのリファクタリング後、既存テストを新しい評価構造に合わせて修正し、実行する。
- **`EquMap` の削除:**
    - 命令ハンドラのリファクタリング完了後、`Pass1` 構造体の古い `EquMap` と関連コードを削除する。
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

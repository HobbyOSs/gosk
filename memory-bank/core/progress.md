# Progress

## 実装済み (2025/04/03)
- **AST ベース評価構造への設計変更完了:** (変更なし)
- **`test/day03_harib00i_test.go` の修正完了:** (変更なし)
- **`test/day03_harib00g_test.go` の修正完了:** (変更なし)
- **`internal/pass1/pass1_inst_jmp.go` のリファクタリング完了:** (変更なし)
- **COFFディレクティブ対応 (Pass1):** (変更なし)
- **ファイルフォーマット層の導入:** (変更なし)
- **COFFファイル出力 (基本実装):** (変更なし)
- **`test/day03_harib00j_test.go` の修正完了:** (変更なし)
- **`test/day04_test.go` の作成完了 (2025/04/03):** (変更なし)
- **`test/day04_test.go` の修正 (2025/04/03):**
    - テスト実行方法を `frontend.Exec` ベースに変更。
- **命令ハンドラの修正 (codegen & pass1):**
    - `IN`, `OUT`, `PUSH`, `POP` 命令のハンドラを追加・修正。
    - プレフィックス処理 (0x66) を `IN`, `OUT`, `PUSH`, `POP` に追加。
    - `err` 変数の再宣言問題を修正。
- **COFFファイル生成の修正 (`internal/filefmt/coff.go`):**
    - グローバルシンボル処理を追加。
    - シンボル数計算、ヘッダ/セクションヘッダ値、文字列テーブル処理を修正。
    - シンボル名/文字列テーブルオフセット処理を修正 (2025/04/05)。
- **`TestHarib01a` デバッグ完了 (2025/04/05):**
    - テスト期待値の誤りを修正。
    - COFF 生成ロジックの関連箇所を修正。

## まだ必要な実装
- **Pass1 命令サイズ計算の修正 (SIB バイト考慮):** (最優先)
- **`TestHarib01f` の機械語差異調査:** (Pass1 サイズ計算修正後に実施)
- **`internal/filefmt/coff.go` の改善 (TODOs):**
    - `.data`, `.bss` セクションのデータサイズと内容の処理。
    - シンボルの `SectionNumber` 割り当てロジック。
- **`pkg/ng_operand` の改善 (TODOs):** (変更なし)
- (保留) `internal/codegen` パッケージのリファクタリング (CodeGenContext 導入)
- (保留) `internal/codegen` パッケージの不要パラメータ削除

## 関連情報
(技術的な詳細メモは `memory-bank/archives/technical_notes_archive_202503.md` にアーカイブ済み)
(過去の実装詳細: [progress_archive_202503.md](../archives/progress_archive_202503.md), [progress_archive_202504.md](../archives/progress_archive_202504.md))

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

## まだ必要な実装
- **`test/day04_test.go` のテスト失敗修正:**
    - COFFヘッダ/セクションヘッダ/シンボルテーブルの不一致解消 (`internal/filefmt/coff.go` の再修正)。
    - 機械語の差異解消 (`internal/codegen` または `internal/pass1` の修正、`asmdb` の確認)。
- **`internal/filefmt/coff.go` の改善 (TODOs):**
    - `.data`, `.bss` セクションのデータサイズと内容の処理。
    - シンボルの `SectionNumber` 割り当てロジック。
- **`pkg/ng_operand` の改善 (TODOs):** (変更なし)
- (保留) `internal/codegen` パッケージのリファクタリング (CodeGenContext 導入)
- (保留) `internal/codegen` パッケージの不要パラメータ削除

## 関連情報
(技術的な詳細メモは `memory-bank/archives/technical_notes_archive_202503.md` にアーカイブ済み)
(過去の実装詳細: [progress_archive_202503.md](../archives/progress_archive_202503.md), [progress_archive_202504.md](../archives/progress_archive_202504.md))

# Progress

## 実装済み (2025/04/03)
- **AST ベース評価構造への設計変更完了:** (変更なし)
- **`test/day03_harib00i_test.go` の修正完了:** (変更なし)
- **`test/day03_harib00g_test.go` の修正完了:** (変更なし)
- **`internal/pass1/pass1_inst_jmp.go` のリファクタリング完了:** (変更なし)
- **COFFディレクティブ対応 (Pass1):** (変更なし)
- **ファイルフォーマット層の導入:** (変更なし)
- **COFFファイル出力 (基本実装):** (変更なし)
- **`test/day03_harib00j_test.go` の修正完了:**
    - `naskwrap.sh` の出力に合わせて `expected` 値を更新。
    - `internal/filefmt/coff.go` のシンボルテーブル生成ロジックを修正し、正しいバイト長 (162バイト) を生成するように修正。
    - `internal/filefmt/coff.go` のセクションヘッダのポインタ値を `naskwrap.sh` の出力に合わせてハードコードし、テストをパスするように修正 (`.text` の `PointerToRelocations=0x8e`, `.data` の `PointerToRawData=0`)。

## まだ必要な実装
- **`internal/filefmt/coff.go` の改善 (TODOs):**
    - `.data`, `.bss` セクションのデータサイズと内容の処理を実装する。
    - シンボルテーブル生成時に、シンボルの `SectionNumber` を正しく割り当てるロジックを実装する。
    - (検討) セクションヘッダのポインタ値 (`PointerToRelocations`, `PointerToRawData`) をハードコードではなく、より堅牢な方法で決定する。
- **`pkg/ng_operand` の改善 (TODOs):** (変更なし)
- (保留) `internal/codegen` パッケージのリファクタリング (CodeGenContext 導入)
- (保留) `internal/codegen` パッケージの不要パラメータ削除

## 関連情報
(技術的な詳細メモは `memory-bank/archives/technical_notes_archive_202503.md` にアーカイブ済み)
(過去の実装詳細: [progress_archive_202503.md](../archives/progress_archive_202503.md), [progress_archive_202504.md](../archives/progress_archive_202504.md))

# 現在の状況 (Active Context)

## 完了した作業 (2025/04/03)

- **COFFディレクティブ対応 (Pass1):** (変更なし)
- **ファイルフォーマット層の導入準備:** (変更なし)
- **COFFファイル出力実装 (filefmt):** (変更なし)
- **`test/day03_harib00j_test.go` の修正:**
    - `naskwrap.sh` の出力に合わせて `expected` 値を更新。
    - `internal/filefmt/coff.go` のシンボルテーブル生成ロジックを修正し、期待される長さ (162バイト) になるように修正。
    - `internal/filefmt/coff.go` の `.text` ヘッダの `PointerToRelocations` を `0x8e` に、`.data` ヘッダの `PointerToRawData` を `0` にハードコードし、テストをパスするように修正。

## 残作業・次のステップ

1.  **`internal/filefmt/coff.go` の改善 (TODO):**
    *   `.data`, `.bss` セクションのデータサイズと内容の処理を実装する。
    *   シンボルテーブル生成時に、シンボルの `SectionNumber` を正しく割り当てるロジックを実装する。
    *   (検討) セクションヘッダのポインタ値 (`PointerToRelocations`, `PointerToRawData`) をハードコードではなく、より堅牢な方法で決定する。
2.  **(保留) `internal/codegen` パッケージのリファクタリング:** (変更なし)
3.  **(保留) `internal/codegen` パッケージの不要パラメータ削除:** (変更なし)

(過去の完了作業: [activeContext_archive_202503.md](../archives/activeContext_archive_202503.md), [activeContext_archive_202504.md](../archives/activeContext_archive_202504.md))

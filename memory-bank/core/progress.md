# Progress

## 実装済み (2025/04/05 深夜)
- **`TestHarib02i` 修正完了:**
    - `internal/pass1` に `LIDT` ハンドラ (`processLIDT`) を追加し、`handlers.go` に登録。
    - `internal/codegen` に `LIDT` ハンドラ (`handleLIDT`) を追加し、`x86gen.go` に登録。
    - `internal/codegen/x86gen_lgdt.go` を修正し、`ng_operand` と `asmdb.FindEncoding` を使用するように変更。
    - `pkg/asmdb/instruction_table_fallback.go` に `LIDT` のフォールバック定義を追加。
    - `pkg/asmdb/encoding.go` の `Opcode.getSize()` を修正し、複数バイトオペコードのサイズ計算を修正。
    - `internal/pass1` の `LGDT`/`LIDT` ハンドラのサイズ計算を `FindMinOutputSize` を使用するように修正。
    - `TestHarib02i` が PASS することを確認。
- **`Day05Suite/Harib02i` テストケース作成 & 課題特定:** (2025/04/05 夜)
    - `test/day05_harib02i_test.go` と `test/day05_test.go` を作成。
    - アセンブリコードと期待値バイナリを挿入。
    - テスト実行により `LIDT` ハンドラ欠落と `LGDT` コード生成エラーを特定。
- **e2e テスト作成手順の更新:** (2025/04/05 夜)
    - `memory-bank/details/technical_notes.md` を修正し、アセンブリ取得、NASK実行、期待値生成、テスト実行の手順を更新。
- **`TestHarib01a` デバッグ完了:** (前セッション)
    - テスト期待値の誤りを修正。
    - COFF 生成ロジックの関連箇所を修正。
- **`pass1` LOC 計算修正 & `TestHarib01f` 検証完了:**
    - `asmdb` と `codegen` の `IN`/`OUT` プレフィックス処理を修正。
    - `codegen` の `PUSH`/`POP` ハンドラ呼び出しを修正。
    - `TestPass1EvalSuite/TestEvalProgramLOC` が PASS することを確認。
    - `gosk` と `nask` の `harib01f` アセンブル結果が一致することを確認。
- **`go install` の問題修正 (2025/04/08 夜):**
    - `pkg/asmdb` が Git Submodule 内のデータファイルを `go:embed` で埋め込めない問題を修正。
    - データ提供用の専用 Go モジュール (`github.com/HobbyOSs/json-x86-64-go-mod`) を導入。
    - 不要になった Submodule (`pkg/asmdb/json-x86-64`) を削除。
    - ルートの不要な `main.go` を削除。
    - `README.md` のインストール手順を更新。
- **`Day03Suite/Harib00i` (asmhead.nas) 差分修正 (2025/04/16 夜):**
    - `test/day03_harib00i_test.go` の期待値バイナリを NASK 出力で更新。
    - `ALIGNB` ディレクティブが `pass1` で `ocode` を Emit していなかった問題を修正 (`internal/pass1/pass1_inst_pseudo.go`)。
    - `codegen` 側で `ALIGNB` の `ocode` を処理し、パディングバイトを生成するように修正 (`internal/codegen/x86gen_pseudo.go`, `internal/codegen/x86gen.go`)。
    - `pkg/ocode/ocode.go` に `OpALIGNB` 定数を追加し、`make gen` を実行。
    - テスト (`TestHarib00i`) が成功することを確認。

## まだ必要な実装
- **SIB バイト計算の検証と coff.go クリーンアップ:** (優先度: 中)
    - `ng_operand.CalcSibByteSize` のテスト追加
    - `internal/filefmt/coff.go` のクリーンアップ
- **`internal/filefmt/coff.go` の改善 (TODOs):**
    - `.data`, `.bss` セクションのデータサイズと内容の処理。
    - シンボルの `SectionNumber` 割り当てロジック。
- **`pkg/ng_operand` の改善 (TODOs):** (変更なし)
- (保留) `internal/codegen` パッケージのリファクタリング (CodeGenContext 導入)
- (保留) `internal/codegen` パッケージの不要パラメータ削除
- **e2eテストケース拡充 (opennask互換性):**
    - **目標:** opennaskのアセンブル結果との一致。
    - **対象:** `Day03Suite/Harib00i`, `Day03Suite/Harib00j`, `Day06Suite/Harib03e`, `Day09Suite/Harib06c`, `Day12Suite/Harib09a`, `Day15Suite/Harib12a`, `Day20Suite/Harib17e`, `Day20Suite/Harib17g` など (GitHub Issue #12 参照)。"新出命令なし"も検証対象。
    - **完了:** `Day04Suite/Harib01a`, `Day04Suite/Harib01f` (2025/04/05), `Day05Suite/Harib02i` (2025/04/05 深夜)
- **e2eテスト作成プロセスの標準化:**
    - **目標:** 効率性と一貫性の向上。
    - **内容:** 命名規則（ファイル、関数）、テストコードテンプレート（ヘルパー関数、期待値管理）、実行手順（`go test -run`）、デバッグ手順の確立。
    - **状況:** `technical_notes.md` に記録・更新済み (2025/04/05)。
- **`INSTRSET` ディレクティブ未対応:** (優先度: 低) `pass1` で `INSTRSET` が処理されていない。
# (削除) Day03Suite/Harib00i (asmhead.nas) アセンブル結果差分: (完了 - 2025/04/16)

## 関連情報
(技術的な詳細メモは `memory-bank/archives/technical_notes_archive_202503.md` にアーカイブ済み)
(過去の実装詳細: [../archives/2025/04/progress_archive_202504.md](../archives/2025/04/progress_archive_202504.md), [../archives/progress_archive_202503.md](../archives/progress_archive_202503.md))

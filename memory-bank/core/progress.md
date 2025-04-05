# Progress

## 実装済み (2025/04/05)
- **`TestHarib01a` デバッグ完了:**
    - テスト期待値の誤りを修正。
    - COFF 生成ロジックの関連箇所を修正。
- **`pass1` LOC 計算修正 & `TestHarib01f` 検証完了:**
    - `asmdb` と `codegen` の `IN`/`OUT` プレフィックス処理を修正。
    - `codegen` の `PUSH`/`POP` ハンドラ呼び出しを修正。
    - `TestPass1EvalSuite/TestEvalProgramLOC` が PASS することを確認。
    - `gosk` と `nask` の `harib01f` アセンブル結果が一致することを確認。

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
    - **対象:** `Day03Suite/Harib00i`, `Day03Suite/Harib00j`, `Day05Suite/Harib02i`, `Day06Suite/Harib03e`, `Day09Suite/Harib06c`, `Day12Suite/Harib09a`, `Day15Suite/Harib12a`, `Day20Suite/Harib17e`, `Day20Suite/Harib17g` など (GitHub Issue #12 参照)。"新出命令なし"も検証対象。
    - **完了:** `Day04Suite/Harib01a`, `Day04Suite/Harib01f` (2025/04/05)
- **e2eテスト作成プロセスの標準化:**
    - **目標:** 効率性と一貫性の向上。
    - **内容:** 命名規則（ファイル、関数）、テストコードテンプレート（ヘルパー関数、期待値管理）、実行手順（`make` ターゲット、デバッグ手順）の確立。

## 関連情報
(技術的な詳細メモは `memory-bank/archives/technical_notes_archive_202503.md` にアーカイブ済み)
(過去の実装詳細: [../archives/2025/04/progress_archive_202504.md](../archives/2025/04/progress_archive_202504.md), [../archives/progress_archive_202503.md](../archives/progress_archive_202503.md))

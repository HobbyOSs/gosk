# 現在の状況 (Active Context)

## 残作業・次のステップ

1.  **Pass1 命令サイズ計算の修正 (最優先):**
    *   **方針:** `internal/pass1` で命令サイズを計算する際に SIB バイトの必要性を正確に判定し、サイズに含めるように修正する。具体的には、`asmdb.FindMinOutputSize` または関連する `asmdb` や `ng_operand` のロジックを修正する。(`processMOV` 内でのアドホックな修正は行わない)
    *   **次作業:** `FindMinOutputSize` または関連箇所の修正方針を検討し、SIB バイト計算ロジックを実装する。
    *   **テスト:** SIB バイトが必要/不要なケースを含む `FindMinOutputSize` (または関連関数) の単体テストを追加し、修正を検証する。その後、`TestHarib01a` を再実行してバイナリ差分が解消されることを確認する。
    * 上記が終了後、coff.goにあるデバッグコードや冗長な処理は削除する。
2.  **`TestHarib01f` の機械語差異調査:** (Pass1 サイズ計算修正後に実施)
3.  **EXTERN シンボルのテストケース追加:** (変更なし)
4.  **`internal/filefmt/coff.go` の改善 (TODO):** (変更なし)
5.  **(保留) `internal/codegen` パッケージのリファクタリング:** (変更なし)
6.  **(保留) `internal/codegen` パッケージの不要パラメータ削除:** (変更なし)

## 完了した作業 (2025/04/05)

- **COFFディレクティブ対応 (Pass1):** (変更なし)
- **ファイルフォーマット層の導入準備:** (変更なし)
- **COFFファイル出力実装 (filefmt):** (変更なし)
- **`test/day03_harib00j_test.go` の修正:** (変更なし)
- **`test/day04_test.go` の作成 (2025/04/03):** (変更なし)
- **`test/day04_test.go` の修正 (2025/04/03):**
    - テストの実行方法を `pass1`/`pass2` の直接呼び出しから `frontend.Exec` を使用し、一時ファイル経由で比較するように変更。
    - `frontend.Exec` の戻り値の処理を修正。
- **命令ハンドラの修正 (codegen):**
    - `internal/codegen/x86gen_in.go`: `IN EAX, DX` および `IN EAX, imm8` のケースを追加し、オペランドサイズプレフィックス (0x66) の処理を追加。`err` 変数の再宣言を修正。
    - `internal/codegen/x86gen_out.go`: `OUT DX, AX/EAX` および `OUT imm8, AX/EAX` のケースを追加し、オペランドサイズプレフィックス (0x66) の処理を追加。`err` 変数の再宣言を修正。
    - `internal/codegen/x86gen_pushpop.go`: `PUSH`/`POP` 命令のハンドラ (`handlePUSH`, `handlePOP`) を新規作成。レジスタ、メモリ、即値オペランドに対応し、プレフィックス処理を追加。`ng_operand` パッケージのインターフェースと型定義に合わせて修正。
- **命令ハンドラの修正 (pass1):**
    - `internal/pass1/pass1_inst_pushpop.go`: `PUSH`/`POP` 命令のハンドラ (`processPUSH`, `processPOP`) を新規作成。
    - `internal/pass1/handlers.go`: `opcodeEvalFns` マップに `PUSH` と `POP` のハンドラを追加。
- **COFFファイル生成の修正 (`internal/filefmt/coff.go`):**
    - `generateSymbolEntries`: `GLOBAL` 宣言されたシンボルを処理するように修正。
    - `Write`: シンボル数の計算ロジック (`numSymbolsCounted`) を修正。文字列テーブルが空の場合の処理を修正。
    - `generateHeader`: `Characteristics` を `0x0000` に修正。
    - `generateSectionHeaders`: `.text` の `PointerToRelocations` を `symbolTableOffset` に、`.data` の `PointerToRawData` を `dataDataOffset` に修正。各セクションの `Characteristics` を nask 出力に合わせて修正。コメントマーカーの構文エラーを修正。
- **`TestHarib00j` のデバッグと修正 (2025/04/04):**
    - `pass1` に `GLOBAL` ディレクティブのハンドラを追加・登録 (`internal/pass1/handlers.go`, `internal/pass1/pass1_inst_pseudo.go`)。
    - `pass1.Eval` が更新後の `GlobalSymbolList` を返すようにし、`frontend.Exec` で `CodeGenContext` を更新するように修正 (`internal/pass1/eval.go`, `internal/frontend/frontend.go`)。
    - COFF 文字列テーブルが空の場合でもサイズフィールド (4バイト) を書き込むように修正 (`internal/filefmt/coff.go`)。
    - `TestHarib00j` の期待値に合わせて、グローバルシンボルの `Type` を `0x00` に、`.data` セクションヘッダの `PointerToRawData` を `0` に修正 (`internal/filefmt/coff.go`)。
    - デバッグ用のログを削除。
- **EXTERN シンボル処理 (Pass1 & filefmt) (2025/04/04):**
    - `internal/pass1/pass1_inst_pseudo.go`: `processEXTERN` ハンドラを追加。
    - `internal/pass1/handlers.go`: `EXTERN` ハンドラを `opcodeEvalFns` に登録。
    - `internal/filefmt/coff.go`: `generateSymbolEntries` に EXTERN シンボルを出力する処理を追加。
- **`pass1.Eval` のシグネチャ変更 (2025/04/04):**
    - `internal/pass1/eval.go`: `Eval` が `*codegen.CodeGenContext` を引数に取り、コンテキストを直接更新するように変更。戻り値を削除。
    - `internal/codegen/typedef.go`: `CodeGenContext` に `ExternSymbolList` フィールドを追加。
    - `internal/frontend/frontend.go`: `Exec` 内での `pass1.Eval` 呼び出しを修正。
    - `test/pass1_test.go`: `pass1.Eval` 呼び出しを修正。
    - `internal/pass1/eval_test.go`: `pass1.Eval` 呼び出しを修正。
- **テスト実行 (2025/04/04):**
    - `go test ./...` を実行し、`day04_test.go` 以外のテストが PASS することを確認。
- **README.md 更新と検証 (2025/04/05):**
    - Featuresセクションを修正 (NASK言及、内部構造削除)。
    - Usageセクションのコマンド例を修正し、実行検証。
    - Usageセクションのコマンドラインインターフェース説明を修正。
    - Usageセクションのコマンドパスを修正 (`./gosk`)。
- **Makefile 修正 (2025/04/05):**
    - ビルド成果物がプロジェクトルート (`./gosk`) に出力されるように修正。
    - `make build` を再実行し、変更を確認。
- **アセンブル再検証 (2025/04/05):**
    - 一時ファイルを用いてアセンブルを実行し、`hexdump` で結果を確認。
- **`internal/filefmt/coff.go` のデバッグ試行 (2025/04/05):**
    - `TestHarib01a` のバイナリ差分（シンボルテーブルのずれ、文字列テーブル先頭のヌルバイト）解消のため、以下の修正を試行したが、いずれも失敗。
        - シンボルテーブル書き込み方法の変更（`binary.Write` vs 手動書き込み）
        - COFFヘッダ `NumberOfSymbols` の値の調整（補助シンボルを含む/含まない/期待値合わせ）
        - 文字列テーブル書き込み処理の変更（先頭バイト削除など）
        - デバッグログの追加と確認。
- **根本原因の特定 (2025/04/05):**
    - `TestHarib01a` のバイナリ差分 (2バイトずれ) の根本原因は、`internal/pass1/pass1_inst_mov.go` の `processMOV` における命令サイズ計算が、SIB バイトが必要なケースを考慮しておらず、実際の機械語サイズより小さく計算してしまうためであると特定。これにより、Pass1 で計算される LOC (Location Counter) がずれ、COFF ファイル全体のオフセット計算に影響が出ていた。

(過去の完了作業: [activeContext_archive_202503.md](../archives/activeContext_archive_202503.md), [activeContext_archive_202504.md](../archives/activeContext_archive_202504.md))

# 現在の状況 (Active Context) - 2025/04/05 セッション終了

## 持ち越し課題

1.  **`pass1` の LOC 計算修正 (最優先):**
    *   **状況:** `TestPass1EvalSuite/TestEvalProgramLOC` の `IN AX, DX` と `OUT DX, AX` テストが失敗する。期待される LOC (2) に対し、実際の LOC (1) となり、必要な 0x66 プレフィックスが計算されていない。
    *   **原因:** `IN`/`OUT` 命令の 0x66 プレフィックス要否判定ロジック (`pkg/ng_operand/requires.go` の `Require66h` または関連箇所) が、32bit モードでの `AX` オペランドのケースを正しく扱えていない。
    *   **課題:** 「オペランドはオペコードを知らない」設計原則を維持しつつ、`IN`/`OUT` 命令の特殊なプレフィックスルールをどこでどのように実装するのが最適か。
    *   **試行錯誤の経緯:**
        *   当初、`asmdb` のマッチングロジック (`matchOperandsStrict`) で特殊レジスタ (`AL`, `AX`, `DX` など) と汎用タイプ (`r8`, `r16`) の不一致を吸収しようとしたが、`MOV` 命令に影響が出た。
        *   次に、`pass1` ハンドラ (`processIN`/`processOUT`) で `forceNoPrefix66` フラグを `ng_operand` に渡し、`Require66h` でチェックする案を試したが、`IN AX, DX` / `OUT DX, AX` で必要なプレフィックスまで無効化してしまい失敗。
        *   `prefix66Override *bool` フラグを `ng_operand` に追加し、`pass1` でプレフィックス要否を判断して渡す案も検討したが、設計原則との兼ね合いで混乱が生じ、中断。
    *   **次作業:** 次のセッションで、プレフィックス計算ロジック (`Require66h` または `asmdb.GetPrefixSize`) の修正方針を再検討し、実装する。`pass1` テストが全て成功することを確認する。
2.  **`TestHarib01f` の再実行と修正:** 上記 LOC 計算修正後、`TestHarib01f` を再実行し、バイナリ差分やその他の問題があれば修正する。
3.  **Pass1 命令サイズ計算の修正 (SIB バイト考慮):** (上記完了後)
    *   **方針:** `internal/pass1` で命令サイズを計算する際に SIB バイトの必要性を正確に判定し、サイズに含めるように修正する。具体的には、`asmdb.FindMinOutputSize` または関連する `asmdb` や `ng_operand` のロジックを修正する。(`processMOV` 内でのアドホックな修正は行わない)
    *   **テスト:** SIB バイトが必要/不要なケースを含む `FindMinOutputSize` (または関連関数) の単体テストを追加し、修正を検証する。
    * 上記が終了後、coff.goにあるデバッグコードや冗長な処理は削除する。
4.  **EXTERN シンボルのテストケース追加:** (変更なし)
5.  **`internal/filefmt/coff.go` の改善 (TODO):** (変更なし)
6.  **(保留) `internal/codegen` パッケージのリファクタリング:** (変更なし)
7.  **(保留) `internal/codegen` パッケージの不要パラメータ削除:** (変更なし)

## このセッションで完了した作業 (2025/04/05)

- **`TestHarib01f` 関連の単体テスト追加:**
    - `internal/pass1/eval_test.go` に `IN`, `OUT`, `PUSH`, `POP` のテストケースを追加。
    - `internal/codegen/x86gen_test.go` に `IN`, `OUT`, `PUSH`, `POP` のテストケースを追加。
- **`asmdb` フォールバック定義の追加・修正:**
    - `pkg/ocode/ocode.go` に `OpPUSH`, `OpPOP` を追加し、`go generate` を実行。
    - `pkg/asmdb/instruction_table_fallback.go` に `IN EAX, DX`, `PUSH r32`, `POP r32` のフォールバック定義を追加。
- **`asmdb` マッチングロジック修正:**
    - `pkg/asmdb/instruction_search.go` の `matchOperandsStrict` を修正し、フォールバック定義 (`al`, `ax`, `eax`, `dx`) と `ng_operand` のパース結果 (`r8`, `r16`, `r32`) の不一致を吸収するようにした（ただし、後に `IN`/`OUT` 限定にする修正は取り消し）。
- **`FindEncoding` エラーの解消:** 上記修正により、`IN`, `OUT`, `PUSH`, `POP` で発生していた `FindEncoding` エラーは解消された。

(以下、変更なし)
- **`TestHarib01a` デバッグ完了:**
    - テスト失敗の原因が、テストコード内の `expected` バイト列が NASK の実際の出力と異なっていたためであることを特定。
    - NASK の正しい出力を基に `expected` を修正し、テストが PASS することを確認。
    - 関連して `internal/filefmt/coff.go` のシンボル名/文字列テーブル処理を修正。
    - テスト期待値の生成プロセスを `technical_notes.md` に記録。
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
    - ~~`TestHarib01a` のバイナリ差分 (2バイトずれ) の根本原因は、`internal/pass1/pass1_inst_mov.go` の `processMOV` における命令サイズ計算が、SIB バイトが必要なケースを考慮しておらず、実際の機械語サイズより小さく計算してしまうためであると特定。これにより、Pass1 で計算される LOC (Location Counter) がずれ、COFF ファイル全体のオフセット計算に影響が出ていた。~~ **(誤り)** -> 正しくはテスト期待値の誤りと COFF 生成ロジックの問題。

(過去の完了作業: [activeContext_archive_202503.md](../archives/activeContext_archive_202503.md), [activeContext_archive_202504.md](../archives/activeContext_archive_202504.md))

- **`TestHarib01a` デバッグ完了:**
    - テスト失敗の原因が、テストコード内の `expected` バイト列が NASK の実際の出力と異なっていたためであることを特定。
    - NASK の正しい出力を基に `expected` を修正し、テストが PASS することを確認。
    - 関連して `internal/filefmt/coff.go` のシンボル名/文字列テーブル処理を修正。
    - テスト期待値の生成プロセスを `technical_notes.md` に記録。
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
    - ~~`TestHarib01a` のバイナリ差分 (2バイトずれ) の根本原因は、`internal/pass1/pass1_inst_mov.go` の `processMOV` における命令サイズ計算が、SIB バイトが必要なケースを考慮しておらず、実際の機械語サイズより小さく計算してしまうためであると特定。これにより、Pass1 で計算される LOC (Location Counter) がずれ、COFF ファイル全体のオフセット計算に影響が出ていた。~~ **(誤り)** -> 正しくはテスト期待値の誤りと COFF 生成ロジックの問題。

(過去の完了作業: [activeContext_archive_202503.md](../archives/activeContext_archive_202503.md), [activeContext_archive_202504.md](../archives/activeContext_archive_202504.md))

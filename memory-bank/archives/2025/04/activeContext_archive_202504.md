# Active Context Archive - 2025/04

## 2025/04/04 以前の完了作業

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
- **COFFディレクティブ対応 (Pass1):** (変更なし)
- **ファイルフォーマット層の導入準備:** (変更なし)
- **COFFファイル出力実装 (filefmt):** (変更なし)
- **`test/day03_harib00j_test.go` の修正:** (変更なし)

(さらに古い履歴は [activeContext_archive_202503.md](../activeContext_archive_202503.md) を参照)

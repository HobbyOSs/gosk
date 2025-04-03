# 現在の状況 (Active Context)

## 完了した作業 (2025/04/03)

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

## 残作業・次のステップ

1.  **`test/day04_test.go` のテスト失敗原因調査と修正:**
    *   COFFヘッダ/セクションヘッダの値の不一致 (`NumberOfSymbols`, `Characteristics` など) を修正する。(`internal/filefmt/coff.go` の再修正が必要。`TestHarib00j` とは異なる値になる可能性あり)
    *   シンボルテーブルの内容/長さの不一致を修正する。(`internal/filefmt/coff.go` の `generateSymbolEntries` の再修正が必要。グローバルシンボルの `Type` など)
    *   `TestHarib01f` での機械語の差異 (`PUSH`/`POP` など) の原因を調査し、`internal/codegen` または `internal/pass1` を修正する。(`asmdb` の確認も必要になる可能性あり)
2.  **`internal/filefmt/coff.go` の改善 (TODO):**
    *   `.data`, `.bss` セクションのデータサイズと内容の処理を実装する。
    *   シンボルテーブル生成時に、シンボルの `SectionNumber` を正しく割り当てるロジックを実装する。
3.  **EXTERN シンボルの処理実装:**
    *   `pass1` で `EXTERN` ディレクティブを処理し、`ExternSymbolList` に追加する。
    *   `filefmt/coff.go` の `generateSymbolEntries` で `EXTERN` シンボルを適切に処理する (例: `SectionNumber = 0`, `StorageClass = IMAGE_SYM_CLASS_EXTERNAL`)。
4.  **(保留) `internal/codegen` パッケージのリファクタリング:** (変更なし)
5.  **(保留) `internal/codegen` パッケージの不要パラメータ削除:** (変更なし)

(過去の完了作業: [activeContext_archive_202503.md](../archives/activeContext_archive_202503.md), [activeContext_archive_202504.md](../archives/activeContext_archive_202504.md))

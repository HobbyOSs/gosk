# 技術ノート

## COFF シンボルテーブルの分析 (2025/04/05)

ユーザー提供の nask 出力 COFF ファイルのシンボルテーブルダンプに基づく分析。

### シンボルテーブルダンプ例 (抜粋)

```
00000160: 0000 0000 0000 0000 0000 005f 696f 5f68  ..........._io_h
00000170: 6c74 0000 0000 0001 0000 0002 005f 696f  lt..........._io
00000180: 5f63 6c69 0002 0000 0001 0000 0002 005f  _cli..........._
00000190: 696f 5f73 7469 0004 0000 0001 0000 0002  io_sti..........
... (以下略) ...
```

### シンボルテーブルエントリ構造 (18バイト)

```
| Bytes | Field              | Example (_io_hlt)        | Notes                                      |
|-------|--------------------|--------------------------|--------------------------------------------|
| 0-7   | Name               | 5f696f5f686c7400         | シンボル名 (8バイト、短い場合はヌル埋め)     |
| 8-11  | Value              | 00000000                 | リロケータブルアドレス (後述)              |
| 12-13 | SectionNumber      | 0100                     | 関連セクション番号 (1 = .text)             |
| 14-15 | Type               | 0000                     | シンボル型 (0x00 = NULL)                   |
| 16    | StorageClass       | 02                       | 格納クラス (0x02 = C_EXT, 外部シンボル)    |
| 17    | NumberOfAuxSymbols | 00                       | 補助シンボル数                             |
```

### `Value` フィールドの解釈 (リロケータブルアドレス)

- シンボルテーブル内の `Value` フィールドは、そのシンボルが定義されているアドレス（`.text` セクション先頭からのオフセット）を示すリロケータブルアドレスである。
- このアドレスは、先行するシンボルからの機械語のバイト数を累積して計算される。

**例:**

```
| Symbol    | Value | Preceding Instructions & Size | Calculation        |
|-----------|-------|-------------------------------|--------------------|
| _io_hlt   | 0     | (Start of .text)              | 0                  |
| _io_cli   | 2     | HLT (1), RET (1)              | 0 + 1 + 1 = 2      |
| _io_sti   | 4     | CLI (1), RET (1)              | 2 + 1 + 1 = 4      |
| (Unnamed) | 6     | STI (1), RET (1)              | 4 + 1 + 1 = 6      |
| _io_in8   | 9     | STI (1), HLT (1), RET (1)     | 6 + 1 + 1 + 1 = 9  |
| ...       | ...   | ...                           | ...                |
```

**実装への影響:**

- `internal/filefmt/coff.go` の `generateSymbolEntries` で `CoffSymbol` を生成する際、`Value` フィールドには `pass1` で計算されたシンボルのアドレス（LOC）を正しく設定する必要がある。現在の実装では `ctx.SymTable[globalName]` から取得しており、これは正しいはずだが、ずれの原因調査においてこの値の正確性も確認する必要がある。

(過去の技術ノートは `memory-bank/archives/technical_notes_archive_YYYYMM.md` にアーカイブされています。)

## テスト期待値 (expected) の生成方法 (2025/04/05)

`go test` で使用するバイナリ期待値 (`expected []byte`) は、基準となるアセンブラ (現在は NASK) の出力を元に生成する。これにより、`gosk` の出力が基準アセンブラと一致するかどうかを正確に検証できる。

**手順:**

1.  **NASK バイナリ生成**: テスト対象のアセンブリコード (`.asm`) を NASK でアセンブルし、バイナリファイル (`.bin`) を生成する。
    ```bash
    # (NASK の実行コマンド例 - 環境に合わせてパス等を調整)
    # wine /path/to/nask.exe input.asm output.bin output.lst
    ```
    *注意*: NASK の実行環境やパスはプロジェクト外の情報であるため、具体的なパスは記録しない。

2.  **Go リテラル生成**: 生成されたバイナリファイルを読み込み、Go の `[]byte` リテラル形式 (`[]byte{0x.., 0x.., ...}`) に変換する Go プログラム (`generate_expected.go` など) を作成・実行する。
    ```go
    // generate_expected.go (抜粋)
    package main
    import ("fmt"; "os"; "strings")
    func main() {
        data, _ := os.ReadFile("output.bin") // NASK が生成したバイナリ
        var builder strings.Builder
        builder.WriteString("[]byte{\n")
        for i, b := range data {
            if i%16 == 0 { builder.WriteString("\t\t") }
            builder.WriteString(fmt.Sprintf("0x%02x,", b))
            if (i+1)%16 == 0 { builder.WriteString(" //\n") } else { builder.WriteString(" ") }
        }
        if len(data)%16 != 0 { builder.WriteString("\n") }
        builder.WriteString("\t}")
        fmt.Println(builder.String())
    }
    ```
    ```bash
    go run generate_expected.go > expected_literal.txt
    ```

3.  **テストコード更新**: 生成された `[]byte` リテラルをコピーし、該当するテストケース (`*_test.go`) 内の `expected` 変数に貼り付ける。

このプロセスにより、`expected` データが常に基準アセンブラの最新の正しい出力を反映するようになり、テストの信頼性が向上する。

## e2e テスト作成プロセスの標準化案 (2025/04/05)

今後の e2e テストケース追加・修正を効率化し、一貫性を保つための標準プロセス案。

### 1. テストファイル命名規則

- **形式:** `test/dayXX_haribYY<suffix>_test.go`
    - `XX`: 対応する「〇日目」の番号 (例: `03`)
    - `YY`: 対応する「harib」の番号 (例: `00`)
    - `<suffix>`: テストケースの識別子 (例: `i`, `j`)
- **例:** `test/day03_harib00i_test.go`

### 2. テスト関数命名規則

- **形式:** `TestDayXXSuite/TestHaribYY<suffix>` (testify/suite を使用する場合)
    - `XX`, `YY`, `<suffix>` はファイル名と同様。
- **例:** `TestDay03Suite/TestHarib00i`

### 3. テストコードの構造とヘルパー関数 (実際のコードに基づく - 2025/04/05)

- **目標:** 定型的な処理を共通化し、テストコードの記述量を削減する。
- **実際のテストコード構造 (例: `test/day03_harib00i_test.go`):**
    1.  **アセンブリコード定義:** テスト対象のアセンブリコード（通常は `haribote-os` のサンプルコードなど、開発者が指定したもの）を Go の文字列リテラルとしてテストコード内に定義する。
    2.  **gosk アセンブル実行:**
        *   `gen.Parse` で上記のアセンブリコード文字列をパースする。
        *   `frontend.Exec` でアセンブルを実行し、結果を一時ファイルに出力する。
        *   `ReadFileAsBytes` (`test/readbin.go` で定義) を使用して一時ファイルの内容を `actual []byte` として読み込む。
    3.  **期待値バイナリ定義:**
        *   `defineHEX` ヘルパー関数 (`test/test_helper.go` で定義) を使用する。この関数は、`DATA` や `FILL` 命令を含むDSL形式の文字列配列を受け取り、期待されるバイナリ (`expected []byte`) を生成する。これにより、バイナリデータを直接記述するよりも可読性が向上する。
        *   DSLの元となる期待値バイナリは、「テスト期待値 (expected) の生成方法」セクションの手順に従い、NASK の出力から生成する。
    4.  **結果比較:**
        *   `github.com/google/go-cmp/cmp` パッケージの `cmp.Diff(expected, actual)` を使用して、期待値と実際の結果を比較する。
        *   差分が存在する場合 (`diff != ""`) は、`DumpDiff` ヘルパー関数 (`test/diff.go` で定義) を使用して詳細な16進ダンプ形式の差分をログに出力し、`t.Fail()` でテストを失敗させる。
        *   バイト列の長さも別途 `len(expected) != len(actual)` で比較する。
- **使用されているヘルパー関数:**
    - `defineHEX(dsl []string) []byte` (`test/test_helper.go`): DSLから期待値バイナリを生成。
    - `ReadFileAsBytes(filePath string) ([]byte, error)` (`test/readbin.go`): ファイル内容をバイト列として読み込む。
    - `DumpDiff(expected, actual []byte, color bool) string` (`test/diff.go`): バイト列の差分を16進ダンプ形式で表示。
- **期待値 (`expected`) の管理:**
    - 現在は `defineHEX` を用いたDSL形式でテストコード内に記述されている。
    - 将来的には、期待値バイナリを `.golden` ファイルとして分離し、テスト実行時にファイルを比較する方法も検討可能。

### 4. テスト実行手順

- **Makefile ターゲット:**
    - `make test-e2e`: 全ての e2e テストを実行する。
    - `make test-e2e-dayXX`: 特定の日のテストスイートを実行する (例: `make test-e2e-day03`)。
        - `go test ./test -run TestDayXXSuite` を実行。
    - `make test-e2e-haribYY<suffix>`: 特定のテストケースを実行する (例: `make test-e2e-harib00i`)。
        - `go test ./test -run TestDayXXSuite/TestHaribYY<suffix>` を実行。
- **デバッグ手順:**
    1. 失敗したテストケースを特定する (`make test-e2e-haribYY<suffix>` で再実行)。
    2. `assembleWithNask` と `assembleWithGosk` でそれぞれバイナリを生成し、`hexdump -C` や `diff <(hexdump -C nask.bin) <(hexdump -C gosk.bin)` などで差分を詳細に比較する。
    3. `gosk` のデバッグログ (`-v` オプションなど、必要なら追加) を確認する。
    4. 関連する `internal/` パッケージ (pass1, codegen, asmdb など) のコードをデバッグする。

この標準化案は、実装を進めながら改善していく。

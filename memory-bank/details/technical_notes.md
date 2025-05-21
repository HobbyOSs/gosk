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

## テスト期待値 (expected) の生成方法 (2025/04/05 更新)

`go test` で使用するバイナリ期待値 (`expected []byte`) は、基準となるアセンブラ (現在は NASK) の出力を元に生成する。これにより、`gosk` の出力が基準アセンブラと一致するかどうかを正確に検証できる。

**手順:**

1.  **NASK バイナリ生成**: テスト対象のアセンブリコード (`.nas` ファイル) を NASK でアセンブルし、バイナリファイル (`.bin`) を生成する。
    ```bash
    # NASK 実行コマンド (ユーザー環境固有 - 実行前にユーザーに nask.exe のパスを確認すること)
    # wine <nask.exeへのパス> input.asm output.bin output.lst
    # 例: wine nask.exe input.asm output.bin output.lst
    ```

2.  **Go リテラル生成**: 生成されたバイナリファイル (`output.bin`) を元に、Go の `[]byte` リテラル形式 (`[]byte{0x.., 0x.., ...}`) を生成する。**`hexdump` からの手動転記は間違いが多いため、必ず以下の Go プログラム (`generate_expected.go`) を使用して生成すること。**
    ```go
    // generate_expected.go (使い捨てスクリプトとして作成・実行)
    package main
    import ("fmt"; "os"; "strings")
    func main() {
        data, err := os.ReadFile("output.bin") // NASK が生成したバイナリ
        if err != nil { panic(err) }
        var builder strings.Builder
        builder.WriteString("[]byte{\n")
        for i, b := range data {
            if i%16 == 0 { builder.WriteString("\t") } // インデント調整
            builder.WriteString(fmt.Sprintf("0x%02x,", b))
            if (i+1)%16 == 0 { builder.WriteString(" //\n") } else { builder.WriteString(" ") }
        }
        if len(data)%16 != 0 { builder.WriteString("\n") }
        builder.WriteString("}") // 末尾のタブ削除
        fmt.Println(builder.String())
    }
    ```
    ```bash
    # スクリプトを一時ファイル (例: gen.go) に保存してから実行
    go run gen.go > expected_literal.txt
    ```

3.  **テストコード更新**: 生成された `expected_literal.txt` の内容 (`[]byte{...}`) をコピーし、該当するテストケース (`*_test.go`) 内の `expected` 変数に貼り付ける。

このプロセスにより、`expected` データが常に基準アセンブラの最新の正しい出力を反映するようになり、テストの信頼性が向上する。

## e2e テスト作成プロセスの標準化案 (2025/04/05 更新)

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
    1.  **アセンブリコード定義:** テスト対象のアセンブリコード（通常は外部の `.asm` ファイル）の内容を、必要に応じて文字コード変換 (`cat file.asm | nkf --ic=CP932 --oc=UTF-8` など) を行った上で、Go の文字列リテラルとしてテストコード内に定義する (`const naskStatements = \`...\``)。
    2.  **gosk アセンブル実行:**
        *   `gen.Parse` で上記のアセンブリコード文字列をパースする。
        *   `frontend.Exec` でアセンブルを実行し、結果を一時ファイルに出力する。
        *   `ReadFileAsBytes` (`test/readbin.go` で定義) を使用して一時ファイルの内容を `actual []byte` として読み込む。
    3.  **期待値バイナリ定義:**
        *   「テスト期待値 (expected) の生成方法」セクションの手順に従い、NASK の出力から生成した `[]byte` リテラルを `expected` 変数に直接代入する。(`defineHEX` DSL は使用しない)。
    4.  **結果比較:**
        *   `github.com/google/go-cmp/cmp` パッケージの `cmp.Diff(expected, actual)` を使用して、期待値と実際の結果を比較する。
        *   差分が存在する場合 (`diff != ""`) は、`DumpDiff` ヘルパー関数 (`test/diff.go` で定義) を使用して詳細な16進ダンプ形式の差分をログに出力し、`t.Fail()` でテストを失敗させる。
        *   バイト列の長さも別途 `len(expected) != len(actual)` で比較する。
- **使用されているヘルパー関数:**
    - `ReadFileAsBytes(filePath string) ([]byte, error)` (`test/readbin.go`): ファイル内容をバイト列として読み込む。
    - `DumpDiff(expected, actual []byte, color bool) string` (`test/diff.go`): バイト列の差分を16進ダンプ形式で表示。
- **期待値 (`expected`) の管理:**
    - 現在は `[]byte` リテラルとしてテストコード内に直接記述されている。
    - 将来的には、期待値バイナリを `.golden` ファイルとして分離し、テスト実行時にファイルを比較する方法も検討可能。

### 4. テスト実行手順

- **テスト実行:**
    - `go test` コマンドを使用し、`-run` フラグで対象のテストスイートやテストメソッドを指定する。
    - 全ての e2e テストを実行: `go test ./test -run ^TestDay` (スイート名が `TestDayXXSuite` であることを仮定)
    - 特定の日のテストスイートを実行: `go test ./test -run TestDayXXSuite` (例: `go test ./test -run TestDay05Suite`)
    - 特定のテストケースを実行: `go test ./test -run TestDayXXSuite/TestHaribYY<suffix>` (例: `go test ./test -run TestDay05Suite/TestHarib02i`)
- **デバッグ手順:**
    1. 失敗したテストケースを特定する (`go test -run ...` で再実行)。
    2. NASK と gosk でそれぞれバイナリを生成する (NASK: `wine ...`, gosk: テストコード内の `frontend.Exec` 部分を一時的に変更してファイル出力するなど)。
    3. 生成されたバイナリを `hexdump -C` や `diff <(hexdump -C nask.bin) <(hexdump -C gosk.bin)` などで差分を詳細に比較する。
    4. `gosk` のデバッグログ (`-v` オプションなど、必要なら追加) を確認する。
    5. 関連する `internal/` パッケージ (pass1, codegen, asmdb など) のコードをデバッグする。

この標準化案は、実装を進めながら改善していく。

## PEGパーサーのデバッグとメモリオペランドのパースに関する注意点 (2025/05/21)

- **問題:** `[EBX*4]`のような「インデックスレジスタとスケールファクタのみ」のメモリオペランドが`pkg/ng_operand/operand_grammar.peg`で正しくパースできない問題が発生した。
- **原因:**
    - 従来の`IndexScaleDisp`ルールがディスプレースメントを必須としていたため、ディスプレースメントがない形式にマッチしなかった。
    - `IndexOnly`ルールはスケールファクタをオプションとしていたが、`*4`のような形式を正しくパースするアクションが不足していた。
    - `pigeon`が生成するGoコードにおいて、PEGファイルでオプションにした要素をアクションブロックで適切に`nil`チェックせずにアクセスしようとすると、`undefined`エラーが発生した。
- **解決策:**
    - `IndexScaleOnly`という新しいルールを導入し、インデックスレジスタとスケールファクタのみを持つメモリオペランドを明示的にパースするようにした。
    - `MemoryBody`のルールにおいて、`IndexScaleOnly`を既存のルール群に追加し、適切な順序で評価されるようにした。
    - `pigeon`が生成するGoコードのビルドエラーを回避するため、PEGアクションブロック内でオプション要素にアクセスする際は、必ず`nil`チェックを行うようにした（今回はルール分離で回避）。
- **教訓:**
    - PEGのルールは上から順にマッチを試みるため、より具体的なルールを先に、より一般的なルールを後に配置することが重要。
    - オプション要素を扱う際は、生成されるGoコードでの`nil`チェックの必要性を常に意識すること。
    - 複雑なパースルールは、より単純なサブルールに分割することで、デバッグとメンテナンスが容易になる。

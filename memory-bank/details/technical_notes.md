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

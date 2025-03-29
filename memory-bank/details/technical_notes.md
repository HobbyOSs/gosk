# 技術ノート

## 技術ノートテンプレート

### 概要
- 簡単な問題概要

### 調査の経緯
- 問題発生の経緯、調査のステップ

### 原因分析
- 詳細な原因分析、根本原因の特定

### 問題点
- 構造的な問題点、技術的な課題

### 解決策/設計提案
- 提案する解決策、改善設計

### 今後の課題/開発方針
- 今後の開発における課題、対策、方針

### 関連ファイル/情報
- 関連するファイルパス、参考情報へのリンクなど

---

## オペランド受け渡しフローと CodegenClient.Emit インターフェースの問題点 (2025/03/29)

### 調査の経緯

`test/day03_harib00i_test.go` のテスト実行時に、`MOV ECX, [EBX + 16]` のようなメモリオペランドを持つ命令で `Failed to parse operand string 'ECX[ EBX + 16 ]'` というエラーが発生し、それに伴い `Failed to find encoding` エラーも発生していました。

### 原因分析

1.  **`pass1` での Ocode 生成**: `internal/pass1/pass1_inst_*.go` や `internal/pass1/handlers.go` では、各命令のトークンを処理し、`CodegenClient.Emit` を呼び出して Ocode を生成していました。当初の実装では、`Emit` メソッドは単一の文字列 (`"MOV ECX,[ EBX + 16 ]"`) を受け取るシグネチャ (`Emit(string)`) でした。
2.  **`ocode_client` での Ocode 格納**: `internal/ocode_client/client.go` の `Emit` 実装は、`pass1` から受け取った単一文字列をパースして `ocode.Ocode` 構造体に格納していました。この実装において、オペランドは一旦カンマ区切りで分割され、`Operands` フィールド (`[]string`) に格納されます。
3.  **`codegen` での Ocode 処理**: `internal/codegen/x86gen_*.go` (例: `handleMOV`) では、`ocode.Ocode` の `Operands` フィールド (`[]string`) を受け取り、`asmdb.FindEncoding` を呼び出すために `operand.Operands` インターフェースを生成する必要があります。この`operand.Operands` の生成処理 (`operand.NewOperandFromString`) において、`codegen` は受け取ったオペランドの `[]string` を `strings.Join(operands, ",")` で**再度単一の文字列に結合**していました。**重要な点として、codegen では operand を `pkg/operand` パッケージのパーサーで再度パース処理を行っています。**
4.  **`pkg/operand` でのパースエラー**: `operand.NewOperandFromString` に渡された結合文字列 (`"ECX,[ EBX + 16 ]"`) は、`pkg/operand` のパーサー (`participle` ベース) が期待する形式 (カンマ区切りの完全な命令文字列) と完全には一致しません。`asmdb.FindEncoding` が内部で `OperandTypes()` を呼び出し、さらにその内部で `getInternalParsed()` がこの結合された文字列をパースしようとした結果、`Failed to parse operand string 'ECX[ EBX + 16 ]'` エラーが発生していました。

### 問題点

- **`CodegenClient.Emit` のインターフェース**: `Emit(string)` というシグネチャが、オペランド情報を構造化して渡す上で不適切でした。`pass1` でパースされたオペランド情報は、単一文字列にシリアライズされるべきではありませんでした。
- **`codegen` での再結合**: `codegen` 側でオペランドスライスを再度文字列に結合していたことが、`pkg/operand` パーサーのエラーを引き起こす直接的な原因でした。
- **モジュール間の結合度**: `pass1`, `ocode_client`, `codegen`, `pkg/operand` の間で、オペランド情報の受け渡し方法に関する暗黙的な依存関係があり、変更が困難になっていました。

### 試みた修正と中断

`CodegenClient.Emit` のシグネチャを `Emit(op string, operands []string)` に変更し、`pass1` から `codegen` までオペランドを `[]string` として渡すように修正を試みました。しかし、関連するファイルが多く、修正が広範囲に及び複雑化したため、ユーザー指示により中断しました。

### 今後の課題

- **オペランド受け渡し方法のリファクタリング**: `pass1` から `codegen` まで、オペランド情報をより構造化された形で (例: `[]operand.ParsedOperand` や専用の構造体) 受け渡すように、関連モジュール全体のリファクタリングが必要です。
- **`CodegenClient.Emit` インターフェース再設計**: オペランド情報を適切に渡せるようなインターフェースを再設計する必要があります。
- **`pkg/operand` と `asmdb` の連携改善**: `asmdb.FindEncoding` が `operand.Operands` インターフェース (単一文字列前提) に依存している点を解消し、より柔軟なオペランド情報 (例: `[]operand.OperandType`) を受け入れられるように改善する必要があります。

### あるべき設計 (提案)

#### 1. `CodegenClient.Emit` インターフェースの変更

`CodegenClient.Emit` のシグネチャを、より構造化されたオペランド情報を扱えるように変更します。

**変更案:** `Emit(op ocode.Ocode)`

- 引数を `string` 型から `ocode.Ocode` 型に変更します。
- `ocode.Ocode` オブジェクトは、opcode とパース済みのオペランド情報 (`[]operand.ParsedOperand`) を保持します。
- `pass1` は、オペランドを `pkg/operand` パーサーでパースし、`ocode.Ocode` オブジェクトを生成して `Emit` を呼び出すように変更します。
- `ocode_client` は、受け取った `ocode.Ocode` オブジェクトをそのまま `Ocodes` リストに格納します。
- `codegen` は、`ocode.Ocode` オブジェクトから opcode とオペランド情報を取得して処理します。

**メリット:**

- オペランド情報を構造化して受け渡せるため、型安全性が向上します。
- `pass1`, `ocode_client`, `codegen` の間のインターフェースが明確になり、疎結合な設計になります。
- `ocode_client` の `Emit` メソッドがシンプルになり、保守性が向上します。
- `codegen` でのオペランド再パース処理が不要になり、効率的な処理が可能になります。

#### 2. `pkg/operand` パッケージの修正

`pkg/operand` パッケージの `participle` パーサー定義 (`operandLexer`) を修正し、複雑なメモリオペランド (`ECX[ EBX + 16 ]` など) を正しくパースできるようにします。

- `IndirectMem` ルールの正規表現を見直し、より複雑なメモリオペランドに対応できるように修正します。
- 必要に応じて、レキサー定義の順序やルールを追加します。
- `operand_impl_test.go` のテストケースを拡充し、修正後のパーサーを検証します。

#### 3. `asmdb.FindEncoding` の改善

`asmdb.FindEncoding` が `operand.Operands` インターフェース (単一文字列前提) に依存している点を解消し、より柔軟なオペランド情報 (例: `[]operand.OperandType` や `[]operand.ParsedOperand`) を受け入れられるように改善します。

- `asmdb.FindEncoding` の引数を `operand.Operands` から `[]operand.OperandType` または `[]operand.ParsedOperand` に変更します。
- `asmdb` が `pkg/operand` に依存しないように、インターフェースを再設計します。

#### 4. オペランド受け渡しフローの改善

`pass1` から `codegen` まで、オペランド情報を `ocode.Ocode` オブジェクトを介して構造化された形で受け渡すように、関連モジュール全体をリファクタリングします。

- `pass1` でオペランドをパースし、`ocode.Ocode` オブジェクトを生成する処理を実装します。
- `ocode_client` で `ocode.Ocode` オブジェクトを格納し、`codegen` に渡す処理を実装します。
- `codegen` で `ocode.Ocode` オブジェクトからオペランド情報を取得し、機械語生成処理を行うように修正します。

### 今後の開発方針

1.  **`CodegenClient.Emit(ocode.Ocode)` インターフェースへの変更**: まず、`Emit` インターフェースの変更を行い、オペランド情報を構造化して受け渡せるようにします。
2.  **`pkg/operand` パーサーの修正**: 次に、`pkg/operand` パッケージのパーサーを修正し、複雑なメモリオペランドを正しくパースできるようにします。
3.  **`asmdb.FindEncoding` の改善**: その後、`asmdb.FindEncoding` を改善し、より柔軟なオペランド情報を受け入れられるようにします。
4.  **オペランド受け渡しフローのリファクタリング**: 最後に、関連モジュール全体をリファクタリングし、オペランド受け渡しフローを改善します。

## jq を用いた JSON データ調査 (2025/03/29)

### json-x86-64/x86_64.json の構造

`json-x86-64/x86_64.json` は、x86-64 アーキテクチャの命令セットに関する詳細な情報を含む JSON ファイルです。このファイルは、命令の名前、オペランド、エンコーディング、属性などの情報を提供します。

ファイルは、トップレベルで `instructions` というキーを持つオブジェクトを含み、`instructions` は命令名 (例: "ADD", "MOV", "IMUL" など) をキーとするオブジェクトの配列です。

各命令オブジェクトは、以下のプロパティを持つ `forms` 配列を含みます。

- `forms`: 命令のエンコーディング形式の配列。各要素はエンコーディング形式に関する情報を持つオブジェクトです。
  - `encodings`: エンコーディングの詳細情報の配列。通常、最初の要素 (`encodings[0]`) が主要なエンコーディング情報です。
    - `opcode`: オペコードに関する情報を持つオブジェクト
      - `byte`: オペコードのバイト表現 (16進数文字列)
    - `operands`: オペランドに関する情報の配列
      - `type`: オペランドのタイプ (例: "r8", "r16", "r32", "r64", "m8", "m16", "m32", "m64" など)

### jq を用いた JSON データ調査

`jq` は、JSON データを処理するための軽量で柔軟なコマンドライン JSON プロセッサです。`jq` を使用することで、`json-x86-64/x86_64.json` から特定の命令やエンコーディングに関する情報を効率的に抽出できます。

**jq コマンド例:**

```bash
cat pkg/asmdb/json-x86-64/x86_64.json | jq -r '.instructions["IMUL"].forms[] | [.encodings[0].opcode.byte, .operands[0].type, .operands[1].type ] | @tsv'
```

このコマンドは、`IMUL` 命令のエンコーディング形式を抽出し、オペコードとオペランドタイプを TSV 形式で出力します。

**出力例:**

```tsv
F6	r8
F7	r16
F7	r32
F7	r64
... (省略) ...
6B	r64	r64
69	r64	r64
```

この出力は、`IMUL` 命令の様々なエンコーディング形式と、対応するオペコード、オペランドタイプを示しています。

### 今後の活用

`jq` を用いることで、`json-x86-64/x86_64.json` のようなJSONデータを効率的に調査し、必要な情報を抽出できることがわかりました。今後、命令実装に必要な情報をJSONファイルから取得する際に、`jq` を積極的に活用します。

## Goリファクタタスクでの反省点 (2025/03/29)

### 発生した問題点

- Goの構文知識不足により、コンパイラエラーが発生した。
- エラーメッセージを十分に理解できず、修正に時間がかかった。
- 単体テストによる動作確認が不足していた。
- 問題解決を**自分自身で**試みる前に、ユーザーフィードバックに依存してしまった。

### 次回の改善策

- **Goのコード例を積極的に参照し、Goのコーディング規約を学習する。**エラーメッセージの分析能力を高める。
- テスト優先の開発プロセスを導入し、より**小さな**単位での動作確認を徹底する。
- **因果推論や反実仮想が必要となる場面（例: 「この変更が他の箇所に影響を与えないか？」）など、不明な点は積極的に `ask_followup_question` ツールを使って質問し、**ユーザーフィードバックを**効果的に活用する**。
- 過去の失敗事例を分析し、**同様の**問題の再発防止に努める。
- 関数をより小さく分割し、**リファクタリングしやすい**構造を心がける。

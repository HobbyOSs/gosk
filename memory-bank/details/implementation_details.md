## アセンブラ命令実装のルーチン

### テストスイートの構造
1. **テストの段階的実装**
   - day01からday20までのテストケースが存在
   - 各dayで新しいアセンブラ命令が追加される
   - day20以降は新しい命令の追加は少ない
   - 各dayのテストは独立して実行可能

2. **テストケースの内容**
   - OSイメージの生成テスト
   - アセンブリコードからバイナリイメージを生成
   - 期待値との比較検証
   - バイナリレベルでの正確性確認

### Pass1での命令実装
1. **機械語サイズの計算**
   - `pkg/asmdb`を使用して命令のサイズを計算
   - プレフィックスバイトのサイズを考慮
   - オフセットサイズを計算に含める

2. **Ocodeの生成**
   - `env.Client.Emit`を使用してOcodeを出力
   - オペランドの正確な解析と変換

### Ocode実装

1.  **命令の定義**
    -   `pkg/ocode/ocode.go`にiota定数として命令を追加
    -   命令の種類に応じた適切な名前付け
2.  **機械語生成の実装**
    -   `internal/codegen/x86gen.go`に命令の処理を追加
    -   asmdbを使用して正確な機械語を生成
    -   ModR/Mの適切な生成
    -   機械語生成の標準パターン：

        ```go
        // 1. 空のスライスから開始
        machineCode := make([]byte, 0)

        // 2. プレフィックスの追加（必要な場合）
        if ops.Require66h() {
            machineCode = append(machineCode, 0x66)
        }

        // 3. オペコードの追加
        opcodeByte, err := strconv.ParseUint(encoding.Opcode.Byte, 16, 8)
        if err != nil {
            return nil, fmt.Errorf("failed to parse opcode byte: %v", err)
        }
        machineCode = append(machineCode, byte(opcodeByte))

        // 4. ModR/Mの追加（必要な場合）
        if encoding.ModRM != nil {
            modrm := generateModRMByte(encoding.ModRM)
            machineCode = append(machineCode, modrm)
        }

        // 5. 即値の追加（必要な場合）
        if encoding.Immediate != nil {
            if imm, err := getImmediateValue(operands[1], encoding.Immediate.Size); err == nil {
                machineCode = append(machineCode, imm...)
            }
        }
        ```

### 実装手順

1.  **Pass1の実装**

    ```go
    // handlers.goに関数を登録
    opcodeEvalFns["NEW_INST"] = processNEW_INST

    // 命令処理関数の実装
    func processNEW_INST(env *Pass1, tokens []*token.ParseToken) {
        // 1. オペランドの解析
        // 2. asmdbでサイズ計算
    	// 3. ラベルの仮登録: SymTable にラベルを登録し、仮アドレスを割り当てる。
        // 4. Ocodeの生成: env.Client.Emit を使用してOcodeを出力。ジャンプ先アドレスはプレースホルダーとする。
    }
    ```
2.  **Ocodeの定義**

    ```go
    // ocode.goに追加
    const (
    	// 既存の定義
    	OpExisting OcodeKind = iota
    	// 新規命令
    	OpNEW_INST
    )
    ```
3.  **Pass2でのラベル解決とプレースホルダー置換**

    -   **`SymTable` に登録されたラベルのアドレスを確定する。**
    -   **Ocode 内のプレースホルダーを実際のアドレスに置き換える。**
4.  **機械語生成**

    ```go
    // x86gen.goに実装
    func processOcode(oc ocode.Ocode, ctx *CodeGenContext) ([]byte, error) {
    	switch oc.Kind {
    	case ocode.OpNEW_INST:
    		return handleNEW_INST(oc.Operands), nil
    	}
    }
    ```


### 実装時の注意点

1.  **asmdbの活用**
    -   命令の正確な機械語表現の取得
    -   プレフィックスバイトの適切な処理
    -   オペランドサイズの正確な計算
    -   複数の機械語パターンがある場合は最小サイズのものを選択
    -   FindEncodingを使用して最小サイズのエンコーディングを一意に決定
    -   セグメントレジスタ（sreg）をr16として扱う
    -   JSONにない命令パターンのフォールバック処理
    -   オペランドタイプの互換性チェック（例：sreg→r16の変換）
2.  **テスト駆動開発**
    -   各命令の基本的なテストケース作成
    -   エッジケースの考慮
    -   他の命令との相互作用の確認
3.  **ドキュメント化**
    -   実装した命令の仕様と制限事項の記録
    -   テストケースの説明
    -   既知の問題点の記録
4.  **pkg/asmdbに存在しない命令への対応**
    -   特殊なOS命令など、json-x86-64/x86_64.jsonに存在しない命令が発生する場合がある
    -   このような場合は、フォールバック実装が必要
    -   フォールバック実装を行う前に、必ず人間に指示を仰ぐこと
5.  **開発ルーチン**
    -   実装完了後は必ず`make test`を実行し、全体のテストを確認
    -   テストのスキップは明示的な許可がある場合のみ実施
    -   テストエラーが発生した場合は、修正を行うか、スキップの許可を得る

## オペランドパーサー移行 (participle -> pigeon) (2025/03/29)

### 背景
- 既存のオペランドパーサー (`pkg/operand`) は `participle` ライブラリを使用していた。
- メモリオペランドの複雑な形式や、レジスタ名とラベル名の曖昧性などにより、パーサーの実装が複雑化し、メンテナンス性や拡張性に課題があった。
- 特に `test/day03_harib00i_test.go` で発生していたエンコーディングエラーの原因の一つとして、オペランドパースの不正確さが疑われていた。

### 移行方針
- より厳密で表現力の高い PEG (Parsing Expression Grammar) を採用し、パーサジェネレータ `pigeon` を使用して新しいパーサーを実装する。
- 新しいパーサーは `pkg/ng_operand` パッケージとして開発を進め、最終的に既存の `pkg/operand` を置き換える。
- 既存の `pkg/operand` のインターフェース (`Operands`) とテストコードを流用し、互換性を保ちながら段階的に移行する。

### 実装状況 (2025/03/29現在)
1.  **peg文法定義 (`pkg/ng_operand/operand_grammar.peg`)**:
    -   オペランド文字列 (`OperandString`) をパースするルールを定義。
    -   `Operand` ルールで `MemoryAddress`, `Register`, `Immediate`, `Label`, `SegmentRegister` を選択。
    -   `MemoryAddress` ルールでデータ型 (`DataType`)、ジャンプタイプ (`JumpType`)、セグメントオーバーライド (`SegmentRegisterName :`)、メモリ本体 (`MemoryBody`) をパース。
    -   `MemoryBody` ルールでベースレジスタ、インデックスレジスタ、スケール、ディスプレースメントの組み合わせをパース。
    -   各ルールのアクションでパース結果を `ParsedOperandPeg` および `MemoryInfo` 構造体に格納。
2.  **型定義 (`pkg/ng_operand/operand_types.go`)**:
    -   パース結果を格納する `ParsedOperandPeg`, `MemoryInfo` 構造体を定義。
    -   オペランドの種類を表す `OperandType` 定義 (既存のものを流用)。
    -   pegアクションで使用するヘルパー関数 (`getRegisterType`, `getImmediateSizeType` など) を定義。
3.  **パーサー生成 (`pkg/ng_operand/generate.go`, `operand_grammar.go`)**:
    -   `go:generate` ディレクティブを使用して `pigeon` コマンドを実行し、`operand_grammar.go` を自動生成。
4.  **ラッパー関数 (`pkg/ng_operand/parser.go`)**:
    -   `ParseOperandString(text string) (*ParsedOperandPeg, error)` を実装し、生成された `Parse` 関数を呼び出す。
5.  **インターフェース実装 (`pkg/ng_operand/operand_impl.go`)**:
    -   既存の `Operands` インターフェースを実装する `OperandPegImpl` 構造体を定義。
    -   `FromString` コンストラクタと、`With*`, `Get*` などの基本的なメソッドを実装。
    -   `OperandTypes`, `Require66h`, `Require67h` の基本的なロジックを実装。
6.  **テスト (`pkg/ng_operand/operand_test.go`)**:
    -   既存の `pkg/operand/operand_test.go` をコピーし、パッケージ名を変更。
    -   `TestRequire66h`, `TestRequire67h` を有効化し、`pkg/ng_operand` の実装でテストが通ることを確認済み。

### 今後の課題
-   `OperandTypes` のサイズ解決ロジックの実装。
-   `CalcOffsetByteSize`, `DetectImmediateSize` の実装。
-   `Require66h`, `Require67h` の詳細化。
-   複数オペランド (カンマ区切り) への対応 (peg文法と `InternalStrings` メソッドの修正)。
-   既存の `pkg/operand` 利用箇所の置き換え。

### 過去の作業

-   オペコード生成処理の改善 (`internal/codegen/x86gen_utils.go`)
    -   `ResolveOpcode` 関数と `GetRegisterNumber` 関数を追加
    -   `x86gen_mov.go` と `x86gen_arithmetic.go` でこれらの関数を使用するように修正
-   `GenerateX86INT` 関数を `handleINT` にリネーム
-   `GenerateModRM` 関数と `GetRegisterNumber` 関数を修正

### JMP系命令 (Jcc命令) の実装 (2025/03/19)

-   `pkg/ocode/ocode.go` に JMP系命令 (JA, JAE, JB, JBE, JC, JE, JG, JGE, JL, JLE, JNA, JNAE, JNB, JNBE, JNC, JNE, JNG, JNGE, JNL, JNLE, JNO, JNP, JNS, JNZ, JO, JP, JPE, JPO, JS, JZ) を追加
-   `internal/codegen/x86gen.go` の `processOcode` 関数に JMP系命令のケースを追加し、`handleJcc` 関数を呼び出すように修正
-   `internal/codegen/x86gen_jmp.go` に `handleJcc` 関数を実装し、各JMP系命令の機械語生成処理を実装
    -   `generateJMPCode` 関数は JMP 命令のみに使用するように変更
    -   各JMP系命令のオペコードは、`handleJcc` 関数内で定義
    -   オフセットは `rel8` (1バイト) のみ対応

### JE命令の実装 (2025/03/14)

-   `internal/codegen/x86gen_jmp.go` に `generateJMPCode` 関数を追加し、JMP命令とJE命令の共通処理を実装
-   `handleJMP` 関数と `handleJE` 関数から `generateJMPCode` 関数を呼び出すように修正
-   `internal/codegen/x86gen.go` の `processOcode` 関数に `ocode.OpJE` のcaseを追加し、`handleJE` 関数を呼び出すように修正
-   `pkg/ocode/ocode.go` に `OpJE` を追加し、`enumer` を実行して `OcodeKind` を再生成

### OUT命令の実装 (2025/03/26)

-   `internal/pass1/pass1_inst_out.go` にOUT命令の処理を追加
-   `internal/pass1/handlers.go` に `processOUT` 関数を登録
-   `internal/codegen/x86gen_out.go` を作成し、OUT命令のエンコード処理を実装
    -   `handleOUT` 関数を実装
    -   現状は `OUT imm8, AL` のみに対応
-   `internal/codegen/x86gen_test.go` にOUT命令のテストケースを追加
-   `internal/codegen/x86gen.go` に `processOcode` 関数に `ocode.OpOUT` のケースを追加

[過去の作業履歴はこちら](./archives/implementation_details_archive_20250313.md)

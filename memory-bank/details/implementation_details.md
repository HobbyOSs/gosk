## アセンブラ命令実装のルーチン

### 背景

#### テストスイートの構造
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

#### Ocode実装の概要

### Pass1とng_operand間の設計的境界 (抽象概念)

Pass1とng_operandは、アセンブラ処理における責務を断固として分担しています。
責務分担に違反する実装は、断じて避けるべきです。

- **Pass1の責務**:
    - AST (抽象構文木) の評価: マクロ展開、定数畳み込みなど
    - 中間表現 (ocode) の生成: 機械語生成に必要な情報を抽出・構造化
    - 機械語サイズの計算

- **ng_operandの責務**:
    - Pass1で還元・Serializeされたオペランドの文字列表現を解析し、`ng_operand.Operands`を生成する
    - 生成された `ng_operand.Operands` は、asmdbを使用した機械語サイズ計算やエンコーディングに使用される
    - オペランドの文字列表現の提供 (デバッグ、エラーメッセージ用など)

この断固たる責務分担により、Pass1の処理が複雑化することを防ぎ、機械語生成に必要な情報のみをng_operandに確実に渡します。
特に、Pass1のオペコード処理 (例: `processMOV`) では、以下の設計原則を厳守します。

#### オペコード処理 (processMOVなど) の設計

`processMOV` などのオペコード処理関数は、以下の引数を受け取ります。

- `env *Pass1`: Pass1の実行環境 (シンボルテーブル、マクロテーブルなど)
- `operands []ast.Exp`: 命令のオペランド (ASTノードの配列)

`operands` の各要素は `ast.Exp` インターフェースを実装しており、`Serialize` (または `TokenLiteral` など、文字列表現を取得するメソッド) を持つことを想定しています。
Pass1では、`operands` を評価・還元した後、`Serialize` メソッド (または `TokenLiteral`) で文字列表現を取得し、**ng_operandには文字列表現の配列としてオペランド情報を厳格に渡します**。

**決して責務を曖昧にせず、重要な制約として、`Serialize` される `[]ast.Exp` は、マクロや計算式が完全に還元されている必要**があります。
もし `Serialize` 時点でマクロや計算式が残存している場合、ng_operand (または機械語生成以降の後続処理) で、責務分担に違反するASTTraverseの再実行が引き起こされる危険性があります。
これは、Pass1とng_operandの責務分担を根底から覆し、コードの複雑化と予期せぬバグの発生を招く行為です。

#### Serializeされた []ast.Exp の制約 (サンプルコード)

断じて責務を侵犯しない実装を行うために、例として以下のマクロ定義とMOV命令を再度確認しましょう。

```assembly
BUFFER_SIZE EQU 1024
MOV AX, BUFFER_SIZE + 5
```

Pass1の `processMOV` 関数に渡されるオペランドは、ASTノード配列 `[ast.RegisterExp{Name: "AX"}, ast.AddExp{Lhs: ast.IdentifierExp{Name: "BUFFER_SIZE"}, Rhs: ast.NumberExp{Value: 5}}]` です。
`BUFFER_SIZE` はマクロ定義であるため、`ast.IdentifierExp` はその値を参照し、`BUFFER_SIZE + 5` は計算式なので、Pass1の段階で評価し、定数 `1029` に還元します。

したがって、`processMOV` 関数内で `operands` をSerializeする際には、**マクロ `BUFFER_SIZE` は必ず展開され、計算式 `BUFFER_SIZE + 5` は評価済みの `ast.NumberExp{Value: 1029}` に還元されていることを保証**しなければなりません。
責務違反の実装は、絶対に許容されません。

```go
func processMOV(env *Pass1, operands []ast.Exp) {
	if len(operands) != 2 {
		log.Printf("Error: MOV instruction requires exactly two operands.")
		return
	}

	// 1. オペランドを評価・還元 (マクロ展開、定数畳み込みなど)
	//    Eval() メソッドを呼び出し、評価済みの AST ノードを取得する
	evaluatedOperands := make([]ast.Exp, len(operands))
	for i, operand := range operands {
		evaluatedOperands[i] = operand.Eval(env) // Eval() で評価・還元
	}

	// 2. Serialize されたオペランド (文字列表現) を取得
	//    TokenLiteral() メソッド (または Serialize() メソッド) を呼び出し、
	//    評価済み AST ノードから文字列表現を取得する
	serializedOperands := make([]string, len(evaluatedOperands))
	for i, operand := range evaluatedOperands {
		serializedOperands[i] = operand.TokenLiteral() // TokenLiteral() で文字列表現を取得
	}

	// 3. ng_operand を生成 (文字列表現から)
	//    ng_operand.FromString() 関数を呼び出し、
	//    Serialize された文字列表現から ng_operand.Operands を生成する
	operandString := serializedOperands[0] + "," + serializedOperands[1] // オペランドをカンマ区切りで結合
	ngOperands, err := ng_operand.FromString(operandString)             // ng_operand.FromString() で ng_operand.Operands を生成
	if err != nil {
		log.Printf("Error creating operands from string '%s' in MOV: %v", operandString, err)
		return // エラーが発生したら処理を中断 (責務違反を検知)
	}

	// 4. BitMode, ForceRelAsImm を設定 (ng_operand の設定)
	//    必要に応じて ng_operand の設定を行う (例: BitMode, ForceRelAsImm)
	ngOperands = ngOperands.WithBitMode(env.BitMode).
		WithForceRelAsImm(true) // Force relative symbols (like labels) to be treated as immediates for size calculation

	// 5. 機械語サイズの計算 (ng_operand を使用)
	//    env.AsmDB.FindMinOutputSize() 関数を呼び出し、
	//    ng_operand.Operands から機械語サイズを計算する
	size, err := env.AsmDB.FindMinOutputSize("MOV", ngOperands)
	if err != nil {
		// エラーログ出力 (詳細なエラー情報を出力)
		log.Printf("Error finding min output size for MOV (op1: '%s', op2: '%s'): %v", serializedOperands[0], serializedOperands[1], err)
		return // エラーが発生したら処理を中断 (責務違反を検知)
	}
	env.LOC += int32(size) // LOC (Location Counter) を更新

	// 6. Ocodeの生成とEmit (ng_operand.Serialize() を使用)
	//    env.Client.Emit() 関数を呼び出し、Ocode を生成・出力する
	//    ngOperands.Serialize() で ng_operand.Operands を文字列表現に変換する
	env.Client.Emit(fmt.Sprintf("MOV %s", ngOperands.Serialize())) // ngOperands.Serialize() で文字列表現を取得
}
```

#### processMOV サンプルコード (擬似コード)


### 実装レベル

#### Pass1での命令実装
1. **機械語サイズの計算**
   - `pkg/asmdb`を使用して命令のサイズを計算
   - プレフィックスバイトのサイズを考慮
   - オフセットサイズを計算に含める

#### Ocodeの生成
   - `env.Client.Emit`を使用してOcodeを出力
   - オペランドの正確な解析と変換

#### Ocode実装

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

#### 実装手順

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


#### 実装時の注意点

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

### オペランドパーサー移行 (participle -> pigeon)

#### 背景
- 既存のオペランドパーサー (`pkg/operand`) は `participle` ライブラリを使用していた。
- メモリオペランドの複雑な形式や、レジスタ名とラベル名の曖昧性などにより、パーサーの実装が複雑化し、メンテナンス性や拡張性に課題があった。
- 特に `test/day03_harib00i_test.go` で発生していたエンコーディングエラーの原因の一つとして、オペランドパースの不正確さが疑われていた。

#### 移行方針
- より厳密で表現力の高い PEG (Parsing Expression Grammar) を採用し、パーサジェネレータ `pigeon` を使用して新しいパーサーを実装する。
- 新しいパーサーは `pkg/ng_operand` パッケージとして開発を進め、最終的に既存の `pkg/operand` を置き換える。
- 既存の `pkg/operand` のインターフェース (`Operands`) とテストコードを流用し、互換性を保ちながら段階的に移行する。

#### 実装状況 (2025/03/29現在)
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
    -   `pkg/ng_operand/operand_impl.go`**:
    -   既存の `Operands` インターフェースを実装する `OperandPegImpl` 構造体を定義。
    -   `FromString` コンストラクタと、`With*`, `Get*` などの基本的なメソッドを実装。
    -   `OperandTypes`, `Require66h`, `Require67h` の基本的なロジックを実装。
6.  **テスト (`pkg/ng_operand/operand_test.go`)**:
    -   既存の `pkg/operand/operand_test.go` をコピーし、パッケージ名を変更。
    -   `TestRequire66h`, `TestRequire67h` を有効化し、`pkg/ng_operand` の実装でテストが通ることを確認済み。

#### 今後の課題
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

-   JMP系命令 (Jcc命令) の実装 (2025/03/19)
-   JE命令の実装 (2025/03/14)
-   OUT命令の実装 (2025/03/26)
-   [過去の作業履歴はこちら](./archives/implementation_details_archive_20250313.md)

### RESB 命令の処理 (2025/04/01)

`RESB` 命令は、指定されたバイト数だけ領域を予約する擬似命令です。オペランドには式 (例: `0x7dfe - $`) を指定できます。

- **Pass1 (`internal/pass1`)**:
    - `TraverseAST` は `RESB` 命令 (`ast.MnemonicStmt`) のオペランド式を `Eval` メソッドで評価します。
    - オペランド内の `$` シンボルは `ast.IdentFactor.Eval` によって現在の `LOC` (ロケーションカウンタ) の値 (`env.GetLOC()`) に解決されます。
    - 式全体 (例: `0x7dfe - $`) が評価され、予約すべきバイト数 (size) が `ast.NumberExp` として計算されます。
    - `processRESB` ハンドラ (`pass1_inst_pseudo.go`) は、この評価済みの `size` を受け取り、`env.LOC += size` を実行して LOC を更新します。
    - 最後に `env.Client.Emit(fmt.Sprintf("RESB %d", size))` を呼び出し、計算されたサイズを Ocode として Codegen に渡します。
- **Codegen (`internal/codegen`)**:
    - `handleRESB` 関数 (`x86gen_pseudo.go`) は、Pass1 から渡された Ocode (`RESB <size>`) を受け取ります。
    - オペランド (`<size>`) を数値としてパースします。
    - `make([]byte, size)` を実行し、指定されたバイト数分の 0x00 で埋められたバイトスライスを生成して返します。

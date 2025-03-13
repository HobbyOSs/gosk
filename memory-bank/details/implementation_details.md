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

1. **asmdbの活用**
   - 命令の正確な機械語表現の取得
   - プレフィックスバイトの適切な処理
   - オペランドサイズの正確な計算
   - 複数の機械語パターンがある場合は最小サイズのものを選択
   - FindEncodingを使用して最小サイズのエンコーディングを一意に決定
   - セグメントレジスタ（sreg）をr16として扱う
   - JSONにない命令パターンのフォールバック処理
   - オペランドタイプの互換性チェック（例：sreg→r16の変換）
2. **テスト駆動開発**
   - 各命令の基本的なテストケース作成
   - エッジケースの考慮
   - 他の命令との相互作用の確認
3. **ドキュメント化**
   - 実装した命令の仕様と制限事項の記録
   - テストケースの説明
   - 既知の問題点の記録
4. **pkg/asmdbに存在しない命令への対応**
   - 特殊なOS命令など、json-x86-64/x86_64.jsonに存在しない命令が発生する場合がある
   - このような場合は、フォールバック実装が必要
   - フォールバック実装を行う前に、必ず人間に指示を仰ぐこと
5. **開発ルーチン**
   - 実装完了後は必ず`make test`を実行し、全体のテストを確認
   - テストのスキップは明示的な許可がある場合のみ実施
   - テストエラーが発生した場合は、修正を行うか、スキップの許可を得る

### 過去の作業

- オペコード生成処理の改善 (`internal/codegen/x86gen_utils.go`)
  - `ResolveOpcode` 関数と `GetRegisterNumber` 関数を追加
  - `x86gen_mov.go` と `x86gen_arithmetic.go` でこれらの関数を使用するように修正
- `GenerateX86INT` 関数を `handleINT` にリネーム
- `GenerateModRM` 関数と `GetRegisterNumber` 関数を修正

### JE命令の実装 (2025/03/14)
- `internal/codegen/x86gen_jmp.go` に `generateJMPCode` 関数を追加し、JMP命令とJE命令の共通処理を実装
- `handleJMP` 関数と `handleJE` 関数から `generateJMPCode` 関数を呼び出すように修正
- `internal/codegen/x86gen.go` の `processOcode` 関数に `ocode.OpJE` のcaseを追加し、`handleJE` 関数を呼び出すように修正
- `pkg/ocode/ocode.go` に `OpJE` を追加し、`enumer` を実行して `OcodeKind` を再生成

[過去の作業履歴はこちら](./archives/implementation_details_archive_20250313.md)

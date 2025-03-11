# Active Context

## 現在の作業の焦点

### オペランド処理の改善
- [x] `pkg/operand` パッケージの `OperandImpl` 構造体に `ForceImm8` フィールドを追加
- [x] `NewOperandFromString` 関数を修正し、`ForceImm8` フィールドを初期化 (デフォルト: false)
- [x] `WithForceImm8` メソッドを追加し、`ForceImm8` フィールドを設定可能に (レシーバを直接変更する実装)
- [x] `WithBitMode` メソッドを修正 (レシーバを直接変更する実装)
- [x] `Operands` インターフェースに `WithForceImm8` メソッドを追加
- [x] `OperandTypes` メソッドを修正し、`ForceImm8` フラグが true の場合は即値のタイプを `CodeIMM8` に設定
- [x] `TestBaseOperand_OperandType` 関数に、`ForceImm8` フラグをテストするための新しいテストケースを追加
- [x] 既存の `"Immediate Value", "SI,1", []OperandType{CodeR16, CodeIMM8}}` テストケースの名前を `"Immediate Value force imm8"` に変更し、`ForceImm8` フィールドを true に設定
- [x] `resolveOperandSizes` 関数をレシーバメソッドに変更し、`ForceImm8` フラグを考慮するように修正

## Day02実装計画

- [x] 1. 基本命令の実装
    - [ ] MOV命令
        - [ ] レジスタ間転送
        - [ ] 即値のロード
        - [ ] セグメントレジスタの設定
    - [x] ADD命令
        - [x] 即値加算の実装
        - [ ] フラグ更新の処理
    - [ ] CMP命令
        - [ ] 比較演算の実装
        - [ ] フラグ設定の処理
- [x] 2. 制御フロー命令の実装
    - [x] JMP命令 (JMP rel8)
        - [x] Pass1: processJMP 関数実装 (internal/pass1/pass1_inst_jmp.go)
        - [x] Ocode: OpJMP 定義 (pkg/ocode/ocode.go)
        - [x] 機械語生成: handleJMP 関数実装 (internal/codegen/x86gen_jmp.go)
        - [ ] 相対アドレス計算
        - [ ] ジャンプ先ラベル解決
- [ ] 3. メモリ操作命令の実装
    - [ ] メモリアドレッシングモードの実装
    - [ ] レジスタ-メモリ間のデータ転送
    - [ ] ModR/Mバイトの生成
- [ ] 4. 制御フロー命令の実装 (続き)
    - [ ] JE命令
        - [ ] フラグに基づく分岐
        - [ ] オフセット計算
- [x] 5. システム命令の実装
    - [x] INT命令
        - [x] 割り込み番号の処理
        - [x] BIOS呼び出しの対応
    - [x] HLT命令
        - [x] CPU停止状態の生成
- [ ] 6. テスト有効化と検証
    - [ ] 各命令のユニットテスト実行
    - [ ] TestHelloos3のスキップ解除
    - [ ] バイナリ出力の検証

## 直近の変更点
- `pkg/operand` パッケージの `OperandImpl` 構造体と関連メソッド、テストを修正 (`WithForceImm8` メソッド、`WithBitMode` メソッドの変更、`Operands` インターフェースへの `WithForceImm8` メソッド追加を含む)

## 次のステップ
- Memory Bankの更新 (`progress.md`の更新)

## アクティブな決定事項と考慮事項
- `ForceImm8` フラグを追加することで、特殊な即値オペランドの扱いを制御する
- テストケースを追加し、`ForceImm8` フラグの動作を確認する

## 関連情報
[implementation_details.md](./implementation_details.md)
[technical_notes.md](./technical_notes.md)

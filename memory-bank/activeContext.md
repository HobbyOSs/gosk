# Active Context

## 現在の作業の焦点

### オペコード生成処理の改善
- [x] オペコード生成処理を`x86gen_utils.go`に集約
  - [x] `ResolveOpcode`関数の実装
    - オペコードとレジスタ番号から最終的なオペコードを生成
    - `Opcode.Addend`に基づいてレジスタ番号を選択
  - [x] `GetRegisterNumber`関数の実装
    - レジスタ名からレジスタ番号（0-7）への変換
    - 8/16/32ビットレジスタの対応付け

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

- `internal/codegen/x86gen_utils.go`に`ResolveOpcode`関数を追加
  - オペコードとレジスタ番号から最終的なオペコードを生成
  - `Opcode.Addend`に基づいてレジスタ番号を選択
  - エラーハンドリングを実装
- `internal/codegen/x86gen_utils.go`に`GetRegisterNumber`関数を追加
  - レジスタ名からレジスタ番号（0-7）への変換
  - 8/16/32ビットレジスタの対応付け
  - エラーハンドリングを実装
- `internal/codegen/x86gen_mov.go`を改善
  - レジスタ名から番号への変換ロジックを`GetRegisterNumber`に移動
  - `ResolveOpcode`関数を使用してオペコードを生成
  - コードの可読性と保守性を向上
- `internal/codegen/x86gen_arithmetic.go`を改善
  - `handleADD` 関数内で `ResolveOpcode` 関数と `GetRegisterNumber` 関数を使用するように修正
  - `x86gen.go` の `processOcode` 関数に `ocode.OpADD` のケースを追加し、`handleADD` 関数を呼び出すように修正
- `internal/codegen/x86gen_int.go`の`GenerateX86INT`関数を`handleINT`にリネーム

## 次のステップ

- `go vet` で検出されたエラーの修正
  - `pkg/operand/operand_impl.go`: struct field tag の構文エラー
  - `pkg/asmdb/instruction_search_test.go`: `db.FindInstruction` が未定義
  - `internal/gen/grammar_test.go`: struct literal で unkeyed fields を使用している
  - `test/pass1_test.go`: struct literal で unkeyed fields を使用している

## アクティブな決定事項と考慮事項

- オペコード生成処理を`x86gen_utils.go`に集約することで、コードの重複を避け、保守性を向上
- レジスタ名から番号への変換を共通化し、一貫性のある処理を実現
- エラーハンドリングを適切に実装し、デバッグ情報を提供

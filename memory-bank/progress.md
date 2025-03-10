# Progress

## 実装済みの機能

- アセンブラ命令実装のルーチンの定義
    - Pass1での命令実装手順
    - Ocodeの実装手順
    - 機械語生成の実装手順
    - 実装時の注意点の整理
- テストスイートの構造化
    - day01からday20までのテストケース構造の確立
    - 各dayでの新しい命令の追加パターンの明確化
    - テストケースの内容と検証方法の整理
- x86アセンブラの解析
- ASTの生成
- 基本的なコード生成
- `GetOutputSize` 関数の引数に `OutputSizeOptions` 構造体を追加
- PEGを用いた構文解析の最適化
- Pass1の基本的な評価処理
- トークン解析の基本実装
- オペランドの基本実装
  - 16ビット汎用レジスタSP, BPのパースを実装
  - メモリアドレッシングでのSP, BPの使用をサポート
  - テストケースの追加と検証
- x86のprefix bytes判定機能の実装
    - オペランドサイズプレフィックス(66h)の判定
    - アドレスサイズプレフィックス(67h)の判定
- 機械語サイズ計算機能の完成
    - プレフィックスバイトのサイズ計算を実装
    - オフセットとプレフィックスを含めた正確なサイズ計算
    - 複数の機械語パターンがある場合の最適な選択機能を実装
        - FindFormをFindEncodingに変更し、最小サイズのエンコーディングを直接選択
        - 複数の機械語パターンから最適なエンコーディングを一意に決定（例：MOV r16, imm16）
- 算術命令の基本構造の実装開始
    - 共通処理を行うヘルパー関数の実装
    - 各命令の基本的なハンドラー関数の定義
    - ADD命令の機械語生成を実装
        - プレフィックスバイト(66h)の処理
        - オペコードの生成
        - ModR/Mバイトの生成
        - 即値の処理
- システム命令の実装完了
    - INT命令（BIOS呼び出し）の実装
        - 2バイト命令（0xCD + 割り込み番号）
        - Ocodeの定義と機械語生成
    - HLT命令（CPU停止）の実装
        - 1バイト命令（0xF4）
        - パラメータなし命令としての処理
- Ocodeの使い方のドキュメント化
    - パラメータなし命令の実装パターン
    - パラメータあり命令の実装パターン
    - Emit関数の使用方法
    - 機械語サイズの計算方法
- 制御フロー命令の実装
    - JMP命令 (JMP rel8)
        - Pass1: processJMP 関数実装 (internal/pass1/pass1_inst_jmp.go)
        - Ocode: OpJMP 定義 (pkg/ocode/ocode.go)
        - 機械語生成: handleJMP 関数実装 (internal/codegen/x86gen_jmp.go)
- asmdbのセグメントレジスタ対応の改善
  - MOV r16, Sreg, MOV Sreg, r16のサポート
  - セグメントレジスタをr16として扱う機能を実装
  - matchOperands関数でsreg→r16の変換を処理
  - FindEncoding関数のリファクタリングとコメントの追加
  - コードの品質維持のためのlintチェックとテスト実行の徹底
- オペコード生成処理の改善
  - `ResolveOpcode`関数の実装
    - オペコードとレジスタ番号から最終的なオペコードを生成
    - `Opcode.Addend`に基づいてレジスタ番号を選択
  - `GetRegisterNumber`関数の実装
    - レジスタ名からレジスタ番号（0-7）への変換
    - 8/16/32ビットレジスタの対応付け
  - コードの可読性と保守性の向上
- `internal/codegen/x86gen_arithmetic.go`の改善
  - `handleADD` 関数内で `ResolveOpcode` 関数と `GetRegisterNumber` 関数を使用するように修正
  - `x86gen.go` の `processOcode` 関数に `ocode.OpADD` のケースを追加し、`handleADD` 関数を呼び出すように修正
- `internal/codegen/x86gen_int.go`の`GenerateX86INT`関数を`handleINT`にリネーム

## まだ構築が必要な部分

- Day02の実装
    - 基本命令
        - [ ] MOV命令（レジスタ間、即値）
        - [ ] ADD命令（フラグ更新）
        - [ ] CMP命令（比較演算、フラグ設定）
    - メモリ操作
        - [ ] メモリアドレッシング
        - [ ] レジスタ-メモリ間転送
        - [ ] ModR/M生成
    - 制御フロー
        - [ ] JE命令（条件分岐）
        - [ ] JMP命令
            - [ ] 相対アドレス計算 (Pass2での実装)
            - [ ] ジャンプ先ラベル解決 (Pass2での実装)
    - システム命令
- day03からday20までの命令実装
    - 各dayで追加される新しい命令の実装
    - テストケースの有効化と検証
- 算術命令の完全な実装
    - Ocodeの定義
    - 機械語生成の実装
    - ModR/Mの生成
    - テストケースの修正と有効化
- 高度なコード生成機能
- ユーザーインターフェースの改善
- `asmdb`と`operand`の依存関係の整理と設計修正
- Pass1の評価処理の完全なテストカバレッジ
- トークン解析の最適化と完全性の確保
- オペランド実装の完成度向上
- スタックマシンベースの設計の最適化

## 現在の進捗状況

- `internal/codegen/x86gen_arithmetic.go` の `handleADD` 関数を `x86gen_mov.go` の `handleMOV` 関数と平仄を合わせた
- `GenerateX86INT` 関数を `handleINT` にリネーム
- `internal/codegen/x86gen_utils.go`の`GenerateModRM`関数を修正
  - オペランドを引数として受け取るように変更
  - `GetRegisterNumber`関数を使用してレジスタ番号を取得
- `internal/codegen/x86gen_utils.go`の`GetRegisterNumber`関数を修正
  - セグメントレジスタに対応する番号を返すように修正
  - `case`を番号ごとにまとめ、可読性を向上
- `internal/codegen/x86gen_mov.go`の`handleMOV`関数を修正
  - `GenerateModRM`関数の代わりに`getModRMFromOperands`関数を呼び出すように変更
- `internal/codegen/x86gen_arithmetic.go`の`handleADD`関数を修正
  - `GenerateModRM`関数の代わりに`getModRMFromOperands`関数を呼び出すように変更

## 現在の進捗状況
- オペコード生成処理の改善と関連関数の修正が完了

## 既知の問題

- day02以降のテストケースが未実装
- 算術命令の実装が進行中（ADD命令は完了）
- トークン解析、オペランド判定、スタックマシン関連の構造に改善の余地あり

## 関連情報
[implementation_details.md](./implementation_details.md)
[technical_notes.md](./technical_notes.md)

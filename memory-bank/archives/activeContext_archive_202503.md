# Active Context

## 現在の作業の焦点
- `internal/codegen` パッケージ内の関数で、`CodeGenContext` をパラメータオブジェクトとして使用するようにリファクタリング

## 直近の変更点
- day02までの実装完了
  - システム命令 (INT, HLT)
  - 算術命令(一部)
  - CMP命令の呼び出し修正とテストケース追加
  - `pkg/operand/operand.go`に`InternalStrings() []string`メソッドを追加
  - `pkg/asmdb/instruction_search.go`の`matchOperands`関数を修正
  - JE命令、MOV命令 (レジスタ間, 即値)、ADD命令 (フラグ更新)、JMP命令のラベル解決を追加
  - CodegenClientインターフェースの拡張 (GetOcodes/SetOcodesメソッドを追加)
  - CodeGenContextへのBitModeの移動
  - `internal/ocode_client/client.go` の `NewCodegenClient` 関数を修正し、エラーハンドリングを追加

## 次のステップ
- `internal/codegen` パッケージのリファクタリング完了 (CodeGenContextパラメータオブジェクト化)
- `internal/codegen` パッケージ内の不要になったパラメータを削除
- メモリアドレッシングの実装

## 関連情報
[technical_notes.md](../details/technical_notes.md)

---
## 直近の変更点 - 2025/03/27 アーカイブ
- `internal/codegen/x86gen_utils.go`:
    - `GetRegisterNumber` 関数を修正し、制御レジスタ (CR0, CR2, CR3, CR4) に対応。
    - `ModRMByOperand` 関数を修正し、`bitMode` に基づいて 16bit/32bit メモリオペランド処理を分岐。
        - 16bit モード処理を改善し、単純なレジスタ間接参照 (`[SI]`, `[DI]`) および直接アドレス (`[imm16]`) の ModR/M とディスプレースメントを生成するように修正。これにより `TestGenerateX86` スイートのデグレを解消。
    - 未使用の `regexp` インポートを削除。
    - `operand.ParseNumeric` の代わりにローカルヘルパー関数 `parseNumeric` を追加・使用。
- `internal/codegen/x86gen_utils.go` のリファクタリング:
    - `modeStr` の switch 文を共通関数 `parseMode` として切り出し。
    - `ModRMByOperand` および `ModRMByValue` がメモリオペランド解析に `pkg/operand.ParseMemoryOperand` を使用するように修正。
    - 冗長な16bitモードの手動解析ロジック、`parseNumeric` 関数、`encoding/binary` インポートを削除。
    - 英語コメントを日本語に翻訳。

---
## 2025/03/29 アーカイブ (pkg/operand パーサー基本修正後)

## 現在の作業の焦点
- `pkg/operand` のパーサー (`participle` ベース) の基本的な問題を修正。
- 引き続き `test/day03_harib00i_test.go` のエラー (エンコーディングエラー、バイナリ長不一致) 対応。

## 直近の変更点
- **`pkg/operand` パーサー修正 (2025/03/29):**
    - `pkg/operand` のテスト (`TestBaseOperand_OperandType`) で失敗していたオペランド解析の問題を特定。
    - TDD アプローチで `pkg/operand/operand_impl_test.go` にテストケースを追加。
    - `pkg/operand/operand_impl.go` のレキサールール (`Reg`, `Seg` の順序) を修正。
    - `pkg/operand/operand.go` のパーサー定義 (`Instruction`, `CommaOperand` 構造体) を修正。
    - `pkg/operand/operand_impl.go` の `OperandTypes` ロジック (ラベルの扱い) を修正。
    - `pkg/operand/operand_impl_test.go` のテストケース (`MOV r32, label`) で `ForceRelAsImm=true` を設定。
    - 上記修正により `pkg/operand` のテスト (`TestBaseOperand_OperandType`) が成功。

## 次のステップ
- `test/day03_harib00i_test.go` を再実行し、エラー内容を確認する。
- エラー内容に基づき、エンコーディングエラーやバイナリ長不一致の原因を調査・修正する。
    - `asmdb` (JSON, fallback) のエンコーディング定義確認
    - `codegen` (MOV, ADD ハンドラ) のロジック確認
- **`Require67h` の TODO コメント解消**:
    - `[disp32]` や `[0x12345678]` のケースを正しく判定できるように `requireAddressSizePrefix` 関数を改善する。(`CalcOffsetByteSize` の改善または個別のオペランドサイズ計算が必要)

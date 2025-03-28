# Active Context

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
- **`internal/codegen` パッケージのリファクタリング (CodeGenContext 導入)**
- **`internal/codegen` パッケージの不要パラメータ削除**
- RESBの計算処理の実装

## 関連情報
[technical_notes.md](../details/technical_notes.md)
[implementation_details.md](../details/implementation_details.md)
(過去の変更点: [activeContext_archive_202503.md](../archives/activeContext_archive_202503.md))

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
- **`CodegenClient.Emit` インターフェース変更試行 (2025/03/29):**
    - `test/day03_harib00i_test.go` のエンコーディングエラー (`Failed to parse operand string 'ECX[ EBX + 16 ]'`) の原因調査。
    - 原因は `pass1` が Ocode 生成時にオペランドを単一文字列として `Emit` に渡し、`codegen` がそれを再結合して `pkg/operand` パーサーに渡していたためと特定。
    - `CodegenClient.Emit` のシグネチャを `Emit(string)` から `Emit(string, []string)` に変更し、関連ファイル (`ocode_client`, `pass1` 各所) を修正しようとした。
    - しかし、修正が広範囲に及び複雑化したため、ユーザー指示により中断。コード変更はユーザーが手動で revert する。

## 次のステップ
- **コード Revert 待ち**: ユーザーによるコード変更の revert を待つ。
- Revert 後、再度 `test/day03_harib00i_test.go` を実行し、エラー状況を確認する。
- **根本解決の検討**: `CodegenClient.Emit` インターフェースの変更を含む、オペランド受け渡し方法のリファクタリングを検討する (文書化後、PLAN MODE で再計画)。
- (保留) `Require67h` の TODO コメント解消
- (保留) RESBの計算処理の実装

## 関連情報
- [オペランド受け渡しフローと CodegenClient.Emit インターフェースの問題点 (2025/03/29)](memory-bank/details/technical_notes.md#オペランド受け渡しフローと-codegenclientemit-インターフェースの問題点-20250329)
[technical_notes.md](../details/technical_notes.md)
[implementation_details.md](../details/implementation_details.md)
(過去の変更点: [activeContext_archive_202503.md](../archives/activeContext_archive_202503.md))

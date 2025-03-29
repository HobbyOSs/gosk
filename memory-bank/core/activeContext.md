# Active Context

## 現在の作業の焦点
- **オペランドパーサー移行 (`pkg/ng_operand`)**:
    - pigeon peg ベースの新しいパーサー (`pkg/ng_operand`) の実装を進める。
    - `pkg/operand` から移植したテストコード (`operand_impl_test.go`) が通るように、`operand_impl.go` の `OperandPegImpl` 構造体とメソッド（特に `OperandTypes`）を実装・修正する。

## 直近の変更点
- **オペランドテストコード移植 (2025/03/29):**
    - `pkg/operand/operand_impl_test.go` を `pkg/ng_operand/operand_impl_test.go` にコピー。
    - `pkg/ng_operand/parser.go` に複数オペランド対応の `ParseOperands` 関数を追加（戻り値は `[]*ParsedOperandPeg`）。
    - `pkg/ng_operand/operand_grammar.peg` を修正し、セグメントオーバーライド、SHORT/FAR PTR、64bitレジスタなどに対応。
    - `go generate` でパーサーコードを更新。
    - `operand_impl_test.go` のテストロジックを `ParseOperands` の戻り値に合わせて修正。
    - テストを実行し、パースエラーは解消されたが、型解決ロジックの不完全性によりサイズ関連のテスト (`TestOperandPegImpl_OperandType`) が多数失敗することを確認。
- **オペランドパーサー移行開始 (2025/03/29):**
    - `participle` ベースのパーサー (`pkg/operand`) の複雑さが開発のボトルネックになっているため、`pigeon` peg パーサーへの移行を決定。
    - オペランド専用のpeg文法 (`operand_grammar.peg`) を設計。
    - 新しいパーサー用のパッケージ `pkg/ng_operand` を作成。
    - `pkg/ng_operand/operand_grammar.peg` (peg定義), `pkg/ng_operand/operand_types.go` (型定義), `pkg/ng_operand/parser.go` (ラッパー関数) を作成。
    - `go:generate` を使用して `pigeon` コマンドを実行し、パーサーコード (`operand_grammar.go`) を生成する仕組みを導入 (`pkg/ng_operand/generate.go`)。
    - 既存の `pkg/operand` のインターフェース (`Operands`) とテストコード (`operand_test.go`) を `pkg/ng_operand` にコピーし、パッケージ名を修正。
    - `pkg/ng_operand/operand_impl.go` を作成し、`Operands` インターフェースの基本的な実装 (`OperandPegImpl`) を開始。
    - `OperandTypes` と `Require66h`, `Require67h` の基本的なロジックを実装し、コピーしたテスト (`TestRequire66h`, `TestRequire67h`) が通ることを確認。
- (過去の変更点) **`pkg/operand` パーサー修正 (2025/03/29):**
    - participleベースのパーサーの基本的な問題を修正 (`TestBaseOperand_OperandType` が成功)。

## 次のステップ
- **`pkg/ng_operand/operand_impl.go` の実装**:
    - `OperandPegImpl` 構造体の設計見直し（`[]*ParsedOperandPeg` を保持するように変更するか検討）。
    - `OperandTypes` メソッドに、ビットモードやオペランド間の依存関係を考慮した完全な型解決ロジックを実装する（`imm`/`m` のサイズ決定）。
    - `CalcOffsetByteSize`, `DetectImmediateSize`, `Require66h`, `Require67h` などのメソッドを実装・修正する。
    - `operand_impl_test.go` のテストが成功するように実装を進める。
- **段階的置換**: `pkg/ng_operand` のテストが安定したら、`internal/pass1` などから `pkg/operand` の利用箇所を `pkg/ng_operand` に置き換えていく。

## 関連情報
- [オペランドパーサー移行 (participle -> pigeon)](../details/implementation_details.md#オペランドパーサー移行-participle---pigeon-20250329)
- [オペランドサイズ決定の複雑さについて](../details/technical_notes.md#オペランドサイズ決定の複雑さについて-20250329)
- (過去) [オペランド受け渡しフローと CodegenClient.Emit インターフェースの問題点 (2025/03/29)](memory-bank/details/technical_notes.md#オペランド受け渡しフローと-codegenclientemit-インターフェースの問題点-20250329) - Revert済み
[technical_notes.md](../details/technical_notes.md)
[implementation_details.md](../details/implementation_details.md)
(過去の変更点: [activeContext_archive_202503.md](../archives/activeContext_archive_202503.md))

# Active Context

## 現在の作業の焦点
- **オペランドパーサー移行 (`pkg/ng_operand`)**:
    - `pkg/ng_operand/operand_impl.go` の `OperandTypes` メソッドにおけるサイズ解決ロジックの修正。
    - `pkg/ng_operand/operand_type_test.go` のテスト (`TestOperandPegImpl_OperandType`) をパスさせることを目指す。

## 直近の変更点
- **コードクリーンアップ (2025/03/29):**
    - `go fmt ./...` と `go mod tidy` を実行。
- **テストファイル分割 (2025/03/29):**
    - `pkg/ng_operand/operand_impl_test.go` を以下のファイルに分割:
        - `operand_type_test.go` (`TestOperandPegImpl_OperandType`)
        - `detect_immediate_size_test.go` (`TestOperandPegImpl_DetectImmediateSize`)
        - `parse_operands_test.go` (`TestParseOperands_FromString`)
- **`requires` 関数のリファクタリング (2025/03/29):**
    - `pkg/ng_operand/requires.go` を作成し、`Require66h` と `Require67h` を `OperandPegImpl` のメソッドとして移動。
    - `operand_impl.go` から該当メソッドを削除。
    - `operand_test.go` (分割前のファイル) の呼び出し箇所を修正。
- **`OperandPegImpl` 構造体変更 (2025/03/29):**
    - `operand_impl.go` の `OperandPegImpl` 構造体が単一の `*ParsedOperandPeg` ではなく、`[]*ParsedOperandPeg` を保持するように変更。
    - 関連するメソッド (`NewOperandPegImpl`, `FromString`, `InternalString`, `InternalStrings`, `Serialize`, `OperandTypes`, `CalcOffsetByteSize`, `DetectImmediateSize`) を修正。
- **`OperandTypes` サイズ解決ロジック修正 (2025/03/29):**
    - `operand_impl.go` の `OperandTypes`, `resolveMemorySize`, `resolveDependentSizes` メソッドを複数回修正し、即値・メモリ・ラベルのサイズ解決ロジックを改善。
    - `operand_type_test.go` 内で `ParseOperands` 呼び出し時に適切なビットモードを設定するように修正。
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
- **`pkg/ng_operand/operand_impl.go` の `OperandTypes` 修正**:
    - `operand_type_test.go` の失敗ケースに基づき、`OperandTypes`, `resolveMemorySize`, `resolveDependentSizes` のサイズ解決ロジックをさらに修正する。
    - 特に、単独即値の `imm32` 解決、`forceImm8` の優先適用、セグメントオーバーライド付きレジスタのメモリ解決を重点的に見直す。
- **テスト実行と修正**:
    - `operand_type_test.go` のテスト (`TestOperandPegImpl_OperandType`) を実行し、パスするまで修正を繰り返す。
    - `detect_immediate_size_test.go` (`TestOperandPegImpl_DetectImmediateSize`) と `parse_operands_test.go` (`TestParseOperands_FromString`) も実行し、必要に応じて修正する。
- **段階的置換**: `pkg/ng_operand` のテストが安定したら、`internal/pass1` などから `pkg/operand` の利用箇所を `pkg/ng_operand` に置き換えていく。

## 関連情報
- [オペランドパーサー移行 (participle -> pigeon)](../details/implementation_details.md#オペランドパーサー移行-participle---pigeon-20250329)
- [オペランドサイズ決定の複雑さについて](../details/technical_notes.md#オペランドサイズ決定の複雑さについて-20250329)
- (過去) [オペランド受け渡しフローと CodegenClient.Emit インターフェースの問題点 (2025/03/29)](memory-bank/details/technical_notes.md#オペランド受け渡しフローと-codegenclientemit-インターフェースの問題点-20250329) - Revert済み
[technical_notes.md](../details/technical_notes.md)
[implementation_details.md](../details/implementation_details.md)
(過去の変更点: [activeContext_archive_202503.md](../archives/activeContext_archive_202503.md))

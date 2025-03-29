# Progress

## 実装済み機能 (要約)
- x86アセンブラ解析基盤 (pass1, pass2)
- Ocode・PEGパーサ等の基礎部分
- 主要命令実装 (システム命令, 算術命令, CMP, JE, MOV, ADD, JMP)
- CodegenClient 関連機能
- EQU命令の展開
- `getModRMFromOperands`の返り値の型変更 (`uint32` -> `[]byte`)
- InstructionForm の Encoding を ModRM の要否で振り分け
- OUT命令の実装
- CALL命令の実装
- 論理命令の実装 (AND, OR, XOR, NOT)
- 論理シフト/算術シフト命令の実装 (SHR, SHL, SAR) (一部テストはコメントアウト)
- IN命令の実装 (pass1, codegen, fallback table)
- RET命令の実装 (pass1, codegen, test)
- ModR/M 生成ロジックの一部修正 (制御レジスタ対応、16bit/32bit 分岐改善)
- `internal/codegen/x86gen_utils.go` のリファクタリング (`modeStr`共通化、`ParseMemoryOperand`利用、コメント翻訳等)
- **ModR/M 生成エラー (`unsupported 16bit mem operand`) の部分修正 (`pkg/operand/modrm_address.go`)**
- **`BitMode` 伝達ロジックの修正 (`client`, `ocode_client`, `pass1`)**
- **`codegen` の MOV, ADD ハンドラに `.WithBitMode()` 追加**
- **JMP rel32 / Jcc rel32 の実装**: `internal/codegen/x86gen_jmp.go` に `rel32` オフセット計算とエンコーディングを追加。
- **`pkg/operand/requires.go` の修正 (2025/03/28)**:
    - `Require66h` (オペランドサイズプレフィックス): 16bitモードでの32bit即値判定を `ParsedOperands()` ベースに修正。
    - `Require67h` (アドレスサイズプレフィックス): 実効アドレスサイズに基づいて判定するようにロジックを修正。
- **`pkg/operand/requires.go` のリファクタリング (2025/03/28)**:
    - `Require66h`, `Require67h` 関数を小さく分割し、可読性と保守性を向上。
    - `is32bitRegInIndirectMem` 関数で正規表現を使用するように修正。
- **`internal/codegen/x86gen_test.go` の修正 (2025/03/28)**:
    - `TestGenerateX86/MOV_SI_a_label` テストケースを `TestGenerateX86/MOV SI, 0x0000` に修正。
- **`pkg/operand` パーサー修正 (2025/03/29):**
    - `participle` ベースのパーサーにおけるレキサールール、パーサー定義、型決定ロジックの基本的な問題を修正。 (`TestBaseOperand_OperandType` が成功)
- **`pkg/ng_operand` パッケージの基本構造作成 (2025/03/29):**
    - pigeon peg ベースの新しいオペランドパーサー用パッケージを作成。
    - peg文法 (`operand_grammar.peg`), 型定義 (`operand_types.go`), パーサー生成 (`generate.go`, `operand_grammar.go`), ラッパー関数 (`parser.go`) を実装。
    - 既存の `pkg/operand` からインターフェース (`operand.go`) とテスト (`operand_test.go`) をコピーし、基本的なテスト (`TestRequire66h`, `TestRequire67h`) が通るように `operand_impl.go` を部分的に実装。
- **`pkg/ng_operand` テストコード移植 (2025/03/29):**
    - `pkg/operand/operand_impl_test.go` を `pkg/ng_operand/operand_impl_test.go` に移植。
    - PEG 文法 (`operand_grammar.peg`) を修正し、パースエラーを解消。
    - テストコード内の型解決ロジックを暫定的に実装（一部テストは失敗中）。

## まだ必要な実装
- **`pkg/ng_operand` の実装完了**:
    - **型解決ロジックの実装**: `operand_impl.go` の `OperandTypes` メソッド等に、ビットモードやオペランド間の依存関係を考慮した完全な型解決ロジックを実装する（特に `imm` および `m` のサイズ決定）。
    - **メソッド実装**: `CalcOffsetByteSize`, `DetectImmediateSize`, `Require66h`, `Require67h` などのメソッドを完全に実装する。
    - **テスト修正/成功**: `operand_impl_test.go` のテストがすべて成功するように実装を修正する。
    - **構造体設計検討**: `OperandPegImpl` が複数のオペランド (`[]*ParsedOperandPeg`) を保持できるように設計を見直すことを検討する。
- **`pkg/ng_operand` への段階的置換**:
    - `internal/pass1`, `internal/codegen`, `pkg/asmdb` など、`pkg/operand` を利用している箇所を `pkg/ng_operand` に置き換える。
    - 最終的に `pkg/operand` を削除し、`pkg/ng_operand` を `pkg/operand` にリネームする。
- **`test/day03_harib00i_test.go` のエラー対応 (継続):**
    - オペランドパーサー移行後に再度確認・対応。
- (保留) RESBの計算処理の実装
- (保留) `internal/codegen` パッケージのリファクタリング (CodeGenContext 導入)
- (保留) `internal/codegen` パッケージの不要パラメータ削除

## 関連情報
[technical_notes.md](../details/technical_notes.md)
(過去の実装詳細: [progress_archive_202503.md](../archives/progress_archive_202503.md))

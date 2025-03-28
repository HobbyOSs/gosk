# Active Context

## 現在の作業の焦点
- `pkg/operand/requires.go` のリファクタリングと `internal/codegen` のテスト `TestGenerateX86/MOV_SI_a_label` の対応が完了。
- `test/day03_harib00i_test.go` の残存エラー (エンコーディングエラー、バイナリ長不一致) 対応。

## 直近の変更点
- `pkg/operand/modrm_address.go` の `ParseMemoryOperand` を修正し、16ビットモードで32ビットレジスタ (`[ESI]` など) が使われた場合の R/M ビット解決に対応。これにより `unsupported 16bit mem operand` エラーは解消。
- `Pass1` から `CodeGenClient` への `BitMode` 伝達ロジックを修正 (`SetBitMode` インターフェース追加と呼び出し)。
- `codegen` の MOV, ADD ハンドラに `.WithBitMode()` 呼び出しを追加。
- **リファクタリング (pkg/operand, pkg/asmdb)**:
    - `internal/ast/bit_mode.go` を `pkg/operand/bit_mode.go` に移動。
    - `internal/ast/support_cpu.go` を `pkg/asmdb/support_cpu.go` に移動。
    - 関連するインポートパスと型参照を修正。
    - 循環参照を解消 (`BitMode` を `pkg/operand` に配置)。
- **`pkg/operand/requires.go` の修正 (2025/03/28)**:
    - `Require66h` (オペランドサイズプレフィックス): 16bitモードでの32bit即値判定ロジックを `ParsedOperands()` を使うように修正。
    - `Require67h` (アドレスサイズプレフィックス): `x86_prefix.cc` を参考に、メモリオペランド文字列内のレジスタ名に基づいて実効アドレスサイズを判定するロジックに修正。これにより `TestGenerateX86/MOV_DWORD_[_0x0ff8_]0x000a0000` のテスト失敗は解消。
- **`pkg/operand/requires.go` のリファクタリング (2025/03/28)**:
    - `Require66h`, `Require67h` 関数を ছোটさく分割し、可読性と保守性を向上。
    - `is32bitRegInIndirectMem` 関数で正規表現を使用するように修正。
- **`internal/codegen/x86gen_test.go` の修正 (2025/03/28)**:
    - `TestGenerateX86/MOV_SI_a_label` テストケースを、`codegen` の段階ではラベル参照が解決済みであることを考慮し、具体的なアドレス値を使用する `TestGenerateX86/MOV SI, 0x0000` に修正。

## 残存エラー (day03_harib00i) - 2025/03/28 更新
- **MOV/ADD エンコーディングエラー (継続)**:
    - `MOV r32, imm32/label/m32` および `ADD r32, r32/imm` 形式での `Failed to find encoding` エラーは依然として残存している可能性が高い (`pkg/operand` パーサーの問題)。
- **バイナリ長不一致 (継続)**:
    - 上記のエンコーディングエラーにより、依然としてバイナリ長が不足している可能性が高い。

## 次のステップ
- **エンコーディングエラー調査 (継続)**:
    - `pkg/operand` パーサーの問題は一旦保留。
    - **`asmdb` の確認**: エラーが発生している `MOV r32, imm32/label/m32` や `ADD r32, r32/imm` の形式に対応するエンコーディングが `pkg/asmdb/json-x86-64/x86_64.json` または `pkg/asmdb/instruction_table_fallback.go` に正しく定義されているか確認する。
    - **`codegen` ロジックの確認**: `internal/codegen/x86gen_mov.go` や `internal/codegen/x86gen_arithmetic.go` (ADD命令) のハンドラが、`asmdb.FindEncoding` を呼び出す際に、(仮にパーサーが正しく動作した場合に期待される) オペランドタイプを正しく渡しているか確認する。
- **`Require67h` の TODO コメント解消**:
    - `[disp32]` や `[0x12345678]` のケースを正しく判定できるように `requireAddressSizePrefix` 関数を改善する。(`CalcOffsetByteSize` の改善または個別のオペランドサイズ計算が必要)

## 関連情報
[technical_notes.md](../details/technical_notes.md)
[implementation_details.md](../details/implementation_details.md)
(過去の変更点: [activeContext_archive_202503.md](../archives/activeContext_archive_202503.md))

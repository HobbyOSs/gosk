# Active Context

## 現在の作業の焦点
- `test/day03_harib00i_test.go` の残存エラー対応。

## 直近の変更点
- `pkg/operand/modrm_address.go` の `ParseMemoryOperand` を修正し、16ビットモードで32ビットレジスタ (`[ESI]` など) が使われた場合の R/M ビット解決に対応。これにより `unsupported 16bit mem operand` エラーは解消。
- `Pass1` から `CodeGenClient` への `BitMode` 伝達ロジックを修正 (`SetBitMode` インターフェース追加と呼び出し)。
- `codegen` の MOV, ADD ハンドラに `.WithBitMode()` 呼び出しを追加。

## 残存エラー (day03_harib00i)
- **MOV エンコーディングエラー**: `MOV r32, imm32/m32/label` 形式で多数発生。
- **ADD エンコーディングエラー**: `ADD r32, r32/imm` 形式で発生。
- **JMP rel32 未実装エラー**: `JMP DWORD ...` で発生。
- **バイナリ長不一致**: `expected length 304 actual length 231`。

## 次のステップ
- **エンコーディングエラー調査**:
    - `pkg/operand/operand_impl.go` の `OperandTypes()` メソッドを調査する。
    - 特に、サイズプレフィックスがないメモリオペランド (`[EBX+16]`) や、即値/ラベル (`bootpack`, `0x00280000`) のオペランドタイプ (`CodeM32`, `CodeIMM32` など) がどのように決定されているかを確認する。
    - `asmdb.FindEncoding` がこれらの型を正しく使ってエンコーディングを検索できているか確認する。
- **JMP rel32 の実装**: `internal/codegen/x86gen_jmp.go` の `handleJMP` 関数に `rel32` の処理を追加する。

## 関連情報
[technical_notes.md](../details/technical_notes.md)
[implementation_details.md](../details/implementation_details.md)
(過去の変更点: [activeContext_archive_202503.md](../archives/activeContext_archive_202503.md))

# Active Context

## 現在の作業の焦点
- `test/day03_harib00i_test.go` の残存エラー対応。

## 直近の変更点
- `pkg/operand/modrm_address.go` の `ParseMemoryOperand` を修正し、16ビットモードで32ビットレジスタ (`[ESI]` など) が使われた場合の R/M ビット解決に対応。これにより `unsupported 16bit mem operand` エラーは解消。
- `Pass1` から `CodeGenClient` への `BitMode` 伝達ロジックを修正 (`SetBitMode` インターフェース追加と呼び出し)。
- `codegen` の MOV, ADD ハンドラに `.WithBitMode()` 呼び出しを追加。

## 残存エラー (day03_harib00i) - 2025/03/27 更新
- **MOV/ADD エンコーディングエラー**:
    - `MOV r32, imm32/label/m32` および `ADD r32, r32/imm` 形式で `Failed to find encoding` / `Failed to find encoding for ADD` が多数発生。
    - **原因**: `pkg/operand` のパーサー (`participle`利用) が、複数のオペランド (特にレジスタ、即値、ラベルの組み合わせ) を含む文字列を正しく解析・分類できていない。Lexerルールの調整（順序変更、厳密化）、キャッシュ無効化、`Instruction`構造体のタグ変更を試みたが解決せず。パーサーの問題は根深いと判断し、一旦調査を保留。
    - 例: `MOV EDI, 0x00280000` が `[imm32 imm32]` と誤判定される。`ADD ESI, EBX` が `[unknown]` となる。
    - 問題再現用テストコード: `pkg/operand/operand_impl_test.go`
- **バイナリ長不一致**: `expected length 304 actual length 236` (68バイト不足)。
    - JMP rel32 実装により5バイト増加したが、上記のMOV/ADDエンコーディングエラーにより多数の命令が生成されず、依然として長さが不足している。

## 次のステップ
- **エンコーディングエラー調査 (方針転換)**:
    - `pkg/operand` パーサーの問題は一旦保留。
    - **`asmdb` の確認**: エラーが発生している `MOV r32, imm32/label/m32` や `ADD r32, r32/imm` の形式に対応するエンコーディングが `pkg/asmdb/json-x86-64/x86_64.json` または `pkg/asmdb/instruction_table_fallback.go` に正しく定義されているか確認する。
    - **`codegen` ロジックの確認**: `internal/codegen/x86gen_mov.go` や `internal/codegen/x86gen_arithmetic.go` (ADD命令) のハンドラが、`asmdb.FindEncoding` を呼び出す際に、(仮にパーサーが正しく動作した場合に期待される) オペランドタイプを正しく渡しているか確認する。

## 関連情報
[technical_notes.md](../details/technical_notes.md)
[implementation_details.md](../details/implementation_details.md)
(過去の変更点: [activeContext_archive_202503.md](../archives/activeContext_archive_202503.md))

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

## まだ必要な実装
- **`test/day03_harib00i_test.go` の残存エラー対応:**
    - エンコーディング未発見エラー (`Failed to find encoding: no matching encoding found`) の修正 (複数の `MOV`, `ADD` 命令)。
    - `Failed to process ocode: not implemented: JMP rel32` エラーの修正 (`JMP DWORD 2*8:0x0000001b`)。
- **ModR/M 生成ロジックのリファクタリング検討:**
    - `internal/codegen/x86gen_utils.go` 内の複雑な手動パースは解消されたが、`pkg/operand` 側に `bitMode` を考慮した統一的なメモリオペランド解析・ModR/M 生成機能 (`ParseMemoryOperand` の改善または新規関数) を実装する検討は継続。
- JMP系命令 (Jcc命令) のrel32オフセット対応 (上記 JMP rel32 と関連)
- RESBの計算処理の実装
- メモリアドレッシング (エンコーディング未発見エラーと関連)
- `internal/codegen` パッケージのリファクタリング (CodeGenContext 導入)
- `internal/codegen` パッケージの不要パラメータ削除

## 関連情報
[technical_notes.md](../details/technical_notes.md)
(過去の実装詳細: [progress_archive_202503.md](../archives/progress_archive_202503.md))

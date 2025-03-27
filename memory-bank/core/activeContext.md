# Active Context

## 現在の作業の焦点
- ModR/M 生成ロジックのリファクタリング検討。
- `test/day03_harib00i_test.go` の残存エラー対応。

## 次のステップ
- **ModR/M 生成ロジックのリファクタリング検討:**
    - `internal/codegen/x86gen_utils.go` 内の複雑な手動パースは解消されたが、`pkg/operand` パッケージ側に、`bitMode` を考慮した統一的なメモリオペランド解析・ModR/M 生成機能 (`ParseMemoryOperand` の改善または新規関数) を実装し、`codegen` 側はそれを呼び出すだけにするリファクタリングを引き続き検討する。
- **`test/day03_harib00i_test.go` の残存エラー対応:**
    - エンコーディング未発見エラー (`Failed to find encoding: no matching encoding found`) の修正 (複数の `MOV`, `ADD` 命令)。
        - 特にラベル (`bootpack`) や `[ EBX + offset ]` 形式のメモリオペランドを含む命令のエンコーディング選択ロジックを確認・修正。(`handleMOV` などが正しい `bitMode` を渡しているか、`asmdb` の検索ロジック自体に問題はないか)
    - `Failed to process ocode: not implemented: JMP rel32` エラーの修正 (`JMP DWORD 2*8:0x0000001b`)。
        - `internal/codegen/x86gen_jmp.go` の `handleJMP` 関数に `rel32` の処理を実装。

## 関連情報
[technical_notes.md](../details/technical_notes.md)
[implementation_details.md](../details/implementation_details.md)
(過去の変更点: [activeContext_archive_202503.md](../archives/activeContext_archive_202503.md))

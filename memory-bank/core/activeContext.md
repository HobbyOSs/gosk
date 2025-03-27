# Active Context

## 現在の作業の焦点
- ModR/M 生成ロジックのリファクタリング検討。
- `test/day03_harib00i_test.go` の残存エラー対応。

## 直近の変更点
- `internal/codegen/x86gen_utils.go`:
    - `GetRegisterNumber` 関数を修正し、制御レジスタ (CR0, CR2, CR3, CR4) に対応。
    - `ModRMByOperand` 関数を修正し、`bitMode` に基づいて 16bit/32bit メモリオペランド処理を分岐。
        - 16bit モード処理を改善し、単純なレジスタ間接参照 (`[SI]`, `[DI]`) および直接アドレス (`[imm16]`) の ModR/M とディスプレースメントを生成するように修正。これにより `TestGenerateX86` スイートのデグレを解消。
    - 未使用の `regexp` インポートを削除。
    - `operand.ParseNumeric` の代わりにローカルヘルパー関数 `parseNumeric` を追加・使用。
- `internal/codegen/x86gen_utils.go` のリファクタリング:
    - `modeStr` の switch 文を共通関数 `parseMode` として切り出し。
    - `ModRMByOperand` および `ModRMByValue` がメモリオペランド解析に `pkg/operand.ParseMemoryOperand` を使用するように修正。
    - 冗長な16bitモードの手動解析ロジック、`parseNumeric` 関数、`encoding/binary` インポートを削除。
    - 英語コメントを日本語に翻訳。

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

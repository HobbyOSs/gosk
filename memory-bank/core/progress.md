# Progress

## 実装済み
- **`JMP DWORD ptr` のエンコーディング:** `JMP ptr16:32` (Opcode `66 EA cd`) の処理を実装 (`pass1` および `codegen`)。単体テスト成功。
- **`JMP rel16` のプレフィックス問題修正:** 16bitモードで不要な `0x66` が付与される問題を修正。

## まだ必要な実装
- **`test/day03_harib00i_test.go` のエラー対応 (継続):**
    - バイナリ差分と長さの不一致 (例: expected 304, actual 305)。
    - ~~**`IMUL` のエンコーディング/ModR/M生成:**~~ (修正済み)
    - ~~**`JMP DWORD ptr` のエンコーディング:**~~ (実装済み)
    - ~~**`JMP rel16` のプレフィックス問題:**~~ (修正済み)
    - **ラベル/定数/データ定義:** `bootpack` ラベル等のアドレス解決、定数計算 (`DSKCAC0+512` 等)、`ALIGNB`/`RESB`/`DW`/`DD` の処理にずれがある可能性。
    - **LOC計算のずれ調査:** 上記のデータ定義や他の命令のLOC計算にずれがある可能性。
- **`pkg/ng_operand` の改善 (TODOs):** (一部対応済み、残りは以下)
    - パーサー: NEAR/FAR PTR, `ParseOperands` での各種フラグ考慮, 文字列リテラル内カンマ対応。
    - 実装: QWORD サポート, `CalcOffsetByteSize` / `DetectImmediateSize` の複数オペランド対応。
    - テスト: 複数オペランド, 未対応ケース, FAR/NEAR PTR, `ParseOperands` 拡充。
- (保留) RESBの計算処理の実装
- (保留) `internal/codegen` パッケージのリファクタリング (CodeGenContext 導入)
- (保留) `internal/codegen` パッケージの不要パラメータ削除

## 関連情報
(技術的な詳細メモは `memory-bank/archives/technical_notes_archive_202503.md` にアーカイブ済み)
(過去の実装詳細: [progress_archive_202503.md](../archives/progress_archive_202503.md))

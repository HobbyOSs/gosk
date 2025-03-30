# Progress

## 実装済み
- **`JMP DWORD ptr` のエンコーディング:** `JMP ptr16:32` (Opcode `66 EA cd`) の処理を実装 (`pass1` および `codegen`)。単体テスト成功。
- **`JMP rel16` のプレフィックス問題修正:** 16bitモードで不要な `0x66` が付与される問題を修正。
- **`ALIGNB` のLOC計算修正:** LOCが既に境界に揃っている場合に不要なパディングを追加しないように修正。
- **`DD`/`DW`/`DB` のラベル解決実装:** 識別子がラベルの場合にアドレス解決を行うように修正。
- **`TestHarib00i` 期待値修正:** `ALIGNB` 修正に伴い不要となった `"FILL 11"` を削除。

## まだ必要な実装
- **`test/day03_harib00i_test.go` の相対ジャンプオフセットずれ調査:**
    - `CALL`, `JMP`, `Jcc` 命令の相対オフセットが期待値とずれている。
    - オフセット計算ロジック (`ターゲットアドレス - (現在の命令のアドレス + 命令のサイズ)`) の見直しが必要。 (ユーザー側でアセンブラダンプマスタを使用して調査予定)
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

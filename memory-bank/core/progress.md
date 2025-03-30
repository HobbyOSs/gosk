# Progress

## 実装済み機能 (要約 - 2025/03)
- x86アセンブラ解析基盤 (pass1, pass2)
- Ocode・PEGパーサ等の基礎部分
- 主要命令実装 (システム命令, 算術命令, CMP, JE, MOV, ADD, JMP, OUT, CALL, 論理, シフト, IN, RET, IMUL, SUB)
- CodegenClient 関連機能
- EQU命令の展開
- ModR/M 生成ロジックの改善 (制御レジスタ, 16/32bit 分岐, 32bit アドレッシング)
- `internal/codegen/x86gen_utils.go` のリファクタリング
- `BitMode` 伝達ロジックの修正
- JMP/CALL rel32 の実装
- `pkg/operand` パーサーの基本的な問題修正と `requires.go` の改善
- `pkg/ng_operand` パッケージの基本構造作成、テスト移植、リファクタリング、バグ修正
- diff出力の改善 (`github.com/akedrou/textdiff` 導入)
- ADD, OUT, MOV, LGDT 命令のエンコーディング/fallback定義修正
- JMP/CALL の pass1 LOC 計算と codegen エンコーディング修正

## まだ必要な実装
- **`test/day03_harib00i_test.go` のエラー対応 (継続):**
    - バイナリ差分と長さの不一致 (expected 304, actual 300)。
    - **`IMUL`/`SUB` のエンコーディング選択:** `asmdb.FindEncoding` が即値サイズに基づいて `imm8` (Opcode `6B`/`83`) と `imm32` (Opcode `69`/`81`) を正しく選択できていない。(`lo.MinBy` の比較ロジック修正中)
    - **`JMP DWORD ptr` のエンコーディング:** `JMP ptr16:32` (Opcode `66 EA cd`) の実装が必要。
    - **ラベル/定数/データ定義:** `bootpack` ラベル等のアドレス解決、定数計算 (`DSKCAC0+512` 等)、`ALIGNB`/`RESB`/`DW`/`DD` の処理にずれがある可能性。
    - (関連) LOC計算のずれ調査。
- **エンコーディング検索アーキテクチャの見直し:**
    - 符号拡張 (`imm8` vs `imm16/32`) の扱いを改善するため、検索方法の見直しを検討 (`relaxImmSearch` アプローチ)。
- **`pkg/ng_operand` への段階的置換**:
    - `internal/pass1`, `internal/codegen`, `pkg/asmdb` 等の利用箇所を置き換え。
    - 最終的に `pkg/operand` を削除し、`pkg/ng_operand` を `pkg/operand` にリネーム。
- **`pkg/ng_operand` の改善 (TODOs):**
    - パーサー: 型推論改善, NEAR/FAR PTR, `ParseOperands` での各種フラグ考慮, 文字列リテラル内カンマ対応。
    - 実装: QWORD サポート, `CalcOffsetByteSize` / `DetectImmediateSize` の複数オペランド対応。
    - テスト: 複数オペランド, 未対応ケース, FAR/NEAR PTR, `ParseOperands` 拡充。
- (保留) RESBの計算処理の実装
- (保留) `internal/codegen` パッケージのリファクタリング (CodeGenContext 導入)
- (保留) `internal/codegen` パッケージの不要パラメータ削除

## 関連情報
[technical_notes.md](../details/technical_notes.md)
(過去の実装詳細: [progress_archive_202503.md](../archives/progress_archive_202503.md))

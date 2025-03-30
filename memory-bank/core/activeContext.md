# 現在の状況 (Active Context) - 2025/03/30

## 現在の焦点
- **`TestDay03Suite/TestHarib00i` のテスト失敗調査 (継続):**
    - `IMUL` の問題は解決したが、依然としてテストが失敗する。
    - 原因として考えられるのは `JMP DWORD ptr` のエンコーディング、ラベル/定数/データ定義 (`ALIGNB`, `RESB`, `DW`, `DD`, 定数計算) の処理のずれなど。

## 次のステップ
1. **`TestDay03Suite/TestHarib00i` の失敗原因調査**:
    - `harib00i.bin` と生成されたバイナリの差分を詳細に比較し、どの命令/データでずれが生じているか特定する。
    - 特に `JMP DWORD ptr` (Opcode `EA cd`)、`ALIGNB`, `RESB`, `DW`, `DD`、ラベルアドレス、定数計算 (`DSKCAC0+512` 等) 周辺の処理を確認する。
2. **`JMP DWORD ptr` の実装確認・修正:** (上記調査結果に応じて)
3. **ラベル/定数/データ定義の問題調査・修正:** (上記調査結果に応じて)

(過去の完了作業: [activeContext_archive_202503.md](../archives/activeContext_archive_202503.md))

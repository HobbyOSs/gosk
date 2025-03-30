# 現在の状況 (Active Context) - 2025/03/30

## 完了した作業
- **`JMP rel16` のプレフィックス問題修正:** 16bitモードで不要な `0x66` が付与される問題を修正。
- **`ALIGNB` のLOC計算修正:** LOCが既に境界に揃っている場合に不要なパディングを追加しないように修正 (`internal/pass1/pass1_inst_pseudo.go`)。
- **`DD`/`DW`/`DB` のラベル解決実装:** 識別子がラベルの場合にアドレス解決を行うように修正 (`internal/pass1/pass1_inst_pseudo.go`)。
- **`TestHarib00i` 期待値修正:** `ALIGNB` 修正に伴い不要となった `"FILL 11"` を削除 (`test/day03_harib00i_test.go`)。

## 現在の焦点
- **`TestDay03Suite/TestHarib00i` の相対ジャンプオフセットずれ調査:**
    - `ALIGNB`, `DD`/`DW`/`DB` の修正により多くの差分は解消されたが、`CALL`, `JMP`, `Jcc` 命令の相対オフセットが期待値とずれている。
    - 原因として、相対オフセット計算ロジック (`ターゲットアドレス - (現在の命令のアドレス + 命令のサイズ)`) に問題がある可能性が高い。

## 次のステップ
1. **相対ジャンプオフセット計算の調査:** (ユーザー側でアセンブラダンプマスタを使用して調査)
    - `internal/pass1/pass1_inst_jmp.go` (Pass1でのサイズ計算)
    - `internal/codegen/x86gen_jmp.go` (Pass2/codegenでのオフセット計算)
    - 上記ファイルのロジックを確認し、ずれの原因を特定する。
2. **相対ジャンプオフセット計算の修正:** (上記調査結果に応じて)

(過去の完了作業: [activeContext_archive_202503.md](../archives/activeContext_archive_202503.md))

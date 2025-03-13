# Active Context

## 現在の作業の焦点
- CMP命令の追加
- テストエラーの修正

## Day02実装計画
- MOV命令 (レジスタ間, 即値)
- ADD命令 (フラグ更新)
- CMP命令 (フラグ更新)
- 他、メモリアドレッシングや制御フロー命令の実装

## 直近の変更点
- `internal/codegen/x86gen.go` の `handleCMP` 関数呼び出しを修正

## 次のステップ
- CMP命令の実装
- メモリアドレッシングモードの実装

(詳細: [implementation_details.md](../details/implementation_details.md))

## 関連情報
[technical_notes.md](../details/technical_notes.md)

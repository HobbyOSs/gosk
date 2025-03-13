# Active Context

## 現在の作業の焦点
- ModR/Mバイト生成の改善
- オペランドのForceImm8対応

## Day02実装計画
- MOV命令 (レジスタ間, 即値)
- ADD命令 (フラグ更新)
- 他、メモリアドレッシングや制御フロー命令の実装

## 直近の変更点
- GenerateModRM関数を2種に分割
- セグメントレジスタ用MOV命令の修正

## 次のステップ
- MOV命令残機能の実装
- メモリアドレッシングモードの実装

(詳細: [implementation_details.md](../details/implementation_details.md))

## 関連情報
[technical_notes.md](../details/technical_notes.md)

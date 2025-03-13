# Progress

## 実装済み機能 (要約)
- x86アセンブラ解析基盤 (pass1, pass2)
- Ocode・PEGパーサ等の基礎部分
- システム命令 (INT, HLT)
- 算術命令(一部)

## まだ必要な実装
- MOV命令 (レジスタ間, 即値)
- ADD命令 (フラグ更新)
- CMP命令
- メモリアドレッシング
- 制御フロー命令 (JE, JMPの相対アドレス計算)
(細かな実装ステップや過去履歴は [implementation_details.md](../details/implementation_details.md) に記載)

## 関連情報
[technical_notes.md](../details/technical_notes.md)

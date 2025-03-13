# Active Context

## 現在の作業の焦点
- メモリアドレッシング

## Day02実装計画
- メモリアドレッシング

## 直近の変更点
- JMP命令のラベル解決を実装
  - pass1でラベルをテンプレート文字列としてEmit
  - pass2でテンプレート文字列をアドレスに置換
- grammar_test.go の HLT 命令関連テストを修正
  - OpcodeStmt の型に合わせてテストを修正
- CodegenClientインターフェースにGetOcodes/SetOcodesメソッドを追加

## 次のステップ
- メモリアドレッシングの実装

(詳細: [implementation_details.md](../details/implementation_details.md))

## 関連情報
[technical_notes.md](../details/technical_notes.md)

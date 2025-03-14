# Progress

## 実装済み機能 (要約)
- x86アセンブラ解析基盤 (pass1, pass2)
- Ocode・PEGパーサ等の基礎部分
- システム命令 (INT, HLT)
- 算術命令(一部)
- `internal/ocode_client/client.go` の `Exec()` メソッドで `BitMode` を呼び出し元で渡せるように修正
- CMP命令の呼び出し修正
- CMP命令のテストケース追加
- `pkg/operand/operand.go`に`InternalStrings() []string`メソッドを追加
- `pkg/operand/operand_impl.go`に`InternalStrings()`メソッドを実装し、キャッシュ処理を追加
- `pkg/asmdb/instruction_search.go`の`matchOperands`関数を修正してアキュムレータレジスタの優先検索を実装

## 実装内容
- `internal/codegen/x86gen.go` の `handleCMP` 関数呼び出しを修正
- `internal/codegen/x86gen_test.go` にCMP命令のテストケースを追加
- `pkg/operand/operand.go`に`InternalStrings() []string`メソッドを追加
- `pkg/operand/operand_impl.go`に`InternalStrings()`メソッドを実装し、キャッシュ処理を追加
- `pkg/asmdb/instruction_search.go`の`matchOperands`関数を修正してアキュムレータレジスタの優先検索を実装

## まだ必要な実装
- MOV命令 (レジスタ間, 即値)
- ADD命令 (フラグ更新)
- メモリアドレッシング
- 制御フロー命令 (JE, JMPの相対アドレス計算)
(細かな実装ステップや過去履歴は [implementation_details.md](../details/implementation_details.md) に記載)

## 関連情報
[technical_notes.md](../details/technical_notes.md)

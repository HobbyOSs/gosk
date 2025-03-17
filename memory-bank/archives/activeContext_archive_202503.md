# Active Context

## 現在の作業の焦点
- `internal/codegen` パッケージ内の関数で、`CodeGenContext` をパラメータオブジェクトとして使用するようにリファクタリング

## 直近の変更点
- day02までの実装完了
  - システム命令 (INT, HLT)
  - 算術命令(一部)
  - CMP命令の呼び出し修正とテストケース追加
  - `pkg/operand/operand.go`に`InternalStrings() []string`メソッドを追加
  - `pkg/asmdb/instruction_search.go`の`matchOperands`関数を修正
  - JE命令、MOV命令 (レジスタ間, 即値)、ADD命令 (フラグ更新)、JMP命令のラベル解決を追加
  - CodegenClientインターフェースの拡張 (GetOcodes/SetOcodesメソッドを追加)
  - CodeGenContextへのBitModeの移動
  - `internal/ocode_client/client.go` の `NewCodegenClient` 関数を修正し、エラーハンドリングを追加

## 次のステップ
- `internal/codegen` パッケージのリファクタリング完了 (CodeGenContextパラメータオブジェクト化)
- `internal/codegen` パッケージ内の不要になったパラメータを削除
- メモリアドレッシングの実装

## 関連情報
[technical_notes.md](../details/technical_notes.md)

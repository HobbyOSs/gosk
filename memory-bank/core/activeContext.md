# Active Context

## 現在の作業の焦点
- パラメータオブジェクト戦略を使ったリファクタリング

## 直近の変更点
- `internal/pass1/pass1_inst_jmp.go` の `processCalcJcc` 関数で、JMP命令のオペランドがラベルの場合に、機械語サイズを2バイト（オペコード + rel8）として計算するように修正
- `internal/codegen/x86gen.go` の `GenerateX86` 関数で、現在のOcodeの直前までの機械語長を計算し、`ctx.MachineCode` の代わりにローカル変数 `machineCode` を使用するように修正
- `internal/codegen/x86gen_jmp.go` の `generateJMPCode` 関数で、現在のアドレスを計算する際に `ctx.MachineCode` の代わりに `ctx.DollarPosition` とローカル変数 `machineCode` の長さを使用するように修正
- `internal/ocode_client/client_test.go`、`internal/frontend/frontend.go`、`internal/pass1/eval_test.go`、`test/pass1_test.go` で `NewCodegenClient` の呼び出しを修正し、`CodeGenContext` を渡すように変更
- `test/pass1_test.go` でimport文の重複を修正
- `internal/codegen/typedef.go` に `BitMode` を追加
- `internal/ocode_client/client.go` で `NewCodeGenContext` を呼び出す際に `bitMode` を渡すように変更
- `internal/ocode_client/client.go` の `Exec()` メソッドで `CodeGenContext` から `bitMode` を取得するように変更
- 関連するテストファイル(`internal/ocode_client/client_test.go`, `test/pass1_test.go`, `internal/frontend/frontend.go`)の`NewCodegenClient`呼び出しを修正

## 次のステップ
- `internal/codegen` パッケージ内の関数で、`CodeGenContext` をパラメータオブジェクトとして使用するようにリファクタリングする。
- 不要になったパラメータを削除する。

(詳細: [implementation_details.md](../details/implementation_details.md))

## 関連情報
[technical_notes.md](../details/technical_notes.md)

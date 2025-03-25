# Active Context

## 現在の作業の焦点
- OUT命令の実装

## 直近の変更点
- `internal/pass1/pass1_inst_out.go` にOUT命令の処理を追加
- `internal/pass1/handlers.go` に `processOUT` 関数を登録
- `internal/codegen/x86gen_out.go` を作成し、OUT命令のエンコード処理を実装
- `internal/codegen/x86gen_test.go` にOUT命令のテストケースを追加
- `internal/codegen/x86gen.go` に `processOcode` 関数に `ocode.OpOUT` のケースを追加

## 次のステップ
- (なし)

## 関連情報
[technical_notes.md](../details/technical_notes.md)
[implementation_details.md](../details/implementation_details.md)

# Active Context

## 現在の作業の焦点
- CALL命令の実装

## 直近の変更点
- `internal/pass1/pass1_inst_call.go` を作成し、CALL命令のpass1での処理を実装
- `internal/codegen/x86gen_call.go` を作成し、CALL命令のcodegenでの処理を実装
- `internal/pass1/handlers.go` を修正し、CALL命令のハンドラを登録
- `internal/codegen/x86gen.go` を修正し、CALL命令のcodegen処理を呼び出すように修正
- `pkg/ocode/ocode.go` を修正し、`OcodeKind` に `OpCALL` を追加
- `internal/codegen/x86gen_test.go` を修正し、CALL命令のテストケースを追加

## 次のステップ
- (なし)

## 関連情報
[technical_notes.md](../details/technical_notes.md)
[implementation_details.md](../details/implementation_details.md)

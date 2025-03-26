# Active Context

## 現在の作業の焦点
- 論理命令 (OR, XOR, NOT) の実装完了

## 直近の変更点
- OR, XOR, NOT 命令の pass1 処理を `internal/pass1/pass1_inst_logical.go` に実装
- OR, XOR, NOT 命令の codegen 処理を `internal/codegen/x86gen_logical.go` に実装 (NOTは新規ハンドラ)
- OR, XOR, NOT 命令のハンドラを `internal/pass1/handlers.go` に登録
- OR, XOR, NOT 命令の codegen 処理を `internal/codegen/x86gen.go` に登録
- `pkg/ocode/ocode.go` に `OpOR`, `OpXOR`, `OpNOT` を追加し、`go generate` を実行
- OR, XOR, NOT 命令のテストケースを `test/logical_test.go` に追加
- REXプレフィックス関連の仮実装を `handleNOT` から削除 (ユーザー指示)
- memory bank (`progress.md`, `activeContext.md`) を更新

## 次のステップ
- (なし)

## 関連情報
[technical_notes.md](../details/technical_notes.md)
[implementation_details.md](../details/implementation_details.md)

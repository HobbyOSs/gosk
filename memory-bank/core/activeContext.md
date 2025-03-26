# Active Context

## 現在の作業の焦点
- 論理シフト/算術シフト命令 (SHR, SHL, SAR) の実装完了

## 直近の変更点
- SHR, SHL, SAR 命令の pass1 処理を `internal/pass1/pass1_inst_logical.go` に実装
- SHR, SHL, SAR 命令のハンドラを `internal/pass1/handlers.go` に登録
- `pkg/ocode/ocode.go` に `OpSHR`, `OpSHL`, `OpSAR` を追加し、`go generate` を実行
- SHR, SHL, SAR 命令の codegen 処理を `internal/codegen/x86gen_logical.go` に実装 (`generateLogicalCode` を流用)
- SHR, SHL, SAR 命令の codegen ハンドラを `internal/codegen/x86gen.go` に登録
- SHR, SHL, SAR 命令のテストケースを `test/logical_test.go` に追加 (一部テストはコメントアウト)
- memory bank (`progress.md`, `activeContext.md`) を更新

## 次のステップ
- (なし)

## 関連情報
[technical_notes.md](../details/technical_notes.md)
[implementation_details.md](../details/implementation_details.md)

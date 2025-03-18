# Active Context

## 現在の作業の焦点
- JMP系命令の実装とmemory bankの更新

## 直近の変更点
- JMP系命令 (Jcc命令) の実装 (rel8のみ, rel32は未実装)
  - `pkg/ocode/ocode.go` にJMP系命令のOcodeを追加
  - `internal/codegen/x86gen.go` にJMP系命令のケースを追加
  - `internal/codegen/x86gen_jmp.go` に `handleJcc` 関数を実装

## 次のステップ
- memory bankの更新
- JMP系命令のテストコード作成
- rel32オフセットの対応 (必要な場合)

## 関連情報
[technical_notes.md](../details/technical_notes.md)
[implementation_details.md](../details/implementation_details.md)

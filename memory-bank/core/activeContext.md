# Active Context

## 現在の作業の焦点
- JMP系命令の実装とmemory bankの更新

## 直近の変更点
- JMP系命令 (Jcc命令) の実装 (rel8のみ, rel32は未実装)
  - `pkg/ocode/ocode.go` にJMP系命令のOcodeを追加
  - `internal/codegen/x86gen.go` にJMP系命令のケースを追加
  - `internal/codegen/x86gen_jmp.go` に `handleJcc` 関数を実装
- `internal/pass1/pass1_inst_arithmetic.go` の `processArithmeticInst` 関数を修正
  - アキュムレータレジスタの判定を正規表現で行うように変更
  - `regexp` パッケージをインポート

## 次のステップ
- JMP系命令のテストコード作成
- rel32オフセットの対応 (必要な場合)

## 関連情報
[technical_notes.md](../details/technical_notes.md)
[implementation_details.md](../details/implementation_details.md)

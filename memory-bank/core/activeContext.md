# Active Context

## 現在の作業の焦点
- EQU命令展開の実装とMemory Bankの更新

## 直近の変更点
- EQU命令展開処理を `internal/pass1/handlers.go` に実装
  - `TraverseAST` 関数内の `case *ast.IdentFactor:` で `env.EquMap` を参照し、EQU定義の値でオペランドを展開
  - `case *ast.MnemonicStmt:` 内のEQU展開処理を削除
- `test/day03_harib00d_test.go` のテストがPASSすることを確認

## 次のステップ
- Memory Bankの更新 (完了後)
- その他の疑似命令の展開方法の検討 (必要に応じて)

## 関連情報
[technical_notes.md](../details/technical_notes.md)
[implementation_details.md](../details/implementation_details.md)

# Active Context

## 現在の作業の焦点
- オペランド名のリファクタリング

## 直近の変更点
- `internal/ast` パッケージに `ExpToString` 関数、`FactorToString` 関数を実装
- `SegmentExp`, `AddExp`, `MultExp` の `TokenLiteral()` メソッドを修正
- `ExpToString` 関数のテストコードを実装し、テストをパス
- `pkg/operand/operand_impl.go` のオペランド名 `Addr` を `DirectMem` に、`Mem` を `IndirectMem` に変更

## 次のステップ
- オペランド名のリファクタリング完了

## 関連情報
[technical_notes.md](../details/technical_notes.md)
[implementation_details.md](../details/implementation_details.md)

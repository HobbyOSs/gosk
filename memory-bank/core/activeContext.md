# Active Context

## Current Task

- `operand` パッケージから `ng_operand` (Pigeonベース) への移行に伴う修正。
- `internal/pass1/eval_test.go` のテストをPASSさせる。

## Focus

- `internal/pass1/eval_test.go` のテスト失敗原因の調査と修正。

## Changes Made

- `go test ./internal/pass1/ -run TestPass1EvalSuite` を実行し、`ADD [BX], AX` と `MOV AL, [SI]` (16bit mode) でLOC計算が期待値(2)ではなく4になる問題を確認。
- `asmdb` のJSON定義 (`pkg/asmdb/json-x86-64/x86_64.json`) を `jq` で調査し、該当する命令形式 (`ADD m16, r16` および `MOV r8, m8`) が正しく定義されていることを確認。
- `asmdb.FindMinOutputSize` の実装 (`pkg/asmdb/instruction_search.go`) を調査。サイズ計算は `encoding.GetOutputSize`, `db.GetPrefixSize`, `operands.CalcOffsetByteSize` の合計で行われている。
- `encoding.GetOutputSize` (`pkg/asmdb/encoding.go`) と `db.GetPrefixSize` (`pkg/asmdb/instruction_search.go`) の実装を確認し、問題がないことを推定。
- `operands.CalcOffsetByteSize` (`pkg/ng_operand/operand_impl.go`) の実装を調査し、16ビットモードでアドレスサイズプレフィックスがない場合に、ディスプレースメントの有無に関わらず常に2バイトを返していたバグを発見。
- `CalcOffsetByteSize` を修正し、直接アドレス指定、ディスプレースメント付き間接アドレス指定、ディスプレースメントなし間接アドレス指定（`[BX]`, `[SI]` など）で正しいオフセットサイズ (0, 1, 2, 4 バイト) を返すように変更。
- 再度 `go test ./internal/pass1/ -run TestPass1EvalSuite` を実行し、テストがPASSすることを確認。

## Next Steps

- `progress.md` を更新する。
- タスク完了として報告する。

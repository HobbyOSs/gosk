# Active Context

## 現在の作業の焦点
- `pass1` のAST traverse処理におけるEQUマクロ展開の点検と修正

## 直近の変更点
- `getModRMFromOperands`の返り値の型変更 (`uint32` -> `[]byte`) に関連する修正は一旦保留
- `test/day03_harib00h_test.go` のEQUマクロ展開問題の調査にフォーカス

## 次のステップ
- `TraverseAST` の `case *ast.IdentFactor` で、EQUマクロ展開された値 (具体的な数値) を持つ `ImmExp` ノードで、元の `IdentFactor` ノードを**完全に置き換える** ように修正を試みる (Act Modeで実施)
- 上記修正後、`test/day03_harib00h_test.go` を実行し、EQUマクロ展開が正しく行われることを確認する
- 必要に応じて、コード生成 (`processOcode`) の点検やコンテキストスタック (`env.Ctx`) の検証を行う

## 関連情報
[technical_notes.md](../details/technical_notes.md)
[implementation_details.md](../details/implementation_details.md)

## 問題分析内容
- `test/day03_harib00h_test.go` のテスト実行時に、EQUマクロで定義されたシンボルがメモリ参照オペランド内で使用された場合に、マクロ展開が正しく行われていない
- `EquMap` へのEQU定義の登録、`IdentFactor` でのEQU値検索は正常に機能していることを確認済み
- `IdentFactor` の解決結果がコード生成 (`processOcode`) まで正しく伝搬されていない可能性
- コード生成 (`processOcode`) でメモリ参照オペランド (`[VMODE]`) を処理するロジックに誤りがある可能性
- コンテキストスタック (`env.Ctx`) のpush/pop操作の不均衡、またはAST traverse処理のロジックに潜在的な問題がある可能性

## 調査記録
- `test/day03_harib00h_test.go`のテスト実行時に以下の問題が発生
  - `MOV [0x0ff0], CH`命令のエンコーディングに失敗
  - `JMP 0xc200`命令のジャンプ先アドレスの解決に失敗  -> 修正済み: `internal/codegen/x86gen_jmp.go`で、`strconv.ParseInt`の基数を10から0に変更
  - 生成されるバイナリデータの長さが期待値より2バイト短い (512バイト、期待値は514バイト)
  - バイナリデータの内容が期待値と異なる
- `internal/pass1/handlers.go` の `TraverseAST` 関数を点検し、EQUマクロ展開処理を調査中

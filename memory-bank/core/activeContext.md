# Active Context

## 現在の作業の焦点
- EQU命令展開の実装とMemory Bankの更新

## 直近の変更点
- EQU命令展開処理を `internal/pass1/handlers.go` に実装
  - `TraverseAST` 関数内の `case *ast.IdentFactor:` で `env.EquMap` を参照し、EQU定義の値でオペランドを展開
  - `case *ast.MnemonicStmt:` 内のEQU展開処理を削除
- `test/day03_harib00d_test.go` のテストがPASSすることを確認

## 次のステップ
- Memory Bankの更新 (完了)
- `test/day03_harib00g_test.go`の修正
  - `MOV [0x0ff0], CH`命令のエンコーディング問題の調査・修正 (優先)
  - `JMP 0xc200`命令のジャンプ先アドレス解決の問題は修正済み
  - 提示された機械語と`defineHex`関数の期待値を比較し、不一致を修正

## 関連情報
[technical_notes.md](../details/technical_notes.md)
[implementation_details.md](../details/implementation_details.md)

## 調査記録
- `test/day03_harib00g_test.go`のテスト実行時に以下の問題が発生
  - `MOV [0x0ff0], CH`命令のエンコーディングに失敗
  - `JMP 0xc200`命令のジャンプ先アドレスの解決に失敗  -> 修正済み: `internal/codegen/x86gen_jmp.go`で、`strconv.ParseInt`の基数を10から0に変更
  - 生成されるバイナリデータの長さが期待値より2バイト短い (512バイト、期待値は514バイト)
  - バイナリデータの内容が期待値と異なる

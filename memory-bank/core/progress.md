# Progress

## 実装済み機能 (要約)
- x86アセンブラ解析基盤 (pass1, pass2)
- Ocode・PEGパーサ等の基礎部分
- 主要命令実装 (システム命令, 算術命令, CMP, JE, MOV, ADD, JMP)
- CodegenClient 関連機能

## まだ必要な実装
- メモリアドレッシング
- `internal/codegen` パッケージのリファクタリング (CodeGenContext 導入)
- `internal/codegen` パッケージの不要パラメータ削除

## 関連情報
[technical_notes.md](../details/technical_notes.md)

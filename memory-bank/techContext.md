# Tech Context

## 使用技術
- PEG (Parsing Expression Grammar) を用いた構文解析
- Go言語
- Gitによるバージョン管理

## 開発環境
- Linux 6.1
- VSCode

## 技術的な制約
- Go言語のバージョンに依存
- プラットフォームに依存しないコード生成

## 依存関係
- Go標準ライブラリ
- 外部ライブラリなし

## Go Test 実行ルーチン
- `go test -list '.'` でテスト関数一覧を取得
- ファイル指定でのテスト実行は非推奨
- パッケージ単位 (`go test ./...` など) で実行
- `-run` でテスト関数フィルタリング

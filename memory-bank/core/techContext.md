# Tech Context

## 使用技術
- Go言語
- Gitによるバージョン管理
- PEG (Parsing Expression Grammar) を用いた構文解析
  - アセンブリ言語全体のパース: `internal/gen` (pigeon)
  - オペランドのパース: `pkg/ng_operand` (pigeon) - 移行中 (旧: `pkg/operand` - participle)

## 開発環境
- Linux 6.1
- VSCode
- `Makefile` を使用したビルドプロセス。`make build` により、実行可能ファイル (`gosk`) がプロジェクトルートに生成される。

## 技術的な制約
- Go言語のバージョンに依存
- プラットフォームに依存しないコード生成

## 依存関係
- Go標準ライブラリ

(さらに詳しい依存ライブラリや特記事項は [technical_notes.md](../details/technical_notes.md))

# System Patterns

## システムアーキテクチャ
- モジュール化されたコンポーネント設計
- 各コンポーネントは独立してテスト可能

### コアアーキテクチャの構成
1. **Pass1とPass2の2段階処理**
   - Pass1: AST走査による評価 (`Eval`)、マクロ展開 (`EQU`)、ラベル定義、機械語サイズ計算（※要リファクタリング）
   - Pass2: 最終的なアドレス解決、機械語生成

2. **AST中心の評価 (完了済み設計変更)**
   - **背景:** 従来のトークンベース評価は、`[ LABEL + 30 ]` のような多項式的なマクロ入りオペランドの解析に課題があったため、ASTベースの評価構造へ移行。
   - **実装:**
     - 各 `ast.Exp` ノードに `Eval(env Env) (ast.Exp, bool)` メソッドを実装。
     - `Eval` は式を再帰的に評価・還元し、可能な限り単純な `ast.Exp` (例: `ast.NumberExp`) を返す。
     - 評価不能な要素（未解決ラベル等）を含む場合は、元のAST構造を保持したまま返す。
     - `internal/pass1/traverse.go` の `TraverseAST` 関数がASTを副作用なしで走査し (`node -> node` 形式)、各ノードの `Eval` を呼び出して評価・変換を行う。
   - **効果:** より複雑な式の評価と構造保持が可能になり、オペランド解析の柔軟性が向上。
   - **廃止:** 従来のスタックマシンベース (`push`/`pop`) の評価は廃止された。

3. **Ocodeによる中間表現**
   - フロントエンドとバックエンド間の橋渡し
   - 命令の抽象化と標準化
   - ラベル解決の柔軟性確保

## 主要な技術的決定
- Go言語を使用した高性能なコード解析
- ASTを用いた中間表現の生成
- asmdbを活用した命令セット情報の管理

### 標準実装パターン
1. **命令実装の基本フロー**
   - オペランドの解析と検証
   - asmdbによる命令情報の取得
   - 機械語サイズの計算
   - Ocodeの生成

2. **機械語生成の標準パターン**
   - プレフィックスの判定と追加
   - オペコードの生成
   - ModR/Mの生成（必要な場合）
   - 即値の追加（必要な場合）

## 使用している設計パターン
- ファクトリーパターンによるオブジェクト生成
- ストラテジーパターンによるアルゴリズムの選択

## プロジェクトのディレクトリ構成
```
.
├── /cmd/                # 各エントリポイント
│   └── codegen/        # コード生成関連のCLIツール (現状は空)
│
├── /internal/           # 内部実装 (外部パッケージには非公開)
│   ├── ast/            # AST (抽象構文木) 関連
│   ├── client/         # CodegenClient インターフェース定義
│   ├── codegen/        # x86 コード生成
│   ├── frontend/       # プログラムのエントリーポイント
│   ├── gen/            # PEG で記述されたパーサ
│   ├── ocode_client/   # OcodeClient 実装
│   ├── pass1/          # AST の１回めの解析（機械語サイズとラベル、マクロ）
│   ├── pass2/          # AST の後処理（ELF,COFFファイルの処理、機械語生成はcodegenで実施）
│   └── token/          # トークン定義とパース処理
│
├── /memory-bank/        # プロジェクトの知識ベース
│   ├── core/           # プロジェクトの中核となる情報
│   ├── details/        # 詳細な実装情報と技術文書
│   └── archives/       # 過去の記録のアーカイブ
│
├── /test/               # テストコード
│
├── go.mod              # Go モジュール定義
└── README.md           # プロジェクト概要 (ルートレベル)
```

## コンポーネント間の関係
- フロントエンドはASTを生成し、バックエンドに渡す
- バックエンドはASTを基に中間表現を生成し、最終コードを出力
- `pkg/asmdb` はx86命令の情報を `github.com/HobbyOSs/json-x86-64-go-mod` モジュールから取得し、オペコード、オペランド、エンコーディングなどの情報を提供する (以前は Git Submodule 内の JSON ファイルから取得)
- **オペランド受け渡し:** `pass1` は解析したオペランド情報を `ocode_client` を介して `codegen` に渡す。`CodegenClient` インターフェース (`internal/client/client.go`) がこの受け渡しを定義する。

(詳細: [implementation_details.md](../details/implementation_details.md) および [technical_notes.md](../details/technical_notes.md))

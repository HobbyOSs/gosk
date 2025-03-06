# Active Context

## 現在の作業の焦点
- 算術命令（ADD, ADC, SUB, SBB, CMP, INC, DEC, NEG, MUL, IMUL, DIV, IDIV）の実装
- Pass1の評価処理の改善とテスト強化
- トークン解析の最適化
- オペランド実装の改善

## 直近の変更点
- 算術命令の基本実装を追加（`internal/pass1/pass1_inst_arithmetic.go`）
- 算術命令のテストを一時的にスキップ（`test/arithmetic_test.go`）
- `internal/pass1/eval_test.go`のテストケース拡充
- `internal/token/parse_token.go`のトークン解析処理の改善
- `pkg/operand/operand_impl.go`にx86のprefix bytes判定機能を実装
  - オペランドサイズプレフィックス(66h)の判定
  - アドレスサイズプレフィックス(67h)の判定
- `pkg/asmdb/instruction_search.go`の機械語サイズ計算機能を改善
  - GetPrefixSize関数を追加してプレフィックスバイトのサイズ計算を実装
  - FindMinOutputSizeメソッドでプレフィックスサイズを計算に含めるように修正

## 次のステップ
1. 算術命令の実装を段階的に進める
   - Ocodeの定義
   - 機械語生成の実装
   - ModR/Mの生成
   - テストケースの修正
2. Pass1の評価処理の網羅的なテスト実装
3. トークン解析の完全性の確認と最適化
4. オペランド実装の完成度向上
5. スタックマシン関連の構造の継続的な改善

## アクティブな決定事項と考慮事項
- アセンブラ命令実装のルーチンを定義し、systemPatterns.mdに記録
  - Pass1での命令実装手順
  - Ocodeの実装手順
  - 機械語生成の実装手順
  - 実装時の注意点
- 算術命令の実装を段階的に進めるため、テストを一時的にスキップ
- テスト駆動開発の継続的な実践
- コードの品質維持のためのlintチェックとテスト実行の徹底
- オペランドの種別判定の精度向上
- スタックマシンベースの設計の最適化
- `ocode`中間言語の実装の継続

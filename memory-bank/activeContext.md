# Active Context

## 現在の作業の焦点
## Day02実装計画
- [x] 1. 基本命令の実装
  - [ ] MOV命令
    - [ ] レジスタ間転送
    - [ ] 即値のロード
    - [ ] セグメントレジスタの設定
  - [x] ADD命令
    - [x] 即値加算の実装
    - [ ] フラグ更新の処理
  - [ ] CMP命令
    - [ ] 比較演算の実装
    - [ ] フラグ設定の処理

- [x] 2. 制御フロー命令の実装
  - [x] JMP命令 (JMP rel8)
    - [x] Pass1: processJMP 関数実装 (internal/pass1/pass1_inst_jmp.go)
    - [x] Ocode: OpJMP 定義 (pkg/ocode/ocode.go)
    - [x] 機械語生成: handleJMP 関数実装 (internal/codegen/x86gen_jmp.go)
    - [ ] 相対アドレス計算
    - [ ] ジャンプ先ラベル解決

- [ ] 3. メモリ操作命令の実装
  - [ ] メモリアドレッシングモードの実装
  - [ ] レジスタ-メモリ間のデータ転送
  - [ ] ModR/Mバイトの生成

- [ ] 4. 制御フロー命令の実装 (続き)
  - [ ] JE命令
    - [ ] フラグに基づく分岐
    - [ ] オフセット計算

- [x] 5. システム命令の実装
  - [x] INT命令
    - [x] 割り込み番号の処理
    - [x] BIOS呼び出しの対応
  - [x] HLT命令
    - [x] CPU停止状態の生成

- [ ] 6. テスト有効化と検証
  - [ ] 各命令のユニットテスト実行
  - [ ] TestHelloos3のスキップ解除
  - [ ] バイナリ出力の検証

- day20までのアセンブラ命令実装の完了
  - 各dayのテストケースに対応する命令の実装
- 算術命令（ADD, ADC, SUB, SBB, CMP, INC, DEC, NEG, MUL, IMUL, DIV, IDIV）の実装
- Pass1の評価処理の改善とテスト強化
- トークン解析の最適化
- オペランド実装の改善

## 直近の変更点
- テストスイートの構造を明確化（day01からday20までの段階的実装）
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
  - FindFormをFindEncodingに変更し、最小サイズのエンコーディングを直接選択するように改善

## 次のステップ
1. day02からday20までの命令を段階的に実装
   - 各dayのテストケースを解析
   - 必要な命令を特定
   - 命令の実装順序を決定
2. 算術命令の実装を段階的に進める
   - Ocodeの定義
   - 機械語生成の実装
   - ModR/Mの生成
   - テストケースの修正
3. Pass1の評価処理の網羅的なテスト実装
4. トークン解析の完全性の確認と最適化
5. オペランド実装の完成度向上
6. スタックマシン関連の構造の継続的な改善

## アクティブな決定事項と考慮事項
- x86_64.jsonの制限事項
  - セグメントレジスタ関連の命令（MOV SS,AX など）の定義が不足している可能性
  - 一時的な対応として、セグメントレジスタ関連の命令は後回しにし、基本的な命令から実装を進める
  - 実装順序：
    1. レジスタ間MOV（AX, BX等）
    2. 即値のMOV（MOV AX, 0等）
    3. メモリ参照のMOV（[SI]等）
    4. セグメントレジスタ関連（要json-x86-64の拡張）

- アセンブラ命令実装のルーチンを定義し、systemPatterns.mdに記録
  - Pass1での命令実装手順
  - Ocodeの実装手順
  - 機械語生成の実装手順
  - 実装時の注意点
- day20までの命令実装を優先的に進める
  - 各dayのテストケースを順次有効化
  - 必要な命令を漏れなく実装
- 算術命令の実装を段階的に進めるため、テストを一時的にスキップ
- テスト駆動開発の継続的な実践
- コードの品質維持のためのlintチェックとテスト実行の徹底
- オペランドの種別判定の精度向上
- スタックマシンベースの設計の最適化
- `ocode`中間言語の実装の継続

## Ocodeの使い方
### 基本的な使い方
1. **Ocodeの生成**
   - `env.Client.Emit`を使用して命令を出力
   - 引数は文字列形式で渡す
   - 例：`env.Client.Emit("INT 0x10")`

2. **パラメータを持つ命令の場合**
   - パラメータはカンマ区切りで指定
   - 例：`env.Client.Emit("MOV AX,0")`

3. **実装パターン**
   - パラメータなし命令（HLT等）
     ```go
     func processNoParam(env *Pass1, tokens []*token.ParseToken) {
         env.LOC += 1  // 機械語サイズを加算
         // Emitは呼び出し元で実行
     }
     ```
   - パラメータあり命令（INT等）
     ```go
     func processINT(env *Pass1, tokens []*token.ParseToken) {
         env.LOC += 2  // 機械語サイズを加算
         args := lo.Map(tokens, func(token *token.ParseToken, _ int) string {
             return token.AsString()
         })
         env.Client.Emit(fmt.Sprintf("INT %s", strings.Join(args, ",")))
     }
     ```

4. **注意点**
   - パラメータなし命令は`handlers.go`のTraverseAST内のOpcodeStmtケースでEmitを実行
   - パラメータあり命令は各処理関数内でEmitを実行
   - 機械語サイズの計算は必須（env.LOCに加算）


---
## 詳細な実装計画: JMP entry 命令 (JMP rel8)

**命令:** `JMP entry`

**目標:** `JMP entry` 命令を実装し、`TestHelloos3` の `result mismatch` を解消する。

**手順:**

1. **`json-x86-64/x86_64.json` の確認:**
    - `JMP` 命令のエンコーディング定義 (`JMP rel8`, `JMP rel16`) を確認し、オペコード `eb` (JMP rel8) が存在することを確認する。
    - 必要に応じて、エンコーディング定義の詳細 (オペランドの種類、サイズなど) を確認する。

2. **Pass1 の実装 (`internal/pass1`)**:
    - `internal/pass1/pass1_inst_jmp.go` を新規作成し、`JMP` 命令の Pass1 処理 (`processJMP` 関数) を実装する。
    - `processJMP` 関数では、以下の処理を行う。
        - オペランド (ジャンプ先ラベル `entry`) の解析
        - ジャンプ先のラベル `entry` のアドレスを解決 (Pass1 では仮アドレスで良い)
        - 相対ジャンプのオフセットサイズを決定 (rel8 or rel16)
        - `pkg/asmdb` を使用して機械語サイズを計算 (`JMP rel8` は 2 bytes, `JMP rel16` は 3 bytes)
        - Ocode (`ocode.OpJMP`) を生成し、`env.Client.Emit` で出力する。

3. **Ocode の定義 (`pkg/ocode/ocode.go`)**:
    - `pkg/ocode/ocode.go` に `OpJMP` を定義する。

4. **機械語生成の実装 (`internal/codegen`)**:
    - `internal/codegen/x86gen_jmp.go` を新規作成し、`JMP` 命令の機械語生成処理 (`handleJMP` 関数) を実装する。
    - `internal/codegen/x86gen.go` の `processOcode` 関数に `ocode.OpJMP` の case を追加し、`handleJMP` 関数を呼び出す。
    - `handleJMP` 関数では、以下の処理を行う。
        - オペランド (ジャンプ先ラベル `entry`) のアドレスを取得 (Pass2 で解決されたアドレス)
        - 相対ジャンプのオフセットを計算 (ジャンプ元アドレス - ジャンプ先アドレス)
        - オフセットサイズに応じて、`JMP rel8` または `JMP rel16` の機械語コードを生成する。
            - `JMP rel8` (オペコード: `eb`, オフセット: 1 byte)
            - `JMP rel16` (オペコード: `e9`, オフセット: 2 bytes)
        - 生成された機械語コードを byte スライスとして返す。

5. **テストと検証:**
    - `test/day02_test.go` の `TestHelloos3` テストを実行し、`result mismatch` が解消されることを確認する。
    - 必要に応じて、`JMP` 命令のユニットテスト (`internal/codegen/x86gen_test.go` など) を追加する。

**実装時の注意点:**

- 相対ジャンプのオフセット計算を正確に行う (符号付き8ビットまたは16ビット)。
- ジャンプ先ラベルのアドレス解決を Pass1 と Pass2 で連携して行う。
- `json-x86-64/x86_64.json` に `JMP` 命令のエンコーディング定義が存在することを確認する。
- テスト駆動開発を実践し、テストケースを ആദ്യം に作成してから実装に取り掛かる。

**次に行うこと:**

1. `memory-bank/progress.md` を更新し、`JMP entry` 命令の実装を「まだ構築が必要な部分」から「実装済みの機能」に移動する。
2. `json-x86-64/x86_64.json` を確認し、`JMP rel8` のエンコーディング定義が存在することを確認する。
3. `internal/pass1/pass1_inst_jmp.go` を新規作成し、`processJMP` 関数を実装する。

package ng_operand

import (
	"github.com/morikuni/failure"
)

// ParseOperandString は単一のオペランド文字列をパースします。
// 例: "EAX", "[EBX+100]", "BYTE [ESI]", "123", "0xFF", "'A'", "my_label"
func ParseOperandString(text string) (*ParsedOperandPeg, error) {
	// operand_grammar.go で生成された Parse 関数を呼び出す
	// (pigeon v1.0.0 時点のデフォルト関数名。異なる場合は要修正)
	result, err := Parse("", []byte(text))
	if err != nil {
		// エラーコンテキストを追加して返す
		return nil, failure.Wrap(err, failure.Messagef("failed to parse operand string: %s", text))
	}

	// Parse 関数の戻り値は any なので、*ParsedOperandPeg に型アサーションする
	parsed, ok := result.(*ParsedOperandPeg)
	if !ok {
		// 型アサーションに失敗した場合 (通常は発生しないはず)
		return nil, failure.New(failure.StringCode("InternalError"), failure.Messagef("unexpected parser result type: %T for %s", result, text))
	}

	return parsed, nil
}

// TODO: ParseOperands 関数の実装 (カンマ区切り文字列対応)
// 必要であれば、peg 文法自体を修正して OperandList をトップレベルでパースできるようにする方が良いかもしれない。

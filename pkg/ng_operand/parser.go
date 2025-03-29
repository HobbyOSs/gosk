package ng_operand

import (
	"strings" // Add import

	"github.com/HobbyOSs/gosk/pkg/cpu" // Add import
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

// ParseOperands はカンマ区切りのオペランド文字列全体をパースします。
// 例: "EAX, EBX", "AL, [SI]", "label, 0x10"
// TODO: BitMode, ForceImm8, ForceRelAsImm をパース処理に反映させる必要がある (現在は引数で受け取るのみ)
// 戻り値を []*ParsedOperandPeg に変更
func ParseOperands(text string, bitMode cpu.BitMode, forceImm8 bool, forceRelAsImm bool) ([]*ParsedOperandPeg, error) {
	// カンマで分割する前に、文字列全体をトリムする
	trimmedText := strings.TrimSpace(text)
	if trimmedText == "" {
		// 空文字列の場合は空のスライスを返す
		return []*ParsedOperandPeg{}, nil
	}

	// カンマでオペランド文字列を分割
	// TODO: 文字列リテラル内のカンマなどを考慮する必要があるかもしれない
	parts := strings.Split(trimmedText, ",")
	parsedOperands := make([]*ParsedOperandPeg, 0, len(parts))

	for _, part := range parts {
		trimmedPart := strings.TrimSpace(part)
		if trimmedPart == "" {
			// 空の部分があればエラー (例: "EAX, , EBX")
			return nil, failure.New(failure.StringCode("InvalidFormat"), failure.Messagef("empty operand found in string: %s", text))
		}
		// 各部分を ParseOperandString でパース
		parsed, err := ParseOperandString(trimmedPart)
		if err != nil {
			// パースエラーが発生したら、エラーコンテキストを追加して返す
			return nil, failure.Wrap(err, failure.Messagef("failed to parse operand part '%s' in string: %s", trimmedPart, text))
		}
		parsedOperands = append(parsedOperands, parsed)
	}

	// パース結果を OperandPegImpl に格納して返す
	// パース結果のスライスを直接返す
	return parsedOperands, nil
}

// Removed OperandPegImpl struct definition and its methods from this file.
// They should reside in operand_impl.go.

package ng_operand

import (
	"strings" // Add import

	"github.com/HobbyOSs/gosk/pkg/cpu" // Add import
	"github.com/morikuni/failure"
)

// ParseOperandString は単一のオペランド文字列をパースします。
// 例: "EAX", "[EBX+100]", "BYTE [ESI]", "123", "0xFF", "'A'", "my_label"
func ParseOperandString(text string) (*ParsedOperandPeg, error) {
	// operand_grammar.go で生成された Parse 関数を呼び出します。
	// (pigeon v1.0.0 時点のデフォルト関数名。異なる場合は要修正)
	result, err := Parse("", []byte(text))
	if err != nil {
		// エラーコンテキストを追加して返します。
		return nil, failure.Wrap(err, failure.Messagef("オペランド文字列のパースに失敗しました: %s", text))
	}

	// Parse 関数の戻り値は any なので、*ParsedOperandPeg に型アサーションします。
	parsed, ok := result.(*ParsedOperandPeg)
	if !ok {
		// 型アサーションに失敗した場合 (通常は発生しないはず)
		return nil, failure.New(failure.StringCode("InternalError"), failure.Messagef("予期しないパーサー結果タイプ: %T for %s", result, text))
	}

	return parsed, nil
}

// ParseOperands はカンマ区切りのオペランド文字列全体をパースします。
// 例: "EAX, EBX", "AL, [SI]", "label, 0x10"
// TODO: BitMode, ForceImm8, ForceRelAsImm をパース処理に反映させる必要があります (現在は引数で受け取るのみ)。
// 戻り値を []*ParsedOperandPeg に変更しました。
func ParseOperands(text string, bitMode cpu.BitMode, forceImm8 bool, forceRelAsImm bool) ([]*ParsedOperandPeg, error) {
	// カンマで分割する前に、文字列全体をトリムします。
	trimmedText := strings.TrimSpace(text)
	if trimmedText == "" {
		// 空文字列の場合は空のスライスを返します。
		return []*ParsedOperandPeg{}, nil
	}

	// カンマでオペランド文字列を分割します。
	// TODO: 文字列リテラル内のカンマなどを考慮する必要があるかもしれません。
	parts := strings.Split(trimmedText, ",")
	parsedOperands := make([]*ParsedOperandPeg, 0, len(parts))

	for _, part := range parts {
		trimmedPart := strings.TrimSpace(part)
		if trimmedPart == "" {
			// 空の部分があればエラー (例: "EAX, , EBX")
			return nil, failure.New(failure.StringCode("InvalidFormat"), failure.Messagef("文字列内に空のオペランドが見つかりました: %s", text))
		}
		// 各部分を ParseOperandString でパースします。
		parsed, err := ParseOperandString(trimmedPart)
		if err != nil {
			// パースエラーが発生したら、エラーコンテキストを追加して返します。
			return nil, failure.Wrap(err, failure.Messagef("文字列 '%s' 内のオペランド部分 '%s' のパースに失敗しました", text, trimmedPart))
		}
		parsedOperands = append(parsedOperands, parsed)
	}

	// パース結果のスライスを直接返します。
	return parsedOperands, nil
}

// OperandPegImpl 構造体の定義とそのメソッドはこのファイルから削除されました。
// それらは operand_impl.go に配置されるべきです。

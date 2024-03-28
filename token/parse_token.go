package token

import "github.com/go-ext/variant"

type TokenType string

const (
	TTIdentifier TokenType = "ttIdentifier"
	TTNumber     TokenType = "ttNumber"
	TTHex        TokenType = "ttHex"
	// その他必要あれば追加する
)

type ParseToken struct {
	TokenType TokenType
	Data      variant.Variant
}

func NewParseToken(tokenType TokenType, v interface{}) *ParseToken {
	return &ParseToken{
		TokenType: tokenType,
		Data:      variant.New(v),
	}
}

// AsString `p.Data.ToString()` のショートハンド
func (p *ParseToken) AsString() string {
	return p.Data.ToString()
}

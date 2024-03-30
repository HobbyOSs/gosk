package token

import (
	"log"
	"strconv"

	"github.com/go-ext/variant"
	"github.com/morikuni/failure"
)

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

func (p *ParseToken) HexAsUInt() uint {
	i, err := strconv.ParseInt(p.Data.ToString()[2:], 16, 32)
	if err != nil {
		log.Fatal(failure.Wrap(err))
	}
	return uint(i)
}

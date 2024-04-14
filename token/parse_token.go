package token

import (
	"log"
	"strconv"

	"github.com/HobbyOSs/gosk/junkjit"
	"github.com/HobbyOSs/gosk/junkjit/x86"
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

func (p *ParseToken) IsNumber() bool {
	return p.TokenType == TTNumber
}

func (p *ParseToken) HexAsUInt() uint {
	i, err := strconv.ParseInt(p.Data.ToString(), 0, 64)
	if err != nil {
		log.Fatal(failure.Wrap(err))
	}
	return uint(i)
}

func (p *ParseToken) AsOperand() junkjit.Operand {
	operand, err := x86.NewX86Operand(p.Data.ToString())
	if err != nil {
		log.Fatal(failure.Wrap(err))
	}
	return operand
}

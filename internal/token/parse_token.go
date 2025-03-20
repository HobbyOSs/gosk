package token

import (
	"fmt"
	"log"
	"strconv"

	"github.com/HobbyOSs/gosk/internal/ast"
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
	Data      ast.Exp
}

func NewParseToken(tokenType TokenType, v ast.Exp) *ParseToken {
	return &ParseToken{
		TokenType: tokenType,
		Data:      v,
	}
}

func (p *ParseToken) AsString() string {
	// 内部的にast.Expの型で分岐する
	return p.Data.TokenLiteral()
}

func (p *ParseToken) IsNumber() bool {
	return p.TokenType == TTNumber
}

func (p *ParseToken) ToInt() int {
	imm, ok := p.Data.(*ast.ImmExp)
	if !ok {
		panic(fmt.Sprintf("token %s should be imm", p.Data.TokenLiteral()))
	}
	switch v := imm.Factor.(type) {
	case *ast.NumberFactor:
		return v.Value
	case *ast.HexFactor:
		return int(p.HexAsUInt())

	default:
		panic(fmt.Sprintf("token %s should be number", p.Data.TokenLiteral()))
	}
}

func (p *ParseToken) ToInt32() int32 {
	imm, ok := p.Data.(*ast.ImmExp)
	if !ok {
		panic(fmt.Sprintf("token %s should be imm", p.Data.TokenLiteral()))
	}
	switch v := imm.Factor.(type) {
	case *ast.NumberFactor:
		return int32(v.Value)
	case *ast.HexFactor:
		return int32(p.HexAsInt())

	default:
		panic(fmt.Sprintf("token %s should be number", p.Data.TokenLiteral()))
	}
}

func (p *ParseToken) ToUInt() uint {
	imm, ok := p.Data.(*ast.ImmExp)
	if !ok {
		panic(fmt.Sprintf("token %s should be imm", p.Data.TokenLiteral()))
	}
	switch v := imm.Factor.(type) {
	case *ast.NumberFactor:
		return uint(v.Value)
	case *ast.HexFactor:
		return p.HexAsUInt()
	default:
		panic(fmt.Sprintf("token %s should be number", p.Data.TokenLiteral()))
	}
}

func (p *ParseToken) HexAsUInt() uint {
	i, err := strconv.ParseInt(p.Data.TokenLiteral(), 0, 64)
	if err != nil {
		log.Fatal(failure.Wrap(err))
	}
	return uint(i)
}

func (p *ParseToken) HexAsInt() int {
	i, err := strconv.ParseInt(p.Data.TokenLiteral(), 0, 64)
	if err != nil {
		log.Fatal(failure.Wrap(err))
	}
	return int(i)
}

package token

import (
	"fmt"
	"log"
	"strconv"

	"github.com/HobbyOSs/gosk/internal/ast" // Restored ast import
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
	Data      ast.Exp // Restored ast.Exp
}

func NewParseToken(tokenType TokenType, v ast.Exp) *ParseToken { // Restored ast.Exp
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
	imm, ok := p.Data.(*ast.ImmExp) // Restored ast.ImmExp
	if !ok {
		panic(fmt.Sprintf("token %s should be imm", p.Data.TokenLiteral()))
	}
	switch v := imm.Factor.(type) {
	case *ast.NumberFactor: // Restored ast.NumberFactor
		return v.Value
	case *ast.HexFactor: // Restored ast.HexFactor
		return int(p.HexAsUInt())

	default:
		panic(fmt.Sprintf("token %s should be number", p.Data.TokenLiteral()))
	}
	// missing return statement added by compiler error, removing it
}

func (p *ParseToken) ToInt32() int32 {
	imm, ok := p.Data.(*ast.ImmExp) // Restored ast.ImmExp
	if !ok {
		panic(fmt.Sprintf("token %s should be imm", p.Data.TokenLiteral()))
	}
	switch v := imm.Factor.(type) {
	case *ast.NumberFactor: // Restored ast.NumberFactor
		return int32(v.Value)
	case *ast.HexFactor: // Restored ast.HexFactor
		return int32(p.HexAsInt())

	default:
		panic(fmt.Sprintf("token %s should be number", p.Data.TokenLiteral()))
	}
	// missing return statement added by compiler error, removing it
}

func (p *ParseToken) ToUInt() uint {
	imm, ok := p.Data.(*ast.ImmExp) // Restored ast.ImmExp
	if !ok {
		panic(fmt.Sprintf("token %s should be imm", p.Data.TokenLiteral()))
	}
	switch v := imm.Factor.(type) {
	case *ast.NumberFactor: // Restored ast.NumberFactor
		return uint(v.Value)
	case *ast.HexFactor: // Restored ast.HexFactor
		return p.HexAsUInt()
	default:
		panic(fmt.Sprintf("token %s should be number", p.Data.TokenLiteral()))
	}
	// missing return statement added by compiler error, removing it
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

// ToExp is added based on compiler error, but its implementation was missing.
// Providing a basic implementation based on the type.
func (t *ParseToken) ToExp() ast.Exp { // Changed return type to ast.Exp
	return t.Data // Assuming Data is already an ast.Exp or compatible
}

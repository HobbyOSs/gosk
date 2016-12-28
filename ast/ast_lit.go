package ast

import (
	"github.com/hangingman/gosk/token"
)

type BinaryLiteral struct {
	Token token.Token
	Value []byte
}

func (bl *BinaryLiteral) expressionNode()      {}
func (bl *BinaryLiteral) TokenLiteral() string { return bl.Token.Literal }
func (bl *BinaryLiteral) String() string       { return bl.Token.Literal }

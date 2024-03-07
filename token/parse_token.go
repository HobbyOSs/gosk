package token

import "github.com/go-ext/variant"

type ParseToken struct {
	Data variant.Variant
}

func NewParseToken(v interface{}) *ParseToken {
	return &ParseToken{Data: variant.New(v)}
}

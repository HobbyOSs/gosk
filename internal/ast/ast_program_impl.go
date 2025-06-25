package ast

import (
	"strings"

	"github.com/samber/lo"
)

//go:generate go tool newc
type Program struct {
	Statements []Statement
}

func (p Program) String() string {
	stmts := lo.Map(p.Statements, func(s Statement, _ int) string {
		return s.TokenLiteral()
	})
	return strings.Join(stmts, "\n")
}

func (p Program) program() {}
func (p Program) TokenLiteral() string {
	return ""
}
func (p Program) Type() string {
	return "Program"
}

package ast

import (
	"strings"

	"github.com/samber/lo"
)

//go:generate newc
type Program struct {
	Statements []Statement
}

func (p Program) String() string {
	stmts := lo.Map(p.Statements, func(s Statement, _ int) string {
		return s.String()
	})
	return strings.Join(stmts, "\n")
}

func (p Program) Type() string {
	return "Program"
}

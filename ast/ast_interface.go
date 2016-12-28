package ast

import (
	"bytes"
	"github.com/hangingman/gosk/token"
	"strings"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
	SetNextNode(Statement)
	SetPrevNode(Statement)
	GetNextNode() Statement
	GetPrevNode() Statement
}

type Expression interface {
	Node
	expressionNode()
}

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) String() string {
	return "{ " + i.Token.Literal + ": " + i.Value + " }"
}

type IdentifierArray struct {
	Tokens []token.Token
	Values []string
}

func (i *IdentifierArray) String() string {
	tokenLit := ""
	if len(i.Tokens) > 0 {
		tokenLit = i.Tokens[0].Literal
	}
	return "{ " + tokenLit + ": " + strings.Join(i.Values, ",") + " }"
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

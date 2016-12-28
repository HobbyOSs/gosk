package ast

import (
	"fmt"
	"github.com/comail/colog"
	"github.com/hangingman/gosk/object"
	"github.com/hangingman/gosk/token"
	"github.com/pk-rawat/gostr/src"
	"log"
	"strconv"
)

func init() {
	colog.Register()
	colog.SetDefaultLevel(colog.LInfo)
	colog.SetMinLevel(colog.LInfo)
	colog.SetFlags(log.Lshortfile)
	colog.SetFormatter(&colog.StdFormatter{Colors: false})
}

// MnemonicStatement は `MOV BX, 15` のような構文を解析する
type MnemonicStatement struct {
	Token    token.Token // OPCODE
	Name     *IdentifierArray
	NextNode Statement
	PrevNode Statement
	Bin      *object.Binary
}

func (m *MnemonicStatement) statementNode()       {}
func (m *MnemonicStatement) TokenLiteral() string { return m.Token.Literal }
func (m *MnemonicStatement) String() string {
	name := "nil"
	if isNotNil(m) && isNotNil(m.Name) {
		name = m.Name.String()
	}
	return "{ " + token.OPCODE + ":" + name + " }"
}
func (m *MnemonicStatement) SetNextNode(stmt Statement) { m.NextNode = stmt }
func (m *MnemonicStatement) SetPrevNode(stmt Statement) { m.PrevNode = stmt }
func (m *MnemonicStatement) GetNextNode() Statement     { return m.NextNode }
func (m *MnemonicStatement) GetPrevNode() Statement     { return m.PrevNode }
func (m *MnemonicStatement) HasOperator() bool {
	result := false
	for _, tok := range m.Name.Tokens {
		if tok.IsOperator() {
			result = true
			break
		}
	}
	return result
}
func (m *MnemonicStatement) PreEval() *MnemonicStatement {

	var start int = -1
	var end int = -1

	for idx, tok := range m.Name.Tokens {
		if tok.IsOperator() {
			if start == -1 {
				start = idx - 1
			}
			end = idx + 1
		}
	}

	evalStr := ""
	for i := start; i <= end; i++ {
		if m.Name.Tokens[i].Type == token.REGISTER ||
			m.Name.Tokens[i].Type == token.DOLLAR {
			return nil
		}
		if m.Name.Tokens[i].Type == token.HEX_LIT {
			hexLit := m.Name.Tokens[i].Literal
			u64v, _ := strconv.ParseUint(hexLit[2:], 16, 64)
			evalStr += fmt.Sprintf("%d", u64v)
		} else {
			evalStr += m.Name.Tokens[i].Literal
		}
	}
	log.Printf("info: expr %s\n", m.Name.Tokens)
	result := gostr.Evaluate(evalStr, nil)

	// new token array
	newTokens := []token.Token{}
	newValues := []string{}

	for i := 0; i < start; i++ {
		newTokens = append(newTokens, m.Name.Tokens[i])
		newValues = append(newValues, m.Name.Values[i])
	}
	newTokens = append(newTokens, token.Token{Type: token.INT, Literal: fmt.Sprintf("%s", result)})
	newValues = append(newValues, fmt.Sprintf("%s", result))
	for i := end + 1; i < len(m.Name.Tokens); i++ {
		newTokens = append(newTokens, m.Name.Tokens[i])
		newValues = append(newValues, m.Name.Values[i])
	}

	m.Name.Tokens = newTokens
	m.Name.Values = newValues
	log.Printf("info: expr %s\n", m.Name.Tokens)

	return nil
}

// SettingStatement は `[FORMAT "WCOFF"]` のような構文を解析する
type SettingStatement struct {
	Token    token.Token // SETTING
	Name     *Identifier
	Value    string
	NextNode Statement
	PrevNode Statement
	Bin      *object.Binary
}

func (s *SettingStatement) statementNode()       {}
func (s *SettingStatement) TokenLiteral() string { return s.Token.Literal }
func (s *SettingStatement) String() string {
	return "{ " + token.SETTING + ":" + s.Name.String() + " }"
}
func (s *SettingStatement) SetNextNode(stmt Statement) { s.NextNode = stmt }
func (s *SettingStatement) SetPrevNode(stmt Statement) { s.PrevNode = stmt }
func (s *SettingStatement) GetNextNode() Statement     { return s.NextNode }
func (s *SettingStatement) GetPrevNode() Statement     { return s.PrevNode }

// LabelStatement は `entry:` のような構文を解析する
type LabelStatement struct {
	Token    token.Token // LABEL
	Name     string
	NextNode Statement
	PrevNode Statement
}

func (s *LabelStatement) statementNode()       {}
func (s *LabelStatement) TokenLiteral() string { return s.Token.Literal }
func (s *LabelStatement) String() string {
	return "{ " + token.LABEL + ": \"" + s.Name + "\" }"
}
func (s *LabelStatement) SetNextNode(stmt Statement) { s.NextNode = stmt }
func (s *LabelStatement) SetPrevNode(stmt Statement) { s.PrevNode = stmt }
func (s *LabelStatement) GetNextNode() Statement     { return s.NextNode }
func (s *LabelStatement) GetPrevNode() Statement     { return s.PrevNode }

// EquStatement は `BOTPAK  EQU  0x00280000` のような構文を解析する
type EquStatement struct {
	Token    token.Token // EQU
	Name     *Identifier
	Value    token.Token
	NextNode Statement
	PrevNode Statement
}

func (s *EquStatement) statementNode()       {}
func (s *EquStatement) TokenLiteral() string { return s.Token.Literal }
func (s *EquStatement) String() string {
	return "{ " + token.EQU + ":" + s.Name.String() + " }"
}
func (s *EquStatement) SetNextNode(stmt Statement) { s.NextNode = stmt }
func (s *EquStatement) SetPrevNode(stmt Statement) { s.PrevNode = stmt }
func (s *EquStatement) GetNextNode() Statement     { return s.NextNode }
func (s *EquStatement) GetPrevNode() Statement     { return s.PrevNode }

type DummyStatement struct{}

func (ds *DummyStatement) String() string       { return "dummy" }
func (ds *DummyStatement) TokenLiteral() string { return "dummy-token-literal" }
func (ds *DummyStatement) statementNode()       {}

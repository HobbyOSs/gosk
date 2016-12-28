package parser

import (
	"fmt"
	"github.com/hangingman/gosk/ast"
	"github.com/hangingman/gosk/lexer"
	"github.com/hangingman/gosk/token"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseSettingStatement(t *testing.T) {
	input := `[INSTRSET "i486p"]`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	// 取得できる Statement は１つ
	assert.Equal(t, len(program.Statements), 1)
	stmt, ok := program.Statements[0].(*ast.SettingStatement)
	// Statement は SettingStatement
	assert.True(t, ok)
	assert.NotNil(t, stmt)
	// トークンの中身
	assert.Equal(t,
		token.Token{Type: token.SETTING, Literal: "INSTRSET"},
		stmt.Token)

	ident := stmt.Name
	assert.Equal(t,
		token.Token{Type: token.SETTING, Literal: "INSTRSET"},
		ident.Token)
	assert.Equal(t,
		"\"i486p\"",
		ident.Value)
}

func TestParseEquStatement(t *testing.T) {
	input := `BOTPAK	EQU		0x00280000		; bootpackのロード先`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	// 取得できる Statement は１つ
	assert.Equal(t, len(program.Statements), 1)
	stmt, ok := program.Statements[0].(*ast.EquStatement)
	// Statement は EquStatement
	assert.True(t, ok)
	assert.NotNil(t, stmt)
	// トークンの中身
	assert.Equal(t,
		token.Token{Type: token.EQU, Literal: "EQU"},
		stmt.Token)
	ident := stmt.Name
	assert.Equal(t,
		token.Token{Type: token.IDENT, Literal: "BOTPAK"},
		ident.Token)
	assert.Equal(t,
		"0x00280000",
		ident.Value)
}

func TestParseLabelStatement(t *testing.T) {
	input := `msg:   `

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	// 取得できる Statement は１つ
	assert.Equal(t, len(program.Statements), 1)
	stmt, ok := program.Statements[0].(*ast.LabelStatement)
	// Statement は EquStatement
	assert.True(t, ok)
	assert.NotNil(t, stmt)
}

func TestParseMOVStatement(t *testing.T) {
	input := `		MOV		AX,0			; レジスタ初期化
		MOV		SS,AX
		MOV		SP,0x7c00
		MOV		DS,AX
		MOV		ES,AX

		MOV		SI,msg
		MOV		AL,[SI]`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	// 取得できる Statement は１つ
	assert.Equal(t, len(program.Statements), 7)
	stmt, ok := program.Statements[0].(*ast.MnemonicStatement)
	// Statement は EquStatement
	assert.True(t, ok)
	assert.NotNil(t, stmt)
	// 中身を検査
	assert.Equal(t, token.Token{Type: token.OPCODE, Literal: "MOV"}, stmt.Token)

	//fmt.Println(stmt.Name.Tokens)

	expectedTokens := [][]token.Token{
		{
			token.Token{Type: token.OPCODE, Literal: "MOV"},
			token.Token{Type: token.REGISTER, Literal: "AX"},
			token.Token{Type: token.INT, Literal: "0"},
		},
		{
			token.Token{Type: token.OPCODE, Literal: "MOV"},
			token.Token{Type: token.SEG_REGISTER, Literal: "SS"},
			token.Token{Type: token.REGISTER, Literal: "AX"},
		},
		{
			token.Token{Type: token.OPCODE, Literal: "MOV"},
			token.Token{Type: token.REGISTER, Literal: "SP"},
			token.Token{Type: token.HEX_LIT, Literal: "0x7c00"},
		},
		{
			token.Token{Type: token.OPCODE, Literal: "MOV"},
			token.Token{Type: token.SEG_REGISTER, Literal: "DS"},
			token.Token{Type: token.REGISTER, Literal: "AX"},
		},
		{
			token.Token{Type: token.OPCODE, Literal: "MOV"},
			token.Token{Type: token.SEG_REGISTER, Literal: "ES"},
			token.Token{Type: token.REGISTER, Literal: "AX"},
		},
		{
			token.Token{Type: token.OPCODE, Literal: "MOV"},
			token.Token{Type: token.REGISTER, Literal: "SI"},
			token.Token{Type: token.IDENT, Literal: "msg"},
		},
		{
			token.Token{Type: token.OPCODE, Literal: "MOV"},
			token.Token{Type: token.REGISTER, Literal: "AL"},
			token.Token{Type: token.LBRACKET, Literal: "["},
			token.Token{Type: token.REGISTER, Literal: "SI"},
			token.Token{Type: token.RBRACKET, Literal: "]"},
		},
	}

	for i, stmt := range program.Statements {
		mnemonic, ok := stmt.(*ast.MnemonicStatement)
		assert.True(t, ok)
		assert.Equal(t,
			expectedTokens[i],
			mnemonic.Name.Tokens,
			fmt.Sprintf("%s should be ", mnemonic.Name.Tokens),
		)
	}
}

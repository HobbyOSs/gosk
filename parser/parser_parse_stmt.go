package parser

import (
	"fmt"
	"github.com/hangingman/gosk/ast"
	"github.com/hangingman/gosk/token"
	"log"
)

// ParseProgram は Parser を受け取ってAST化されたProgramを返す
func (p *Parser) ParseProgram() *ast.Program {
	log.Println("debug: ParseProgram!")
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if isNotNil(stmt) {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	// next, prevを設定する
	for index := range program.Statements {
		if index == len(program.Statements)-1 {
			program.Statements[index].SetNextNode(nil)
		} else {
			nextNode := program.Statements[index+1]

			switch nextNode.(type) {
			case *ast.MnemonicStatement:
				m := nextNode.(*ast.MnemonicStatement)
				program.Statements[index].SetNextNode(m)
			case *ast.SettingStatement:
				m := nextNode.(*ast.SettingStatement)
				program.Statements[index].SetNextNode(m)
			case *ast.LabelStatement:
				m := nextNode.(*ast.LabelStatement)
				program.Statements[index].SetNextNode(m)
			case *ast.EquStatement:
				m := nextNode.(*ast.EquStatement)
				program.Statements[index].SetNextNode(m)
			}

		}
		if index == 0 {
			program.Statements[index].SetPrevNode(nil)
		} else {
			prevNode := program.Statements[index-1]

			switch prevNode.(type) {
			case *ast.MnemonicStatement:
				m := prevNode.(*ast.MnemonicStatement)
				program.Statements[index].SetPrevNode(m)
			case *ast.SettingStatement:
				m := prevNode.(*ast.SettingStatement)
				program.Statements[index].SetPrevNode(m)
			case *ast.LabelStatement:
				m := prevNode.(*ast.LabelStatement)
				program.Statements[index].SetPrevNode(m)
			case *ast.EquStatement:
				m := prevNode.(*ast.EquStatement)
				program.Statements[index].SetPrevNode(m)
			}
		}
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken().Type {
	case token.OPCODE:
		return p.parseMnemonicStatement()
	case token.LABEL:
		return p.parseLabelStatement()
	case token.LBRACKET:
		return p.parseSettingStatement()
	case token.IDENT:
		return p.parseEquStatement()
	default:
		return nil
	}
}

func (p *Parser) parseMnemonicStatement() *ast.MnemonicStatement {

	opcodeParseFn := p.opcodeParseFns[p.curToken().Literal]
	if opcodeParseFn == nil {
		return nil
	}
	stmt := opcodeParseFn()
	log.Println(fmt.Sprintf("info: %s", stmt.String()))
	return stmt
}

func (p *Parser) parseLabelStatement() *ast.LabelStatement {
	stmt := &ast.LabelStatement{
		Token: p.curToken(),
		Name:  p.curToken().Literal,
	}
	log.Println(fmt.Sprintf("info: %s", stmt.String()))
	return stmt
}

func (p *Parser) parseSettingStatement() *ast.SettingStatement {

	p.nextToken()

	stmt := &ast.SettingStatement{
		Token: p.curToken(),
		Name: &ast.Identifier{
			Token: token.Token{Type: token.SETTING, Literal: string(p.curToken().Literal)},
			Value: p.peekToken().Literal,
		},
	}

	for !p.curTokenIs(token.RBRACKET) && !p.curTokenIs(token.EOF) {
		p.nextToken()
	}

	log.Println(fmt.Sprintf("info: %s", stmt.String()))
	return stmt
}

func (p *Parser) parseEquStatement() *ast.EquStatement {

	if !(p.curTokenIs(token.IDENT) && p.peekTokenIs(token.EQU)) {
		return nil
	}

	stmt := &ast.EquStatement{
		Token: token.Token{Type: token.EQU, Literal: "EQU"},
		Name: &ast.Identifier{
			Token: token.Token{Type: p.curToken().Type, Literal: string(p.curToken().Literal)},
		},
	}

	p.nextToken()
	pTok := p.peekToken()
	stmt.Name.Value = pTok.Literal
	stmt.Value = pTok
	log.Println(fmt.Sprintf("info: %s", stmt.String()))
	p.nextToken()

	return stmt
}

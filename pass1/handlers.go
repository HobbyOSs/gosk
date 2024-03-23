package pass1

import (
	"log"

	"github.com/HobbyOSs/gosk/ast"
	"github.com/HobbyOSs/gosk/token"
)

//go:generate newc
type ProgramHandlerImpl struct {
	Visitor *ast.Visitor
	Env     *Pass1
}

//go:generate newc
type DeclareStmtHandlerImpl struct {
	Visitor *ast.Visitor
	Env     *Pass1
}

//go:generate newc
type LabelStmtHandlerImpl struct {
	Visitor *ast.Visitor
	Env     *Pass1
}

//go:generate newc
type ExportSymStmtHandlerImpl struct {
	Visitor *ast.Visitor
	Env     *Pass1
}

//go:generate newc
type ExternSymStmtHandlerImpl struct {
	Visitor *ast.Visitor
	Env     *Pass1
}

//go:generate newc
type ConfigStmtHandlerImpl struct {
	Visitor *ast.Visitor
	Env     *Pass1
}

//go:generate newc
type MnemonicStmtHandlerImpl struct {
	Visitor *ast.Visitor
	Env     *Pass1
}

//go:generate newc
type NumberFactorHandlerImpl struct {
	Visitor *ast.Visitor
	Env     *Pass1
}

//go:generate newc
type StringFactorHandlerImpl struct {
	Visitor *ast.Visitor
	Env     *Pass1
}

//go:generate newc
type HexFactorHandlerImpl struct {
	Visitor *ast.Visitor
	Env     *Pass1
}

//go:generate newc
type IdentFactorHandlerImpl struct {
	Visitor *ast.Visitor
	Env     *Pass1
}

//go:generate newc
type CharFactorHandlerImpl struct {
	Visitor *ast.Visitor
	Env     *Pass1
}

func (p *ProgramHandlerImpl) Program(node *ast.Program) bool {
	log.Println("debug: program handler!!!")
	for _, stmt := range node.Statements {
		p.Visitor.Visit(stmt)
	}
	return true
}

func (p *DeclareStmtHandlerImpl) DeclareStmt(node *ast.DeclareStmt) bool {
	log.Println("debug: declare handler!!!")
	return true
}

func (p *LabelStmtHandlerImpl) LabelStmt(node *ast.LabelStmt) bool {
	log.Println("debug: label handler!!!")
	return true
}

func (p *ExportSymStmtHandlerImpl) ExportSymStmt(node *ast.ExportSymStmt) bool {
	log.Println("debug: export sym handler!!!")
	return true
}

func (p *ExternSymStmtHandlerImpl) ExternSymStmt(node *ast.ExternSymStmt) bool {
	log.Println("debug: extern sym handler!!!")
	return true
}

func (p *ConfigStmtHandlerImpl) ConfigStmt(node *ast.ConfigStmt) bool {
	log.Println("debug: config stmt handler!!!")
	return true
}

func (p *MnemonicStmtHandlerImpl) MnemonicStmt(node *ast.MnemonicStmt) bool {
	log.Println("debug: mnemonic stmt handler!!!")
	p.Visitor.Visit(node.Opcode)
	for _, operand := range node.Operands {
		p.Visitor.Visit(operand)
	}
	return true
}

func (p *NumberFactorHandlerImpl) NumberFactor(node *ast.NumberFactor) bool {
	log.Println("debug: number factor handler!!!")
	p.Env.Ctx.Push(token.NewParseToken(token.TTNumber, node.Value))
	return true
}

func (p *StringFactorHandlerImpl) StringFactor(node *ast.StringFactor) bool {
	log.Println("debug: string factor handler!!!")
	p.Env.Ctx.Push(token.NewParseToken(token.TTIdentifier, node.Value))
	return true
}

func (p *HexFactorHandlerImpl) HexFactor(node *ast.HexFactor) bool {
	log.Println("debug: hex factor handler!!!")
	p.Env.Ctx.Push(token.NewParseToken(token.TTHex, node.Value))
	return true
}

func (p *IdentFactorHandlerImpl) IdentFactor(node *ast.IdentFactor) bool {
	log.Println("debug: ident factor handler!!!")
	p.Env.Ctx.Push(token.NewParseToken(token.TTIdentifier, node.Value))
	return true
}

func (p *CharFactorHandlerImpl) CharFactor(node *ast.CharFactor) bool {
	log.Println("debug: char factor handler!!!")
	p.Env.Ctx.Push(token.NewParseToken(token.TTIdentifier, node.Value))
	return true
}

package ast

import (
	"strings"

	"github.com/samber/lo"
)

//go:generate go tool newc
type DeclareStmt struct {
	BaseStatement
	Id    *IdentFactor
	Value Exp // interfaceはポインタにしない
}

//go:generate go tool newc
type OpcodeStmt struct {
	BaseStatement
	Opcode *IdentFactor
}

func (os OpcodeStmt) statementNode() {}
func (os OpcodeStmt) TokenLiteral() string {
	return os.Opcode.TokenLiteral()
}

func (d DeclareStmt) statementNode() {}
func (d DeclareStmt) TokenLiteral() string {
	return d.Id.TokenLiteral() + " EQU " + d.Value.TokenLiteral()
}

//go:generate go tool newc
type LabelStmt struct {
	BaseStatement
	Label *IdentFactor
}

func (l LabelStmt) statementNode() {}
func (l LabelStmt) TokenLiteral() string {
	return l.Label.TokenLiteral()
}

//go:generate go tool newc
type ExportSymStmt struct {
	BaseStatement
	Symbols []*IdentFactor
}

func (es ExportSymStmt) statementNode() {}
func (es ExportSymStmt) TokenLiteral() string {
	symbols := lo.Map(es.Symbols, func(i *IdentFactor, _ int) string {
		return i.TokenLiteral()
	})
	return "GLOBAL " + strings.Join(symbols, ",")
}

//go:generate go tool newc
type ExternSymStmt struct {
	BaseStatement
	Symbols []*IdentFactor
}

func (es ExternSymStmt) statementNode() {}
func (es ExternSymStmt) TokenLiteral() string {
	symbols := lo.Map(es.Symbols, func(i *IdentFactor, _ int) string {
		return i.TokenLiteral()
	})
	return "EXTERN " + strings.Join(symbols, ",")
}

//go:generate go tool newc
type ConfigStmt struct {
	BaseStatement
	ConfigType ConfigType
	Factor     Factor
}

func (c ConfigStmt) statementNode() {}
func (c ConfigStmt) TokenLiteral() string {
	return string(c.ConfigType) + " " + c.Factor.TokenLiteral()
}

// TODO: go generateで作成できないか
type ConfigType string

const (
	Bits     ConfigType = "BITS"
	InstrSet ConfigType = "INSTRSET"
	Optimize ConfigType = "OPTIMIZE"
	Format   ConfigType = "FORMAT"
	Padding  ConfigType = "PADDING"
	PadSet   ConfigType = "PADSET"
	Section  ConfigType = "SECTION"
	Absolute ConfigType = "ABSOLUTE"
	File     ConfigType = "FILE"
)

var stringToConfigType = map[string]ConfigType{
	"BITS":     Bits,
	"INSTRSET": InstrSet,
	"OPTIMIZE": Optimize,
	"FORMAT":   Format,
	"PADDING":  Padding,
	"PADSET":   PadSet,
	"SECTION":  Section,
	"ABSOLUTE": Absolute,
	"FILE":     File,
}

func NewConfigType(s string) (ConfigType, bool) {
	c, ok := stringToConfigType[s]
	return c, ok
}

//go:generate go tool newc
type MnemonicStmt struct {
	BaseStatement
	Opcode   *IdentFactor
	Operands []Exp
}

func (ms MnemonicStmt) statementNode() {}
func (ms MnemonicStmt) TokenLiteral() string {
	operands := lo.Map(ms.Operands, func(f Exp, _ int) string {
		return f.TokenLiteral()
	})
	return ms.Opcode.TokenLiteral() + " " + strings.Join(operands, ",")
}

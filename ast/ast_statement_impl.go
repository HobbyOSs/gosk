package ast

import (
	"strings"

	"github.com/samber/lo"
)

//go:generate newc
type DeclareStmt struct {
	BaseStatement
	Id    *IdentFactor
	Value Exp // interfaceはポインタにしない
}

func (d DeclareStmt) String() string {
	return d.Id.String() + " EQU " + d.Value.String()
}

//go:generate newc
type LabelStmt struct {
	BaseStatement
	Label *IdentFactor
}

func (l LabelStmt) String() string {
	return l.Label.String()
}

//go:generate newc
type ExportSymStmt struct {
	BaseStatement
	Symbols []*IdentFactor
}

func (es ExportSymStmt) String() string {
	symbols := lo.Map(es.Symbols, func(i *IdentFactor, _ int) string {
		return i.String()
	})
	return "GLOBAL " + strings.Join(symbols, ",")
}

//go:generate newc
type ExternSymStmt struct {
	BaseStatement
	Symbols []*IdentFactor
}

func (es ExternSymStmt) String() string {
	symbols := lo.Map(es.Symbols, func(i *IdentFactor, _ int) string {
		return i.String()
	})
	return "EXTERN " + strings.Join(symbols, ",")
}

//go:generate newc
type ConfigStmt struct {
	BaseStatement
	ConfigType ConfigType
	Factor     Factor
}

func (c ConfigStmt) String() string {
	return string(c.ConfigType) + " " + c.Factor.String()
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

//go:generate newc
type MnemonicStmt struct {
	BaseStatement
	Opcode       *IdentFactor
	MnemonicArgs []Exp
}

func (ms MnemonicStmt) String() string {
	args := lo.Map(ms.MnemonicArgs, func(f Exp, _ int) string {
		return f.String()
	})
	return ms.Opcode.String() + " " + strings.Join(args, ",")
}

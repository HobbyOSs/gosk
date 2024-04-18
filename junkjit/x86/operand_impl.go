package x86

import (
	"regexp"

	"github.com/HobbyOSs/gosk/asmdb"
	"github.com/HobbyOSs/gosk/ast"
)

type X86Operand struct {
	A asmdb.AddressingType
}

// TODO: もし複雑になりそうだったらここでもPEGなどでパースする?
func NewX86Operand(lit string) (*X86Operand, error) {

	if re, err := regexp.Compile(`CR[0-8]`); err == nil && re.MatchString(lit) {
		return &X86Operand{A: asmdb.CodeCRField}, nil
	}
	if re, err := regexp.Compile(`DR[0-3]|DR[6-7]`); err == nil && re.MatchString(lit) {
		return &X86Operand{A: asmdb.CodeDebugField}, nil
	}
	if re, err := regexp.Compile(`CS|DS|ES|SS|FS|GS`); err == nil && re.MatchString(lit) {
		return &X86Operand{A: asmdb.CodeSregField}, nil
	}
	if re, err := regexp.Compile(`(R|E)?(A|B|C|D)X|(R|E)?(SI|DI|SP|BP)|(A|B|C|D)(H|L)|(SI|DI|SP|BP)L`); err == nil && re.MatchString(lit) {
		return &X86Operand{A: asmdb.CodeGeneralReg}, nil
	}

	return &X86Operand{A: asmdb.CodeImmediate}, nil
}

func NewX86OperandByImmExp(imm *ast.ImmExp) (*X86Operand, error) {
	lit := imm.TokenLiteral()
	return NewX86Operand(lit)
}

func (x *X86Operand) AddressingType() asmdb.AddressingType {
	return x.A
}

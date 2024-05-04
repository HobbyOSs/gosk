package x86

import (
	"regexp"

	"github.com/HobbyOSs/gosk/ast"
)

type X86Operand struct {
	A ast.AddressingType
	T ast.OperandType
}

// TODO: もし複雑になりそうだったらここでもPEGなどでパースする?
func NewX86Operand(lit string) (*X86Operand, error) {

	if re, err := regexp.Compile(`CR[0-8]`); err == nil && re.MatchString(lit) {
		return &X86Operand{A: ast.CodeCRField}, nil
	}
	if re, err := regexp.Compile(`DR[0-3]|DR[6-7]`); err == nil && re.MatchString(lit) {
		return &X86Operand{A: ast.CodeDebugField}, nil
	}
	if re, err := regexp.Compile(`CS|DS|ES|SS|FS|GS`); err == nil && re.MatchString(lit) {
		return &X86Operand{A: ast.CodeSregField}, nil
	}
	if re, err := regexp.Compile(`(R|E)?(A|B|C|D)X|(R|E)?(SI|DI|SP|BP)|(A|B|C|D)(H|L)|(SI|DI|SP|BP)L`); err == nil && re.MatchString(lit) {
		return &X86Operand{A: ast.CodeGeneralReg}, nil
	}

	return &X86Operand{A: ast.CodeImmediate}, nil
}

func NewX86OperandByImmExp(imm *ast.ImmExp) (*X86Operand, error) {
	lit := imm.TokenLiteral()
	return NewX86Operand(lit)
}

func (x *X86Operand) AddressingType() ast.AddressingType {
	return x.A
}

func (x *X86Operand) OperandType() ast.OperandType {
	return x.T
}

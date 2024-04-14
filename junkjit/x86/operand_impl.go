package x86

import (
	"github.com/HobbyOSs/gosk/asmdb"
)

type X86Operand struct {
	A asmdb.AddressingType
}

func NewX86Operand(s string) *X86Operand {
	return &X86Operand{A: asmdb.CodeGeneralReg}
}

func (x *X86Operand) AddressingType() asmdb.AddressingType {
	return x.A
}

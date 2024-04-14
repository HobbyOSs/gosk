package x86

import (
	"github.com/HobbyOSs/gosk/asmdb"
)

type X86Operand struct {
	A asmdb.AddressingType
}

// TODO: もし複雑になりそうだったらここでもPEGなどでパースする?
func NewX86Operand(s string) (*X86Operand, error) {
	return &X86Operand{A: asmdb.CodeGeneralReg}, nil
}

func (x *X86Operand) AddressingType() asmdb.AddressingType {
	return x.A
}

package x86

import "github.com/HobbyOSs/gosk/junkjit"

type X86Assembler struct {
	Code *junkjit.CodeHolder
}

func NewX86Assembler(code *junkjit.CodeHolder) *X86Assembler {
	return &X86Assembler{
		Code: code,
	}
}

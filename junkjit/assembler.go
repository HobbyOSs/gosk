package junkjit

type Assembler struct {
	Code *CodeHolder
}

func NewAssembler(code *CodeHolder) *Assembler {
	return &Assembler{
		Code: code,
	}
}

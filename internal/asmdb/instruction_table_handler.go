package asmdb

import (
	"github.com/HobbyOSs/gosk/internal/ast"
	"github.com/samber/lo"
)

type Options struct {
	BitMode   ast.BitMode
	TargetCPU ast.SupCPU
	Operand1  ast.Operand
	Operand2  ast.Operand
}

type Option func(*Options)

func BitMode(b ast.BitMode) Option {
	return func(opts *Options) {
		opts.BitMode = b
	}
}

func TargetCPU(s ast.SupCPU) Option {
	return func(opts *Options) {
		opts.TargetCPU = s
	}
}

// TODO: ここのテストを書く
func (x x86Reference) InstructionsBy(mnem string, options ...Option) []Instruction {
	filteredInst := lo.Filter(x.Instructions, func(i Instruction, _ int) bool {
		return i.Mnemonic == mnem
	})

	opts := &Options{
		BitMode:   ast.ID_16BIT_MODE,
		TargetCPU: ast.SUP_8086,
	}
	for _, option := range options {
		option(opts)
	}

	filteredInst = lo.Filter(filteredInst, func(i Instruction, _ int) bool {
		return i.IsSupported(opts.TargetCPU)
	})

	return filteredInst
}

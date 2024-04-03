package asmdb

import (
	"github.com/HobbyOSs/gosk/ast"
	"github.com/samber/lo"
)

type Options struct {
	BitMode ast.BitMode
}

type Option func(*Options)

func BitMode(b ast.BitMode) Option {
	return func(opts *Options) {
		opts.BitMode = b
	}
}

func (x x86Reference) InstructionsBy(mnem string, options ...Option) []Instruction {
	return lo.Filter(x.Instructions, func(i Instruction, _ int) bool {
		return i.Mnemonic == mnem
	})
}

func (x x86Reference) MachineCodeSize(mnem string, options ...Option) int {
	return 0
}

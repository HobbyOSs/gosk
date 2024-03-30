package junkjit

type Assembler interface {
	DB(x uint8, options ...DOption)
	DW(x uint16, options ...DOption)
	DD(x uint32, options ...DOption)
	DStruct(x any)
	CLI()
}

type DOptions struct {
	Count int
}

type DOption func(*DOptions)

func DCount(count int) DOption {
	return func(opts *DOptions) {
		opts.Count = count
	}
}

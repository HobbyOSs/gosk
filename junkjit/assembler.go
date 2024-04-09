package junkjit

// すべての命令で機械語のサイズを返す
type Assembler interface {
	BufferData() []byte

	DB(x uint8, options ...Option) int
	DW(x uint16, options ...Option) int
	DD(x uint32, options ...Option) int
	DStruct(x any) int
	CLI() int
}

type Options struct {
	Count int
}

type Option func(*Options)

func Count(count int) Option {
	return func(opts *Options) {
		opts.Count = count
	}
}

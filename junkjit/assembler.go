package junkjit

type Assembler struct {
	Code *CodeHolder
}

func NewAssembler(code *CodeHolder) *Assembler {
	return &Assembler{
		Code: code,
	}
}

// ボイラープレートコードには下記を使う
// https://github.com/switchupcb/copygen/blob/main/examples/tmpl/template/generate.tmpl
func (a *Assembler) Cli() {
	// CLI命令のオペコードは0xfa
	opcode := []byte{0xfa}
	// CodeHolderが保持するバイト配列にオペコードを追加
	a.Code.Bytes = append(a.Code.Bytes, opcode...)
}

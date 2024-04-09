package x86

// ボイラープレートコードには下記を使う
// https://github.com/switchupcb/copygen/blob/main/examples/tmpl/template/generate.tmpl
func (a *X86Assembler) CLI() int {
	opcode := []byte{0xfa}
	a.Code.Bytes = append(a.Code.Bytes, opcode...)
	return 1
}

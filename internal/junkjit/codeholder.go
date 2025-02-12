package junkjit

// CodeHolderは機械語のコードを保持する
type CodeHolder struct {
	Bytes []byte
}

// 保持している機械語コードを返す
func (ch *CodeHolder) Buffer() []byte {
	return ch.Bytes
}

// 必要に応じて　https://pkg.go.dev/go4.org/mem　などを使う
func (ch *CodeHolder) Append(b []byte) {
	ch.Bytes = append(ch.Bytes, b...)
}

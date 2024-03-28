package junkjit

// CodeHolderは機械語のコードを保持する
type CodeHolder struct {
	Bytes []byte
}

// 保持している機械語コードを返す
func (ch *CodeHolder) Buffer() []byte {
	return ch.Bytes
}

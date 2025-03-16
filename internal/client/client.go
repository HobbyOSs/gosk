package client

import (
	"github.com/HobbyOSs/gosk/pkg/ocode"
)

// CodegenClient インターフェースの定義
type CodegenClient interface {
	Emit(line string) error
	EmitAll(text string) error
	Exec() ([]byte, error)
	GetOcodes() []ocode.Ocode
	SetOcodes(ocodes []ocode.Ocode)
	SetDollarPosition(pos uint32)
	SetLOC(loc int32) // SetLOCメソッドを追加
}

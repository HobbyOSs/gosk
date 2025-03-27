package client

import (
	"github.com/HobbyOSs/gosk/internal/ast" // Import ast package
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
	SetSymbolTable(map[string]int32)
	SetBitMode(mode ast.BitMode)
}

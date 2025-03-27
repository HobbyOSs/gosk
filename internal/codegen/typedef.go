package codegen

import "github.com/HobbyOSs/gosk/pkg/operand"
import "github.com/HobbyOSs/gosk/pkg/variantstack"

type Byte uint8
type Word uint16
type DWord uint32

func ConvertToByte(value DWord) Byte {
	return Byte(value & 0xFF)
}

func ConvertToWord(value DWord) Word {
	return Word(value & 0xFFFF)
}

type CodeGenContext struct {
	SymTable       map[string]int32
	DollarPosition uint64
	MachineCode    []byte
	VS             *variantstack.VariantStack
	BitMode        operand.BitMode
}

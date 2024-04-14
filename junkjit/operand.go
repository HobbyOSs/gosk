package junkjit

import "github.com/HobbyOSs/gosk/asmdb"

type Operand interface {
	AddressingType() asmdb.AddressingType
}

package pass1

import (
	"log"

	"github.com/HobbyOSs/gosk/token"
)

func processCalcJcc(env *Pass1, tokens []*token.ParseToken) {
	arg := tokens[0]

	if arg.TokenType == token.TTIdentifier {
		env.LOC += 2
		return
	}
	dataSize := checkUintRange(arg.ToUInt())
	env.LOC += int32(dataSize)
}

func checkUintRange(value uint) int {
	switch {
	case value <= uint(^uint8(0)):
		return 2
	case value <= uint(^uint16(0)):
		return 4
	case value <= uint(^uint32(0)):
		return 6
	default:
		log.Fatal("The value is larger than uint32")
	}
	return 0
}

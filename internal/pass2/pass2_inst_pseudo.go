package pass2

import (
	"log"
	"strconv"
	"strings"

	"github.com/HobbyOSs/gosk/internal/junkjit"
	"github.com/HobbyOSs/gosk/internal/token"
	"github.com/morikuni/failure"
)

func processALIGNB(env *Pass2, tokens []*token.ParseToken) {
	//arg := tokens[0]
	//unit := arg.Data.ToInt32()
	//nearestSize := env.LOC/unit + 1
	//loc := nearestSize*unit - env.LOC
	//env.LOC += loc
}

func processDB(env *Pass2, tokens []*token.ParseToken) {
	for _, t := range tokens {
		switch t.TokenType {
		case token.TTNumber:
			env.Asm.DB(uint8(t.ToUInt()))
		case token.TTHex:
			env.Asm.DB(uint8(t.HexAsUInt()))
		case token.TTIdentifier:
			env.Asm.DStruct(t.AsString())
		}
	}
}

func processDW(env *Pass2, tokens []*token.ParseToken) {
	for _, t := range tokens {
		switch t.TokenType {
		case token.TTNumber:
			env.Asm.DW(uint16(t.ToUInt()))
		case token.TTHex:
			env.Asm.DW(uint16(t.HexAsUInt()))
		case token.TTIdentifier:
			env.Asm.DStruct(t.AsString())
		}
	}
}

func processDD(env *Pass2, tokens []*token.ParseToken) {
	for _, t := range tokens {
		switch t.TokenType {
		case token.TTNumber:
			env.Asm.DD(uint32(t.ToUInt()))
		case token.TTHex:
			env.Asm.DD(uint32(t.HexAsUInt()))
		case token.TTIdentifier:
			env.Asm.DStruct(t.AsString())
		}
	}
}

func processORG(env *Pass2, tokens []*token.ParseToken) {
	arg := tokens[0].AsString()
	currentPos, err := strconv.ParseInt(arg, 0, 32)
	if err != nil {
		log.Fatal(failure.Wrap(err))
	}
	env.DollarPos = uint32(currentPos)
}

func processRESB(env *Pass2, tokens []*token.ParseToken) {
	arg := tokens[0].AsString()
	suffix := "-$"

	if strings.Contains(arg, suffix) {
		rangeOfResb := arg[:len(arg)-len(suffix)]
		reserveSize, err := strconv.ParseInt(rangeOfResb, 0, 32)
		if err != nil {
			log.Fatal(failure.Wrap(err))
		}

		currentBufferSize := len(env.Asm.BufferData())
		needToAppendSize := reserveSize - int64(currentBufferSize)

		env.Asm.DB(0x00, junkjit.Count(int(needToAppendSize)))
		return
	}

	reserveSize, err := strconv.ParseInt(arg, 10, 32)
	if err != nil {
		log.Fatal(failure.Wrap(err))
	}
	env.Asm.DB(0x00, junkjit.Count(int(reserveSize)))
}

package pass2

import (
	"github.com/HobbyOSs/gosk/token"
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
			env.Asm.DB(uint8(t.Data.ToUInt()))
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
			env.Asm.DW(uint16(t.Data.ToUInt()))
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
			env.Asm.DD(uint32(t.Data.ToUInt()))
		case token.TTHex:
			env.Asm.DD(uint32(t.HexAsUInt()))
		case token.TTIdentifier:
			env.Asm.DStruct(t.AsString())
		}
	}
}

func processORG(env *Pass2, tokens []*token.ParseToken) {
	//arg := tokens[0].Data.ToString()
	//size, err := strconv.ParseInt(arg[2:], 16, 32)
	//if err != nil {
	// 	log.Fatal(failure.Wrap(err))
	//}
	//env.LOC = int32(size)
}

func processRESB(env *Pass2, tokens []*token.ParseToken) {
	//arg := tokens[0].Data.ToString()
	//suffix := "-$"
	//
	//if strings.Contains(arg, suffix) {
	// 	reserveSize := arg[:len(arg)-len(suffix)]
	// 	size, err := strconv.ParseInt(reserveSize[2:], 16, 32)
	// 	if err != nil {
	// 		log.Fatal(failure.Wrap(err))
	// 	}
	// 	env.LOC += int32(size)
	// 	return
	//}
	//
	//size, err := strconv.ParseInt(arg, 10, 32)
	//if err != nil {
	// 	log.Fatal(failure.Wrap(err))
	//}
	//env.LOC += int32(size)
}

package pass1

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/HobbyOSs/gosk/internal/token"
	"github.com/morikuni/failure"
)

func processALIGNB(env *Pass1, tokens []*token.ParseToken) {
	arg := tokens[0]
	unit := arg.ToInt32()
	nearestSize := env.LOC/unit + 1
	loc := nearestSize*unit - env.LOC
	env.LOC += loc
}

func processDB(env *Pass1, tokens []*token.ParseToken) {
	var loc int32 = 0
	ocodes := []int32{}

	for _, t := range tokens {
		switch t.TokenType {
		case token.TTNumber, token.TTHex:
			loc += 1
			ocodes = append(ocodes, t.ToInt32())
		case token.TTIdentifier:
			loc += int32(len(t.AsString()))
			ocodes = append(stringToOcodes(t.AsString()), ocodes...)
		}
	}

	env.LOC += loc
	emitCommand(env, "DB", ocodes)
}

func processDW(env *Pass1, tokens []*token.ParseToken) {
	var loc int32 = 0
	ocodes := []int32{}

	for _, t := range tokens {
		switch t.TokenType {
		case token.TTNumber, token.TTHex:
			loc += 2
			ocodes = append(ocodes, t.ToInt32())
		case token.TTIdentifier:
			loc += int32(len(t.AsString()))
			ocodes = append(stringToOcodes(t.AsString()), ocodes...)
		}
	}

	env.LOC += loc
	emitCommand(env, "DW", ocodes)
}

func processDD(env *Pass1, tokens []*token.ParseToken) {
	var loc int32 = 0
	ocodes := []int32{}

	for _, t := range tokens {
		switch t.TokenType {
		case token.TTNumber, token.TTHex:
			loc += 4
			ocodes = append(ocodes, t.ToInt32())
		case token.TTIdentifier:
			loc += int32(len(t.AsString()))
			ocodes = append(stringToOcodes(t.AsString()), ocodes...)
		}
	}

	env.LOC += loc
	emitCommand(env, "DD", ocodes)
}

func processORG(env *Pass1, tokens []*token.ParseToken) {
	arg := tokens[0].AsString()
	size, err := strconv.ParseInt(arg, 0, 32)
	if err != nil {
		log.Fatal(failure.Wrap(err))
	}
	env.LOC = int32(size)
}

func processRESB(env *Pass1, tokens []*token.ParseToken) {
	arg := tokens[0].AsString()

	if strings.Contains(arg, `-$`) {
		rangeOfResb := arg[:len(arg)-len(`-$`)]
		reserveSize, err := strconv.ParseInt(rangeOfResb, 0, 32)
		if err != nil {
			log.Fatal(failure.Wrap(err))
		}
		needToAppendSize := reserveSize - int64(env.LOC)
		env.LOC += int32(needToAppendSize)
		env.Client.Emit(fmt.Sprintf("RESB %d-$", reserveSize))
		return
	}

	size, err := strconv.ParseInt(arg, 10, 32)
	if err != nil {
		log.Fatal(failure.Wrap(err))
	}
	env.LOC += int32(size)
	env.Client.Emit(fmt.Sprintf("RESB %d", size))
}

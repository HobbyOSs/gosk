package pass1

import (
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

	for _, t := range tokens {
		switch t.TokenType {
		case token.TTNumber, token.TTHex:
			loc += 1
		case token.TTIdentifier:
			loc += int32(len(t.AsString()))
		}
	}

	env.LOC += loc
}

func processDW(env *Pass1, tokens []*token.ParseToken) {
	var loc int32 = 0

	for _, t := range tokens {
		switch t.TokenType {
		case token.TTNumber, token.TTHex:
			loc += 2
		case token.TTIdentifier:
			loc += int32(len(t.AsString()))
		}
	}

	env.LOC += loc
}

func processDD(env *Pass1, tokens []*token.ParseToken) {
	var loc int32 = 0

	for _, t := range tokens {
		switch t.TokenType {
		case token.TTNumber, token.TTHex:
			loc += 4
		case token.TTIdentifier:
			loc += int32(len(t.AsString()))
		}
	}

	env.LOC += loc
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
	suffix := "-$"

	if strings.Contains(arg, suffix) {
		reserveSize := arg[:len(arg)-len(suffix)]
		size, err := strconv.ParseInt(reserveSize, 0, 32)
		if err != nil {
			log.Fatal(failure.Wrap(err))
		}
		env.LOC += int32(size)
		return
	}

	size, err := strconv.ParseInt(arg, 10, 32)
	if err != nil {
		log.Fatal(failure.Wrap(err))
	}
	env.LOC += int32(size)
}

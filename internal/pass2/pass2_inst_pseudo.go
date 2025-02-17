package pass2

import (
	"log"
	"strconv"

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

func processORG(env *Pass2, tokens []*token.ParseToken) {
	arg := tokens[0].AsString()
	currentPos, err := strconv.ParseInt(arg, 0, 32)
	if err != nil {
		log.Fatal(failure.Wrap(err))
	}
	env.DollarPos = uint32(currentPos)
}

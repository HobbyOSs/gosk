package pass1

import (
	"fmt"

	"github.com/HobbyOSs/gosk/internal/token"
)

func processMOV(env *Pass1, tokens []*token.ParseToken) {
	for _, token := range tokens {
		env.Client.Emit(fmt.Sprintf("L %s\n", token.AsString()))
	}
	//env.LOC += int32(size)
	env.Client.Emit("MOV\n")
}

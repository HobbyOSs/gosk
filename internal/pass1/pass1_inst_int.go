package pass1

import (
	"fmt"

	"github.com/HobbyOSs/gosk/internal/token"
)

func processINT(env *Pass1, tokens []*token.ParseToken) {
	// Ocodeの生成
	if len(tokens) != 1 {
		panic("INT instruction requires one operand")
	}

	arg := tokens[0].AsString()
	if arg == "3" {
		env.LOC += int32(1)
	} else {
		env.LOC += int32(2)
	}

	env.Client.Emit(fmt.Sprintf("INT %s", arg))
}

package pass1

import (
	"fmt" // Keep only one fmt import
	"strings"

	"github.com/HobbyOSs/gosk/internal/token"
	"github.com/HobbyOSs/gosk/pkg/ng_operand" // Use ng_operand
	"github.com/samber/lo"
)

func processMOV(env *Pass1, tokens []*token.ParseToken) {
	args := lo.Map(tokens, func(token *token.ParseToken, _ int) string {
		return token.AsString()
	})

	// Use ng_operand.FromString factory function
	operands, err := ng_operand.FromString(strings.Join(args, ","))
	if err != nil {
		// TODO: より適切なエラーハンドリングを行う
		fmt.Printf("Error creating operands from string in MOV: %v\n", err)
		return // エラーが発生したら処理を中断
	}

	// Set BitMode and ForceRelAsImm
	operands = operands.WithBitMode(env.BitMode).
		WithForceRelAsImm(true) // Keep this flag for MOV

	// Restore LOC calculation
	size, err := env.AsmDB.FindMinOutputSize("MOV", operands)
	if err != nil {
		// TODO: より適切なエラーハンドリングを行う
		fmt.Printf("Error finding min output size for MOV %s: %v\n", strings.Join(args, ","), err)
		return // エラーが発生したら処理を中断
	}
	env.LOC += int32(size)

	deb := fmt.Sprintf("MOV %s\n", strings.Join(args, ","))
	env.Client.Emit(deb)
}

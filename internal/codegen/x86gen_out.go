package codegen

import (
	"fmt"
	"strconv"
)

func handleOUT(params x86genParams, ctx *CodeGenContext) ([]byte, error) {
	// OUT imm8, AL
	if len(params.Operands) != 2 {
		return nil, fmt.Errorf("OUT instruction requires 2 operands")
	}

	imm8, err := strconv.ParseUint(params.Operands[0], 0, 8)
	if err != nil {
		return nil, fmt.Errorf("invalid immediate value for OUT instruction: %v", err)
	}
	if params.Operands[1] != "AL" {
		return nil, fmt.Errorf("invalid operand for OUT instruction: %s", params.Operands[1])
	}

	return []byte{0xE6, byte(imm8)}, nil
}

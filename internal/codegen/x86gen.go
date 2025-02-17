package codegen

import (
	"fmt"

	"github.com/HobbyOSs/gosk/pkg/ocode"
	"github.com/HobbyOSs/gosk/pkg/variantstack"
)

// CodeGenContext はコード生成全体の状態を保持するコンテキストです。
type CodeGenContext struct {
	MachineCode []byte
	VS          *variantstack.VariantStack
}

func GenerateX86(ocodes []ocode.Ocode) []byte {
	ctx := CodeGenContext{
		MachineCode: make([]byte, 0),
		VS:          variantstack.NewVariantStack(),
	}

	for _, ocode := range ocodes {
		fmt.Printf("trace: %s\n", ocode)
		code, err := processOcode(ocode, &ctx)
		if err != nil {

		}
		ctx.MachineCode = append(ctx.MachineCode, code...)
	}

	return ctx.MachineCode
}

// processOcode processes a single ocode and returns the corresponding machine code.
func processOcode(oc ocode.Ocode, ctx *CodeGenContext) ([]byte, error) {
	currentBufferSize := len(ctx.MachineCode)

	switch oc.Kind {
	case ocode.OpL:
		return handleL(oc.Operands, ctx.VS)
	case ocode.OpDB:
		return handleDB(oc.Operands), nil
	case ocode.OpDW:
		return handleDW(oc.Operands), nil
	case ocode.OpDD:
		return handleDD(oc.Operands), nil
	case ocode.OpRESB:
		return handleRESB(oc.Operands, currentBufferSize), nil
	default:
		return handleNoParamOpcode(oc), nil
	}
}

// handleNoParamOpcode handles opcodes that do not require parameters.
func handleNoParamOpcode(ocode ocode.Ocode) []byte {
	if _, exists := opcodeMap[ocode.Kind]; exists {
		return GenerateX86NoParam(ocode)
	}
	return nil
}

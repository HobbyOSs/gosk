package codegen

import (
	"github.com/HobbyOSs/gosk/pkg/ocode"
	"github.com/HobbyOSs/gosk/pkg/variantstack"
)

func GenerateX86(ocodes []ocode.Ocode) []byte {
	var machineCode []byte
	vs := variantstack.NewVariantStack()

	for _, ocode := range ocodes {
		code, err := processOcode(ocode, vs)
		if err != nil {

		}
		machineCode = append(machineCode, code...)
	}

	return machineCode
}

// processOcode processes a single ocode and returns the corresponding machine code.
func processOcode(oc ocode.Ocode, vs *variantstack.VariantStack) ([]byte, error) {
	switch oc.Kind {
	case ocode.OpL:
		return handleL(oc.Operands, vs)
	case ocode.OpDB:
		return handleDB(oc.Operands), nil
	case ocode.OpDW:
		return handleDW(oc.Operands), nil
	case ocode.OpDD:
		return handleDD(oc.Operands), nil
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

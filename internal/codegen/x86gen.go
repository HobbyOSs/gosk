package codegen

import (
	"fmt"
	"log"

	"github.com/HobbyOSs/gosk/internal/ast"
	"github.com/HobbyOSs/gosk/pkg/ocode"
	"github.com/HobbyOSs/gosk/pkg/variantstack"
)

// x86genParams は x86 コード生成用のパラメータをまとめた構造体です。
type x86genParams struct {
	Operands       []string
	SymTable       map[string]int32
	MachineCodeLen int
	OCode          ocode.Ocode
}

func GenerateX86(ocodes []ocode.Ocode, bitMode ast.BitMode, ctx *CodeGenContext) []byte {
	ctx.VS = variantstack.NewVariantStack()
	var machineCode []byte

	log.Printf("debug: === ocode ===\n")
	for _, oc := range ocodes {
		log.Printf("debug: %s\n", oc)
		code, err := processOcode(oc, ctx, &machineCode)
		if err != nil {
			log.Printf("error: Failed to process ocode: %v", err)
		}
		machineCode = append(machineCode, code...)
	}
	log.Printf("debug: === ocode ===\n")
	ctx.MachineCode = machineCode
	return ctx.MachineCode
}

func processOcode(oc ocode.Ocode, ctx *CodeGenContext, machineCode *[]byte) ([]byte, error) {
	params := x86genParams{
		Operands:       oc.Operands,
		SymTable:       ctx.SymTable,
		MachineCodeLen: len(*machineCode),
		OCode:          oc,
	}

	log.Printf("debug: processOcode: %s, operands: %v\n", oc.Kind, oc.Operands)

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
		return handleRESB(oc.Operands, params.MachineCodeLen), nil
	case ocode.OpMOV:
		return handleMOV(params.Operands, ctx), nil
	case ocode.OpINT:
		return handleINT(oc), nil
	case ocode.OpJMP:
		return handleJMP(params, ctx)
	case ocode.OpJE:
		return handleJE(params, ctx)
	case ocode.OpADD:
		return handleADD(params, ctx)
	case ocode.OpCMP:
		return handleCMP(params, ctx)
	case ocode.OpHLT, ocode.OpNOP, ocode.OpRETN:
		return handleNoParamOpcode(oc), nil
	default:
		return nil, fmt.Errorf("not implemented: %v", oc.Kind)
	}
}

func handleNoParamOpcode(ocode ocode.Ocode) []byte {
	log.Printf("debug: handleNoParamOpcode: %s\n", ocode.Kind)
	if _, exists := opcodeMap[ocode.Kind]; exists {
		return GenerateX86NoParam(ocode)
	}
	return nil
}

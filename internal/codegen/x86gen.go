package codegen

import (
	"fmt"
	"log"

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

func GenerateX86(ocodes []ocode.Ocode, ctx *CodeGenContext) []byte {
	ctx.VS = variantstack.NewVariantStack()
	// machineCode を非 nil の空スライスで初期化
	machineCode := make([]byte, 0)

	log.Printf("debug: [codegen] === ocode processing start ===\n")
	for _, oc := range ocodes {
		log.Printf("debug: [codegen] Processing ocode: %s\n", oc)
		code, err := processOcode(oc, ctx, &machineCode)
		if err != nil {
			log.Printf("error: Failed to process ocode: %v", err)
		}
		machineCode = append(machineCode, code...)
	}
	log.Printf("debug: [codegen] === ocode processing end ===\n")
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

	// Check if the instruction is a no-parameter instruction handled by opcodeMap
	if _, exists := opcodeMap[oc.Kind]; exists {
		return handleNoParamOpcode(oc), nil
	}

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
		return handleRESB(oc.Operands, params, ctx), nil
	case ocode.OpMOV:
		return handleMOV(params.Operands, ctx), nil
	case ocode.OpINT:
		return handleINT(oc), nil
	case ocode.OpJMP, ocode.OpJMP_FAR, ocode.OpJE, ocode.OpJA, ocode.OpJAE, ocode.OpJB, ocode.OpJBE, ocode.OpJC, ocode.OpJG, ocode.OpJGE, ocode.OpJL, ocode.OpJLE,
		ocode.OpJNA, ocode.OpJNAE, ocode.OpJNB, ocode.OpJNBE, ocode.OpJNC, ocode.OpJNE, ocode.OpJNG, ocode.OpJNGE, ocode.OpJNL,
		ocode.OpJNLE, ocode.OpJNO, ocode.OpJNP, ocode.OpJNS, ocode.OpJNZ, ocode.OpJO, ocode.OpJP, ocode.OpJPE, ocode.OpJPO,
		ocode.OpJS, ocode.OpJZ:
		return handleJcc(params, ctx)
	case ocode.OpADD:
		return handleADD(params, ctx)
	case ocode.OpCMP:
		return handleCMP(params, ctx)
	case ocode.OpIMUL:
		return handleIMUL(params, ctx)
	case ocode.OpSUB:
		return handleSUB(params, ctx)
	case ocode.OpAND:
		return handleAND(params, ctx)
	case ocode.OpOR:
		return handleOR(params, ctx)
	case ocode.OpXOR:
		return handleXOR(params, ctx)
	case ocode.OpNOT:
		return handleNOT(params, ctx)
	case ocode.OpSHR:
		return handleSHR(params, ctx)
	case ocode.OpSHL:
		return handleSHL(params, ctx)
	case ocode.OpSAR:
		return handleSAR(params, ctx)
	case ocode.OpRET:
		return handleRET(oc)
	case ocode.OpIN:
		return handleIN(params, ctx)
	case ocode.OpOUT:
		return handleOUT(params, ctx)
	case ocode.OpPUSH:
		return handlePUSH(params, ctx)
	case ocode.OpPOP:
		return handlePOP(params, ctx)
	case ocode.OpCALL:
		return handleCALL(params, ctx)
	case ocode.OpLGDT:
		return handleLGDT(params.Operands, ctx)
	case ocode.OpALIGNB:
		return handleALIGNB(params, ctx)
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

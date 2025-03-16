package codegen

import (
	"fmt"
	"log"

	"github.com/HobbyOSs/gosk/internal/ast"
	"github.com/HobbyOSs/gosk/pkg/ocode"
	"github.com/HobbyOSs/gosk/pkg/variantstack"
)

// CodeGenContext はコード生成全体の状態を保持するコンテキストです。
type CodeGenContext struct {
	MachineCode    []byte
	VS             *variantstack.VariantStack
	BitMode        ast.BitMode
	DollarPosition uint32 // エントリーポイントのアドレス
	LOC            int32  // Location Counter
}

func GenerateX86(ocodes []ocode.Ocode, bitMode ast.BitMode, ctx *CodeGenContext) []byte {
	ctx.VS = variantstack.NewVariantStack()
    var machineCode []byte // ローカルのmachineCode変数を追加

	log.Printf("debug: === ocode ===\n")
	for _, ocode := range ocodes {
		log.Printf("debug: %s\n", ocode)
		code, err := processOcode(ocode, ctx)
		if err != nil {
			log.Printf("error: Failed to process ocode: %v", err)
		}
		machineCode = append(machineCode, code...) // ローカル変数に追加
	}
	log.Printf("debug: === ocode ===\n")
    ctx.MachineCode = machineCode // 最後にctx.MachineCodeに結合
	return ctx.MachineCode
}

// processOcode processes a single ocode and returns the corresponding machine code.
func processOcode(oc ocode.Ocode, ctx *CodeGenContext) ([]byte, error) {
	currentBufferSize := len(ctx.MachineCode)

	log.Printf("debug: processOcode: %s, operands: %v\n", oc.Kind, oc.Operands) // デバッグログを追加

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
	case ocode.OpMOV:
		return handleMOV(oc.Operands, ctx), nil
	case ocode.OpINT:
		return handleINT(oc), nil
	case ocode.OpJMP:
		return handleJMP(oc, ctx)
	case ocode.OpJE:
		return handleJE(oc, ctx)
	case ocode.OpADD:
		return handleADD(oc.Operands, ctx)
	case ocode.OpCMP:
		return handleCMP(oc.Operands, ctx)
	case ocode.OpHLT, ocode.OpNOP, ocode.OpRETN:
		return handleNoParamOpcode(oc), nil
	default:
		return nil, fmt.Errorf("not implemented: %v", oc.Kind)
	}
}

// handleNoParamOpcode handles opcodes that do not require parameters.
func handleNoParamOpcode(ocode ocode.Ocode) []byte {
	log.Printf("debug: handleNoParamOpcode: %s\n", ocode.Kind)
	if _, exists := opcodeMap[ocode.Kind]; exists {
		return GenerateX86NoParam(ocode)
	}
	return nil
}

package pass1

import (
	"log" // Import log package

	"github.com/HobbyOSs/gosk/internal/ast" // Add ast import
	"github.com/samber/lo"
)

// opcodeEvalFn defines the function signature for instruction handlers in pass1.
// It takes the Pass1 environment and a slice of evaluated operand expressions.
type opcodeEvalFn func(env *Pass1, operands []ast.Exp)

var (
	// opcodeEvalFns maps mnemonic strings to their corresponding handler functions.
	opcodeEvalFns = make(map[string]opcodeEvalFn, 0)
)

func init() {
	// Assign the refactored handlers to the map.

	// 疑似命令
	opcodeEvalFns["ALIGNB"] = processALIGNB // Refactored
	opcodeEvalFns["DB"] = processDB         // Refactored
	opcodeEvalFns["DD"] = processDD         // Refactored
	opcodeEvalFns["DW"] = processDW         // Refactored
	opcodeEvalFns["ORG"] = processORG       // Refactored
	opcodeEvalFns["RESB"] = processRESB     // Refactored

	// Jump命令
	jmpOps := []string{
		"JA", "JAE", "JB", "JBE", "JC", "JE", "JG", "JGE", "JL", "JLE", "JMP", "JNA", "JNAE",
		"JNB", "JNBE", "JNC", "JNE", "JNG", "JNGE", "JNL", "JNLE", "JNO", "JNP", "JNS", "JNZ",
		"JO", "JP", "JPE", "JPO", "JS", "JZ",
	}
	jmpFns := lo.SliceToMap(
		jmpOps,
		func(op string) (string, opcodeEvalFn) {
			// Assign the refactored processCalcJcc
			localOp := op // Capture loop variable
			return localOp, func(env *Pass1, operands []ast.Exp) {
				processCalcJcc(env, operands, localOp) // Call refactored handler
			}
		},
	)
	opcodeEvalFns = lo.Assign(opcodeEvalFns, jmpFns)

	// No-parameter instructions
	noParamOps := []string{ // Keep noParamOps for the loop below
		"AAA", "AAD", "AAM", "AAS", "CBW", "CDQ", "CDQE", "CLC", "CLD", "CLI", "CLTS", "CMC",
		"CPUID", "CQO", "CS", "CWD", "CWDE", "DAA", "DAS", "DIV", "DS", "EMMS", "ENTER", "ES",
		"F2XM1", "FABS", "FADDP", "FCHS", "FCLEX", "FCOM", "FCOMP", "FCOMPP", "FCOS", "FDECSTP",
		"FDISI", "FDIVP", "FDIVRP", "FENI", "FINCSTP", "FINIT", "FLD1", "FLDL2E", "FLDL2T",
		"FLDLG2", "FLDLN2", "FLDPI", "FLDZ", "FMULP", "FNCLEX", "FNDISI", "FNENI", "FNINIT",
		"FNOP", "FNSETPM", "FPATAN", "FPREM", "FPREM1", "FPTAN", "FRNDINT", "FRSTOR", "FS",
		"FSCALE", "FSETPM", "FSIN", "FSINCOS", "FSQRT", "FSUBP", "FSUBRP", "FTST", "FUCOM",
		"FUCOMP", "FUCOMPP", "FXAM", "FXCH", "FXRSTOR", "FXTRACT", "FYL2X", "FYL2XP1", "GETSEC",
		"GS", "HLT", "ICEBP", "IDIV", "IMUL", "INTO", "INVD", "IRET", "IRETD", "IRETQ", "JMPE",
		"LAHF", "LEAVE", "LFENCE", "LOADALL", "LOCK", "MFENCE", "MONITOR", "MUL", "MWAIT", "NOP",
		"PAUSE", "POPA", "POPAD", "POPF", "POPFD", "POPFQ", "PUSHA", "PUSHAD", "PUSHF", "PUSHFD",
		"PUSHFQ", "RDMSR", "RDPMC", "RDTSC", "RDTSCP", "REP", "REPE", "REPNE", "RETF", "RETN",
		"RSM", "SAHF", "SETALC", "SFENCE", "SS", "STC", "STD", "STI", "SWAPGS", "SYSCALL",
		"SYSENTER", "SYSEXIT", "SYSRET", "TAKEN", "UD2", "VMCALL", "VMLAUNCH", "VMRESUME",
		"VMXOFF", "WAIT", "WBINVD", "WRMSR", "XGETBV", "XRSTOR", "XSETBV",
	}
	// Assign placeholder for no-parameter instructions (processNoParam needs refactoring)
	for _, op := range noParamOps {
		localOp := op
		opcodeEvalFns[localOp] = func(env *Pass1, operands []ast.Exp) {
			// TODO: Refactor processNoParam and call it here.
			log.Printf("TODO: Refactor processNoParam for %s and call it.", localOp)
			// Placeholder logic:
			// size := calculateNoParamInstructionSize(env, localOp) // Placeholder
			// env.Client.Emit(localOp) // Placeholder
			// env.LOC += size
		}
	}
	/* Remove the old SliceToMap block for noParamFns
	noParamFns := lo.SliceToMap( // Remove this variable declaration
		noParamOps,
		// processNoParam のシグネチャが変わるため、一旦コメントアウト。後で修正する。
		// func(op string) (string, opcodeEvalFn) {
		// 	return op, processNoParam
		// },
		func(op string) (string, opcodeEvalFn) {
			// TODO: Implement proper processNoParam with new signature
			return op, func(env *Pass1, operands []ast.Exp) {
				// Placeholder implementation for processNoParam
				log.Printf("TODO: Implement processNoParam for %s with new signature.", op) // Use Printf
				// Original processNoParam logic needs to be adapted here.
				// It likely involves calling env.Client.EmitXXX() without operands
			// and calculating the instruction size.
			// For now, just log and maybe update LOC with a placeholder size.
			// size := calculateNoParamInstructionSize(env, localOp) // Placeholder function
			// env.Client.EmitXXX(localOp) // Placeholder call
			// env.LOC += size
		}
	}
	/* Remove the old SliceToMap block for noParamFns
	noParamFns := lo.SliceToMap( // Remove this variable declaration
		noParamOps,
		// processNoParam のシグネチャが変わるため、一旦コメントアウト。後で修正する。
		// func(op string) (string, opcodeEvalFn) {
		// 	return op, processNoParam
		// },
		func(op string) (string, opcodeEvalFn) {
			// TODO: Implement proper processNoParam with new signature
			return op, func(env *Pass1, operands []ast.Exp) {
				// Placeholder implementation for processNoParam
				log.Printf("TODO: Implement processNoParam for %s with new signature.", op) // Use Printf
				// Original processNoParam logic needs to be adapted here.
				// It likely involves calling env.Client.EmitXXX() without operands
				// and calculating the instruction size.
				// For now, just log and maybe update LOC with a placeholder size.
				// size := calculateNoParamInstructionSize(env, op) // Placeholder function
				// env.Client.EmitXXX(op) // Placeholder call
				// env.LOC += size
		},
	)
	*/ // End of removed SliceToMap block

	// Assign refactored handlers for remaining instructions
	opcodeEvalFns["RET"] = processRET   // Refactored,
	opcodeEvalFns["MOV"] = processMOV   // Refactored,
	opcodeEvalFns["INT"] = processINT   // Refactored,
	opcodeEvalFns["ADD"] = processADD   // Refactored (uses helper),
	opcodeEvalFns["ADC"] = processADC   // Refactored (uses helper),
	opcodeEvalFns["SUB"] = processSUB   // Refactored (uses helper),
	opcodeEvalFns["SBB"] = processSBB   // Refactored (uses helper),
	opcodeEvalFns["CMP"] = processCMP   // Refactored (uses helper),
	opcodeEvalFns["INC"] = processINC   // Refactored (uses helper),
	opcodeEvalFns["DEC"] = processDEC   // Refactored (uses helper),
	opcodeEvalFns["NEG"] = processNEG   // Refactored (uses helper),
	opcodeEvalFns["MUL"] = processMUL   // Refactored (uses helper),
	opcodeEvalFns["IMUL"] = processIMUL // Refactored,
	opcodeEvalFns["DIV"] = processDIV   // Refactored (uses helper),
	opcodeEvalFns["IDIV"] = processIDIV // Refactored (uses helper),
	opcodeEvalFns["AND"] = processAND   // Refactored (uses helper),
	opcodeEvalFns["OR"] = processOR     // Refactored (uses helper),
	opcodeEvalFns["XOR"] = processXOR   // Refactored (uses helper),
	opcodeEvalFns["NOT"] = processNOT   // Refactored,
	opcodeEvalFns["SHR"] = processSHR   // Refactored (uses helper),
	opcodeEvalFns["SHL"] = processSHL   // Refactored (uses helper),
	opcodeEvalFns["SAR"] = processSAR   // Refactored (uses helper),
	opcodeEvalFns["IN"] = processIN     // Refactored,
	opcodeEvalFns["OUT"] = processOUT   // Refactored,
	opcodeEvalFns["CALL"] = processCALL // Refactored,
	opcodeEvalFns["LGDT"] = processLGDT // Refactored,
}

// --- TraverseAST function moved to traverse.go ---

// --- DefineMacro and LookupMacro moved to traverse.go (as methods on Pass1) ---
// Note: They are now methods on Pass1 in traverse.go to implement ast.Env

// TODO: Refactor all processXXX functions (opcodeEvalFn implementations)
// They need to be adapted to the new TraverseAST model:
// - They should not use the stack (env.Ctx).
// - They should accept evaluated operands (likely []ast.Exp).
// - They need to handle size calculation based on evaluated operands.
// - They need to interact with the CodegenClient using evaluated data.

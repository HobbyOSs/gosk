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
	// Assign placeholder functions with the new signature to fix compiler errors.
	// Actual logic implementation will follow.

	// 疑似命令
	assignPlaceholderHandler("ALIGNB")
	assignPlaceholderHandler("DB")
	assignPlaceholderHandler("DD")
	assignPlaceholderHandler("DW")
	assignPlaceholderHandler("ORG")
	assignPlaceholderHandler("RESB")

	// Jump命令 (processCalcJcc needs refactoring)
	jmpOps := []string{
		"JA", "JAE", "JB", "JBE", "JC", "JE", "JG", "JGE", "JL", "JLE", "JMP", "JNA", "JNAE",
		"JNB", "JNBE", "JNC", "JNE", "JNG", "JNGE", "JNL", "JNLE", "JNO", "JNP", "JNS", "JNZ",
		"JO", "JP", "JPE", "JPO", "JS", "JZ",
	}
	jmpFns := lo.SliceToMap(
		jmpOps,
		func(op string) (string, opcodeEvalFn) {
			// Placeholder for Jcc instructions
			return op, func(env *Pass1, operands []ast.Exp) {
				log.Printf("TODO: Implement Jcc handler for %s with new signature.", op) // Use Printf
				// Original processCalcJcc logic needs adaptation.
				// It involves evaluating the target operand (operands[0])
				// and calculating the relative offset based on env.LOC.
				// size := calculateJccInstructionSize(env, op, operands[0]) // Placeholder
				// env.Client.EmitJcc(op, target) // Placeholder
				// env.LOC += size
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
	// Assign placeholders for no-parameter instructions
	for _, op := range noParamOps {
		// Capture op in the closure
		localOp := op
		opcodeEvalFns[localOp] = func(env *Pass1, operands []ast.Exp) {
			// Placeholder implementation for processNoParam
			log.Printf("TODO: Implement processNoParam for %s with new signature.", localOp) // Use Printf
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

	// Assign placeholders for remaining instructions
	assignPlaceholderHandler("RET")
	assignPlaceholderHandler("MOV")
	assignPlaceholderHandler("INT")
	assignPlaceholderHandler("ADD")
	assignPlaceholderHandler("ADC")
	assignPlaceholderHandler("SUB")
	assignPlaceholderHandler("SBB")
	assignPlaceholderHandler("CMP")
	assignPlaceholderHandler("INC")
	assignPlaceholderHandler("DEC")
	assignPlaceholderHandler("NEG")
	assignPlaceholderHandler("MUL")
	assignPlaceholderHandler("IMUL")
	assignPlaceholderHandler("DIV")
	assignPlaceholderHandler("IDIV")
	assignPlaceholderHandler("AND")
	assignPlaceholderHandler("OR")
	assignPlaceholderHandler("XOR")
	assignPlaceholderHandler("NOT")
	assignPlaceholderHandler("SHR")
	assignPlaceholderHandler("SHL")
	assignPlaceholderHandler("SAR")
	assignPlaceholderHandler("IN")
	assignPlaceholderHandler("OUT")
	assignPlaceholderHandler("CALL")
	assignPlaceholderHandler("LGDT")
}

// assignPlaceholderHandler assigns a placeholder function with the new signature.
func assignPlaceholderHandler(mnemonic string) {
	opcodeEvalFns[mnemonic] = func(env *Pass1, operands []ast.Exp) {
		log.Printf("TODO: Implement handler for %s with new signature.", mnemonic) // Use Printf
		// Placeholder: Calculate size based on operands and update LOC
		// size := calculateInstructionSize(env, mnemonic, operands) // Placeholder
		// env.Client.EmitXXX(...) // Placeholder
		// env.LOC += size
	}
}

// --- Stack operations removed ---

// --- TraverseAST function moved to traverse.go ---

// --- DefineMacro and LookupMacro moved to traverse.go (as methods on Pass1) ---
// Note: They are now methods on Pass1 in traverse.go to implement ast.Env

// TODO: Refactor all processXXX functions (opcodeEvalFn implementations)
// They need to be adapted to the new TraverseAST model:
// - They should not use the stack (env.Ctx).
// - They should accept evaluated operands (likely []ast.Exp).
// - They need to handle size calculation based on evaluated operands.
// - They need to interact with the CodegenClient using evaluated data.

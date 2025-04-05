package pass1

import (
	"github.com/HobbyOSs/gosk/internal/ast" // Add ast import
	"github.com/samber/lo"
)

// opcodeEvalFn は pass1 における命令ハンドラの関数シグネチャを定義します。
// Pass1 環境と評価済みオペランド式のスライスを受け取ります。
type opcodeEvalFn func(env *Pass1, operands []ast.Exp)

var (
	// opcodeEvalFns はニーモニック文字列を対応するハンドラ関数にマッピングします。
	opcodeEvalFns = make(map[string]opcodeEvalFn, 0)
)

func init() {
	// リファクタリングされたハンドラをマップに割り当てます。

	// 疑似命令
	opcodeEvalFns["ALIGNB"] = processALIGNB
	opcodeEvalFns["DB"] = processDB
	opcodeEvalFns["DD"] = processDD
	opcodeEvalFns["DW"] = processDW
	opcodeEvalFns["ORG"] = processORG
	opcodeEvalFns["RESB"] = processRESB
	opcodeEvalFns["GLOBAL"] = processGLOBAL // Add GLOBAL handler
	opcodeEvalFns["EXTERN"] = processEXTERN // Add EXTERN handler

	// Jump命令
	jmpOps := []string{
		"JA", "JAE", "JB", "JBE", "JC", "JE", "JG", "JGE", "JL", "JLE", "JMP", "JNA", "JNAE",
		"JNB", "JNBE", "JNC", "JNE", "JNG", "JNGE", "JNL", "JNLE", "JNO", "JNP", "JNS", "JNZ",
		"JO", "JP", "JPE", "JPO", "JS", "JZ",
	}
	jmpFns := lo.SliceToMap(
		jmpOps,
		func(op string) (string, opcodeEvalFn) {
			// リファクタリングされた processCalcJcc を割り当てます
			localOp := op // ループ変数をキャプチャ
			return localOp, func(env *Pass1, operands []ast.Exp) {
				processCalcJcc(env, operands, localOp) // リファクタリングされたハンドラを呼び出す
			}
		},
	)
	opcodeEvalFns = lo.Assign(opcodeEvalFns, jmpFns)

	// パラメータなし命令
	noParamOps := []string{ // 以下のループのために noParamOps を保持
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
	// パラメータなし命令に processNoParam を呼び出すクロージャを割り当てます
	for _, op := range noParamOps {
		localOp := op // ループ変数をキャプチャ
		opcodeEvalFns[localOp] = func(env *Pass1, operands []ast.Exp) {
			processNoParam(env, operands, localOp) // 命令名を渡してハンドラを呼び出す
		}
	}

	// 残りの命令にリファクタリングされたハンドラを割り当てます
	opcodeEvalFns["RET"] = processRET
	opcodeEvalFns["MOV"] = processMOV
	opcodeEvalFns["INT"] = processINT
	opcodeEvalFns["ADD"] = processADD
	opcodeEvalFns["ADC"] = processADC
	opcodeEvalFns["SUB"] = processSUB
	opcodeEvalFns["SBB"] = processSBB
	opcodeEvalFns["CMP"] = processCMP
	opcodeEvalFns["INC"] = processINC
	opcodeEvalFns["DEC"] = processDEC
	opcodeEvalFns["NEG"] = processNEG
	opcodeEvalFns["MUL"] = processMUL
	opcodeEvalFns["IMUL"] = processIMUL
	opcodeEvalFns["DIV"] = processDIV
	opcodeEvalFns["IDIV"] = processIDIV
	opcodeEvalFns["AND"] = processAND
	opcodeEvalFns["OR"] = processOR
	opcodeEvalFns["XOR"] = processXOR
	opcodeEvalFns["NOT"] = processNOT
	opcodeEvalFns["SHR"] = processSHR
	opcodeEvalFns["SHL"] = processSHL
	opcodeEvalFns["SAR"] = processSAR
	opcodeEvalFns["IN"] = processIN
	opcodeEvalFns["OUT"] = processOUT
	opcodeEvalFns["CALL"] = processCALL
	opcodeEvalFns["LGDT"] = processLGDT
	opcodeEvalFns["LIDT"] = processLIDT // Add LIDT handler
	opcodeEvalFns["PUSH"] = processPUSH // Add PUSH handler
	opcodeEvalFns["POP"] = processPOP   // Add POP handler
}

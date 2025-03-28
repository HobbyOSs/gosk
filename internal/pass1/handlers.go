package pass1

import (
	"log"
	"strings"

	"github.com/HobbyOSs/gosk/internal/ast" // Restored ast import
	// "github.com/HobbyOSs/gosk/internal/client" // Removed unused import
	"github.com/HobbyOSs/gosk/internal/token"
	"github.com/HobbyOSs/gosk/pkg/cpu" // Keep cpu import for NewBitMode
	"github.com/morikuni/failure"
	"github.com/samber/lo"
)

type opcodeEvalFn func(*Pass1, []*token.ParseToken)

var (
	opcodeEvalFns = make(map[string]opcodeEvalFn, 0)
)

func init() {
	// 疑似命令
	opcodeEvalFns["ALIGNB"] = processALIGNB
	opcodeEvalFns["DB"] = processDB
	opcodeEvalFns["DD"] = processDD
	opcodeEvalFns["DW"] = processDW
	opcodeEvalFns["ORG"] = processORG
	opcodeEvalFns["RESB"] = processRESB

	// Jump命令
	jmpOps := []string{
		"JA", "JAE", "JB", "JBE", "JC", "JE", "JG", "JGE", "JL", "JLE", "JMP", "JNA", "JNAE",
		"JNB", "JNBE", "JNC", "JNE", "JNG", "JNGE", "JNL", "JNLE", "JNO", "JNP", "JNS", "JNZ",
		"JO", "JP", "JPE", "JPO", "JS", "JZ",
	}
	jmpFns := lo.SliceToMap(
		jmpOps,
		func(op string) (string, opcodeEvalFn) {
			return op, func(env *Pass1, tokens []*token.ParseToken) {
				processCalcJcc(env, tokens, op)
			}
		},
	)
	opcodeEvalFns = lo.Assign(opcodeEvalFns, jmpFns)

	// No-parameter instructions
	noParamOps := []string{
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
	noParamFns := lo.SliceToMap(
		noParamOps,
		func(op string) (string, opcodeEvalFn) {
			return op, processNoParam
		},
	)
	opcodeEvalFns = lo.Assign(opcodeEvalFns, noParamFns)

	// RET
	opcodeEvalFns["RET"] = processRET

	// MOV
	opcodeEvalFns["MOV"] = processMOV

	// Interrupt Instructions
	opcodeEvalFns["INT"] = processINT

	// Arithmetic Instructions
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

	// Logical Instructions
	opcodeEvalFns["AND"] = processAND
	opcodeEvalFns["OR"] = processOR
	opcodeEvalFns["XOR"] = processXOR
	opcodeEvalFns["NOT"] = processNOT
	opcodeEvalFns["SHR"] = processSHR
	opcodeEvalFns["SHL"] = processSHL
	opcodeEvalFns["SAR"] = processSAR

	// IN
	opcodeEvalFns["IN"] = processIN

	// OUT
	opcodeEvalFns["OUT"] = processOUT

	// CALL
	opcodeEvalFns["CALL"] = processCALL

	// LGDT
	opcodeEvalFns["LGDT"] = processLGDT
}

func popAndPush(env *Pass1) {
	ok, t := env.Ctx.Pop()
	if !ok {
		log.Fatal("error: failed to pop token")
	}
	err := env.Ctx.Push(t)
	if err != nil {
		log.Fatal(failure.Wrap(err))
	}
}

func pop(env *Pass1) *token.ParseToken {
	ok, t := env.Ctx.Pop()
	if !ok {
		log.Fatal("error: failed to pop token")
	}
	return t
}

func push(env *Pass1, t *token.ParseToken) {
	err := env.Ctx.Push(t)
	if err != nil {
		log.Fatal(failure.Wrap(err))
	}
}

func TraverseAST(node ast.Node, env *Pass1) { // Restored ast.Node
	switch n := node.(type) {
	case *ast.Program: // Restored ast.Program
		log.Println("trace: program handler!!!")
		for _, stmt := range n.Statements {
			TraverseAST(stmt, env)
		}

	case *ast.DeclareStmt: // Restored ast.DeclareStmt
		log.Println("trace: declare handler!!!")
		TraverseAST(n.Id, env)
		ok, key := env.Ctx.Pop()
		if !ok {
			log.Fatal("error: EQU failed to pop token key")
		}

		TraverseAST(n.Value, env)
		ok, value := env.Ctx.Pop()
		if !ok {
			log.Fatal("error: EQU failed to pop token value")
		}

		env.EquMap[key.AsString()] = value
		log.Printf("debug: EquMap after DeclareStmt: %+v\n", env.EquMap)

	case *ast.LabelStmt: // Restored ast.LabelStmt
		// ラベルが存在するので、シンボルテーブルのラベルのレコードに現在のLOCを設定
		log.Println("trace: label handler!!!")
		TraverseAST(n.Label, env)
		vLabel := pop(env)
		label := strings.TrimSuffix(vLabel.AsString(), ":")
		env.SymTable[label] = env.LOC

	case *ast.MnemonicStmt: // Restored ast.MnemonicStmt
		log.Println("trace: mnemonic stmt handler!!!")
		TraverseAST(n.Opcode, env)
		vOpcode := pop(env)

		vOperands := make([]*token.ParseToken, 0)
		for _, operand := range n.Operands {
			TraverseAST(operand, env)
			ok, vOperand := env.Ctx.Pop()
			if !ok {
				log.Fatal("error: failed to pop operand")
				return
			}
			vOperands = append(vOperands, vOperand)
		}

		if vOpcode.Data == nil {
			log.Fatal("error: opcode is invalid")
		}

		opcode := vOpcode.AsString()
		evalOpcodeFn := opcodeEvalFns[opcode]
		if evalOpcodeFn == nil {
			log.Fatal("error: not registered opcode process function; ", opcode)
		}

		evalOpcodeFn(env, vOperands)
		log.Printf("debug: [pass1] LOC = %d\n", env.LOC)

	case *ast.OpcodeStmt: // Restored ast.OpcodeStmt
		log.Println("trace: opcode stmt handler!!!")
		TraverseAST(n.Opcode, env)
		vOpcode := pop(env)

		if vOpcode.Data == nil {
			log.Fatal("error: opcode is invalid")
		}

		opcode := vOpcode.AsString()
		evalOpcodeFn := opcodeEvalFns[opcode]
		if evalOpcodeFn == nil {
			log.Fatal("error: not registered opcode process function; ", opcode)
		}

		evalOpcodeFn(env, nil)
		env.Client.Emit(opcode) // opcodeFnの中で実行できないので
		log.Printf("debug: [pass1] LOC = %d\n", env.LOC)

	case *ast.ExportSymStmt: // Restored ast.ExportSymStmt
		log.Println("trace: export sym stmt handler!!!")
		//for _, factor := range n.Factors {
		//	TraverseAST(factor, env)
		//}

	case *ast.ExternSymStmt: // Restored ast.ExternSymStmt
		log.Println("trace: extern sym stmt handler!!!")
		//for _, factor := range n.Factors {
		//	TraverseAST(factor, env)
		//}

	case *ast.ConfigStmt: // Restored ast.ConfigStmt
		log.Println("trace: config stmt handler!!!")
		// 使用するbit_modeは機械語サイズに影響するので読み取って設定する
		TraverseAST(n.Factor, env)
		if n.ConfigType == ast.Bits { // Restored ast.Bits
			ok, token := env.Ctx.Pop()
			if !ok {
				log.Fatal("Failed to pop token")
			}
			bitMode, ok := cpu.NewBitMode(token.ToInt()) // Keep cpu.NewBitMode
			if !ok {
				log.Fatal("error: Failed to parse BITS")
			}
			env.BitMode = bitMode
			// Notify the CodeGenClient (and its context) about the bit mode change
			env.Client.SetBitMode(bitMode)
		}

	case *ast.MemoryAddrExp: // Restored ast.MemoryAddrExp
		log.Println("trace: memory addr exp handler!!!")
		// Recursively traverse left and right sides of memory address expression
		TraverseAST(n.Left, env)
		if n.Right != nil {
			TraverseAST(n.Right, env)
		}

		pop(env) // Pop the result of TraverseAST(n.Right) or TraverseAST(n.Left)
		v := token.NewParseToken(token.TTIdentifier, n)
		push(env, v)

	case *ast.SegmentExp: // Restored ast.SegmentExp
		log.Println("trace: segment exp handler!!!")
		TraverseAST(n.Left, env)
		if n.Right != nil {
			TraverseAST(n.Right, env)
		}
		popAndPush(env)

	case *ast.AddExp: // Restored ast.AddExp
		log.Println("trace: add exp handler!!!")
		TraverseAST(n.HeadExp, env)
		vHead := pop(env)
		vTail := make([]*token.ParseToken, 0)
		ops := make([]string, 0)
		tuples := lo.Zip2(n.Operators, n.TailExps)

		for _, t := range tuples {
			ops = append(ops, t.A)
			TraverseAST(t.B, env)
			vTail = append(vTail, pop(env))
		}

		if len(vTail) == 0 {
			push(env, vHead)
			return
		}

		if vHead.TokenType == token.TTHex &&
			ops[0] == "-" &&
			vTail[0].AsString() == "$" {
			// 0xffff - $ という特殊系
			v := token.NewParseToken(token.TTIdentifier,
				ast.NewImmExp(ast.BaseExp{}, ast.NewIdentFactor(ast.BaseFactor{}, vHead.AsString()+"-$")), // Restored ast types
			)
			push(env, v)
			return
		}

		acc := 0
		if vHead.IsNumber() {
			acc = vHead.ToInt()
		} else {
			push(env, vHead)
		}

		sum := lo.Reduce(lo.Zip2(ops, vTail), func(acc int, t lo.Tuple2[string, *token.ParseToken], _ int) int {
			if t.A == "+" && t.B.IsNumber() {
				return acc + t.B.ToInt()
			} else if t.A == "-" && t.B.IsNumber() {
				return acc - t.B.ToInt()
			}
			return acc
		}, acc)

		v := token.NewParseToken(token.TTNumber,
			ast.NewImmExp(ast.BaseExp{}, ast.NewNumberFactor(ast.BaseFactor{}, sum)), // Restored ast types
		)
		push(env, v)
		return

	case *ast.MultExp: // Restored ast.MultExp
		log.Println("trace: mult exp handler!!!")
		TraverseAST(n.HeadExp, env)
		vHead := pop(env)
		vTail := make([]*token.ParseToken, 0)
		ops := make([]string, 0)
		tuples := lo.Zip2(n.Operators, n.TailExps)

		for _, t := range tuples {
			ops = append(ops, t.A)
			TraverseAST(t.B, env)
			vTail = append(vTail, pop(env))
		}

		if len(vTail) == 0 {
			push(env, vHead)
			return
		}

		base := 1
		if vHead.IsNumber() {
			base = vHead.ToInt()
		} else {
			push(env, vHead)
		}

		sum := lo.Reduce(lo.Zip2(ops, vTail), func(acc int, t lo.Tuple2[string, *token.ParseToken], _ int) int {
			if t.A == "*" && t.B.IsNumber() {
				return acc * t.B.ToInt()
			} else if t.A == "/" && t.B.IsNumber() {
				return acc / t.B.ToInt()
			}
			return acc
		}, base)

		v := token.NewParseToken(token.TTNumber,
			ast.NewImmExp(ast.BaseExp{}, ast.NewNumberFactor(ast.BaseFactor{}, sum)), // Restored ast types
		)
		push(env, v)
		return

	case *ast.ImmExp: // Restored ast.ImmExp
		log.Println("trace: imm exp handler!!!")
		TraverseAST(n.Factor, env)

		if ident, ok := n.Factor.(*ast.IdentFactor); ok { // Restored ast.IdentFactor
			if value, ok := env.EquMap[ident.Value]; ok {
				log.Printf("debug: IdentFactor: %s found in EquMap: %+v\n", ident.Value, value)
				// EQU対応; 置き換え対象の Factor を新しい Factor に変更する
				if immExp, ok := value.Data.(*ast.ImmExp); ok { // Restored ast.ImmExp
					n.Factor = immExp.Factor

					pop(env)
					v := token.NewParseToken(token.TTIdentifier, n)
					push(env, v)
					return
				}
			}
		}
		popAndPush(env)
		return

	case *ast.NumberFactor, // Restored ast types
		*ast.StringFactor,
		*ast.HexFactor,
		*ast.IdentFactor,
		*ast.CharFactor:

		log.Printf("trace: %T factor: %+v\n", n, n)
		var t *token.ParseToken
		switch f := n.(type) {
		case *ast.NumberFactor: // Restored ast types
			t = token.NewParseToken(token.TTNumber, ast.NewImmExp(ast.BaseExp{}, f))
		case *ast.StringFactor:
			t = token.NewParseToken(token.TTIdentifier, ast.NewImmExp(ast.BaseExp{}, f))
		case *ast.HexFactor:
			t = token.NewParseToken(token.TTHex, ast.NewImmExp(ast.BaseExp{}, f))
		case *ast.IdentFactor:
			t = token.NewParseToken(token.TTIdentifier, ast.NewImmExp(ast.BaseExp{}, f))
		case *ast.CharFactor:
			t = token.NewParseToken(token.TTIdentifier, ast.NewImmExp(ast.BaseExp{}, f))
		default:
			return
		}

		err := env.Ctx.Push(t)
		if err != nil {
			log.Fatal(failure.Wrap(err))
		}
	default:
		log.Printf("Unknown AST node: %T\n", node)
		return
	}
}

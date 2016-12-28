package eval

import (
	"fmt"
	"github.com/comail/colog"
	"github.com/hangingman/gosk/ast"
	"github.com/hangingman/gosk/object"
	"github.com/hangingman/gosk/token"
	"log"
	"strconv"
	"strings"
)

type (
	opcodeEvalFn func(stmt *ast.MnemonicStatement) object.Object
)

var (
	// 変数格納
	equMap = make(map[string]token.Token)
	// オペコードごとに評価関数を切り替える
	opcodeEvalFns = make(map[string]opcodeEvalFn)
	// '$' が表す現在のポジション
	dollarPosition = 0
	// 現在までで評価されたバイナリ
	curByteSize = 0
	// ラベルとジャンプ命令管理用オブジェクト
	labelManage = LabelManagement{
		opcode:            map[string][]byte{},
		labelBinaryRefMap: map[string]*object.Binary{},
		labelBytesMap:     map[string]int{},
		labelFromMap:      map[string][]int{},
		genBytesFns:       map[string]func(i int) []byte{},
	}
	globalSymbolList = []string{}
	externSymbolList = []string{}
)

func init() {
	colog.Register()
	colog.SetDefaultLevel(colog.LInfo)
	colog.SetMinLevel(colog.LInfo)
	colog.SetFlags(log.Lshortfile)
	colog.SetFormatter(&colog.StdFormatter{Colors: false})

	opcodeEvalFns["AAA"] = evalSingleByteOpcode("AAA", 0x37)
	opcodeEvalFns["AAS"] = evalSingleByteOpcode("AAS", 0x3f)
	opcodeEvalFns["ADD"] = evalADDStatement
	opcodeEvalFns["ALIGNB"] = evalALIGNBStatement
	opcodeEvalFns["AND"] = evalANDStatement
	opcodeEvalFns["CALL"] = evalCallStatement
	opcodeEvalFns["CBW"] = evalSingleByteOpcode("CBW", 0x98)
	opcodeEvalFns["CDQ"] = evalSingleByteOpcode("CDQ", 0x99)
	opcodeEvalFns["CLC"] = evalSingleByteOpcode("CLC", 0xf8)
	opcodeEvalFns["CLD"] = evalSingleByteOpcode("CLD", 0xfc)
	opcodeEvalFns["CLI"] = evalSingleByteOpcode("CLI", 0xfa)
	opcodeEvalFns["CLTS"] = evalSingleWordOpcode("CLTS", []byte{0x0f, 0x06})
	opcodeEvalFns["CMC"] = evalSingleByteOpcode("CMC", 0xf5)
	opcodeEvalFns["CMP"] = evalCMPStatement
	opcodeEvalFns["CPUID"] = evalSingleByteOpcode("CPUID", 0xf8)
	opcodeEvalFns["CWD"] = evalSingleByteOpcode("CWD", 0x99)
	opcodeEvalFns["CWDE"] = evalSingleByteOpcode("CWDE", 0x98)
	opcodeEvalFns["DAA"] = evalSingleByteOpcode("DAA", 0x27)
	opcodeEvalFns["DAS"] = evalSingleByteOpcode("DAS", 0x2f)
	opcodeEvalFns["DB"] = evalDBStatement
	opcodeEvalFns["DD"] = evalDDStatement
	opcodeEvalFns["DW"] = evalDWStatement
	opcodeEvalFns["GLOBAL"] = evalGLOBALStatement
	opcodeEvalFns["JAE"] = evalJumpStatement([]byte{0x73})
	opcodeEvalFns["JB"] = evalJumpStatement([]byte{0x72})
	opcodeEvalFns["JBE"] = evalJumpStatement([]byte{0x76})
	opcodeEvalFns["JC"] = evalJumpStatement([]byte{0x72})
	opcodeEvalFns["JE"] = evalJumpStatement([]byte{0x74})
	opcodeEvalFns["JMP"] = evalJumpStatement([]byte{0xeb})
	opcodeEvalFns["JNC"] = evalJumpStatement([]byte{0x73})
	opcodeEvalFns["JZ"] = evalJumpStatement([]byte{0x74})
	opcodeEvalFns["JNZ"] = evalJumpStatement([]byte{0x75})
	opcodeEvalFns["FWAIT"] = evalSingleByteOpcode("WAIT", 0x9b)
	opcodeEvalFns["HLT"] = evalSingleByteOpcode("HLT", 0xf4)
	opcodeEvalFns["IMUL"] = evalIMULStatement
	opcodeEvalFns["IN"] = evalINStatement
	opcodeEvalFns["INCO"] = evalSingleByteOpcode("INCO", 0xce)
	opcodeEvalFns["INSB"] = evalSingleByteOpcode("INSB", 0x6c)
	opcodeEvalFns["INSD"] = evalSingleByteOpcode("INSD", 0x6d)
	opcodeEvalFns["INSW"] = evalSingleByteOpcode("INSW", 0x6d)
	opcodeEvalFns["INT"] = evalINTStatement
	opcodeEvalFns["INVD"] = evalSingleWordOpcode("INVD", []byte{0x0f, 0x08})
	opcodeEvalFns["IRET"] = evalSingleByteOpcode("IRET", 0xcf)
	opcodeEvalFns["IRETD"] = evalSingleByteOpcode("IRETD", 0xcf)
	opcodeEvalFns["LAHF"] = evalSingleByteOpcode("LAHF", 0x9f)
	opcodeEvalFns["LGDT"] = evalLGDTStatement
	opcodeEvalFns["LEAVE"] = evalSingleByteOpcode("LEAVE", 0xc9)
	opcodeEvalFns["LOCK"] = evalSingleByteOpcode("LOCK", 0xf0)
	opcodeEvalFns["MOV"] = evalMOVStatement
	opcodeEvalFns["NOP"] = evalSingleByteOpcode("NOP", 0x90)
	opcodeEvalFns["OR"] = evalORStatement
	opcodeEvalFns["ORG"] = evalORGStatement
	opcodeEvalFns["OUT"] = evalOUTStatement
	opcodeEvalFns["OUTSB"] = evalSingleByteOpcode("OUTSB", 0x6e)
	opcodeEvalFns["OUTSD"] = evalSingleByteOpcode("OUTSD", 0x6f)
	opcodeEvalFns["OUTSW"] = evalSingleByteOpcode("OUTSW", 0x6f)
	opcodeEvalFns["POPA"] = evalSingleByteOpcode("POPA", 0x61)
	opcodeEvalFns["POPAD"] = evalSingleByteOpcode("POPAD", 0x61)
	opcodeEvalFns["POPF"] = evalSingleByteOpcode("POPF", 0x9d)
	opcodeEvalFns["POPFD"] = evalSingleByteOpcode("POPFD", 0x9d)
	opcodeEvalFns["PUSHA"] = evalSingleByteOpcode("PUSHA", 0x60)
	opcodeEvalFns["PUSHD"] = evalSingleByteOpcode("PUSHD", 0x60)
	opcodeEvalFns["PUSHF"] = evalSingleByteOpcode("PUSHF", 0x9c)
	opcodeEvalFns["RESB"] = evalRESBStatement
	opcodeEvalFns["RET"] = evalSingleByteOpcode("RET", 0xc3)
	opcodeEvalFns["RETF"] = evalSingleByteOpcode("RETF", 0xcb)
	opcodeEvalFns["RSM"] = evalSingleWordOpcode("RSM", []byte{0x0f, 0xaa})
	opcodeEvalFns["SAHF"] = evalSingleByteOpcode("SAHF", 0x9e)
	opcodeEvalFns["SHR"] = evalSHRStatement
	opcodeEvalFns["STC"] = evalSingleByteOpcode("STC", 0xf9)
	opcodeEvalFns["STD"] = evalSingleByteOpcode("STD", 0xfd)
	opcodeEvalFns["STI"] = evalSingleByteOpcode("STI", 0xfb)
	opcodeEvalFns["SUB"] = evalSUBStatement
	opcodeEvalFns["UD2"] = evalSingleWordOpcode("UD2", []byte{0x0f, 0x0b})
	opcodeEvalFns["WAIT"] = evalSingleByteOpcode("WAIT", 0x9b)
	opcodeEvalFns["RDMSR"] = evalSingleWordOpcode("RDMSR", []byte{0x0f, 0x32})
	opcodeEvalFns["RDPMC"] = evalSingleWordOpcode("RDPMC", []byte{0x0f, 0x33})
	opcodeEvalFns["RDTSC"] = evalSingleWordOpcode("RDTSC", []byte{0x0f, 0x31})
	opcodeEvalFns["WBINVD"] = evalSingleWordOpcode("WBINVD", []byte{0x0f, 0x09})
	opcodeEvalFns["WRMSR"] = evalSingleWordOpcode("WRMSR", []byte{0x0f, 0x30})
}

func evalSingleByteOpcode(opcode string, b byte) func(stmt *ast.MnemonicStatement) object.Object {
	return func(stmt *ast.MnemonicStatement) object.Object {
		log.Println(fmt.Sprintf("info: [%s, %x]", opcode, b))
		stmt.Bin = &object.Binary{Value: []byte{b}}
		return stmt.Bin
	}
}

func evalSingleWordOpcode(opcode string, w []byte) func(stmt *ast.MnemonicStatement) object.Object {
	return func(stmt *ast.MnemonicStatement) object.Object {
		log.Println(fmt.Sprintf("info: [%s, %x, %x]", opcode, w[0], w[1]))
		stmt.Bin = &object.Binary{Value: w}
		return stmt.Bin
	}
}

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		dollarPosition = 0
		curByteSize = 0
		labelManage.opcode = map[string][]byte{}
		labelManage.labelBinaryRefMap = map[string]*object.Binary{}
		labelManage.labelBytesMap = map[string]int{}
		labelManage.genBytesFns = map[string]func(i int) []byte{}
		return evalStatements(node.Statements)
	case *ast.MnemonicStatement:
		return evalMnemonicStatement(node)
	case *ast.SettingStatement:
		return evalSettingStatement(node)
	case *ast.LabelStatement:
		return evalLabelStatement(node)
	case *ast.EquStatement:
		return evalEquStatement(node)
	case *ast.BinaryLiteral:
		return &object.Binary{Value: node.Value}
	}
	return nil
}

// evalStatements は文を評価する
func evalStatements(stmts []ast.Statement) object.Object {
	results := object.ObjectArray{}

	// 文を評価して、結果としてobject.ObjectArrayを返す
	for _, stmt := range stmts {
		if IsNotNil(stmt) {
			result := Eval(stmt)
			bin, ok := result.(*object.Binary)
			if ok {
				evalByteSize := len(bin.Value)
				log.Println(fmt.Sprintf("info: evaled byte size: %d", evalByteSize))
				curByteSize += len(bin.Value)
				log.Println(fmt.Sprintf("info: current byte size: %d", curByteSize))
			}
			results = append(results, result)
		}
	}

	// セクションテーブルを最後に加える
	//results = append(results, evalSectionTable())

	return &results
}

func evalMnemonicStatement(stmt *ast.MnemonicStatement) object.Object {
	opcode := stmt.Name.Tokens[0].Literal

	if stmt.HasOperator() {
		stmt.PreEval()
	}
	evalOpcodeFn := opcodeEvalFns[opcode]

	if evalOpcodeFn == nil {
		return nil
	}

	return evalOpcodeFn(stmt)
}

func evalLabelStatement(stmt *ast.LabelStatement) object.Object {
	label := strings.TrimSuffix(stmt.Name, ":")
	// ラベルが見つかったのでコールバックを起動して処理する
	labelManage.Emit(label, curByteSize)
	// 先にラベルが見つかった場合、バイト数を記録しておく
	labelManage.labelBytesMap[label] = curByteSize
	return nil
}

func evalEquStatement(stmt *ast.EquStatement) object.Object {
	// EQUで指定された文字列を置き換える
	equKey := stmt.Name.Token.Literal
	equTok := stmt.Value
	log.Println(fmt.Sprintf("info: %s = %s", equKey, equTok))
	equMap[equKey] = equTok

	nextStmt := stmt.GetNextNode()
	for {
		switch nextStmt.(type) {
		case *ast.MnemonicStatement:
			m := nextStmt.(*ast.MnemonicStatement)
			for idx, tok := range m.Name.Tokens {
				if tok.Type == token.IDENT && tok.Literal == equKey {
					log.Println("info: replace token by EQU specified")
					m.Name.Tokens[idx] = equTok
					m.Name.Values[idx] = equTok.Literal
				}
			}
		default:
			// do nothing
		}
		nextStmt = nextStmt.GetNextNode()
		if nextStmt == nil {
			break
		}
	}

	return nil
}

func makeZeroFill(bs []byte) []byte {
	for i := range bs {
		bs[i] = 0x00
	}
	return bs
}

func makeZeroFilledBytesU64(byteSize uint64) []byte {
	bs := make([]byte, byteSize)
	return makeZeroFill(bs)
}

func makeZeroFilledBytes(byteSize int) []byte {
	bs := make([]byte, byteSize)
	return makeZeroFill(bs)
}

func evalRESBStatement(stmt *ast.MnemonicStatement) object.Object {
	toks := []string{}
	bytes := []byte{}

	for i, tok := range stmt.Name.Tokens {
		if tok.Type == token.INT {
			v, _ := strconv.Atoi(tok.Literal)
			bs := makeZeroFilledBytes(v)
			bytes = append(bytes, bs...)
		} else if tok.Type == token.HEX_LIT {
			// RESB	0x1fe-$ のように hexリテラル値の後に
			// ハイフンとダラーがあることを期待する
			if stmt.Name.Tokens[i+1].Type == token.MINUS &&
				stmt.Name.Tokens[i+2].Type == token.DOLLAR {
				u64v, _ := strconv.ParseUint(tok.Literal[2:], 16, 64)

				log.Println(fmt.Sprintf("info: RESB will fill by zero, upto %d", u64v))
				required := int(u64v) - dollarPosition - curByteSize
				log.Println(fmt.Sprintf("info: RESB required %d zero filled binary", required))
				for i := 0; i < required; i++ {
					bytes = append(bytes, 0x00)
				}
				break
			}
		}
		toks = append(toks, fmt.Sprintf("%s: %s", tok.Type, tok.Literal))
	}

	log.Println(fmt.Sprintf("info: [%s]", strings.Join(toks, ", ")))
	stmt.Bin = &object.Binary{Value: bytes}
	return stmt.Bin
}

func evalALIGNBStatement(stmt *ast.MnemonicStatement) object.Object {
	toks := []string{}
	bytes := []byte{}

	for _, tok := range stmt.Name.Tokens {
		if tok.Type == token.INT {
			unit, _ := strconv.Atoi(tok.Literal)
			nearestSize := curByteSize/unit + 1
			times := nearestSize*unit - curByteSize
			bs := makeZeroFilledBytes(times)
			bytes = append(bytes, bs...)
			log.Println(fmt.Sprintf("info: ALIGNB stores 0x00 %d times", times))
		}
		toks = append(toks, fmt.Sprintf("%s: %s", tok.Type, tok.Literal))
	}

	log.Println(fmt.Sprintf("info: [%s]", strings.Join(toks, ", ")))
	stmt.Bin = &object.Binary{Value: bytes}
	return stmt.Bin
}

func evalORGStatement(stmt *ast.MnemonicStatement) object.Object {
	toks := []string{}

	for _, tok := range stmt.Name.Tokens {
		if tok.Type == token.INT {
			// Go言語のintは常にint64
			v, _ := strconv.Atoi(tok.Literal)
			dollarPosition = v
		} else if tok.Type == token.HEX_LIT {
			u64v, _ := strconv.ParseUint(tok.Literal[2:], 16, 64)
			dollarPosition = int(u64v)
		}
		toks = append(toks, fmt.Sprintf("%s: %s", tok.Type, tok.Literal))
	}
	log.Println(fmt.Sprintf("info: [%s]", strings.Join(toks, ", ")))
	log.Println(fmt.Sprintf("info: ORG = %d", dollarPosition))
	return nil
}

func evalLGDTStatement(stmt *ast.MnemonicStatement) object.Object {
	stmt.Bin = &object.Binary{Value: []byte{}}
	toks := []string{}

	for _, tok := range stmt.Name.Tokens {
		if tok.Type == token.IDENT {
			stmt.Bin.Value = append(stmt.Bin.Value, 0x0f)
			stmt.Bin.Value = append(stmt.Bin.Value, 0x01)
			modrm := generateModRMSlashN(0x0f, RegReg, "["+tok.Literal+"]", "/2")
			stmt.Bin.Value = append(stmt.Bin.Value, modrm)

			if _, ok := labelManage.labelBytesMap[tok.Literal]; ok {
				stmt.Bin.Value = append(stmt.Bin.Value, int2Word(dollarPosition)...)
			} else {
				stmt.Bin.Value = append(stmt.Bin.Value, 0x00)
				stmt.Bin.Value = append(stmt.Bin.Value, 0x00)
				labelManage.AddLabelCallback(
					// CALL自体のバイト数を含まないので +2 しておく
					[]byte{0x0f, 0x01, modrm}, tok.Literal, stmt.Bin, -dollarPosition, int2Word,
				)
			}
		}
		toks = append(toks, fmt.Sprintf("%s: %s", tok.Type, tok.Literal))
	}
	log.Println(fmt.Sprintf("info: [%s]", strings.Join(toks, ", ")))

	return stmt.Bin
}

func evalCallStatement(stmt *ast.MnemonicStatement) object.Object {
	stmt.Bin = &object.Binary{Value: []byte{}}

	for _, tok := range stmt.Name.Tokens {
		if tok.Type == token.IDENT {
			if from, ok := labelManage.labelBytesMap[tok.Literal]; ok {
				// ラベルが見つかっていればバイト数を計算して設定する
				log.Println(fmt.Sprintf("info: already has label %s", tok.Literal))
				log.Println(fmt.Sprintf("info: %d - %d - 3 = %d", from, curByteSize, from-curByteSize-3))
				stmt.Bin.Value = append(stmt.Bin.Value, 0xe8)
				stmt.Bin.Value = append(stmt.Bin.Value, int2Word(from-curByteSize-3)...)
			} else {
				// ラベルが見つかっていないならば
				// callbackを配置し今のバイト数を設定する
				log.Println(fmt.Sprintf("info: no label %s", tok.Literal))
				stmt.Bin.Value = append(stmt.Bin.Value, 0xe8)
				stmt.Bin.Value = append(stmt.Bin.Value, 0x00)
				stmt.Bin.Value = append(stmt.Bin.Value, 0x00)

				labelManage.AddLabelCallback(
					// CALL自体のバイト数を含まないので +3 しておく
					[]byte{0xe8}, tok.Literal, stmt.Bin, curByteSize+3, int2Word,
				)
			}
		}
		log.Println(fmt.Sprintf("info: %s", tok))
	}

	return stmt.Bin
}

func evalJumpStatement(b []byte) func(stmt *ast.MnemonicStatement) object.Object {
	return func(stmt *ast.MnemonicStatement) object.Object {
		stmt.Bin = &object.Binary{Value: []byte{}}

		for idx, tok := range stmt.Name.Tokens {
			if tok.Type == token.IDENT {
				if from, ok := labelManage.labelBytesMap[tok.Literal]; ok {
					// ラベルが見つかっていればバイト数を計算して設定する
					log.Println(fmt.Sprintf("info: already has label %s", tok.Literal))
					log.Println(fmt.Sprintf("info: %d - %d - 2 = %d", from, curByteSize, from-curByteSize-2))
					stmt.Bin.Value = append(stmt.Bin.Value, b...)
					stmt.Bin.Value = append(stmt.Bin.Value, int2Byte(from-curByteSize-2)...)
				} else {
					// ラベルが見つかっていないならば
					// callbackを配置し今のバイト数を設定する
					log.Println(fmt.Sprintf("info: no label %s", tok.Literal))
					stmt.Bin.Value = append(stmt.Bin.Value, b...)
					stmt.Bin.Value = append(stmt.Bin.Value, 0x00)

					labelManage.AddLabelCallback(
						// JMP自体のバイト数を含まないので +2 しておく
						b, tok.Literal, stmt.Bin, curByteSize+2, int2Byte,
					)
				}
			} else if tok.Type == token.HEX_LIT {
				// JMP 0xc200
				// のようにジャンプさせたい時用
				u64v, _ := strconv.ParseUint(tok.Literal[2:], 16, 64)
				stmt.Bin.Value = append(stmt.Bin.Value, 0xe9)
				stmt.Bin.Value = append(stmt.Bin.Value, int2Word(int(u64v)-dollarPosition-curByteSize-3)...)
			} else if tok.Type == token.DATA_TYPE {
				// JMP DWORD 2*8:0x0000001b
				// のようにジャンプさせたい時用
				// 0xEA cd => JMP ptr16:16
				// 0xEA cp => JMP ptr16:32
				// 0xFF /5 => JMP m16:16
				// 0xFF /5 => JMP m16:32
				stmt.Bin.Value = append(stmt.Bin.Value, 0x66)
				stmt.Bin.Value = append(stmt.Bin.Value, 0xea)
				addr := stmt.Name.Tokens[idx+3]
				m16 := stmt.Name.Tokens[idx+1]
				stmt.Bin.Value = append(stmt.Bin.Value, imm32ToDword(addr)...)
				stmt.Bin.Value = append(stmt.Bin.Value, imm16ToWord(m16)...)
				break
			}
			log.Println(fmt.Sprintf("info: %s", tok))
		}

		return stmt.Bin
	}
}

func evalINTStatement(stmt *ast.MnemonicStatement) object.Object {
	bin := &object.Binary{Value: []byte{}}
	toks := []string{}

	for _, tok := range stmt.Name.Tokens {
		if tok.Type == token.INT {
			// Go言語のintは常にint64
			v, _ := strconv.Atoi(tok.Literal)
			bin.Value = append(bin.Value, 0xcd)
			bin.Value = append(bin.Value, int2Byte(v)...)
		} else if tok.Type == token.HEX_LIT {
			u64v, _ := strconv.ParseUint(tok.Literal[2:], 16, 64)
			bin.Value = append(bin.Value, 0xcd)
			bin.Value = append(bin.Value, int2Byte(int(u64v))...)
		}
		toks = append(toks, fmt.Sprintf("%s: %s", tok.Type, tok.Literal))
	}
	log.Println(fmt.Sprintf("info: [%s]", strings.Join(toks, ", ")))

	stmt.Bin = bin
	return stmt.Bin
}

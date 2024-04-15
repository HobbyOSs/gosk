package pass2

import (
	"log"

	"github.com/HobbyOSs/gosk/ast"
	"github.com/HobbyOSs/gosk/token"
	"github.com/morikuni/failure"
	"github.com/samber/lo"
)

//go:generate newc
type ProgramHandlerImpl struct {
	Visitor *ast.Visitor
	Env     *Pass2
}

//go:generate newc
type DeclareStmtHandlerImpl struct {
	Visitor *ast.Visitor
	Env     *Pass2
}

//go:generate newc
type LabelStmtHandlerImpl struct {
	Visitor *ast.Visitor
	Env     *Pass2
}

//go:generate newc
type ExportSymStmtHandlerImpl struct {
	Visitor *ast.Visitor
	Env     *Pass2
}

//go:generate newc
type ExternSymStmtHandlerImpl struct {
	Visitor *ast.Visitor
	Env     *Pass2
}

//go:generate newc
type ConfigStmtHandlerImpl struct {
	Visitor *ast.Visitor
	Env     *Pass2
}

//go:generate newc
type MnemonicStmtHandlerImpl struct {
	Visitor *ast.Visitor
	Env     *Pass2
}

//go:generate newc
type MemoryAddrExpHandlerImpl struct {
	Visitor *ast.Visitor
	Env     *Pass2
}

//go:generate newc
type SegmentExpHandlerImpl struct {
	Visitor *ast.Visitor
	Env     *Pass2
}

//go:generate newc
type AddExpHandlerImpl struct {
	Visitor *ast.Visitor
	Env     *Pass2
}

//go:generate newc
type MultExpHandlerImpl struct {
	Visitor *ast.Visitor
	Env     *Pass2
}

//go:generate newc
type ImmExpHandlerImpl struct {
	Visitor *ast.Visitor
	Env     *Pass2
}

//go:generate newc
type NumberFactorHandlerImpl struct {
	Visitor *ast.Visitor
	Env     *Pass2
}

//go:generate newc
type StringFactorHandlerImpl struct {
	Visitor *ast.Visitor
	Env     *Pass2
}

//go:generate newc
type HexFactorHandlerImpl struct {
	Visitor *ast.Visitor
	Env     *Pass2
}

//go:generate newc
type IdentFactorHandlerImpl struct {
	Visitor *ast.Visitor
	Env     *Pass2
}

//go:generate newc
type CharFactorHandlerImpl struct {
	Visitor *ast.Visitor
	Env     *Pass2
}

type opcodeEvalFn func(*Pass2, []*token.ParseToken)

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
}

func popAndPush(env *Pass2) {
	ok, t := env.Ctx.Pop()
	if !ok {
		log.Fatal("error: failed to pop token")
	}
	err := env.Ctx.Push(t)
	if err != nil {
		log.Fatal(failure.Wrap(err))
	}
}

func pop(env *Pass2) *token.ParseToken {
	ok, t := env.Ctx.Pop()
	if !ok {
		log.Fatal("error: failed to pop token")
	}
	return t
}

func push(env *Pass2, t *token.ParseToken) {
	err := env.Ctx.Push(t)
	if err != nil {
		log.Fatal(failure.Wrap(err))
	}
}

func (p *ProgramHandlerImpl) Program(node *ast.Program) bool {
	log.Println("trace: program handler!!!")
	for _, stmt := range node.Statements {
		p.Visitor.Visit(stmt)
	}
	return true
}

func (p *DeclareStmtHandlerImpl) DeclareStmt(node *ast.DeclareStmt) bool {
	log.Println("trace: declare handler!!!")

	p.Visitor.Visit(node.Id)
	ok, key := p.Env.Ctx.Pop()
	if !ok {
		log.Fatal("error: EQU failed to pop token key")
	}

	p.Visitor.Visit(node.Value)
	ok, value := p.Env.Ctx.Pop()
	if !ok {
		log.Fatal("error: EQU failed to pop token value")
	}
	p.Env.EquMap[key.AsString()] = value

	return true
}

func (p *LabelStmtHandlerImpl) LabelStmt(node *ast.LabelStmt) bool {
	log.Println("trace: label handler!!!")
	return true
}

func (p *ExportSymStmtHandlerImpl) ExportSymStmt(node *ast.ExportSymStmt) bool {
	log.Println("trace: export sym handler!!!")
	return true
}

func (p *ExternSymStmtHandlerImpl) ExternSymStmt(node *ast.ExternSymStmt) bool {
	log.Println("trace: extern sym handler!!!")
	return true
}

func (p *ConfigStmtHandlerImpl) ConfigStmt(node *ast.ConfigStmt) bool {
	// 使用するbit_modeは機械語サイズに影響するので読み取って設定する
	p.Visitor.Visit(node.Factor)

	if node.ConfigType == ast.Bits {
		ok, token := p.Env.Ctx.Pop()
		if !ok {
			log.Fatal("Failed to pop token")
		}
		bitMode, ok := ast.NewBitMode(token.ToInt())
		if !ok {
			log.Fatal("error: Failed to parse BITS")
		}
		p.Env.BitMode = bitMode
	}

	return true
}

func (p *MnemonicStmtHandlerImpl) MnemonicStmt(node *ast.MnemonicStmt) bool {
	log.Println("trace: mnemonic stmt handler!!!")

	// オペコードに応じて機械語を出力する
	p.Visitor.Visit(node.Opcode)
	vOpcode := pop(p.Env)

	vOperands := make([]*token.ParseToken, 0)
	for _, operand := range node.Operands {
		p.Visitor.Visit(operand)
		vOperands = append(vOperands, pop(p.Env))
	}

	if vOpcode.Data == nil {
		log.Fatal("error: opcode is invalid")
	}

	opcode := vOpcode.AsString()
	evalOpcodeFn := opcodeEvalFns[opcode]
	if evalOpcodeFn == nil {
		log.Fatal("error: not registered opcode process function; ", opcode)
	}

	evalOpcodeFn(p.Env, vOperands) // 評価

	return true
}

/**
 * Handling Exp elements
 */
func (p *MemoryAddrExpHandlerImpl) MemoryAddrExp(node *ast.MemoryAddrExp) bool {
	log.Println("trace: memory_addr exp handler!!!")
	log.Printf("trace: %+v", node)
	return true
}

func (p *SegmentExpHandlerImpl) SegmentExp(node *ast.SegmentExp) bool {
	log.Println("trace: segment exp handler!!!")
	p.Visitor.Visit(node.Left)
	if node.Right != nil {
		p.Visitor.Visit(node.Right)
	}
	popAndPush(p.Env)
	return true
}

func (p *AddExpHandlerImpl) AddExp(node *ast.AddExp) bool {
	log.Println("trace: add exp handler!!!")
	p.Visitor.Visit(node.HeadExp)
	vHead := pop(p.Env)

	vTail := make([]*token.ParseToken, 0)
	ops := make([]string, 0)
	tuples := lo.Zip2(node.Operators, node.TailExps)

	for _, t := range tuples {
		ops = append(ops, t.A)
		p.Visitor.Visit(t.B)
		vTail = append(vTail, pop(p.Env))
	}

	if len(vTail) == 0 {
		push(p.Env, vHead)
		return true
	}
	if vHead.TokenType == token.TTHex &&
		ops[0] == "-" &&
		vTail[0].AsString() == "$" {
		// 0xffff - $ という特殊系
		v := token.NewParseToken(token.TTIdentifier,
			ast.NewImmExp(ast.BaseExp{}, ast.NewIdentFactor(ast.BaseFactor{}, vHead.AsString()+"-$")),
		)
		push(p.Env, v)
		return true
	}

	acc := 0
	if vHead.IsNumber() {
		acc = vHead.ToInt()
	} else {
		push(p.Env, vHead)
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
		ast.NewImmExp(ast.BaseExp{}, ast.NewNumberFactor(ast.BaseFactor{}, sum)),
	)
	push(p.Env, v)

	return true
}

func (p *MultExpHandlerImpl) MultExp(node *ast.MultExp) bool {
	log.Println("trace: mult exp handler!!!")
	p.Visitor.Visit(node.HeadExp)
	vHead := pop(p.Env)
	vTail := make([]*token.ParseToken, 0)
	ops := make([]string, 0)
	tuples := lo.Zip2(node.Operators, node.TailExps)

	for _, t := range tuples {
		ops = append(ops, t.A)
		p.Visitor.Visit(t.B)
		vTail = append(vTail, pop(p.Env))
	}

	if len(vTail) == 0 {
		push(p.Env, vHead)
		return true
	}

	base := 1
	if vHead.IsNumber() {
		base = vHead.ToInt()
	} else {
		push(p.Env, vHead)
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
		ast.NewImmExp(ast.BaseExp{}, ast.NewNumberFactor(ast.BaseFactor{}, sum)),
	)
	push(p.Env, v)

	return true
}

func (p *ImmExpHandlerImpl) ImmExp(node *ast.ImmExp) bool {
	log.Println("trace: imm exp handler!!!")
	p.Visitor.Visit(node.Factor)
	popAndPush(p.Env)
	return true
}

/**
 * Handling Factor elements
 */
func (p *NumberFactorHandlerImpl) NumberFactor(node *ast.NumberFactor) bool {
	log.Printf("trace: number factor: %+v\n", node)

	t := token.NewParseToken(token.TTNumber, ast.NewImmExp(ast.BaseExp{}, node))
	err := p.Env.Ctx.Push(t)
	if err != nil {
		log.Fatal(failure.Wrap(err))
	}
	return true
}

func (p *StringFactorHandlerImpl) StringFactor(node *ast.StringFactor) bool {
	log.Printf("trace: string factor: %+v\n", node)

	t := token.NewParseToken(token.TTIdentifier, ast.NewImmExp(ast.BaseExp{}, node))
	err := p.Env.Ctx.Push(t)
	if err != nil {
		log.Fatal(failure.Wrap(err))
	}
	return true

}

func (p *HexFactorHandlerImpl) HexFactor(node *ast.HexFactor) bool {
	log.Printf("trace: hex factor: %+v\n", node)

	t := token.NewParseToken(token.TTHex, ast.NewImmExp(ast.BaseExp{}, node))
	err := p.Env.Ctx.Push(t)
	if err != nil {
		log.Fatal("error: Failed to push token; ", err)
	}
	return true

}

func (p *IdentFactorHandlerImpl) IdentFactor(node *ast.IdentFactor) bool {
	log.Printf("trace: ident factor: %+v\n", node)

	t := token.NewParseToken(token.TTIdentifier, ast.NewImmExp(ast.BaseExp{}, node))
	err := p.Env.Ctx.Push(t)
	if err != nil {
		log.Fatal(failure.Wrap(err))
	}
	return true

}

func (p *CharFactorHandlerImpl) CharFactor(node *ast.CharFactor) bool {
	log.Printf("trace: char factor: %+v\n", node)

	t := token.NewParseToken(token.TTIdentifier, ast.NewImmExp(ast.BaseExp{}, node))
	err := p.Env.Ctx.Push(t)
	if err != nil {
		log.Fatal(failure.Wrap(err))
	}
	return true
}

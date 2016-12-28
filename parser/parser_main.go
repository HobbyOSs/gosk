package parser

import (
	"fmt"
	"github.com/comail/colog"
	"github.com/hangingman/gosk/ast"
	"github.com/hangingman/gosk/lexer"
	"github.com/hangingman/gosk/token"
	"log"
	"reflect"
)

type (
	opcodeParseFn func() *ast.MnemonicStatement
)

type Parser struct {
	l              *lexer.Lexer
	curIndex       int
	lexedTokens    []token.Token
	errors         []string
	opcodeParseFns map[string]opcodeParseFn // オペコードごとに構文解析を切り替える
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:           l,
		curIndex:    0,
		lexedTokens: make([]token.Token, 50),
		errors:      []string{},
	}
	colog.Register()
	colog.SetDefaultLevel(colog.LInfo)
	colog.SetMinLevel(colog.LInfo)
	colog.SetFlags(log.Lshortfile)
	colog.SetFormatter(&colog.StdFormatter{Colors: false})

	// オペコードの構文解析方式を格納するmap
	p.opcodeParseFns = make(map[string]opcodeParseFn)
	p.registerOpecode("AAA", p.parseOnlyOpcodeStatement)
	p.registerOpecode("AAS", p.parseOnlyOpcodeStatement)
	p.registerOpecode("ADD", p.parseGeneralOpcodeStatement)
	p.registerOpecode("ALIGNB", p.parseGeneralOpcodeStatement)
	p.registerOpecode("AND", p.parseGeneralOpcodeStatement)
	p.registerOpecode("CALL", p.parseJMPStatement)
	p.registerOpecode("CBW", p.parseOnlyOpcodeStatement)
	p.registerOpecode("CDQ", p.parseOnlyOpcodeStatement)
	p.registerOpecode("CLC", p.parseOnlyOpcodeStatement)
	p.registerOpecode("CLD", p.parseOnlyOpcodeStatement)
	p.registerOpecode("CLI", p.parseOnlyOpcodeStatement)
	p.registerOpecode("CLTS", p.parseOnlyOpcodeStatement)
	p.registerOpecode("CMC", p.parseOnlyOpcodeStatement)
	p.registerOpecode("CMP", p.parseGeneralOpcodeStatement)
	p.registerOpecode("CPUID", p.parseOnlyOpcodeStatement)
	p.registerOpecode("CWD", p.parseOnlyOpcodeStatement)
	p.registerOpecode("CWDE", p.parseOnlyOpcodeStatement)
	p.registerOpecode("DAA", p.parseOnlyOpcodeStatement)
	p.registerOpecode("DAS", p.parseOnlyOpcodeStatement)
	p.registerOpecode("DB", p.parseGeneralOpcodeStatement)
	p.registerOpecode("DD", p.parseGeneralOpcodeStatement)
	p.registerOpecode("DW", p.parseGeneralOpcodeStatement)
	p.registerOpecode("GLOBAL", p.parseDBLikeStatement)
	p.registerOpecode("JAE", p.parseJMPStatement)
	p.registerOpecode("JB", p.parseJMPStatement)
	p.registerOpecode("JBE", p.parseJMPStatement)
	p.registerOpecode("JC", p.parseJMPStatement)
	p.registerOpecode("JE", p.parseJMPStatement)
	p.registerOpecode("JMP", p.parseGeneralOpcodeStatement)
	p.registerOpecode("JNC", p.parseJMPStatement)
	p.registerOpecode("JNZ", p.parseJMPStatement)
	p.registerOpecode("JZ", p.parseJMPStatement)
	p.registerOpecode("FWAIT", p.parseOnlyOpcodeStatement)
	p.registerOpecode("HLT", p.parseOnlyOpcodeStatement)
	p.registerOpecode("IMUL", p.parseGeneralOpcodeStatement)
	p.registerOpecode("IN", p.parseGeneralOpcodeStatement)
	p.registerOpecode("INCO", p.parseOnlyOpcodeStatement)
	p.registerOpecode("INSB", p.parseOnlyOpcodeStatement)
	p.registerOpecode("INSD", p.parseOnlyOpcodeStatement)
	p.registerOpecode("INSW", p.parseOnlyOpcodeStatement)
	p.registerOpecode("INT", p.parseDBLikeStatement)
	p.registerOpecode("INVD", p.parseOnlyOpcodeStatement)
	p.registerOpecode("IRET", p.parseOnlyOpcodeStatement)
	p.registerOpecode("IRETD", p.parseOnlyOpcodeStatement)
	p.registerOpecode("LAHF", p.parseOnlyOpcodeStatement)
	p.registerOpecode("LEAVE", p.parseOnlyOpcodeStatement)
	p.registerOpecode("LGDT", p.parseLGDTStatement)
	p.registerOpecode("LOCK", p.parseOnlyOpcodeStatement)
	p.registerOpecode("MOV", p.parseMOVStatement)
	p.registerOpecode("NOP", p.parseOnlyOpcodeStatement)
	p.registerOpecode("OR", p.parseGeneralOpcodeStatement)
	p.registerOpecode("ORG", p.parseORGStatement)
	p.registerOpecode("OUT", p.parseGeneralOpcodeStatement)
	p.registerOpecode("OUTSB", p.parseOnlyOpcodeStatement)
	p.registerOpecode("OUTSD", p.parseOnlyOpcodeStatement)
	p.registerOpecode("OUTSW", p.parseOnlyOpcodeStatement)
	p.registerOpecode("POPA", p.parseOnlyOpcodeStatement)
	p.registerOpecode("POPAD", p.parseOnlyOpcodeStatement)
	p.registerOpecode("POPF", p.parseOnlyOpcodeStatement)
	p.registerOpecode("POPFD", p.parseOnlyOpcodeStatement)
	p.registerOpecode("PUSHA", p.parseOnlyOpcodeStatement)
	p.registerOpecode("PUSHD", p.parseOnlyOpcodeStatement)
	p.registerOpecode("PUSHF", p.parseOnlyOpcodeStatement)
	p.registerOpecode("RDMSR", p.parseRESBStatement)
	p.registerOpecode("RESB", p.parseRESBStatement)
	p.registerOpecode("RET", p.parseOnlyOpcodeStatement)
	p.registerOpecode("RETF", p.parseOnlyOpcodeStatement)
	p.registerOpecode("RSM", p.parseOnlyOpcodeStatement)
	p.registerOpecode("SHR", p.parseGeneralOpcodeStatement)
	p.registerOpecode("STD", p.parseOnlyOpcodeStatement)
	p.registerOpecode("STI", p.parseOnlyOpcodeStatement)
	p.registerOpecode("SUB", p.parseGeneralOpcodeStatement)
	p.registerOpecode("UD2", p.parseOnlyOpcodeStatement)
	p.registerOpecode("WAIT", p.parseOnlyOpcodeStatement)
	p.registerOpecode("RDPMC", p.parseRESBStatement)
	p.registerOpecode("RDTSC", p.parseRESBStatement)
	p.registerOpecode("SAHF", p.parseOnlyOpcodeStatement)
	p.registerOpecode("STC", p.parseOnlyOpcodeStatement)
	p.registerOpecode("WBINVD", p.parseOnlyOpcodeStatement)
	p.registerOpecode("WRMSR", p.parseOnlyOpcodeStatement)

	// 初回のTokenを配列に追加
	p.lexedTokens = append(p.lexedTokens, p.curToken())

	// EOFまでトークンを読み込む
	for {
		tok := p.l.NextToken()
		p.lexedTokens = append(p.lexedTokens, tok)
		if tok.Type == token.EOF {
			log.Println(fmt.Sprintf("debug: p.curIndex max = %d", len(p.lexedTokens)))
			break
		}
	}
	return p
}

func (p *Parser) registerOpecode(opcode string, fn opcodeParseFn) {
	p.opcodeParseFns[opcode] = fn
}

func (p *Parser) maxTokenSize() int {
	return len(p.lexedTokens)
}

func (p *Parser) curToken() token.Token {
	if p.curIndex >= p.maxTokenSize() {
		return token.Token{Type: token.EOF, Literal: ""}
	}
	return p.lexedTokens[p.curIndex]
}

func (p *Parser) peekToken() token.Token {
	if p.curIndex+1 >= p.maxTokenSize() {
		return token.Token{Type: token.EOF, Literal: ""}
	}
	return p.lexedTokens[p.curIndex+1]
}

func (p *Parser) lookAhead(n int) token.Token {
	if p.curIndex+n < p.maxTokenSize() {
		return p.lexedTokens[p.curIndex+n]
	}
	return token.Token{Type: token.ILLEGAL, Literal: ""}
}

func (p *Parser) nextToken() {
	p.curIndex++
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken().Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken().Type == t
}

func (p *Parser) lookAheadIs(n int, t token.TokenType) bool {
	return p.lookAhead(n).Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		t, p.peekToken().Type)
	p.errors = append(p.errors, msg)
}

func isNil(x interface{}) bool {
	return x == nil || reflect.ValueOf(x).IsNil()
}

func isNotNil(x interface{}) bool {
	return !isNil(x)
}

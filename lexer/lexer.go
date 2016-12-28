package lexer

import (
	"github.com/comail/colog"
	"github.com/hangingman/gosk/token"
	"log"
)

// Lexer は入力された文字列に対する現状の検査状況を保持します
type Lexer struct {
	input        []rune // 入力
	position     int    // 現在の文字の位置
	readPosition int    // これから読み込む位置
	ch           rune   // 現在検査中の文字
}

// New は与えられた文字列に対するトークンを返します
func New(input string) *Lexer {
	l := &Lexer{
		input: []rune(input),
	}
	colog.Register()
	colog.SetDefaultLevel(colog.LInfo)
	colog.SetMinLevel(colog.LInfo)
	colog.SetFlags(log.Lshortfile)
	colog.ParseFields(true)
	colog.SetFormatter(&colog.StdFormatter{Colors: false})
	l.readChar()
	return l
}

// readChar は文字列を１文字読み進める
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

// peekChar は次の文字を読んで返す（現在位置は進めない）
func (l *Lexer) peekChar() rune {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

// readIdentifier は識別子を読み出して非英字まで読み進める
func (l *Lexer) readIdentifier() string {
	position := l.position
	log.Printf("debug: [%d] = %s", position, string(l.input[l.position]))
	for isLetter(l.ch) {
		l.readChar()
	}
	ident := string(l.input[position:l.position])

	l.position--
	l.readPosition--

	return ident
}

// isLetter は入力バイトが英字＋アンダーバーであればtrueを返す
func isLetter(ch rune) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_' || isDigit(ch)
}

func newToken(tokenType token.TokenType, ch rune) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

// NextToken は入力を１文字読み出してトークンを返します
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	// 文字がスペースならば読み飛ばす
	l.skipWhitespaceWithLF()

	switch l.ch {
	case '=':
		tok = newToken(token.ASSIGN, l.ch)
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '!':
		tok = newToken(token.BANG, l.ch)
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case ':':
		tok = newToken(token.COLON, l.ch)
	case ';':
		l.skipUntilNextLF() // 行コメント
		return l.NextToken()
	case '#':
		l.skipUntilNextLF() // 行コメント
		return l.NextToken()
	case '<':
		tok = newToken(token.LT, l.ch)
	case '>':
		tok = newToken(token.GT, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case '[':
		tok = newToken(token.LBRACKET, l.ch)
	case ']':
		tok = newToken(token.RBRACKET, l.ch)
	case '"':
		tok.Type = token.STR_LIT
		tok.Literal = l.readDoubleQuotedString()
	case '$':
		tok = newToken(token.DOLLAR, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isHexNotation(l.ch, l.peekChar()) {
			tok.Type = token.HEX_LIT
			tok.Literal = l.readHex()
		} else if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
		} else if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)

			if tok.Type == token.IDENT && l.ch == ':' {
				// ':' で終わる識別子は基本的にラベルとみなす
				l.readChar()
				tok.Literal = tok.Literal + ":"
				tok.Type = token.LABEL
			}
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}
	l.readChar()
	return tok
}

func (l *Lexer) skipUntilNextLF() {
	for !(l.ch == '\n' || l.ch == '\r' || l.ch == 0) {
		l.readChar()
	}
	l.position--
	l.readPosition--
}

func (l *Lexer) readDoubleQuotedString() string {
	position := l.position

	l.readChar() // l.ch を " の次へ
	for !(l.ch == '"' || l.ch == 0) {
		l.readChar()
	}
	return string(l.input[position : l.position+1])
}

func (l *Lexer) skipWhitespaceWithLF() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	number := string(l.input[position:l.position])
	l.position--
	l.readPosition--
	return number
}

func (l *Lexer) readHex() string {
	position := l.position

	l.readChar()
	l.readChar() // l.ch を 0x の次へ

	for isDigit(l.ch) || isHexAlpha(l.ch) {
		l.readChar()
	}
	hex := string(l.input[position:l.position])
	l.position--
	l.readPosition--
	return hex
}

func isHexNotation(ch1 rune, ch2 rune) bool {
	return ch1 == '0' && ch2 == 'x'
}

func isHexAlpha(ch rune) bool {
	return ('a' <= ch && ch <= 'f') || ('A' <= ch && ch <= 'F')
}

func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

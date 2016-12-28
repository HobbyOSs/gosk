package eval

import (
	"fmt"
	"github.com/hangingman/gosk/lexer"
	"github.com/hangingman/gosk/object"
	"github.com/hangingman/gosk/parser"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

func TestEvalMov(t *testing.T) {
	input := `MOV [0x0ff0], CH
MOV BYTE[0x0ff2],8
MOV WORD[0x0ff4],320
MOV WORD[0x0ff6],200
MOV DWORD[0x0ff8],0x000a0000
MOV [0x0ff1] AL
`
	tests := []struct {
		Value []byte
	}{
		{[]byte{0x88, 0x2e, 0xf0, 0x0f}},
		{[]byte{0xc6, 0x06, 0xf2, 0x0f, 0x08}},
		{[]byte{0xc7, 0x06, 0xf4, 0x0f, 0x40, 0x01}},
		{[]byte{0xc7, 0x06, 0xf6, 0x0f, 0xc8, 0x00}},
		{[]byte{0x66, 0xc7, 0x06, 0xf8, 0x0f, 0x00, 0x00, 0x0a, 0x00}},
		{[]byte{0xa2, 0xf1, 0x0f}},
	}

	l := lexer.New(input)
	p := parser.New(l)

	// プログラムの解析と評価
	program := p.ParseProgram()
	evaluated := Eval(program)
	objArray, _ := evaluated.(*object.ObjectArray)
	testTarget := strings.Split(input, "\n")
	// 結果を１つずつ見てみる
	assert.Equal(t, 6, len(*objArray))
	for i, obj := range *objArray {
		bin, _ := obj.(*object.Binary)

		assert.Equal(t,
			tests[i].Value,
			bin.Value,
			fmt.Sprintf("Opcode: [%s] should be...", testTarget[i]))
	}
}

func TestEvalSingleOpcode(t *testing.T) {
	input := `AAA ;; 0x37
AAS ;; 0x3f
CBW ;; 0x98
CDQ ;; 0x99
CLC ;; 0xf8
CLD ;; 0xfc
CLI ;; 0xfa
CLTS ;; []byte{0x0f, 0x06}
CMC ;; 0xf5
CPUID ;; 0xf8
CWD ;; 0x99
CWDE ;; 0x98
DAA ;; 0x27
DAS ;; 0x2f
WAIT ;; 0x9b
HLT ;; 0xf4
INCO ;; 0xce
INSB ;; 0x6c
INSD ;; 0x6d
INSW ;; 0x6d
INVD ;; []byte{0x0f, 0x08}
IRET ;; 0xcf
IRETD ;; 0xcf
LAHF ;; 0x9f
LEAVE ;; 0xc9
LOCK ;; 0xf0
NOP ;; 0x90
OUTSB ;; 0x6e
OUTSD ;; 0x6f
OUTSW ;; 0x6f
POPA ;; 0x61
POPAD ;; 0x61
POPF ;; 0x9d
POPFD ;; 0x9d
PUSHA ;; 0x60
PUSHD ;; 0x60
PUSHF ;; 0x9c
RET ;; 0xc3
RETF ;; 0xcb
RSM ;; []byte{0x0f, 0xaa}
SAHF ;; 0x9e
STC ;; 0xf9
STD ;; 0xfd
STI ;; 0xfb
UD2 ;; []byte{0x0f, 0x0b}
WAIT ;; 0x9b
RDMSR ;; []byte{0x0f, 0x32}
RDPMC ;; []byte{0x0f, 0x33}
RDTSC ;; []byte{0x0f, 0x31}
WBINVD ;; []byte{0x0f, 0x09}
WRMSR ;; []byte{0x0f, 0x30}`

	tests := []struct {
		Value []byte
	}{
		{[]byte{0x37}},
		{[]byte{0x3f}},
		{[]byte{0x98}},
		{[]byte{0x99}},
		{[]byte{0xf8}},
		{[]byte{0xfc}},
		{[]byte{0xfa}},
		{[]byte{0x0f, 0x06}},
		{[]byte{0xf5}},
		{[]byte{0xf8}},
		{[]byte{0x99}},
		{[]byte{0x98}},
		{[]byte{0x27}},
		{[]byte{0x2f}},
		{[]byte{0x9b}},
		{[]byte{0xf4}},
		{[]byte{0xce}},
		{[]byte{0x6c}},
		{[]byte{0x6d}},
		{[]byte{0x6d}},
		{[]byte{0x0f, 0x08}},
		{[]byte{0xcf}},
		{[]byte{0xcf}},
		{[]byte{0x9f}},
		{[]byte{0xc9}},
		{[]byte{0xf0}},
		{[]byte{0x90}},
		{[]byte{0x6e}},
		{[]byte{0x6f}},
		{[]byte{0x6f}},
		{[]byte{0x61}},
		{[]byte{0x61}},
		{[]byte{0x9d}},
		{[]byte{0x9d}},
		{[]byte{0x60}},
		{[]byte{0x60}},
		{[]byte{0x9c}},
		{[]byte{0xc3}},
		{[]byte{0xcb}},
		{[]byte{0x0f, 0xaa}},
		{[]byte{0x9e}},
		{[]byte{0xf9}},
		{[]byte{0xfd}},
		{[]byte{0xfb}},
		{[]byte{0x0f, 0x0b}},
		{[]byte{0x9b}},
		{[]byte{0x0f, 0x32}},
		{[]byte{0x0f, 0x33}},
		{[]byte{0x0f, 0x31}},
		{[]byte{0x0f, 0x09}},
		{[]byte{0x0f, 0x30}},
	}

	l := lexer.New(input)
	p := parser.New(l)

	// プログラムの解析と評価
	program := p.ParseProgram()
	evaluated := Eval(program)
	// リフレクションで結果をチェック
	assert.Equal(t, "*object.ObjectArray", reflect.TypeOf(evaluated).String())
	// キャストをやる
	objArray, ok := evaluated.(*object.ObjectArray)
	assert.True(t, ok)
	testTarget := strings.Split(input, `\n`)
	// 結果を１つずつ見てみる
	for i, obj := range *objArray {
		if len(testTarget) > i {
			assert.NotEqual(t, obj, nil)
			bin, ok := obj.(*object.Binary)
			assert.True(t, ok)
			assert.Equal(t,
				tests[i].Value,
				bin.Value, fmt.Sprintf("Opcode: [%s] should be...", testTarget[i]))
		}
	}
}

func TestEvalAsmHead(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	asmheadPath := path.Join(path.Dir(filename), "..", "testdata", "helloos.nas")
	err := os.Chdir("../../")

	if err != nil {
		panic(err)
	}
	b, err := ioutil.ReadFile(asmheadPath)
	if err != nil {
		fmt.Print(err)
	}
	input := string(b)

	l := lexer.New(input)
	p := parser.New(l)

	// プログラムの解析と評価
	program := p.ParseProgram()
	evaluated := Eval(program)
	// リフレクションで結果をチェック
	assert.Equal(t, "*object.ObjectArray", reflect.TypeOf(evaluated).String())
	// キャストをやる
	objArray, ok := evaluated.(*object.ObjectArray)
	assert.True(t, ok)
	// 結果を１つずつ見てみる
	for _, obj := range *objArray {
		if obj != nil {
			assert.Equal(t, "*object.Binary", reflect.TypeOf(obj).String())
		}
	}
}

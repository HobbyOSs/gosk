package operand

import (
	"testing"

	"github.com/HobbyOSs/gosk/internal/ast"
)

func TestBaseOperand_OperandType(t *testing.T) {
	tests := []struct {
		name      string
		internal  string
		expected  []OperandType
		forceImm8 bool
	}{
		{"Multiple Operands", "EAX, EBX", []OperandType{CodeR32, CodeR32}, false},
		{"Multiple Operands with Memory", "EAX, [EBX]", []OperandType{CodeR32, CodeM32}, false},
		{"Multiple Operands with Immediate", "EAX, 0x10", []OperandType{CodeR32, CodeIMM32}, false},
		{"Multiple Operands with Different Types", "[EAX], EBX", []OperandType{CodeM32, CodeR32}, false},
		{"General Register", " EAX", []OperandType{CodeR32}, false},
		{"Memory Address", "[EBX]", []OperandType{CodeM32}, false},
		{"Memory Address", "AL,[ SI ]", []OperandType{CodeR8, CodeM8}, false},
		{"Immediate Value", "0x10", []OperandType{CodeIMM32}, false},
		{"Immediate Value force imm8", "SI,1", []OperandType{CodeR16, CodeIMM8}, true},
		{"Segment Register", " CS", []OperandType{CodeSREG}, false},
		{"Segmented Address", " DS:SI", []OperandType{CodeM16}, false},
		{"Segmented Address", "DS:BX ", []OperandType{CodeM16}, false},
		{"Segmented Address", "ES:DI", []OperandType{CodeM16}, false},
		{"Segmented Address", "ES:CX", []OperandType{CodeM16}, false},
		{"Relative Offset", "LABEL", []OperandType{CodeREL32}, false},
		{"Relative Offset", "SHORT label", []OperandType{CodeREL8}, false},
		{"Direct Address (with imm and size prefix)", "BYTE [0x1234], 8", []OperandType{CodeM8, CodeIMM8}, false},
		// TODO: サイズプレフィックスがないのでこのパターンはないかも
		//{"Direct Address (simple)", "[0x1234]", []OperandType{CodeM32}, false},
		//{"Direct Address (with imm)", "[0x1234], 8", []OperandType{CodeM32, CodeIMM16}, false},
		{"8-bit Register", "AL", []OperandType{CodeR8}, false},
		{"16-bit Register", "AX", []OperandType{CodeR16}, false},
		{"SP Register", "SP", []OperandType{CodeR16}, false},
		{"BP Register", "BP", []OperandType{CodeR16}, false},
		{"CL Register", "CL", []OperandType{CodeR8}, false},
		{"Complex Memory", "[RAX+4]", []OperandType{CodeM32}, false},
		{"Memory with SP", "[SP+2]", []OperandType{CodeM32}, false},
		{"Memory with BP", "[BP-4]", []OperandType{CodeM32}, false},
		// TODO: まだ未対応
		// {"DWORD PTR", "DWORD PTR [ECX]", []OperandType{CodeM32}, false},
		{"Immediate 10", "10", []OperandType{CodeIMM32}, false},
		{"Immediate Hex", "0xFF", []OperandType{CodeIMM32}, false},
		{"Negative Immediate", "-128", []OperandType{CodeIMM32}, false},
		{"Control Register", "CR0", []OperandType{CodeCR}, false},
		{"Control Register", "CR2", []OperandType{CodeCR}, false},
		{"Control Register", "CR4", []OperandType{CodeCR}, false},
		{"Debug Register", "DR0", []OperandType{CodeDR}, false},
		{"Debug Register", "DR3", []OperandType{CodeDR}, false},
		{"Debug Register", "DR7", []OperandType{CodeDR}, false},
		{"Test Register", "TR3", []OperandType{CodeTR}, false},
		{"Test Register", "TR5", []OperandType{CodeTR}, false},
		{"Test Register", "TR6", []OperandType{CodeTR}, false},
		{"MMX Register", "MM0", []OperandType{CodeMM}, false},
		{"MMX Register", "MM5", []OperandType{CodeMM}, false},
		{"MMX Register", "MM7", []OperandType{CodeMM}, false},
		{"XMM Register", "XMM1", []OperandType{CodeXMM}, false},
		{"XMM Register", "XMM3", []OperandType{CodeXMM}, false},
		{"YMM Register", "YMM4", []OperandType{CodeYMM}, false},
		{"YMM Register", "YMM1", []OperandType{CodeYMM}, false},
		/*
		   TODO: FAR/NEAR PTR の扱いを正しく区別する

		   - FAR PTR (16:16モード)     : 4バイト (セグメント2 + オフセット2)
		   - FAR PTR (16:32モード)     : 6バイト (セグメント2 + オフセット4)
		       - 上記２つはラベルの場合も同様
		   - NEAR PTR / PTR (16ビット) : 2バイト (オフセットのみ)
		   - NEAR PTR / PTR (32ビット) : 4バイト (オフセットのみ)

		   いま CodeM32 としている項目に本来は FAR ポインタか NEAR ポインタかを考慮して、
		   オペランドタイプを区別する必要がある。
		*/
		{"Far Pointer", "FAR PTR [0x5678]", []OperandType{CodeM32}, false},
		{"Pointer + Direct Address", "PTR [0x1234]", []OperandType{CodeM32}, false},
		{"Far Pointer + Direct Address", "FAR PTR [0x5678]", []OperandType{CodeM32}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &OperandImpl{Internal: tt.internal, BitMode: ast.MODE_16BIT, ForceImm8: tt.forceImm8}
			if got := b.OperandTypes(); !equalOperandTypes(got, tt.expected) {
				t.Errorf("OperandType() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestOperandImpl_DetectImmediateSize(t *testing.T) {
	tests := []struct {
		name     string
		internal string
		expected int
	}{
		{"Immediate 8 bit", "0x7f", 1},
		{"Immediate 16 bit", "0x7fff", 2},
		{"Immediate 32 bit", "0x7fffffff", 4},
		// TODO: 負の数のテストがうまくいってない
		//{"Immediate negative 8 bit", "-128", 1},
		//{"Immediate negative 16 bit", "-32768", 2},
		//{"Immediate negative 32 bit", "-2147483648", 4},
		{"No Immediate", "EAX", 0},
		// TODO: 複数のオペランドがある場合のテストがうまくいってない
		//{"Multiple Operands with Immediate", "EAX, 0x10", 4},
		// {"Multiple Operands with Different Immediate Sizes", "EAX, 0x7f", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := OperandImpl{Internal: tt.internal}
			if got := b.DetectImmediateSize(); got != tt.expected {
				t.Errorf("DetectImmediateSize() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestNewOperandFromString(t *testing.T) {
	text := "EAX, EBX"
	expected := "EAX, EBX"
	operand := NewOperandFromString(text)
	if got := operand.InternalString(); got != expected {
		t.Errorf("NewOperandFromString() = %v, want %v", got, expected)
	}
}

func TestOperandImpl_InternalString(t *testing.T) {
	operand := OperandImpl{Internal: "EAX, EBX"}
	expected := "EAX, EBX"
	if got := operand.InternalString(); got != expected {
		t.Errorf("InternalString() = %v, want %v", got, expected)
	}
}

func TestOperandImpl_Serialize(t *testing.T) {
	operand := OperandImpl{Internal: "EAX, EBX"}
	expected := "EAX, EBX"
	if got := operand.Serialize(); got != expected {
		t.Errorf("Serialize() = %v, want %v", got, expected)
	}
}

func TestOperandImpl_FromString(t *testing.T) {
	operand := OperandImpl{}
	text := "EAX, EBX"
	expected := "EAX, EBX"
	newOperand := operand.FromString(text)
	if got := newOperand.InternalString(); got != expected {
		t.Errorf("FromString() = %v, want %v", got, expected)
	}
}

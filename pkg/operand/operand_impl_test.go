package operand

import (
	"testing"

	"github.com/HobbyOSs/gosk/pkg/cpu" // Added import
)

func TestBaseOperand_OperandType(t *testing.T) {
	tests := []struct {
		name          string
		internal      string
		expected      []OperandType
		forceImm8     bool
		forceRelAsImm bool // Add forceRelAsImm field
	}{
		{"Multiple Operands", "EAX, EBX", []OperandType{CodeR32, CodeR32}, false, false},
		{"Multiple Operands with Memory", "EAX, [EBX]", []OperandType{CodeR32, CodeM32}, false, false},
		{"Multiple Operands with Immediate", "EAX, 0x10", []OperandType{CodeR32, CodeIMM32}, false, false},
		{"Multiple Operands with Different Types", "[EAX], EBX", []OperandType{CodeM32, CodeR32}, false, false},
		{"General Register", " EAX", []OperandType{CodeR32}, false, false},
		{"Memory Address", "[EBX]", []OperandType{CodeM32}, false, false},
		{"Memory Address", "AL,[ SI ]", []OperandType{CodeR8, CodeM8}, false, false},
		{"Immediate Value", "0x10", []OperandType{CodeIMM32}, false, false},
		{"Immediate Value force imm8", "SI,1", []OperandType{CodeR16, CodeIMM8}, true, false},
		{"Segment Register", " CS", []OperandType{CodeSREG}, false, false},
		{"Segmented Address", " DS:SI", []OperandType{CodeM16}, false, false},
		{"Segmented Address", "DS:BX ", []OperandType{CodeM16}, false, false},
		{"Segmented Address", "ES:DI", []OperandType{CodeM16}, false, false},
		{"Segmented Address", "ES:CX", []OperandType{CodeM16}, false, false},
		{"Relative Offset", "LABEL", []OperandType{CodeREL32}, false, false}, // JMP/CALL context
		{"Relative Offset", "SHORT label", []OperandType{CodeREL8}, false, false}, // JMP/CALL context
		{"Direct Address (with imm and size prefix)", "BYTE [0x1234], 8", []OperandType{CodeM8, CodeIMM8}, false, false},
		// TODO: サイズプレフィックスがないのでこのパターンはないかも
		//{"Direct Address (simple)", "[0x1234]", []OperandType{CodeM32}, false, false},
		//{"Direct Address (with imm)", "[0x1234], 8", []OperandType{CodeM32, CodeIMM16}, false, false},
		{"8-bit Register", "AL", []OperandType{CodeR8}, false, false},
		{"16-bit Register", "AX", []OperandType{CodeR16}, false, false},
		{"SP Register", "SP", []OperandType{CodeR16}, false, false},
		{"BP Register", "BP", []OperandType{CodeR16}, false, false},
		{"CL Register", "CL", []OperandType{CodeR8}, false, false},
		{"Complex Memory", "[RAX+4]", []OperandType{CodeM32}, false, false},
		{"Memory with SP", "[SP+2]", []OperandType{CodeM32}, false, false},
		{"Memory with BP", "[BP-4]", []OperandType{CodeM32}, false, false},
		// TODO: まだ未対応
		// {"DWORD PTR", "DWORD PTR [ECX]", []OperandType{CodeM32}, false, false},
		{"Immediate 10", "10", []OperandType{CodeIMM32}, false, false},
		{"Immediate Hex", "0xFF", []OperandType{CodeIMM32}, false, false},
		{"Negative Immediate", "-128", []OperandType{CodeIMM32}, false, false},
		{"Control Register", "CR0", []OperandType{CodeCR}, false, false},
		{"Control Register", "CR2", []OperandType{CodeCR}, false, false},
		{"Control Register", "CR4", []OperandType{CodeCR}, false, false},
		{"Debug Register", "DR0", []OperandType{CodeDR}, false, false},
		{"Debug Register", "DR3", []OperandType{CodeDR}, false, false},
		{"Debug Register", "DR7", []OperandType{CodeDR}, false, false},
		{"Test Register", "TR3", []OperandType{CodeTR}, false, false},
		{"Test Register", "TR5", []OperandType{CodeTR}, false, false},
		{"Test Register", "TR6", []OperandType{CodeTR}, false, false},
		{"MMX Register", "MM0", []OperandType{CodeMM}, false, false},
		{"MMX Register", "MM5", []OperandType{CodeMM}, false, false},
		{"MMX Register", "MM7", []OperandType{CodeMM}, false, false},
		{"XMM Register", "XMM1", []OperandType{CodeXMM}, false, false},
		{"XMM Register", "XMM3", []OperandType{CodeXMM}, false, false},
		{"YMM Register", "YMM4", []OperandType{CodeYMM}, false, false},
		{"YMM Register", "YMM1", []OperandType{CodeYMM}, false, false},
		// Harib00i 問題箇所
		{"MOV r32, label", "ESI, bootpack", []OperandType{CodeR32, CodeIMM32}, false, true}, // ForceRelAsImm = true
		{"MOV r32, imm32", "EDI, 0x00280000", []OperandType{CodeR32, CodeIMM32}, false, true}, // ForceRelAsImm = true (immなので影響ないはず)
		{"MOV r32, m (no size prefix)", "ECX, [EBX+16]", []OperandType{CodeR32, CodeM32}, false, false},
		{"ADD r32, r32", "ESI, EBX", []OperandType{CodeR32, CodeR32}, false, false},
		{"ADD r32, imm (small)", "ESI, 4", []OperandType{CodeR32, CodeIMM32}, false, false}, // Expect IMM32 even for small imm in ADD
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
		{"Far Pointer", "FAR PTR [0x5678]", []OperandType{CodeM32}, false, false},
		{"Pointer + Direct Address", "PTR [0x1234]", []OperandType{CodeM32}, false, false},
		{"Far Pointer + Direct Address", "FAR PTR [0x5678]", []OperandType{CodeM32}, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Changed MODE_16BIT to cpu.MODE_16BIT
			// Apply ForceRelAsImm from test case
			b := &OperandImpl{Internal: tt.internal, BitMode: cpu.MODE_16BIT, ForceImm8: tt.forceImm8, ForceRelAsImm: tt.forceRelAsImm}
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

package operand

import (
	"testing"
)

func TestBaseOperand_OperandType(t *testing.T) {
	tests := []struct {
		name     string
		internal string
		expected []OperandType
	}{
		{"Multiple Operands", "EAX, EBX", []OperandType{CodeR32, CodeR32}},
		{"Multiple Operands with Memory", "EAX, [EBX]", []OperandType{CodeR32, CodeM32}},
		{"Multiple Operands with Immediate", "EAX, 0x10", []OperandType{CodeR32, CodeIMM32}},
		{"Multiple Operands with Different Types", "[EAX], EBX", []OperandType{CodeM32, CodeR32}},
		{"General Register", "EAX", []OperandType{CodeR32}},
		{"Memory Address", "[EBX]", []OperandType{CodeM32}},
		{"Immediate Value", "0x10", []OperandType{CodeIMM32}},
		{"Segment Register", "CS", []OperandType{CodeR16}},
		{"Segmented Address", "DS:SI", []OperandType{CodeM16}},
		{"Segmented Address", "DS:BX", []OperandType{CodeM16}},
		{"Segmented Address", "ES:DI", []OperandType{CodeM16}},
		{"Segmented Address", "ES:CX", []OperandType{CodeM16}},
		{"Relative Offset", "LABEL", []OperandType{CodeREL32}},
		{"Relative Offset", "SHORT label", []OperandType{CodeREL8}},
		// TODO: 下記に対応してない
		// FAR PTR label → m16:16 or m16:32
		{"Direct Address", "[0x1234]", []OperandType{CodeM32}},
		{"Direct Address", "BYTE [0x1234], 8", []OperandType{CodeM32, CodeIMM8}},
		{"Direct Address", "[0x1234], 8", []OperandType{CodeM32, CodeIMM16}},
		{"8-bit Register", "AL", []OperandType{CodeR8}},
		{"16-bit Register", "AX", []OperandType{CodeR16}},
		{"CL Register", "CL", []OperandType{CodeR8}},
		{"Complex Memory", "[RAX+4]", []OperandType{CodeM32}},
		// TODO: まだ未対応
		// {"DWORD PTR", "DWORD PTR [ECX]", []OperandType{CodeM32}},
		{"Immediate 10", "10", []OperandType{CodeIMM32}},
		{"Immediate Hex", "0xFF", []OperandType{CodeIMM32}},
		{"Negative Immediate", "-128", []OperandType{CodeIMM32}},
		{"Control Register", "CR0", []OperandType{CodeCR}},
		{"Control Register", "CR2", []OperandType{CodeCR}},
		{"Control Register", "CR4", []OperandType{CodeCR}},
		{"Debug Register", "DR0", []OperandType{CodeDR}},
		{"Debug Register", "DR3", []OperandType{CodeDR}},
		{"Debug Register", "DR7", []OperandType{CodeDR}},
		{"Test Register", "TR3", []OperandType{CodeTR}},
		{"Test Register", "TR5", []OperandType{CodeTR}},
		{"Test Register", "TR6", []OperandType{CodeTR}},
		{"MMX Register", "MM0", []OperandType{CodeMM}},
		{"MMX Register", "MM5", []OperandType{CodeMM}},
		{"MMX Register", "MM7", []OperandType{CodeMM}},
		{"XMM Register", "XMM1", []OperandType{CodeXMM}},
		{"XMM Register", "XMM3", []OperandType{CodeXMM}},
		{"YMM Register", "YMM4", []OperandType{CodeYMM}},
		{"YMM Register", "YMM1", []OperandType{CodeYMM}},
		{"Far Pointer", "FAR PTR [0x5678]", []OperandType{CodeM32}},
		{"Direct Address", "PTR [0x1234]", []OperandType{CodeM32}},
		{"Far Pointer", "FAR PTR [0x5678]", []OperandType{CodeM32}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := OperandImpl{Internal: tt.internal}
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

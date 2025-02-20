package operand

import (
	"testing"
)

func TestBaseOperand_OperandType(t *testing.T) {
	tests := []struct {
		name     string
		internal string
		expected OperandType
	}{
		{"General Register", "EAX", CodeR32},
		{"Memory Address", "[EBX]", CodeM32},
		{"Immediate Value", "0x10", CodeIMM32},
		{"Segment Register", "CS", CodeR16},
		{"Segmented Address", "DS:BX", CodeM16},
		{"Segmented Address", "ES:DI", CodeM16},
		{"Segmented Address", "ES:CX", CodeM16},
		{"Relative Offset", "LABEL", CodeREL32},
		{"Relative Offset", "SHORT label", CodeREL32},
		{"Direct Address", "[0x1234]", CodeM32},
		{"8-bit Register", "AL", CodeR8},
		{"16-bit Register", "AX", CodeR16},
		{"64-bit Register", "RAX", CodeR32},
		{"CL Register", "CL", CodeR8},
		{"Complex Memory", "[RAX+4]", CodeM32},
		{"DWORD PTR", "DWORD PTR [ECX]", CodeM32},
		{"Immediate 10", "10", CodeIMM32},
		{"Immediate Hex", "0xFF", CodeIMM32},
		{"Negative Immediate", "-128", CodeIMM32},
		{"Control Register", "CR0", CodeK},
		{"Control Register", "CR2", CodeK},
		{"Control Register", "CR4", CodeK},
		{"Debug Register", "DR0", CodeK},
		{"Debug Register", "DR3", CodeK},
		{"Debug Register", "DR7", CodeK},
		{"Test Register", "TR3", CodeK},
		{"Test Register", "TR5", CodeK},
		{"Test Register", "TR6", CodeK},
		{"MMX Register", "MM0", CodeMM},
		{"MMX Register", "MM5", CodeMM},
		{"MMX Register", "MM7", CodeMM},
		{"XMM Register", "XMM1", CodeXMM},
		{"XMM Register", "XMM3", CodeXMM},
		{"YMM Register", "YMM4", CodeYMM},
		{"YMM Register", "YMM1", CodeYMM},
		{"XMM Memory", "XMMWORD PTR [RAX]", CodeM128},
		{"YMM Memory", "YMMWORD PTR [RAX]", CodeM256},
		{"Far Pointer", "FAR PTR [0x5678]", CodeM32},
		{"Direct Address", "PTR [1234H]", CodeM32},
		{"Far Pointer", "FAR PTR [5678H]", CodeM32},
		{"Segmented Address", "DS:SI", CodeM16},
		// {"Moffs Address", "MOV AL, [0x1234]", CodeM32},
		// {"Moffs Address", "MOV EAX, [0xABCD]", CodeM32},
		// {"Moffs Address", "MOV RAX, [0x1000]", CodeM32},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := OperandImpl{Internal: tt.internal}
			if got := b.OperandType(); got != tt.expected {
				t.Errorf("OperandType() = %v, want %v", got, tt.expected)
			}
		})
	}
}

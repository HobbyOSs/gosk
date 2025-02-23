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
		{"Segmented Address", "DS:BX", []OperandType{CodeM16}},
		{"Segmented Address", "ES:DI", []OperandType{CodeM16}},
		{"Segmented Address", "ES:CX", []OperandType{CodeM16}},
		{"Relative Offset", "LABEL", []OperandType{CodeREL32}},
		{"Relative Offset", "SHORT label", []OperandType{CodeREL32}},
		{"Direct Address", "[0x1234]", []OperandType{CodeM32}},
		{"8-bit Register", "AL", []OperandType{CodeR8}},
		{"16-bit Register", "AX", []OperandType{CodeR16}},
		{"CL Register", "CL", []OperandType{CodeR8}},
		{"Complex Memory", "[RAX+4]", []OperandType{CodeM32}},
		{"DWORD PTR", "DWORD PTR [ECX]", []OperandType{CodeM32}},
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
		{"XMM Memory", "XMMWORD PTR [RAX]", []OperandType{CodeM128}},
		{"YMM Memory", "YMMWORD PTR [RAX]", []OperandType{CodeM256}},
		{"Far Pointer", "FAR PTR [0x5678]", []OperandType{CodeM32}},
		{"Direct Address", "PTR [1234H]", []OperandType{CodeM32}},
		{"Far Pointer", "FAR PTR [5678H]", []OperandType{CodeM32}},
		{"Segmented Address", "DS:SI", []OperandType{CodeM16}},
		// {"Moffs Address", "MOV AL, [0x1234]", CodeM32},
		// {"Moffs Address", "MOV EAX, [0xABCD]", CodeM32},
		// {"Moffs Address", "MOV RAX, [0x1000]", CodeM32},
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

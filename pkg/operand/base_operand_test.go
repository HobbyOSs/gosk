package operand_test

import (
	"testing"

	"github.com/HobbyOSs/gosk/pkg/operand"
)

func TestBaseOperand_AddressingType(t *testing.T) {
	tests := []struct {
		name     string
		internal string
		expected operand.AddressingType
	}{
		{"General Register", "EAX", operand.CodeGeneralReg},
		{"Memory Address", "[EBX]", operand.CodeModRMAddress},
		{"Immediate Value", "0x10", operand.CodeImmediate},
		{"Segment Register", "CS", operand.CodeSregField},
		{"Relative Offset", "JMP LABEL", operand.CodeRelativeOffset},
		{"Direct Address", "[0x1234]", operand.CodeDirectAddress},
		{"AL", "AL", operand.CodeGeneralReg},
		{"EAX", "EAX", operand.CodeGeneralReg},
		{"RDX", "RDX", operand.CodeGeneralReg},
		{"CL", "CL", operand.CodeGeneralReg},
		{"[EBX]", "[EBX]", operand.CodeModRMAddress},
		{"[RAX+4]", "[RAX+4]", operand.CodeModRMAddress},
		{"DWORD PTR [ECX]", "DWORD PTR [ECX]", operand.CodeModRMAddress},
		{"10", "10", operand.CodeImmediate},
		{"0xFF", "0xFF", operand.CodeImmediate},
		{"-128", "-128", operand.CodeImmediate},
		{"MM0", "MM0", operand.CodeModRM_MMX},
		{"MM5", "MM5", operand.CodeModRM_MMX},
		{"MM7", "MM7", operand.CodeModRM_MMX},
		{"XMM1", "XMM1", operand.CodeXmmRegField},
		{"XMM3", "XMM3", operand.CodeXmmRegField},
		{"XMMWORD PTR [RAX]", "XMMWORD PTR [RAX]", operand.CodeModRMAddress},
		{"YMM1", "YMM1", operand.CodeXmmRMField},
		{"YMM4", "YMM4", operand.CodeXmmRMField},
		{"YMMWORD PTR [RAX]", "YMMWORD PTR [RAX]", operand.CodeModRMAddress},
		{"TR3", "TR3", operand.CodeRegFieldTest},
		{"TR5", "TR5", operand.CodeRegFieldTest},
		{"TR6", "TR6", operand.CodeRegFieldTest},
		{"CR0", "CR0", operand.CodeCRField},
		{"CR2", "CR2", operand.CodeCRField},
		{"CR4", "CR4", operand.CodeCRField},
		{"DR0", "DR0", operand.CodeDebugField},
		{"DR3", "DR3", operand.CodeDebugField},
		{"DR7", "DR7", operand.CodeDebugField},
		{"DS:(E)SI", "DS:(E)SI", operand.CodeMemoryAddressX},
		{"DS:(E)BX", "DS:(E)BX", operand.CodeMemoryAddressX},
		{"ES:(E)DI", "ES:(E)DI", operand.CodeMemoryAddressY},
		{"ES:(E)CX", "ES:(E)CX", operand.CodeMemoryAddressY},
		{"PTR [1234H]", "PTR [1234H]", operand.CodeDirectAddress},
		{"FAR PTR [5678H]", "FAR PTR [5678H]", operand.CodeDirectAddress},
		{"NEAR PTR [0x200]", "NEAR PTR [0x200]", operand.CodeDirectAddress},
		{"CS", "CS", operand.CodeSregField},
		{"DS", "DS", operand.CodeSregField},
		{"SS", "SS", operand.CodeSregField},
		{"JMP SHORT label", "JMP SHORT label", operand.CodeRelativeOffset},
		{"CALL label", "CALL label", operand.CodeRelativeOffset},
		{"JNZ label", "JNZ label", operand.CodeRelativeOffset},
		{"MOV AL, [0x1234]", "MOV AL, [0x1234]", operand.CodeModRMAddressMoffs},
		{"MOV EAX, [0xABCD]", "MOV EAX, [0xABCD]", operand.CodeModRMAddressMoffs},
		{"MOV RAX, [0x1000]", "MOV RAX, [0x1000]", operand.CodeModRMAddressMoffs},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := operand.BaseOperand{Internal: tt.internal}
			if got := b.AddressingType(); got != tt.expected {
				t.Errorf("AddressingType() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// safeSlice returns a substring of the given string from start to end.
// If the string is shorter than the requested slice, it returns as much as possible.
func safeSlice(s string, start, end int) string {
	if len(s) < start {
		return ""
	}
	if len(s) < end {
		return s[start:]
	}
	return s[start:end]
}

func TestBaseOperand_OperandType(t *testing.T) {
	tests := []struct {
		name     string
		internal string
		expected operand.OperandType
	}{
		{"General Register", "EAX", operand.CodeDoubleword},
		{"Memory Address", "[EBX]", operand.CodeDoubleword},
		{"Immediate Value", "0x10", operand.CodeDoublewordInteger},
		{"Segment Register", "CS", operand.CodeWord},
		{"Relative Offset", "JMP LABEL", operand.CodeWord},
		{"Direct Address", "[0x1234]", operand.CodeDoubleword},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := operand.BaseOperand{Internal: tt.internal}
			if got := b.OperandType(); got != tt.expected {
				t.Errorf("OperandType() = %v, want %v", got, tt.expected)
			}
		})
	}
}

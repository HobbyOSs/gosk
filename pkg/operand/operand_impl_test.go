package operand

import (
	"testing"
)

// {"General Register", "EAX", CodeGeneralReg},
// {"Memory Address", "[EBX]", CodeModRMAddress},
// {"Immediate Value", "0x10", CodeImmediate},
// {"Segment Register", "CS", CodeSregField},
//{"Relative Offset", "LABEL", CodeRelativeOffset},
//{"Direct Address", "[0x1234]", CodeDirectAddress},
// {"AL", "AL", CodeGeneralReg},
// {"EAX", "EAX", CodeGeneralReg},
// {"RDX", "RDX", CodeGeneralReg},
// {"CL", "CL", CodeGeneralReg},
// {"[EBX]", "[EBX]", CodeModRMAddress},
// {"[RAX+4]", "[RAX+4]", CodeModRMAddress},
//{"DWORD PTR [ECX]", "DWORD PTR [ECX]", CodeModRMAddress},
// {"10", "10", CodeImmediate},
// {"0xFF", "0xFF", CodeImmediate},
// {"-128", "-128", CodeImmediate},
// {"MM0", "MM0", CodeModRM_MMX},
// {"MM5", "MM5", CodeModRM_MMX},
// {"MM7", "MM7", CodeModRM_MMX},
// {"XMM1", "XMM1", CodeXmmRegField},
// {"XMM3", "XMM3", CodeXmmRegField},
//{"XMMWORD PTR [RAX]", "XMMWORD PTR [RAX]", CodeModRMAddress},
// {"YMM1", "YMM1", CodeXmmRMField},
// {"YMM4", "YMM4", CodeXmmRMField},
//{"YMMWORD PTR [RAX]", "YMMWORD PTR [RAX]", CodeModRMAddress},
// {"TR3", "TR3", CodeRegFieldTest},
// {"TR5", "TR5", CodeRegFieldTest},
// {"TR6", "TR6", CodeRegFieldTest},
// {"CR0", "CR0", CodeCRField},
// {"CR2", "CR2", CodeCRField},
// {"CR4", "CR4", CodeCRField},
// {"DR0", "DR0", CodeDebugField},
// {"DR3", "DR3", CodeDebugField},
// {"DR7", "DR7", CodeDebugField},
//{"PTR [1234H]", "DWORD PTR [0x1234]", CodeDirectAddress},
//{"FAR PTR [5678H]", "FAR PTR [0x5678]", CodeDirectAddress},
// {"DS:SI", "DS:SI", CodeMemoryAddressX},
// {"DS:BX", "DS:BX", CodeMemoryAddressX},
// {"ES:DI", "ES:DI", CodeMemoryAddressY},
// {"ES:CX", "ES:CX", CodeMemoryAddressY},
//{"SHORT label", "SHORT label", CodeRelativeOffset},
//{"MOV AL, [0x1234]", "[0x1234]", CodeModRMAddressMoffs},
//{"MOV EAX, [0xABCD]", "[0xABCD]", CodeModRMAddressMoffs},
//{"MOV RAX, [0x1000]", "[0x1000]", CodeModRMAddressMoffs},

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
		expected OperandType
	}{
		{"General Register", "EAX", CodeDoubleword},
		{"Memory Address", "[EBX]", CodeDoubleword},
		{"Immediate Value", "0x10", CodeDoublewordInteger},
		{"Segment Register", "CS", CodeWord},
		//{"Relative Offset", "LABEL", CodeWord},
		{"Direct Address", "[0x1234]", CodeDoubleword},
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

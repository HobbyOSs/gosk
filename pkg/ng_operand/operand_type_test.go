package ng_operand

import (
	"reflect" // Add reflect for deep equal comparison
	"testing"

	"github.com/HobbyOSs/gosk/pkg/cpu" // Added import
)

// Helper function
func equalOperandTypes(a, b []OperandType) bool {
	return reflect.DeepEqual(a, b)
}

func TestOperandPegImpl_OperandType(t *testing.T) {
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
		{"Relative Offset", "LABEL", []OperandType{CodeREL32}, false, false},      // JMP/CALL context
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
		{"Control Register", "CR0", []OperandType{CodeCREG}, false, false}, // Use CodeCREG
		{"Control Register", "CR2", []OperandType{CodeCREG}, false, false}, // Use CodeCREG
		{"Control Register", "CR4", []OperandType{CodeCREG}, false, false}, // Use CodeCREG
		{"Debug Register", "DR0", []OperandType{CodeDREG}, false, false},   // Use CodeDREG
		{"Debug Register", "DR3", []OperandType{CodeDREG}, false, false},   // Use CodeDREG
		{"Debug Register", "DR7", []OperandType{CodeDREG}, false, false},   // Use CodeDREG
		{"Test Register", "TR3", []OperandType{CodeTREG}, false, false},    // Use CodeTREG
		{"Test Register", "TR5", []OperandType{CodeTREG}, false, false},    // Use CodeTREG
		{"Test Register", "TR6", []OperandType{CodeTREG}, false, false},    // Use CodeTREG
		{"MMX Register", "MM0", []OperandType{CodeMM}, false, false},
		{"MMX Register", "MM5", []OperandType{CodeMM}, false, false},
		{"MMX Register", "MM7", []OperandType{CodeMM}, false, false},
		{"XMM Register", "XMM1", []OperandType{CodeXMM}, false, false},
		{"XMM Register", "XMM3", []OperandType{CodeXMM}, false, false},
		{"YMM Register", "YMM4", []OperandType{CodeYMM}, false, false},
		{"YMM Register", "YMM1", []OperandType{CodeYMM}, false, false},
		// Harib00i 問題箇所
		{"MOV r32, label", "ESI, bootpack", []OperandType{CodeR32, CodeIMM32}, false, true},   // ForceRelAsImm = true
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
			// Determine the appropriate bit mode for the test case
			bitMode := cpu.MODE_16BIT // Default to 16-bit
			// Infer 32-bit mode if expected types contain 32-bit codes
			for _, expType := range tt.expected {
				if expType == CodeR32 || expType == CodeM32 || expType == CodeIMM32 || expType == CodeREL32 {
					bitMode = cpu.MODE_32BIT
					break
				}
			}
			// Override for specific test cases if needed (e.g., if a 16-bit test expects 32-bit context implicitly)
			// Example: if strings.Contains(tt.name, "some specific case") { bitMode = cpu.MODE_32BIT }

			// Call ParseOperands with the determined bit mode and flags
			parsedOperands, err := ParseOperands(tt.internal, bitMode, tt.forceImm8, tt.forceRelAsImm)
			if err != nil {
				t.Fatalf("ParseOperands failed for %q with mode %v: %v", tt.internal, bitMode, err)
			}

			// Create OperandPegImpl and get the resolved types
			opImpl := NewOperandPegImpl(parsedOperands)
			// Set the same bit mode and flags used for parsing to ensure consistency in OperandTypes resolution
			opImpl.WithBitMode(bitMode)
			opImpl.WithForceImm8(tt.forceImm8)
			opImpl.WithForceRelAsImm(tt.forceRelAsImm)

			// Get the resolved types from the implementation
			gotTypes := opImpl.OperandTypes()

			// Compare the resolved types with the expected types
			if !equalOperandTypes(gotTypes, tt.expected) {
				t.Errorf("OperandTypes() with mode %v = %v, want %v", bitMode, gotTypes, tt.expected)
			}
		})
	}
}

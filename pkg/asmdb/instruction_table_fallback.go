package asmdb

import (
	"github.com/samber/lo" // Add import for lo package
)

// Bool is a helper function that returns a pointer to a boolean value.
func Bool(b bool) *bool {
	return &b
}

// addImulFallbackEncodings adds fallback encodings for IMUL instructions.
// This function is called from instruction_table.go:init().
func addImulFallbackEncodings() {
	instructionData.Instructions["IMUL"] = Instruction{
		Summary: "Multiply",
		Forms: []InstructionForm{
			{
				Operands: &[]Operand{
					{Type: "r16", Input: Bool(true), Output: Bool(true)},
					{Type: "imm8", Input: Bool(true), Output: Bool(false)},
				},
				Encodings: []Encoding{
					{
						Opcode:    Opcode{Byte: "6B"},
						ModRM:     &Modrm{Mode: "11", Rm: "#0", Reg: "#0"}, // Corrected: Reg should refer to the destination register (#0)
						Immediate: &Immediate{Size: 1, Value: "#1"},
					},
				},
			},
			{
				Operands: &[]Operand{
					{Type: "r32", Input: Bool(true), Output: Bool(true)},
					{Type: "imm8", Input: Bool(true), Output: Bool(false)},
				},
				Encodings: []Encoding{
					{
						Opcode:    Opcode{Byte: "6B"},
						ModRM:     &Modrm{Mode: "11", Rm: "#0", Reg: "#0"}, // Corrected: Reg should refer to the destination register (#0)
						Immediate: &Immediate{Size: 1, Value: "#1"},
					},
				},
			},
			{
				Operands: &[]Operand{
					{Type: "r16", Input: Bool(true), Output: Bool(true)},
					{Type: "imm16", Input: Bool(true), Output: Bool(false)},
				},
				Encodings: []Encoding{
					{
						Opcode:    Opcode{Byte: "69"},
						ModRM:     &Modrm{Mode: "11", Rm: "#0", Reg: "#1"},
						Immediate: &Immediate{Size: 2, Value: "#1"},
					},
				},
			},
			{
				Operands: &[]Operand{
					{Type: "r32", Input: Bool(true), Output: Bool(true)},
					{Type: "imm32", Input: Bool(true), Output: Bool(false)},
				},
				Encodings: []Encoding{
					{
						Opcode:    Opcode{Byte: "69"},
						ModRM:     &Modrm{Mode: "11", Rm: "#0", Reg: "#1"},
						Immediate: &Immediate{Size: 4, Value: "#1"},
					},
				},
			},
		},
	}
}

// addOutFallbackEncodings adds fallback encodings for OUT instructions.
// This function is called from instruction_table.go:init().
func addOutFallbackEncodings() {
	instructionData.Instructions["OUT"] = Instruction{
		Summary: "Output to Port",
		Forms: []InstructionForm{
			// OUT imm8, AL (Intel syntax: OUT port, data)
			{
				Operands: &[]Operand{
					{Type: "imm8", Input: Bool(true), Output: Bool(false)}, // Port
					{Type: "al", Input: Bool(true), Output: Bool(false)},   // Data
				},
				Encodings: []Encoding{
					{
						Opcode:    Opcode{Byte: "E6"},
						Immediate: &Immediate{Size: 1, Value: "#0"}, // imm8 is the first operand (#0)
					},
				},
			},
			// OUT imm8, AX
			{
				Operands: &[]Operand{
					{Type: "imm8", Input: Bool(true), Output: Bool(false)}, // Port
					{Type: "ax", Input: Bool(true), Output: Bool(false)},   // Data
				},
				Encodings: []Encoding{
					{
						Opcode:    Opcode{Byte: "E7"},
						Immediate: &Immediate{Size: 1, Value: "#0"}, // imm8 is the first operand (#0)
					},
				},
			},
			// OUT imm8, EAX
			{
				Operands: &[]Operand{
					{Type: "imm8", Input: Bool(true), Output: Bool(false)}, // Port
					{Type: "eax", Input: Bool(true), Output: Bool(false)},  // Data
				},
				Encodings: []Encoding{
					{
						Opcode:    Opcode{Byte: "E7"},
						Immediate: &Immediate{Size: 1, Value: "#0"}, // imm8 is the first operand (#0)
					},
				},
			},
			{
				Operands: &[]Operand{
					{Type: "dx", Input: Bool(false), Output: Bool(false)},
					{Type: "al", Input: Bool(true), Output: Bool(false)},
				},
				Encodings: []Encoding{
					{
						Opcode: Opcode{Byte: "EE"},
					},
				},
			},
			{
				Operands: &[]Operand{
					{Type: "dx", Input: Bool(false), Output: Bool(false)},
					{Type: "ax", Input: Bool(true), Output: Bool(false)},
				},
				Encodings: []Encoding{
					{
						Opcode: Opcode{Byte: "EF"},
					},
				},
			},
			{
				Operands: &[]Operand{
					{Type: "dx", Input: Bool(false), Output: Bool(false)},
					{Type: "eax", Input: Bool(true), Output: Bool(false)},
				},
				Encodings: []Encoding{
					{
						Opcode: Opcode{Byte: "EF"},
					},
				},
			},
		},
	}
}

func addMovFallbackEncodings() {
	// "MOV" 命令の既存の Forms を取得
	currentMOVInstructionForms := instructionData.Instructions["MOV"].Forms
	// 新しい Forms を追加
	newMOVInstructionForms := append(currentMOVInstructionForms,
		InstructionForm{
			Operands: &[]Operand{
				{Type: "r16", Input: Bool(false), Output: Bool(true)},
				{Type: "sreg", Input: Bool(true), Output: Bool(false)},
			},
			Encodings: []Encoding{
				{
					Opcode: Opcode{Byte: "8C"},
					ModRM:  &Modrm{Mode: "11", Reg: "#1", Rm: "#1"},
				},
			},
		},
		InstructionForm{
			Operands: &[]Operand{
				{Type: "sreg", Input: Bool(false), Output: Bool(true)},
				{Type: "r16", Input: Bool(true), Output: Bool(false)},
			},
			Encodings: []Encoding{
				{
					Opcode: Opcode{Byte: "8E"},
					ModRM:  &Modrm{Mode: "11", Reg: "#0", Rm: "#1"},
				},
			},
		},
	)

	// https://www.felixcloutier.com/x86/mov
	// A0 	MOV AL, moffs8 	（セグメント：オフセット）のバイトをALに転送します
	// A1 	MOV AX, moffs16 	（セグメント：オフセット）のワードをAXに転送します
	// A1 	MOV EAX, moffs32 	（セグメント：オフセット）のダブルワードをEAXに転送します
	// A2 	MOV moffs8, AL 	ALを（セグメント：オフセット）に転送します
	// A3 	MOV moffs16, AX 	AXを（セグメント：オフセット）に転送します
	// A3 	MOV moffs32, EAX 	EAXを（セグメント：オフセット）に転送します
	newMOVInstructionForms = append(newMOVInstructionForms,
		InstructionForm{
			Operands: &[]Operand{
				{Type: "al", Input: Bool(false), Output: Bool(true)},
				{Type: "m8", Input: Bool(true), Output: Bool(false)},
			},
			Encodings: []Encoding{
				{
					Opcode: Opcode{Byte: "A0"},
				},
			},
		},
		InstructionForm{
			Operands: &[]Operand{
				{Type: "ax", Input: Bool(false), Output: Bool(true)},
				{Type: "m16", Input: Bool(true), Output: Bool(false)},
			},
			Encodings: []Encoding{
				{
					Opcode: Opcode{Byte: "A1"},
				},
			},
		},
		InstructionForm{
			Operands: &[]Operand{
				{Type: "eax", Input: Bool(false), Output: Bool(true)},
				{Type: "m32", Input: Bool(true), Output: Bool(false)},
			},
			Encodings: []Encoding{
				{
					Opcode: Opcode{Byte: "A1"},
				},
			},
		},
		InstructionForm{
			Operands: &[]Operand{
				{Type: "m8", Input: Bool(false), Output: Bool(true)},
				{Type: "al", Input: Bool(true), Output: Bool(false)},
			},
			Encodings: []Encoding{
				{
					Opcode: Opcode{Byte: "A2"},
				},
			},
		},
		InstructionForm{
			Operands: &[]Operand{
				{Type: "m16", Input: Bool(false), Output: Bool(true)},
				{Type: "ax", Input: Bool(true), Output: Bool(false)},
			},
			Encodings: []Encoding{
				{
					Opcode: Opcode{Byte: "A3"},
				},
			},
		},
		InstructionForm{
			Operands: &[]Operand{
				{Type: "m32", Input: Bool(false), Output: Bool(true)},
				{Type: "eax", Input: Bool(true), Output: Bool(false)},
			},
			Encodings: []Encoding{
				{
					Opcode: Opcode{Byte: "A3"},
				},
			},
		},
		// Add MOV r32, CR0 (0F 20 /r)
		InstructionForm{
			Operands: &[]Operand{
				{Type: "r32", Input: Bool(false), Output: Bool(true)},
				{Type: "creg", Input: Bool(true), Output: Bool(false)}, // Use "creg" type defined in operand_types.go
			},
			Encodings: []Encoding{
				{
					Opcode: Opcode{Byte: "0F20"},
					ModRM:  &Modrm{Mode: "11", Reg: "#1", Rm: "#0"}, // reg should be CR number (0 for CR0), rm is r32
				},
			},
		},
		// Add MOV CR0, r32 (0F 22 /r)
		InstructionForm{
			Operands: &[]Operand{
				{Type: "creg", Input: Bool(false), Output: Bool(true)}, // Use "creg" type defined in operand_types.go
				{Type: "r32", Input: Bool(true), Output: Bool(false)},
			},
			Encodings: []Encoding{
				{
					Opcode: Opcode{Byte: "0F22"},
					ModRM:  &Modrm{Mode: "11", Reg: "#0", Rm: "#1"}, // reg should be CR number (0 for CR0), rm is r32
				},
			},
		},
	)

	// 更新された Forms で "MOV" 命令を更新
	instructionData.Instructions["MOV"] = Instruction{
		Summary: instructionData.Instructions["MOV"].Summary,
		Forms:   newMOVInstructionForms,
	}
}

// addInFallbackEncodings adds fallback encodings for IN instructions.
// This function is called from instruction_table.go:init().
func addInFallbackEncodings() {
	instructionData.Instructions["IN"] = Instruction{
		Summary: "Input from Port",
		Forms: []InstructionForm{
			{
				Operands: &[]Operand{
					{Type: "al", Input: Bool(false), Output: Bool(true)},   // Destination
					{Type: "imm8", Input: Bool(true), Output: Bool(false)}, // Source (Port)
				},
				Encodings: []Encoding{
					{
						Opcode:    Opcode{Byte: "E4"},
						Immediate: &Immediate{Size: 1, Value: "#1"}, // imm8 is the second operand (#1)
					},
				},
			},
			{
				Operands: &[]Operand{
					{Type: "ax", Input: Bool(false), Output: Bool(true)},   // Destination
					{Type: "imm8", Input: Bool(true), Output: Bool(false)}, // Source (Port)
				},
				Encodings: []Encoding{
					{
						Opcode:    Opcode{Byte: "E5"},
						Immediate: &Immediate{Size: 1, Value: "#1"}, // imm8 is the second operand (#1)
					},
				},
			},
			// Note: EAX form might exist for 32-bit, but sticking to AL/AX for now based on Wikipedia 8086 info
			{
				Operands: &[]Operand{
					{Type: "al", Input: Bool(false), Output: Bool(true)}, // Destination
					{Type: "dx", Input: Bool(true), Output: Bool(false)}, // Source (Port)
				},
				Encodings: []Encoding{
					{
						Opcode: Opcode{Byte: "EC"},
					},
				},
			},
			{
				Operands: &[]Operand{
					{Type: "ax", Input: Bool(false), Output: Bool(true)}, // Destination
					{Type: "dx", Input: Bool(true), Output: Bool(false)}, // Source (Port)
				},
				Encodings: []Encoding{
					{
						Opcode: Opcode{Byte: "ED"},
					},
				},
			},
			// Add IN EAX, DX form (Opcode ED, same as AX but no 66h prefix in 32bit)
			{
				Operands: &[]Operand{
					{Type: "eax", Input: Bool(false), Output: Bool(true)}, // Destination
					{Type: "dx", Input: Bool(true), Output: Bool(false)},  // Source (Port)
				},
				Encodings: []Encoding{
					{
						Opcode: Opcode{Byte: "ED"},
					},
				},
			},
		},
	}
}

// addLgdtFallbackEncodings adds fallback encodings for LGDT instruction.
// This function is called from instruction_table.go:init().
func addLgdtFallbackEncodings() {
	instructionData.Instructions["LGDT"] = Instruction{
		Summary: "Load Global Descriptor Table Register",
		Forms: []InstructionForm{
			{
				Operands: &[]Operand{
					{Type: "mem", Input: Bool(true), Output: Bool(false)}, // Change type from "m" to "mem"
				},
				Encodings: []Encoding{
					{
						Opcode: Opcode{Byte: "0F01"},
						ModRM:  &Modrm{Mode: "#0", Reg: "02", Rm: "#0"}, // Reg: 02 (/2)
					},
				},
			},
		},
	}
}

func init() {
	addImulFallbackEncodings()
	addOutFallbackEncodings()
	addMovFallbackEncodings()
	addLgdtFallbackEncodings()
	addInFallbackEncodings() // Add call to the new function
}

// addPushPopFallbackEncodings adds fallback encodings for PUSH/POP r32 instructions.
// JSON data seems to be missing r32 forms.
func addPushPopFallbackEncodings() {
	// PUSH r32 (Opcode 50+rd)
	// Check if PUSH already exists, if so append, otherwise create
	pushInst, pushExists := instructionData.Instructions["PUSH"]
	if !pushExists {
		pushInst = Instruction{Summary: "Push Word, Doubleword or Quadword Onto the Stack"}
	}
	pushInst.Forms = append(pushInst.Forms, InstructionForm{
		Operands: &[]Operand{
			{Type: "r32", Input: Bool(true), Output: Bool(false)},
		},
		Encodings: []Encoding{
			{
				Opcode: Opcode{Byte: "50", Addend: lo.ToPtr("#0")}, // 50 + register number
			},
			// Also add FF /6 form for r/m32? JSON has it for r/m16, r/m64
			// {
			// 	Opcode: Opcode{Byte: "FF"},
			// 	ModRM:  &Modrm{Mode: "11", Rm: "#0", Reg: "6"},
			// },
		},
	})
	instructionData.Instructions["PUSH"] = pushInst

	// POP r32 (Opcode 58+rd)
	// Check if POP already exists, if so append, otherwise create
	popInst, popExists := instructionData.Instructions["POP"]
	if !popExists {
		popInst = Instruction{Summary: "Pop a Value from the Stack"}
	}
	popInst.Forms = append(popInst.Forms, InstructionForm{
		Operands: &[]Operand{
			{Type: "r32", Input: Bool(false), Output: Bool(true)},
		},
		Encodings: []Encoding{
			{
				Opcode: Opcode{Byte: "58", Addend: lo.ToPtr("#0")}, // 58 + register number
			},
			// Also add 8F /0 form for r/m32? JSON has it for r/m16, r/m64
			// {
			// 	Opcode: Opcode{Byte: "8F"},
			// 	ModRM:  &Modrm{Mode: "11", Rm: "#0", Reg: "0"},
			// },
		},
	})
	instructionData.Instructions["POP"] = popInst
}

func init() {
	addImulFallbackEncodings()
	addOutFallbackEncodings()
	addMovFallbackEncodings()
	addLgdtFallbackEncodings()
	addInFallbackEncodings()      // Add call to the new function
	addPushPopFallbackEncodings() // Add call for PUSH/POP
}

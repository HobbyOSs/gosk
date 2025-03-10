package asmdb

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
						ModRM:     &Modrm{Mode: "11", Rm: "#0", Reg: "#1"},
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
						ModRM:     &Modrm{Mode: "11", Rm: "#0", Reg: "#1"},
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
			{
				Operands: &[]Operand{
					{Type: "al", Input: Bool(true), Output: Bool(false)},
					{Type: "imm8", Input: Bool(true), Output: Bool(false)},
				},
				Encodings: []Encoding{
					{
						Opcode:    Opcode{Byte: "E6"},
						Immediate: &Immediate{Size: 1, Value: "#0"},
					},
				},
			},
			{
				Operands: &[]Operand{
					{Type: "ax", Input: Bool(true), Output: Bool(false)},
					{Type: "imm8", Input: Bool(true), Output: Bool(false)},
				},
				Encodings: []Encoding{
					{
						Opcode:    Opcode{Byte: "E7"},
						Immediate: &Immediate{Size: 1, Value: "#0"},
					},
				},
			},
			{
				Operands: &[]Operand{
					{Type: "eax", Input: Bool(true), Output: Bool(false)},
					{Type: "imm8", Input: Bool(true), Output: Bool(false)},
				},
				Encodings: []Encoding{
					{
						Opcode:    Opcode{Byte: "E7"},
						Immediate: &Immediate{Size: 1, Value: "#0"},
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

func addMovSegmentRegisterEncodings(instructionData *InstructionData) {
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
					ModRM:  &Modrm{Mode: "11", Reg: "#2", Rm: "#1"},
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
					ModRM:  &Modrm{Mode: "11", Reg: "#2", Rm: "#0"},
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

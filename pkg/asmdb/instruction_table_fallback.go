package asmdb

// addImulFallbackEncodings adds fallback encodings for IMUL instructions.
// This function is called from instruction_table.go:init().
func addImulFallbackEncodings() {
	instructionData.Instructions["IMUL r16, imm8"] = Instruction{
		Summary: "Multiply r16 by imm8",
		Forms: []InstructionForm{
			{
				Operands: &[]Operand{
					{Type: "r16", Input: new(bool), Output: new(bool)},
					{Type: "imm8", Input: new(bool), Output: new(bool)},
				},
				Encodings: []Encoding{
					{
						Opcode:    Opcode{Byte: "6B"},
						ModRM:     &Modrm{Mode: "11", Rm: "#0", Reg: "#1"},
						Immediate: &Immediate{Size: 1, Value: "#1"},
					},
				},
			},
		},
	}
	*(*instructionData.Instructions["IMUL r16, imm8"].Forms[0].Operands)[0].Output = false
	*(*instructionData.Instructions["IMUL r16, imm8"].Forms[0].Operands)[1].Input = false

	instructionData.Instructions["IMUL r32, imm8"] = Instruction{
		Summary: "Multiply r32 by imm8",
		Forms: []InstructionForm{
			{
				Operands: &[]Operand{
					{Type: "r32", Input: new(bool), Output: new(bool)},
					{Type: "imm8", Input: new(bool), Output: new(bool)},
				},
				Encodings: []Encoding{
					{
						Opcode:    Opcode{Byte: "6B"},
						ModRM:     &Modrm{Mode: "11", Rm: "#0", Reg: "#1"},
						Immediate: &Immediate{Size: 1, Value: "#1"},
					},
				},
			},
		},
	}
	*(*instructionData.Instructions["IMUL r32, imm8"].Forms[0].Operands)[0].Output = false
	*(*instructionData.Instructions["IMUL r32, imm8"].Forms[0].Operands)[1].Input = false

	instructionData.Instructions["IMUL r16, imm16"] = Instruction{
		Summary: "Multiply r16 by imm16",
		Forms: []InstructionForm{
			{
				Operands: &[]Operand{
					{Type: "r16", Input: new(bool), Output: new(bool)},
					{Type: "imm16", Input: new(bool), Output: new(bool)},
				},
				Encodings: []Encoding{
					{
						Opcode:    Opcode{Byte: "69"},
						ModRM:     &Modrm{Mode: "11", Rm: "#0", Reg: "#1"},
						Immediate: &Immediate{Size: 2, Value: "#1"},
					},
				},
			},
		},
	}
	*(*instructionData.Instructions["IMUL r16, imm16"].Forms[0].Operands)[0].Output = false
	*(*instructionData.Instructions["IMUL r16, imm16"].Forms[0].Operands)[1].Input = false

	instructionData.Instructions["IMUL r32, imm32"] = Instruction{
		Summary: "Multiply r32 by imm32",
		Forms: []InstructionForm{
			{
				Operands: &[]Operand{
					{Type: "r32", Input: new(bool), Output: new(bool)},
					{Type: "imm32", Input: new(bool), Output: new(bool)},
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
	*(*instructionData.Instructions["IMUL r32, imm32"].Forms[0].Operands)[0].Output = false
	*(*instructionData.Instructions["IMUL r32, imm32"].Forms[0].Operands)[1].Input = false
}

// addOutFallbackEncodings adds fallback encodings for OUT instructions.
// This function is called from instruction_table.go:init().
func addOutFallbackEncodings() {
	instructionData.Instructions["OUT imm8, al"] = Instruction{
		Summary: "Output to Port",
		Forms: []InstructionForm{
			{
				Operands: &[]Operand{
					{Type: "imm8", Input: new(bool), Output: new(bool)},
					{Type: "al", Input: new(bool), Output: new(bool)},
				},
				Encodings: []Encoding{
					{
						Opcode:    Opcode{Byte: "E6"},
						Immediate: &Immediate{Size: 1, Value: "#0"},
					},
				},
			},
		},
	}
	*(*instructionData.Instructions["OUT imm8, al"].Forms[0].Operands)[0].Input = false
	*(*instructionData.Instructions["OUT imm8, al"].Forms[0].Operands)[1].Output = false

	instructionData.Instructions["OUT imm8, ax"] = Instruction{
		Summary: "Output to Port",
		Forms: []InstructionForm{
			{
				Operands: &[]Operand{
					{Type: "imm8", Input: new(bool), Output: new(bool)},
					{Type: "ax", Input: new(bool), Output: new(bool)},
				},
				Encodings: []Encoding{
					{
						Opcode:    Opcode{Byte: "E7"},
						Immediate: &Immediate{Size: 1, Value: "#0"},
					},
				},
			},
		},
	}
	*(*instructionData.Instructions["OUT imm8, ax"].Forms[0].Operands)[0].Input = false
	*(*instructionData.Instructions["OUT imm8, ax"].Forms[0].Operands)[1].Output = false

	instructionData.Instructions["OUT imm8, eax"] = Instruction{
		Summary: "Output to Port",
		Forms: []InstructionForm{
			{
				Operands: &[]Operand{
					{Type: "imm8", Input: new(bool), Output: new(bool)},
					{Type: "eax", Input: new(bool), Output: new(bool)},
				},
				Encodings: []Encoding{
					{
						Opcode:    Opcode{Byte: "E7"},
						Immediate: &Immediate{Size: 1, Value: "#0"},
					},
				},
			},
		},
	}
	*(*instructionData.Instructions["OUT imm8, eax"].Forms[0].Operands)[0].Input = false
	*(*instructionData.Instructions["OUT imm8, eax"].Forms[0].Operands)[1].Output = false

	instructionData.Instructions["OUT dx, al"] = Instruction{
		Summary: "Output to Port",
		Forms: []InstructionForm{
			{
				Operands: &[]Operand{
					{Type: "dx", Input: new(bool), Output: new(bool)},
					{Type: "al", Input: new(bool), Output: new(bool)},
				},
				Encodings: []Encoding{
					{
						Opcode: Opcode{Byte: "EE"},
					},
				},
			},
		},
	}
	*(*instructionData.Instructions["OUT dx, al"].Forms[0].Operands)[0].Output = false
	*(*instructionData.Instructions["OUT dx, al"].Forms[0].Operands)[1].Output = false

	instructionData.Instructions["OUT dx, ax"] = Instruction{
		Summary: "Output to Port",
		Forms: []InstructionForm{
			{
				Operands: &[]Operand{
					{Type: "dx", Input: new(bool), Output: new(bool)},
					{Type: "ax", Input: new(bool), Output: new(bool)},
				},
				Encodings: []Encoding{
					{
						Opcode: Opcode{Byte: "EF"},
					},
				},
			},
		},
	}
	*(*instructionData.Instructions["OUT dx, ax"].Forms[0].Operands)[0].Output = false
	*(*instructionData.Instructions["OUT dx, ax"].Forms[0].Operands)[1].Output = false

	instructionData.Instructions["OUT dx, eax"] = Instruction{
		Summary: "Output to Port",
		Forms: []InstructionForm{
			{
				Operands: &[]Operand{
					{Type: "dx", Input: new(bool), Output: new(bool)},
					{Type: "eax", Input: new(bool), Output: new(bool)},
				},
				Encodings: []Encoding{
					{
						Opcode: Opcode{Byte: "EF"},
					},
				},
			},
		},
	}
	*(*instructionData.Instructions["OUT dx, eax"].Forms[0].Operands)[0].Output = false
	*(*instructionData.Instructions["OUT dx, eax"].Forms[0].Operands)[1].Output = false
}

func addMovSegmentRegisterEncodings(instructionData *InstructionData) {
	// "MOV" 命令の既存の Forms を取得
	currentMOVInstructionForms := instructionData.Instructions["MOV"].Forms
	// 新しい Forms を追加
	newMOVInstructionForms := append(currentMOVInstructionForms,
		InstructionForm{
			Operands: &[]Operand{
				{Type: "r16", Input: new(bool), Output: new(bool)},
				{Type: "sreg", Input: new(bool), Output: new(bool)},
			},
			Encodings: []Encoding{
				{
					Opcode: Opcode{Byte: "8C"},
					ModRM:  &Modrm{Mode: "11", Reg: "#1", Rm: "#0"},
				},
			},
		},
		InstructionForm{
			Operands: &[]Operand{
				{Type: "sreg", Input: new(bool), Output: new(bool)},
				{Type: "r16", Input: new(bool), Output: new(bool)},
			},
			Encodings: []Encoding{
				{
					Opcode: Opcode{Byte: "8E"},
					ModRM:  &Modrm{Mode: "11", Reg: "#0", Rm: "#1"},
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

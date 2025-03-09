package asmdb

func addSegmentRegisterEncodings() {
	// セグメントレジスタのエンコーディングを追加
	instructionData.Instructions["MOV r16, Sreg"] = Instruction{
		Summary: "Move segment register to r16",
		Forms: []InstructionForm{
			{
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
		},
	}
    *(*instructionData.Instructions["MOV r16, Sreg"].Forms[0].Operands)[0].Output = false
	*(*instructionData.Instructions["MOV r16, Sreg"].Forms[0].Operands)[1].Input = false

	instructionData.Instructions["MOV Sreg, r16"] = Instruction{
		Summary: "Move r16 to segment register",
		Forms: []InstructionForm{
			{
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
		},
	}
	*(*instructionData.Instructions["MOV Sreg, r16"].Forms[0].Operands)[0].Output = false
	*(*instructionData.Instructions["MOV Sreg, r16"].Forms[0].Operands)[1].Input = false
}

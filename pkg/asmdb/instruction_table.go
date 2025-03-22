package asmdb

import (
	"encoding/json"
	"log"
)

type Isa struct {
	ID string `json:"id"`
}

type DataOffset struct {
	Size  int    `json:"size"`
	Value string `json:"value"`
}

type CodeOffset struct {
	Size  int    `json:"size"`
	Value string `json:"value"`
}

type Prefix struct {
	Mandatory bool   `json:"mandatory"`
	Byte      string `json:"byte"`
}

type Rex struct {
	Mandatory bool    `json:"mandatory"`
	W         *string `json:"W,omitempty"`
	R         *string `json:"R,omitempty"`
	B         *string `json:"B,omitempty"`
	X         *string `json:"X,omitempty"`
}

type Vex struct {
	Mmmmm *string `json:"mmmmm,omitempty"`
	Pp    *string `json:"pp,omitempty"`
	W     *string `json:"W,omitempty"`
	L     *string `json:"L,omitempty"`
	R     *string `json:"R,omitempty"`
	X     *string `json:"X,omitempty"`
}

type Modrm struct {
	Mode string `json:"mode"` // 想定値 "#0", "#1", "#2", "11"
	Rm   string `json:"rm"`   // 想定値 "#0", "#1"  rm用のoperand 0-indexed
	Reg  string `json:"reg"`  // 想定値 "#0", "#1"  reg用のoperand 0-indexed
}

type Operand struct {
	Type         string `json:"type"`
	Input        *bool  `json:"input,omitempty"`
	Output       *bool  `json:"output,omitempty"`
	ExtendedSize *int   `json:"extended_size,omitempty"`
}

type ImplicitOperand struct {
	ID     string `json:"id"`
	Input  *bool  `json:"input,omitempty"`
	Output *bool  `json:"output,omitempty"`
}

type Immediate struct {
	Size  int    `json:"size"`
	Value string `json:"value"`
}

type Opcode struct {
	Byte   string  `json:"byte"`
	Addend *string `json:"addend,omitempty"`
}

type Encoding struct {
	Opcode     Opcode      `json:"opcode"`
	Prefix     *Prefix     `json:"prefix,omitempty"`
	REX        *Rex        `json:"REX,omitempty"`
	VEX        *Vex        `json:"VEX,omitempty"`
	ModRM      *Modrm      `json:"ModRM,omitempty"`
	Immediate  *Immediate  `json:"immediate,omitempty"`
	DataOffset *DataOffset `json:"data_offset,omitempty"`
	CodeOffset *CodeOffset `json:"code_offset,omitempty"`
}

type InstructionForm struct {
	Encodings        []Encoding         `json:"encodings"`
	Operands         *[]Operand         `json:"operands,omitempty"`
	ImplicitOperands *[]ImplicitOperand `json:"implicit_operands,omitempty"`
	XmmMode          *string            `json:"xmm_mode,omitempty"`
	CancellingInputs *bool              `json:"cancelling_inputs,omitempty"`
	Isa              *[]Isa             `json:"isa,omitempty"`
}

type Instruction struct {
	Summary string            `json:"summary"`
	Forms   []InstructionForm `json:"forms"`
}

type InstructionData struct {
	InstructionSet string                 `json:"instruction_set"`
	Instructions   map[string]Instruction `json:"instructions"`
}

var instructionData InstructionData

func init() {
	data, err := decompressGzip(compressedJSON)
	if err != nil {
		log.Fatalf("Failed to decompress JSON: %v", err)
	}

	// Unmarshal the JSON data into the instructions map
	if err := json.Unmarshal(data, &instructionData); err != nil {
		log.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// fallback
	addMovFallbackEncodings(&instructionData)
	addImulFallbackEncodings()
	addOutFallbackEncodings()
}

func GetInstructionByOpcode(opcode string) (*Instruction, error) {
	instr, exists := instructionData.Instructions[opcode]
	if !exists {
		return nil, nil
	}
	return &instr, nil
}

func X86Instructions() map[string]Instruction {
	return instructionData.Instructions
}

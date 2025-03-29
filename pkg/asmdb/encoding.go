package asmdb

import (
	"fmt"
	"log"
	"strings"
)

type OutputSizeOptions struct {
	ImmSize int // 即値サイズ
}

func (e *Encoding) GetOutputSize(options *OutputSizeOptions) int {
	outputSize := 0
	var builder strings.Builder

	// Calculate size based on REX
	if e.REX != nil && e.REX.Mandatory {
		size := e.REX.getSize()
		builder.WriteString(fmt.Sprintf(" rex:%d", size))
		outputSize += size
	}

	// Calculate size based on VEX
	if e.VEX != nil {
		size := e.VEX.getSize()
		builder.WriteString(fmt.Sprintf(" vex:%d", size))
		outputSize += size
	}

	// Calculate size based on Opcode
	sizeOpcode := e.Opcode.getSize()
	builder.WriteString(fmt.Sprintf(" opcode:%d", sizeOpcode))
	outputSize += sizeOpcode

	// Calculate size based on ModRM
	if e.ModRM != nil {
		size := e.ModRM.getSize()
		builder.WriteString(fmt.Sprintf(" modrm:%d", size))
		outputSize += size
	}

	// Calculate size based on Immediate
	if e.Immediate != nil {
		var immSize int = e.Immediate.Size
		if options != nil && options.ImmSize > 0 {
			immSize = options.ImmSize
		}
		builder.WriteString(fmt.Sprintf(" immediate:%d", immSize))
		outputSize += immSize
	}

	// Calculate size based on DataOffset
	if e.DataOffset != nil {
		size := e.DataOffset.Size
		builder.WriteString(fmt.Sprintf(" data_offset:%d", size))
		outputSize += size
	}

	// Calculate size based on CodeOffset
	if e.CodeOffset != nil {
		size := e.CodeOffset.Size
		builder.WriteString(fmt.Sprintf(" code_offset:%d", size))
		outputSize += size
	}

	log.Printf("trace: [pass1] output_size_details:%s total:%d", builder.String(), outputSize)
	return outputSize
}

// Implement getSize methods for REX, VEX, Opcode, ModRM, etc.
func (r *Rex) getSize() int {
	return 1 // REX prefix size
}

func (v *Vex) getSize() int {
	// Implement logic to calculate VEX size
	return 1 // Placeholder
}

func (o *Opcode) getSize() int {
	return 1 // Opcode size
}

func (m *Modrm) getSize() int {
	return 1 // ModRM size
}

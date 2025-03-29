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
	if e.REX != nil { // Check if REX is not nil
		if e.REX.Mandatory { // Check Mandatory only if REX is not nil
			size := e.REX.getSize()
			builder.WriteString(fmt.Sprintf(" rex:%d", size))
			outputSize += size
		}
	}


	// Calculate size based on VEX
	if e.VEX != nil { // Check if VEX is not nil
		size := e.VEX.getSize()
		builder.WriteString(fmt.Sprintf(" vex:%d", size))
		outputSize += size
	}

	// Calculate size based on Opcode
	// Check if Opcode is not the zero value before getting its size
	var zeroOpcode Opcode // Create a zero value Opcode for comparison
	if e.Opcode != zeroOpcode {
		sizeOpcode := e.Opcode.getSize()
		builder.WriteString(fmt.Sprintf(" opcode:%d", sizeOpcode))
		outputSize += sizeOpcode
	} else {
		log.Printf("warn: Encoding has zero Opcode field: %+v", e) // Log warning if Opcode is zero
	}


	// Calculate size based on ModRM
	if e.ModRM != nil {
		size := e.ModRM.getSize()
		builder.WriteString(fmt.Sprintf(" modrm:%d", size))
		outputSize += size
	}

	// Calculate size based on Immediate
	if e.Immediate != nil {
		// Use the size defined by the encoding itself.
		// options.ImmSize might be used for selection logic elsewhere,
		// but the encoding's defined size determines the actual output size here.
		immSize := e.Immediate.Size
		builder.WriteString(fmt.Sprintf(" immediate:%d", immSize))
		outputSize += immSize
	}

	// Calculate size based on DataOffset
	if e.DataOffset != nil { // Check if DataOffset is not nil
		size := e.DataOffset.Size
		builder.WriteString(fmt.Sprintf(" data_offset:%d", size))
		outputSize += size
	}

	// Calculate size based on CodeOffset
	if e.CodeOffset != nil { // Check if CodeOffset is not nil
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

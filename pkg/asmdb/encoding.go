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
		immSize := e.Immediate.Size // Default to encoding definition
		if options != nil && options.ImmSize > 0 {
			// If options provide a specific ImmSize, use it instead.
			// This allows pass1 to calculate size based on the actual immediate value.
			immSize = options.ImmSize
		}
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
	// Calculate size based on the length of the Byte string (hex representation)
	if len(o.Byte)%2 != 0 {
		log.Printf("warn: Opcode byte string has odd length: %s", o.Byte)
		// Handle error or return default? Returning length/2 might be misleading.
		return 0 // Or perhaps return an error? For now, return 0 on odd length.
	}
	return len(o.Byte) / 2 // Each byte is represented by 2 hex characters
}

func (m *Modrm) getSize() int {
	return 1 // ModRM size
}

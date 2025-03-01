package asmdb

import (
	"log"
)

type OutputSizeOptions struct {
	ImmSize int // 即値サイズ
}

func (e *Encoding) GetOutputSize(options *OutputSizeOptions) int {
	outputSize := 0

	// Log the start of the calculation
	log.Printf("debug: [pass1] --- get_output_size ---")

	// Calculate size based on REX
	if e.REX != nil && e.REX.Mandatory {
		log.Printf("debug: [pass1] bytes rex %d", e.REX.getSize())
		outputSize += e.REX.getSize()
	}

	// Calculate size based on VEX
	if e.VEX != nil {
		log.Printf("debug: [pass1] bytes vex %d", e.VEX.getSize())
		outputSize += e.VEX.getSize()
	}

	// Calculate size based on Opcode
	log.Printf("debug: [pass1] bytes opcode %d", e.Opcode.getSize())
	outputSize += e.Opcode.getSize()

	// Calculate size based on ModRM
	if e.ModRM != nil {
		log.Printf("debug: [pass1] bytes modrm %d", e.ModRM.getSize())
		outputSize += e.ModRM.getSize()
	}

	// Calculate size based on Immediate
	if e.Immediate != nil {
		var immSize int = e.Immediate.Size
		if options != nil && options.ImmSize > 0 {
			immSize = options.ImmSize
		}
		log.Printf("debug: [pass1] bytes immediate %d", immSize)
		outputSize += immSize
	}

	// Calculate size based on DataOffset
	if e.DataOffset != nil {
		log.Printf("debug: [pass1] bytes data offset %d", e.DataOffset.Size)
		outputSize += e.DataOffset.Size
	}

	// Calculate size based on CodeOffset
	if e.CodeOffset != nil {
		log.Printf("debug: [pass1] bytes code offset %d", e.CodeOffset.Size)
		outputSize += e.CodeOffset.Size
	}

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

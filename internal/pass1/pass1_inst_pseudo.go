package pass1

import (
	"fmt"
	"log"

	"github.com/HobbyOSs/gosk/internal/ast" // Add ast import
)

// processALIGNB handles the ALIGNB pseudo-instruction.
func processALIGNB(env *Pass1, operands []ast.Exp) {
	if len(operands) != 1 {
		log.Printf("Error: ALIGNB directive requires exactly one numeric operand.")
		return
	}

	numExp, ok := operands[0].(*ast.NumberExp)
	if !ok {
		log.Printf("Error: ALIGNB directive requires a numeric operand, got %T.", operands[0])
		return
	}

	unit := int32(numExp.Value) // Value is int64, cast to int32 for unit
	if unit <= 0 {
		log.Printf("Error: ALIGNB unit must be positive, got %d.", unit)
		return
	}

	var padding int32 = 0
	if env.LOC%unit != 0 { // 現在のLOCが境界に揃っていない場合のみ計算
		nearestSize := (env.LOC + unit - 1) / unit // Calculate next multiple
		padding = nearestSize*unit - env.LOC
	}
	env.LOC += padding
	// ALIGNB は ocode を発行しない (LOC調整のみ)
}

// processDB handles the DB pseudo-instruction.
// It processes evaluated operands (numbers, strings, labels).
func processDB(env *Pass1, operands []ast.Exp) {
	var loc int32 = 0
	var ocodes []int32 // Use int32 for consistency with emitCommand, though DB is bytes

	for _, operand := range operands {
		switch op := operand.(type) {
		case *ast.NumberExp:
			val := op.Value // Value is int64
			if val < 0 || val > 255 {
				log.Printf("Warning: Value %d out of range for DB, truncating.", val)
			}
			loc += 1
			ocodes = append(ocodes, int32(val&0xFF)) // Append lower byte
		case *ast.ImmExp: // Handles strings, identifiers (labels), etc. via Factor
			switch factor := op.Factor.(type) {
			case *ast.StringFactor:
				strVal := factor.Value // Assuming StringFactor has a Value field
				loc += int32(len(strVal))
				for _, char := range []byte(strVal) {
					ocodes = append(ocodes, int32(char))
				}
			case *ast.IdentFactor:
				ident := factor.Value                        // Assuming IdentFactor has a Value field
				if labelLOC, ok := env.SymTable[ident]; ok { // Check if it's a known label
					loc += 1                               // DB stores 1 byte for label address (lower byte)
					ocodes = append(ocodes, labelLOC&0xFF) // Append lower byte of label address
				} else {
					// Handle unresolved identifier - this should ideally be an error
					log.Printf("Error: Unresolved identifier '%s' in DB directive.", ident)
					// Decide how to handle this: skip, add placeholder, or halt?
					// For now, let's skip adding to ocodes and loc.
				}
			// case *ast.NumberFactor, *ast.HexFactor, *ast.CharFactor:
			// These should have been evaluated to *ast.NumberExp by TraverseAST/Eval
			// If they appear here, it might indicate an issue in evaluation logic.
			// log.Printf("Warning: Unexpected Factor type %T within ImmExp in DB.", factor)
			// Handle defensively if needed, e.g., try to get value and add as byte.
			default:
				log.Printf("Error: Unsupported Factor type %T within ImmExp in DB directive.", factor)
				// Decide how to handle this: skip, add placeholder, or halt?
				// For now, let's skip adding to ocodes and loc.
			}
		default:
			// This case should ideally not happen if TraverseAST evaluates correctly
			log.Printf("Error: Unsupported operand type %T in DB directive.", operand)
		}
	}

	env.LOC += loc
	emitCommand(env, "DB", ocodes) // TODO: Adjust emitCommand if needed for []int32
}

// processDW handles the DW pseudo-instruction.
func processDW(env *Pass1, operands []ast.Exp) {
	var loc int32 = 0
	var ocodes []int32

	for _, operand := range operands {
		switch op := operand.(type) {
		case *ast.NumberExp:
			val := op.Value                  // Value is int64
			if val < -32768 || val > 65535 { // Check signed/unsigned 16-bit range
				log.Printf("Warning: Value %d out of range for DW, truncating.", val)
			}
			loc += 2
			ocodes = append(ocodes, int32(val&0xFFFF)) // Append lower 16 bits
		case *ast.ImmExp:
			switch factor := op.Factor.(type) {
			case *ast.IdentFactor:
				ident := factor.Value
				if labelLOC, ok := env.SymTable[ident]; ok { // Check if it's a known label
					loc += 2                          // DW stores 2 bytes for label address
					ocodes = append(ocodes, labelLOC) // Append label address (int32)
				} else {
					log.Printf("Error: Unresolved identifier '%s' in DW directive.", ident)
					// Skip adding to ocodes and loc.
				}
			case *ast.StringFactor:
				log.Printf("Error: String literal '%s' is not a valid operand for DW directive.", factor.Value)
				// Skip adding to ocodes and loc.
			default:
				log.Printf("Error: Unsupported Factor type %T within ImmExp in DW directive.", factor)
				// Skip adding to ocodes and loc.
			}
		default:
			log.Printf("Error: Unsupported operand type %T in DW directive.", operand)
		}
	}

	env.LOC += loc
	emitCommand(env, "DW", ocodes)
}

// processDD handles the DD pseudo-instruction.
func processDD(env *Pass1, operands []ast.Exp) {
	var loc int32 = 0
	var ocodes []int32

	for _, operand := range operands {
		switch op := operand.(type) {
		case *ast.NumberExp:
			val := op.Value // Value is int64
			// Check 32-bit range (optional, as int32 conversion handles it)
			// if val < -2147483648 || val > 4294967295 {
			// 	log.Printf("Warning: Value %d out of 32-bit range for DD, truncating.", val)
			// }
			loc += 4
			ocodes = append(ocodes, int32(val)) // Append value as int32
		case *ast.ImmExp:
			switch factor := op.Factor.(type) {
			case *ast.IdentFactor:
				ident := factor.Value
				if labelLOC, ok := env.SymTable[ident]; ok { // Check if it's a known label
					loc += 4                          // DD stores 4 bytes for label address
					ocodes = append(ocodes, labelLOC) // Append label address (int32)
				} else {
					log.Printf("Error: Unresolved identifier '%s' in DD directive.", ident)
					// Skip adding to ocodes and loc.
				}
			case *ast.StringFactor:
				log.Printf("Error: String literal '%s' is not a valid operand for DD directive.", factor.Value)
				// Skip adding to ocodes and loc.
			default:
				log.Printf("Error: Unsupported Factor type %T within ImmExp in DD directive.", factor)
				// Skip adding to ocodes and loc.
			}
		default:
			log.Printf("Error: Unsupported operand type %T in DD directive.", operand)
		}
	}

	env.LOC += loc
	emitCommand(env, "DD", ocodes)
}

// processORG handles the ORG pseudo-instruction.
func processORG(env *Pass1, operands []ast.Exp) {
	if len(operands) != 1 {
		log.Printf("Error: ORG directive requires exactly one numeric operand.")
		return
	}

	numExp, ok := operands[0].(*ast.NumberExp)
	if !ok {
		log.Printf("Error: ORG directive requires a numeric operand, got %T.", operands[0])
		return
	}

	size := numExp.Value // Value is int64
	env.LOC = int32(size)
	env.DollarPosition += uint32(size) // エントリーポイントのアドレスを加算
	// ORG does not emit ocode
}

// processRESB handles the RESB pseudo-instruction.
// It expects a single evaluated NumberExp operand representing the size.
func processRESB(env *Pass1, operands []ast.Exp) {
	if len(operands) != 1 {
		log.Printf("Error: RESB directive requires exactly one numeric operand.")
		return
	}

	numExp, ok := operands[0].(*ast.NumberExp)
	if !ok {
		// The expression might involve '$' which needs special handling during Eval.
		// If it didn't evaluate to a number, log an error.
		log.Printf("Error: RESB directive requires a numeric operand (potentially involving '$' evaluated earlier), got %T.", operands[0])
		return
	}

	size := numExp.Value // Value is int64
	if size < 0 {
		log.Printf("Error: RESB size cannot be negative (%d).", size)
		return
	}

	env.LOC += int32(size)
	// Emit the RESB command with the calculated size.
	// The original source might have been 'X-$', but we emit the final size.
	env.Client.Emit(fmt.Sprintf("RESB %d", size))
}

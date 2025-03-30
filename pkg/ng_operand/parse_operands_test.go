package ng_operand

import (
	"testing"

	"github.com/HobbyOSs/gosk/pkg/cpu"
)

// Renamed and adapted TestOperandImpl_FromString to test ParseOperands directly
// This test now verifies if ParseOperands correctly parses the string into expected structures.
func TestParseOperands_FromString(t *testing.T) {
	text := "EAX, EBX"
	// expected := "EAX, EBX" // Original string check is less useful now

	parsedOperands, err := ParseOperands(text, cpu.MODE_16BIT, false) // Assuming default flags (forceRelAsImm=false)
	if err != nil {
		t.Fatalf("ParseOperands failed for %q: %v", text, err)
	}

	// Verify the number of parsed operands
	if len(parsedOperands) != 2 {
		t.Fatalf("Expected 2 operands, got %d", len(parsedOperands))
	}

	// Verify the details of each parsed operand (example)
	if parsedOperands[0] == nil || parsedOperands[0].Type != CodeR32 || parsedOperands[0].Register != "EAX" {
		t.Errorf("First operand mismatch: got %+v, want R32 EAX", parsedOperands[0])
	}
	if parsedOperands[1] == nil || parsedOperands[1].Type != CodeR32 || parsedOperands[1].Register != "EBX" {
		t.Errorf("Second operand mismatch: got %+v, want R32 EBX", parsedOperands[1])
	}
	// Add more detailed checks as needed based on ParsedOperandPeg structure
}

// TODO: Add more test cases for ParseOperands covering various operand combinations and edge cases.

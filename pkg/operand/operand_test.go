package operand_test

import (
	"testing"

	"github.com/HobbyOSs/gosk/pkg/operand"
)

func TestOperandBuilder(t *testing.T) {
	op1 := operand.OperandBuilder{}.Reg("EAX")
	if op1.Serialize() != "EAX" {
		t.Errorf("Expected EAX, got %s", op1.Serialize())
	}

	op2 := operand.OperandBuilder{}.Imm(42)
	if op2.Serialize() != "#42" {
		t.Errorf("Expected #42, got %s", op2.Serialize())
	}

	op3 := operand.OperandBuilder{}.Mem("EBX", "ECX", 2, 8)
	if op3.Serialize() != "[EBX ECX*2 +8]" {
		t.Errorf("Expected [EBX ECX*2 +8], got %s", op3.Serialize())
	}
}

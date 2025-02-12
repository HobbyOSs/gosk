package variantstack

import (
	"testing"
)

func TestVariantStack_PushPop(t *testing.T) {
	stack := VariantStack{}

	// Test pushing different types
	stack.Push(42)
	stack.Push("Hello")
	stack.Push(3.14)
	stack.Push(true)

	if stack.Size() != 4 {
		t.Errorf("Expected stack size 4, got %d", stack.Size())
	}

	// Test popping and type assertion
	if val, ok := stack.Pop(); !ok || val != true {
		t.Errorf("Expected true, got %v", val)
	}

	if val, ok := stack.Pop(); !ok || val != 3.14 {
		t.Errorf("Expected 3.14, got %v", val)
	}

	if val, ok := stack.Pop(); !ok || val != "Hello" {
		t.Errorf("Expected 'Hello', got %v", val)
	}

	if val, ok := stack.Pop(); !ok || val != 42 {
		t.Errorf("Expected 42, got %v", val)
	}

	if stack.Size() != 0 {
		t.Errorf("Expected stack size 0, got %d", stack.Size())
	}
}

func TestVariantStack_Clear(t *testing.T) {
	stack := VariantStack{}
	stack.Push(1)
	stack.Push(2)
	stack.Clear()

	if stack.Size() != 0 {
		t.Errorf("Expected stack size 0 after clear, got %d", stack.Size())
	}
}

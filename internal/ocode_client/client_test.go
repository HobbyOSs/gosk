package client

import (
	"testing"
)

func TestEmit(t *testing.T) {
	client := NewCodegenClient()

	err := client.Emit("L 5")
	if err != nil {
		t.Errorf("Emit failed: %v", err)
	}
}

func TestExec(t *testing.T) {
	client := NewCodegenClient()

	err := client.EmitAll("L 5\nL EAX\nADD")
	if err != nil {
		t.Fatalf("EmitAll failed: %v", err)
	}

	_, err = client.Exec()
	if err != nil {
		t.Errorf("Exec failed: %v", err)
	}
}

func TestEmitAll(t *testing.T) {
	client := NewCodegenClient()

	err := client.EmitAll("L 5\nL EAX\nADD")
	if err != nil {
		t.Errorf("EmitAll failed: %v", err)
	}
}
func TestExecWithDBDWDD(t *testing.T) {
	client := NewCodegenClient()

	err := client.Emit("DB 2,224")
	if err != nil {
		t.Fatalf("Emit failed: %v", err)
	}

	err = client.Emit("DW 4660")
	if err != nil {
		t.Fatalf("Emit failed: %v", err)
	}

	err = client.Emit("DD 305419896")
	if err != nil {
		t.Fatalf("Emit failed: %v", err)
	}

	// Execで結果を取得
	// 期待される結果を定義
	expected := []byte{0x02, 0xe0, 0x34, 0x12, 0x78, 0x56, 0x34, 0x12}
	result, err := client.Exec()
	if err != nil {
		t.Fatalf("Exec failed: %v", err)
	}

	// 結果を検証
	if !equal(result, expected) {
		t.Errorf("got %v, expected %v", result, expected)
	}
}

// equal関数はx86gen_test.goからコピーするか、共通のユーティリティとして定義してください
func equal(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

package client

import (
	"testing"

	"github.com/HobbyOSs/gosk/internal/codegen"
	"github.com/HobbyOSs/gosk/internal/pass1"
	"github.com/stretchr/testify/assert"
)

func TestNewCodegenClient(t *testing.T) {
	client := NewCodegenClient(&codegen.CodeGenContext{}, nil)
	assert.NotNil(t, client)
}

func TestEmit(t *testing.T) {
	client := NewCodegenClient(&codegen.CodeGenContext{}, nil)
	err := client.Emit("MOV AX, 1")
	assert.NoError(t, err)
}

func TestEmitAll(t *testing.T) {
	client := NewCodegenClient(&codegen.CodeGenContext{}, nil)
	err := client.EmitAll("MOV AX, 1\nMOV BX, 2")
	assert.NoError(t, err)
}

func TestExec(t *testing.T) {
	pass1 := &pass1.Pass1{
		SymTable: make(map[string]int32),
	}
	client := NewCodegenClient(&codegen.CodeGenContext{}, pass1)
	err := client.EmitAll("MOV AX, 1\nMOV BX, 2")
	assert.NoError(t, err)

	machineCode, err := client.Exec()
	assert.NoError(t, err)
	assert.NotNil(t, machineCode)
}

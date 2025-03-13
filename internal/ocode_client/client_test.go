package client

import (
	"testing"

	"github.com/HobbyOSs/gosk/internal/ast"
	"github.com/stretchr/testify/assert"
)

func TestNewCodegenClient(t *testing.T) {
	client := NewCodegenClient(ast.MODE_32BIT)
	assert.NotNil(t, client)
}

func TestEmit(t *testing.T) {
	client := NewCodegenClient(ast.MODE_32BIT)
	err := client.Emit("MOV AX, 1")
	assert.NoError(t, err)
}

func TestEmitAll(t *testing.T) {
	client := NewCodegenClient(ast.MODE_32BIT)
	err := client.EmitAll("MOV AX, 1\nMOV BX, 2")
	assert.NoError(t, err)
}

func TestExec(t *testing.T) {
	client := NewCodegenClient(ast.MODE_32BIT)
	err := client.EmitAll("MOV AX, 1\nMOV BX, 2")
	assert.NoError(t, err)

	machineCode, err := client.Exec()
	assert.NoError(t, err)
	assert.NotNil(t, machineCode)
}

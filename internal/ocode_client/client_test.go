package client

import (
	"testing"

	"github.com/HobbyOSs/gosk/internal/codegen"
	"github.com/stretchr/testify/assert"
)

func TestNewCodegenClient(t *testing.T) {
	// ctx == nil の場合、エラーが発生することを確認
	clientNilCtx, errNilCtx := NewCodegenClient(nil)
	assert.Nil(t, clientNilCtx)
	assert.Error(t, errNilCtx)

	// ctx != nil の場合、エラーが発生しないことを確認
	clientValidCtx, errValidCtx := NewCodegenClient(&codegen.CodeGenContext{})
	assert.NotNil(t, clientValidCtx)
	assert.NoError(t, errValidCtx)
}

func TestEmit(t *testing.T) {
	client, err := NewCodegenClient(&codegen.CodeGenContext{})
	assert.NoError(t, err)
	errEmit := client.Emit("MOV AX, 1")
	assert.NoError(t, errEmit)
}

func TestEmitAll(t *testing.T) {
	client, err := NewCodegenClient(&codegen.CodeGenContext{})
	assert.NoError(t, err)
	errEmitAll := client.EmitAll("MOV AX, 1\nMOV BX, 2")
	assert.NoError(t, errEmitAll)
}

func TestExec(t *testing.T) {
	client, err := NewCodegenClient(&codegen.CodeGenContext{})
	assert.NoError(t, err)
	errEmitAll := client.EmitAll("MOV AX, 1\nMOV BX, 2")
	assert.NoError(t, errEmitAll)

	machineCode, errExec := client.Exec()
	assert.NoError(t, errExec)
	assert.NotNil(t, machineCode)
}

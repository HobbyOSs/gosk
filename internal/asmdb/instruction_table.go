package asmdb

import (
	_ "embed"
	"log"

	"github.com/HobbyOSs/gosk/internal/ast"
	"github.com/goccy/go-yaml"
)

//go:embed x86reference.yml
var x86referenceYml string
var x86Ref x86Reference

type x86Reference struct {
	Instructions []Instruction
}

type Operand struct {
	Destination struct {
		Str     string `yaml:"operand_s"`
		Address string `yaml:"a,omitempty"`
		Type    string `yaml:"t,omitempty"`
	} `yaml:"dst,omitempty"`
	Source struct {
		Str     string `yaml:"operand_s"`
		Address string `yaml:"a,omitempty"`
		Type    string `yaml:"t,omitempty"`
	} `yaml:"src,omitempty"`
}

type Instruction struct {
	Mnemonic    string   `yaml:"mnem"`
	Opcode      string   `yaml:"opcd"`
	Operand1    *Operand `yaml:"op1,omitempty"`
	Operand2    *Operand `yaml:"op2,omitempty"`
	Proc        string   `yaml:"proc"`
	Description string   `yaml:"desc,omitempty"`
}

func (i Instruction) IsSupported(targetCPU ast.SupCPU) bool {
	supportCPURange, ok := ast.NewSupCPUByCode(i.Proc)
	if !ok {
		return false
	}
	return supportCPURange.IsSupported(targetCPU)
}

func init() {
	if err := yaml.Unmarshal([]byte(x86referenceYml), &x86Ref.Instructions); err != nil {
		log.Fatalf("error: %v", err)
	}
}

func X86Reference() *x86Reference {
	return &x86Ref
}

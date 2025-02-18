package asmdb

import (
	_ "embed"
	"bytes"
	"compress/gzip"
	"io"
	"log"

	"github.com/tidwall/gjson"
)

	//go:embed json-x86-64/x86_64.json.gz
	var compressedJSON []byte

type Instruction struct {
	Mnemonic    string   `yaml:"mnem"`
	Opcode      string   `yaml:"opcd"`
	Proc        string   `yaml:"proc"`
	Description string   `yaml:"desc,omitempty"`
}

var instructions []Instruction

func decompressGzip(data []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	return io.ReadAll(reader)
}

func init() {
	data, err := decompressGzip(compressedJSON)
	if err != nil {
		log.Fatalf("Failed to decompress JSON: %v", err)
	}

	gjson.ParseBytes(data).ForEach(func(key, value gjson.Result) bool {
		instructions = append(instructions, Instruction{
			Mnemonic:    value.Get("mnem").String(),
			Opcode:      value.Get("opcd").String(),
			Proc:        value.Get("proc").String(),
			Description: value.Get("desc").String(),
		})
		return true
	})
}

func X86Instructions() []Instruction {
	return instructions
}

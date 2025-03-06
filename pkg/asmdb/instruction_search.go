package asmdb

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/HobbyOSs/gosk/pkg/operand"
	"github.com/tidwall/gjson"
)

var jsonData []byte

func init() {
	var err error
	jsonData, err = decompressGzip(compressedJSON)
	if err != nil {
		log.Fatalf("Failed to decompress JSON: %v", err)
	}
}

type InstructionDB struct {
	jsonData []byte
}

func NewInstructionDB() *InstructionDB {
	return &InstructionDB{jsonData: jsonData}
}

func (db *InstructionDB) FindInstruction(opcode string) (*Instruction, bool) {
	path := `instructions.` + opcode
	result := gjson.GetBytes(db.jsonData, path)
	if !result.Exists() {
		return nil, false
	}

	var instr Instruction
	if err := json.Unmarshal([]byte(result.Raw), &instr); err != nil {
		return nil, false
	}
	return &instr, true
}

func (db *InstructionDB) FindEncoding(opcode string, operands operand.Operands) (*Encoding, error) {
	instr, found := db.FindInstruction(opcode)
	if !found {
		return nil, errors.New("instruction not found")
	}

	var (
		minEncoding *Encoding
		minSize     = -1
	)

	for i := range instr.Forms {
		form := &instr.Forms[i]
		if !matchOperands(form.Operands, operands) {
			continue
		}

		// 各エンコーディングのサイズを計算し、最小のものを見つける
		for j := range form.Encodings {
			e := &form.Encodings[j]
			options := &OutputSizeOptions{
				ImmSize: operands.DetectImmediateSize(),
			}
			size := e.GetOutputSize(options)

			// より小さいサイズのエンコーディングを見つけた場合に更新
			if minEncoding == nil || size < minSize {
				minEncoding = e
				minSize = size
			}
		}
	}

	if minEncoding == nil {
		return nil, errors.New("no matching encoding found")
	}

	return minEncoding, nil
}

// GetPrefixSize はプレフィックスバイトのサイズを計算します
func (db *InstructionDB) GetPrefixSize(operands operand.Operands) int {
	size := 0

	// operand size prefix (0x66)のみ必要
	if operands.Require66h() {
		size += 1
	}

	return size
}

func (db *InstructionDB) FindMinOutputSize(opcode string, operands operand.Operands) (int, error) {
	encoding, err := db.FindEncoding(opcode, operands)
	if err != nil {
		return 0, err
	}

	options := &OutputSizeOptions{
		ImmSize: operands.DetectImmediateSize(),
	}
	size := encoding.GetOutputSize(options)

	// プレフィックスとオフセットのサイズを加算
	return size + db.GetPrefixSize(operands) + operands.CalcOffsetByteSize(), nil
}

func matchOperands(formOperands *[]Operand, queryOperands operand.Operands) bool {
	if formOperands == nil || len(*formOperands) != len(queryOperands.OperandTypes()) {
		return false
	}
	for i, operand := range *formOperands {
		if operand.Type != queryOperands.OperandTypes()[i].String() {
			return false
		}
	}
	return true
}

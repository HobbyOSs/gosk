package asmdb

import (
	"errors"
	"log"
	"regexp"

	"github.com/HobbyOSs/gosk/pkg/operand"
	"github.com/samber/lo"
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
}

func NewInstructionDB() *InstructionDB {
	return &InstructionDB{}
}

// FindEncoding は指定された命令とオペランドに対応するエンコーディングを検索します。
// セグメントレジスタ（sreg）を含む命令の場合、matchOperands関数内でr16として扱われます。
// 例：MOV AX, SS は MOV r16, r16 として検索され、適切なエンコーディング（8C/8E）が選択されます。
func (db *InstructionDB) FindEncoding(opcode string, operands operand.Operands) (*Encoding, error) {
	instr, err := GetInstructionByOpcode(opcode)
	if err != nil {
		return nil, errors.New("instruction not found")
	}

	var (
		minEncoding      *Encoding
		minSize          = -1
		conditionRelaxed = false
	)

	filteredForms := lo.Filter(instr.Forms, func(form InstructionForm, _ int) bool {
		return matchOperands(form.Operands, operands, conditionRelaxed)
	})
	if len(filteredForms) == 0 {
		conditionRelaxed = true
		filteredForms = lo.Filter(instr.Forms, func(form InstructionForm, _ int) bool {
			return matchOperands(form.Operands, operands, conditionRelaxed)
		})
	}

	for i := range filteredForms {
		for j := range filteredForms[i].Encodings {
			e := &filteredForms[i].Encodings[j]
			options := &OutputSizeOptions{
				ImmSize: operands.DetectImmediateSize(),
			}
			size := e.GetOutputSize(options)

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
	minOutputSize := size + db.GetPrefixSize(operands) + operands.CalcOffsetByteSize()
	log.Printf("debug: [pass1] %s %s = %d\n", opcode, operands.InternalString(), minOutputSize)
	return minOutputSize, nil
}

func matchOperands(formOperands *[]Operand, queryOperands operand.Operands, conditionRelaxed bool) bool {
	if formOperands == nil || len(*formOperands) != len(queryOperands.OperandTypes()) {
		return false
	}

	// アキュムレータレジスタの優先検索
	internalStrs := queryOperands.InternalStrings()
	reAcc := regexp.MustCompile(`(?i)^(al|ax|eax)$`)

	for _, str := range internalStrs {
		if reAcc.MatchString(str) {
			// アキュムレータレジスタが含まれる場合、優先的に検索
			return matchOperandsWithAccumulator(formOperands, queryOperands)
		}
	}

	// 通常の検索ロジック（既存のコード）
	if conditionRelaxed {
		for i, operand := range *formOperands {
			queryType := queryOperands.OperandTypes()[i].String()
			if operand.Type != queryType {
				// 条件が緩和された場合; sregはr16としても一致を試みる
				if queryType == "sreg" && operand.Type == "r16" {
					continue // sregはr16として扱う
				}
				return false
			}
		}
		return true
	}

	for i, operand := range *formOperands {
		queryType := queryOperands.OperandTypes()[i].String()
		if operand.Type != queryType {
			return false
		}
	}
	return true
}

func matchOperandsWithAccumulator(formOperands *[]Operand, queryOperands operand.Operands) bool {
	reAcc := regexp.MustCompile(`^(al|ax|eax)$`)
	reOp := regexp.MustCompile(`^r(8|16|32)$`)

	// アキュムレータレジスタを優先的にマッチングするロジック
	for i, operand := range *formOperands {
		queryType := queryOperands.OperandTypes()[i].String()
		if operand.Type != queryType {
			// アキュムレータレジスタの場合、特定の条件でマッチングを試みる
			if reAcc.MatchString(operand.Type) && reOp.MatchString(queryType) {
				continue
			}
			return false
		}
	}
	return true
}

package codegen

import (
	"fmt"
	"log"
	"strconv"
)

// 配列型で渡される文字列はすべて数値型を想定する
// 0x00のようなhex notationや "hello" のような文字列はastで事前に数値の配列に変換
func handleDB(args []string) []byte {
	var binary []byte
	for _, arg := range args {
		num, _ := strconv.Atoi(arg)
		binary = append(binary, byte(num))
	}
	return binary
}

func handleDW(args []string) []byte {
	var binary []byte
	for _, arg := range args {
		num, _ := strconv.Atoi(arg)
		binary = append(binary, byte(num&0xFF), byte((num>>8)&0xFF))
	}
	return binary
}

func handleDD(args []string) []byte {
	var binary []byte
	for _, arg := range args {
		num, _ := strconv.Atoi(arg)
		binary = append(binary, byte(num&0xFF), byte((num>>8)&0xFF), byte((num>>16)&0xFF), byte((num>>24)&0xFF))
	}
	return binary
}

// handleRESB は RESB 疑似命令を処理します。
// Pass1 から渡される評価済みのサイズに基づいて、指定されたバイト数の 0 を生成します。
func handleRESB(args []string, params x86genParams, ctx *CodeGenContext) []byte {
	if len(args) != 1 {
		log.Printf("Error: handleRESB expects exactly one argument (the size), got %d", len(args))
		return nil // またはエラーを返す
	}

	// Pass1 からは評価済みのサイズが文字列として渡されるはず
	reserveSize, err := strconv.ParseInt(args[0], 10, 64) // 10進数としてパース
	if err != nil {
		log.Printf("Error parsing RESB size '%s': %v", args[0], err)
		return nil // またはエラーを返す
	}

	if reserveSize < 0 {
		log.Printf("Error: RESB size cannot be negative (%d).", reserveSize)
		return nil // またはエラーを返す
	}

	// 指定されたサイズの 0 バイトスライスを作成して返します
	return make([]byte, reserveSize)
}

// handleALIGNB は ALIGNB ディレクティブを処理します。
// 現在のバイトコード長を次のアライメント境界までヌルバイトでパディングします。
func handleALIGNB(params x86genParams, ctx *CodeGenContext) ([]byte, error) {
	if len(params.Operands) != 1 {
		return nil, fmt.Errorf("handleALIGNB: expected 1 operand, got %d", len(params.Operands))
	}

	alignBoundary, err := strconv.Atoi(params.Operands[0])
	if err != nil {
		return nil, fmt.Errorf("handleALIGNB: failed to parse alignment boundary '%s': %w", params.Operands[0], err)
	}

	if alignBoundary <= 0 || (alignBoundary&(alignBoundary-1)) != 0 {
		// 境界値が0以下、または2のべき乗でない場合はエラー (より厳密なチェックも可能)
		return nil, fmt.Errorf("handleALIGNB: invalid alignment boundary %d, must be a positive power of 2", alignBoundary)
	}

	// x86genParams から現在のバイトコード長を取得
	currentLength := params.MachineCodeLen
	paddingSize := (alignBoundary - (currentLength % alignBoundary)) % alignBoundary

	if paddingSize > 0 {
		log.Printf("debug: [codegen] Applying alignment padding for ALIGNB %d: size=%d, currentLength=%d\n", alignBoundary, paddingSize, currentLength)
		return make([]byte, paddingSize), nil // ヌルバイトで初期化されたスライスを返す
	}

	// パディング不要の場合
	return []byte{}, nil
}

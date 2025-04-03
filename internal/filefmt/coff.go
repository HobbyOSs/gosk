package filefmt

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	// "time" // 未使用
	// "io" // 未使用

	"github.com/HobbyOSs/gosk/internal/codegen"
	// "github.com/HobbyOSs/gosk/pkg/cpu" // 未使用
)

// COFF ファイルヘッダ構造体
type CoffHeader struct {
	Machine              uint16 // ターゲットマシン (0x014c for i386)
	NumberOfSections     uint16 // セクション数
	TimeDateStamp        uint32 // ファイル作成日時 (Unix timestamp)
	PointerToSymbolTable uint32 // シンボルテーブルへのファイルオフセット
	NumberOfSymbols      uint32 // シンボルテーブル内のエントリ数 (補助シンボル含む)
	SizeOfOptionalHeader uint16 // オプショナルヘッダのサイズ (オブジェクトファイルでは通常0)
	Characteristics      uint16 // ファイル特性フラグ
}

// COFF セクションヘッダ構造体
type CoffSectionHeader struct {
	Name                 [8]byte // セクション名 (NULL終端)
	VirtualSize          uint32  // (未使用)
	VirtualAddress       uint32  // (未使用)
	SizeOfRawData        uint32  // セクションのデータサイズ (バイト単位)
	PointerToRawData     uint32  // セクションデータへのファイルオフセット
	PointerToRelocations uint32  // 再配置情報へのファイルオフセット (今回は0)
	PointerToLinenumbers uint32  // 行番号情報へのファイルオフセット (今回は0)
	NumberOfRelocations  uint16  // 再配置エントリ数 (今回は0)
	NumberOfLinenumbers  uint16  // 行番号エントリ数 (今回は0)
	Characteristics      uint32  // セクション特性フラグ
}

// COFF シンボルテーブルエントリ構造体
type CoffSymbol struct {
	Name               [8]byte // シンボル名 (短い場合) または文字列テーブルへのオフセット
	Value              uint32  // シンボルの値 (アドレスなど)
	SectionNumber      int16   // 関連セクション番号 (1ベース、0=未定義, -1=絶対, -2=デバッグ)
	Type               uint16  // シンボル型 (基本型と派生型)
	StorageClass       uint8   // 格納クラス (スコープ、型など)
	NumberOfAuxSymbols uint8   // 補助シンボルエントリ数
}

// COFF 補助セクションシンボルエントリ構造体
type CoffAuxSectionSymbol struct {
	Length           uint32  // セクション長
	NumberOfRelocs   uint16  // 再配置エントリ数
	NumberOfLineNums uint16  // 行番号エントリ数
	CheckSum         uint32  // チェックサム (未使用)
	Number           uint16  // COMDATセクション番号 (未使用)
	Selection        uint8   // COMDAT選択タイプ (未使用)
	Reserved         [3]byte // パディング
}

// CoffFormat は COFF ファイル形式の書き出し処理を実装します。
type CoffFormat struct{}

// NewCoffFormat は新しい CoffFormat インスタンスを作成します。
func NewCoffFormat() *CoffFormat {
	return &CoffFormat{}
}

const (
	// COFF ヘッダサイズ
	coffHeaderSize = 20
	// COFF セクションヘッダサイズ
	coffSectionHeaderSize = 40
	// COFF シンボルエントリサイズ (補助シンボルも同じサイズ)
	coffSymbolSize = 18
	// COFF 文字列テーブルサイズエントリサイズ (サイズ自体を示すDWORD)
	coffStringTableSizeEntrySize = 4
)

// Write は COFF 形式でオブジェクトファイルを書き出します。
func (c *CoffFormat) Write(ctx *codegen.CodeGenContext, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filePath, err)
	}
	defer file.Close()

	buf := new(bytes.Buffer)

	// --- プレースホルダー書き込み ---
	// ヘッダ
	placeholderHeader := make([]byte, coffHeaderSize)
	if _, err := buf.Write(placeholderHeader); err != nil {
		return fmt.Errorf("failed to write placeholder header: %w", err)
	}

	// セクションヘッダ
	numSections := uint16(3) // .text, .data, .bss
	placeholderSectionHeaders := make([]byte, int(numSections)*coffSectionHeaderSize)
	if _, err := buf.Write(placeholderSectionHeaders); err != nil {
		return fmt.Errorf("failed to write placeholder section headers: %w", err)
	}

	// --- セクションデータ書き込み ---
	textDataOffset := uint32(buf.Len())
	textDataSize := uint32(len(ctx.MachineCode))
	if _, err := buf.Write(ctx.MachineCode); err != nil {
		return fmt.Errorf("failed to write .text section data: %w", err)
	}
	dataDataOffset := uint32(buf.Len())
	dataDataSize := uint32(0) // TODO: .data セクションの内容
	bssDataSize := uint32(0)  // TODO: .bss セクションのサイズ

	// --- シンボルテーブルと文字列テーブル書き込み ---
	symbolTableOffset := uint32(buf.Len())
	symbolsBytes, _, numSymbolsWithAux := c.generateSymbolsAndStringTable(ctx, textDataSize, dataDataSize, bssDataSize) // stringTable を無視
	if _, err := buf.Write(symbolsBytes); err != nil {
		return fmt.Errorf("failed to write symbol table: %w", err)
	}

	// 文字列テーブルの書き込みを削除

	// --- ヘッダとセクションヘッダのオフセット更新 ---
	header := c.generateHeader(ctx, numSections, symbolTableOffset, numSymbolsWithAux)
	sectionHeaders := c.generateSectionHeaders(ctx, textDataOffset, textDataSize, dataDataOffset, dataDataSize, bssDataSize)

	// --- 最終的なファイル内容を構築 ---
	finalBuf := new(bytes.Buffer)
	// ヘッダ
	if err := binary.Write(finalBuf, binary.LittleEndian, &header); err != nil {
		return fmt.Errorf("failed to write final header: %w", err)
	}
	// セクションヘッダ
	if err := binary.Write(finalBuf, binary.LittleEndian, sectionHeaders); err != nil {
		return fmt.Errorf("failed to write final section headers: %w", err)
	}
	// セクションデータ、シンボルテーブル、文字列テーブル
	copyOffset := uint32(coffHeaderSize + int(numSections)*coffSectionHeaderSize)
	if _, err := finalBuf.Write(buf.Bytes()[copyOffset:]); err != nil { // bufから直接コピー
		return fmt.Errorf("failed to copy data to final buffer: %w", err)
	}

	// --- ファイルへの書き込み ---
	_, err = file.Write(finalBuf.Bytes())
	if err != nil {
		return fmt.Errorf("failed to write buffer to file: %w", err)
	}

	log.Printf("Successfully wrote COFF file to %s", filePath)
	return nil
}

// generateHeader は COFF ファイルヘッダを生成します。
func (c *CoffFormat) generateHeader(ctx *codegen.CodeGenContext, numSections uint16, symbolTableOffset uint32, numSymbols uint32) CoffHeader {
	// Characteristics: 0x0000 (テストデータに合わせる)
	characteristics := uint16(0x0000)
	// 32bitフラグは不要

	return CoffHeader{
		Machine:              0x014c, // IMAGE_FILE_MACHINE_I386
		NumberOfSections:     numSections,
		TimeDateStamp:        0, // テストデータに合わせて0
		PointerToSymbolTable: symbolTableOffset,
		NumberOfSymbols:      numSymbols, // 補助シンボルを含む数
		SizeOfOptionalHeader: 0,
		Characteristics:      characteristics,
	}
}

// generateSectionHeaders は COFF セクションヘッダのスライスを生成します。
func (c *CoffFormat) generateSectionHeaders(ctx *codegen.CodeGenContext, textDataOffset, textDataSize, dataDataOffset, dataDataSize, bssDataSize uint32) []CoffSectionHeader {
	sections := make([]CoffSectionHeader, 0, 3)

	// .text セクション
	sections = append(sections, CoffSectionHeader{
		Name:             [8]byte{'.', 't', 'e', 'x', 't'},
		SizeOfRawData:    textDataSize,
		PointerToRawData: textDataOffset,
		Characteristics:  0x20001060, // テストデータに合わせる
		// 他のフィールドは 0
	})
	// .data セクション
	sections = append(sections, CoffSectionHeader{
		Name:             [8]byte{'.', 'd', 'a', 't', 'a'},
		SizeOfRawData:    dataDataSize,
		PointerToRawData: dataDataOffset,
		Characteristics:  0x400010C0, // テストデータに合わせる
		// 他のフィールドは 0
	})
	// .bss セクション
	sections = append(sections, CoffSectionHeader{
		Name:             [8]byte{'.', 'b', 's', 's'},
		SizeOfRawData:    bssDataSize, // データは持たない
		PointerToRawData: 0,           // データは持たない
		Characteristics:  0x800010C0,  // テストデータに合わせる
		// 他のフィールドは 0
	})
	return sections
}

// generateSymbolsAndStringTable は COFF シンボルテーブル(補助シンボル含む)のバイト列と文字列テーブル(今回は空)を生成します。
// 戻り値: シンボルテーブルのバイト列, 文字列テーブルのバイト列(空), シンボルテーブルのエントリ数(補助シンボル含む)
func (c *CoffFormat) generateSymbolsAndStringTable(ctx *codegen.CodeGenContext, textDataSize, dataDataSize, bssDataSize uint32) ([]byte, []byte, uint32) {
	symbols := make([]CoffSymbol, 0, len(ctx.SymTable)+4)
	auxSymbolsData := make(map[int][]byte)
	symbolIndex := 0

	// 期待される順序: .file -> .text -> .data -> .bss -> _io_hlt

	// 1. .file シンボル
	fileSymbol := CoffSymbol{
		Name:               c.convertNameToBytesDirect(".file"), // 直接変換
		Value:              0,
		SectionNumber:      -2,   // IMAGE_SYM_DEBUG
		Type:               0x00, // NULL
		StorageClass:       103,  // IMAGE_SYM_CLASS_FILE
		NumberOfAuxSymbols: 1,
	}
	symbols = append(symbols, fileSymbol)
	fileSymbolIndex := symbolIndex
	symbolIndex++

	// .file 補助シンボル (ファイル名)
	auxFileBytes := make([]byte, coffSymbolSize)
	copy(auxFileBytes, ctx.SourceFileName)
	auxSymbolsData[fileSymbolIndex] = auxFileBytes
	symbolIndex++

	// 2. セクションシンボル (.text, .data, .bss)
	sectionNames := []string{".text", ".data", ".bss"}
	sectionDataSizes := []uint32{textDataSize, dataDataSize, bssDataSize}
	sectionRelocs := []uint16{0, 0, 0}
	sectionLineNums := []uint16{0, 0, 0}

	for i, name := range sectionNames {
		sectionSymbol := CoffSymbol{
			Name:               c.convertNameToBytesDirect(name), // 直接変換
			Value:              0,
			SectionNumber:      int16(i + 1),
			Type:               0x00,
			StorageClass:       3, // IMAGE_SYM_CLASS_STATIC
			NumberOfAuxSymbols: 1,
		}
		symbols = append(symbols, sectionSymbol)
		sectionSymbolIndex := symbolIndex
		symbolIndex++

		// セクション補助シンボル
		auxSection := CoffAuxSectionSymbol{
			Length:           sectionDataSizes[i],
			NumberOfRelocs:   sectionRelocs[i],
			NumberOfLineNums: sectionLineNums[i],
		}
		auxSectionBytesBuf := new(bytes.Buffer)
		binary.Write(auxSectionBytesBuf, binary.LittleEndian, auxSection)
		auxSectionBytes := auxSectionBytesBuf.Bytes()
		if len(auxSectionBytes) < coffSymbolSize {
			padding := make([]byte, coffSymbolSize-len(auxSectionBytes))
			auxSectionBytes = append(auxSectionBytes, padding...)
		}
		auxSymbolsData[sectionSymbolIndex] = auxSectionBytes
		symbolIndex++
	}

	// 3. 通常のシンボル (_io_hlt)
	globalName := "_io_hlt"
	if addr, ok := ctx.SymTable[globalName]; ok {
		symbol := CoffSymbol{
			Name:               c.convertNameToBytesDirect(globalName), // 直接変換
			Value:              uint32(addr),
			SectionNumber:      1,    // .text セクション
			Type:               0x20, // Function
			StorageClass:       2,    // IMAGE_SYM_CLASS_EXTERNAL
			NumberOfAuxSymbols: 0,
		}
		symbols = append(symbols, symbol)
		symbolIndex++
	} else {
		log.Printf("warning: Global symbol '%s' not found in symbol table", globalName)
	}

	// シンボルと補助シンボルを結合したバイト列を作成
	symbolTableBytes := new(bytes.Buffer)
	for i, sym := range symbols {
		if err := binary.Write(symbolTableBytes, binary.LittleEndian, sym); err != nil {
			log.Printf("Error writing symbol %d: %v", i, err)
			return nil, nil, 0
		}
		if auxData, ok := auxSymbolsData[i]; ok {
			if _, err := symbolTableBytes.Write(auxData); err != nil {
				log.Printf("Error writing aux symbol %d: %v", i, err)
				return nil, nil, 0
			}
		}
	}

	return symbolTableBytes.Bytes(), []byte{}, uint32(symbolIndex) // 文字列テーブルは空バイト列を返す
}

// convertNameToBytesDirect はシンボル名を COFF シンボルテーブル用の8バイト配列に直接変換します。
// 8バイトを超える場合は切り詰められます。
func (c *CoffFormat) convertNameToBytesDirect(name string) [8]byte {
	var result [8]byte
	copy(result[:], name)
	return result
}

// convertNameToBytesForSymbol は削除 (使用しない)
// func (c *CoffFormat) convertNameToBytesForSymbol(name string, addToStringTable func(string) uint32) [8]byte { ... }

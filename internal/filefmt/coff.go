package filefmt

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"

	"github.com/HobbyOSs/gosk/internal/codegen"
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

// SymbolEntry はメインシンボルと対応する補助シンボルデータを保持します。
type SymbolEntry struct {
	Main CoffSymbol
	Aux  []byte // 補助シンボルがない場合は nil
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

	buf := new(bytes.Buffer) // ファイル全体のバッファ

	// --- プレースホルダー書き込み (後で上書き) ---
	// COFFヘッダ
	placeholderHeader := make([]byte, coffHeaderSize)
	if _, err := buf.Write(placeholderHeader); err != nil {
		return fmt.Errorf("failed to write placeholder header: %w", err)
	}

	// COFFセクションヘッダ (3セクション分)
	numSections := uint16(3) // .text, .data, .bss
	placeholderSectionHeaders := make([]byte, int(numSections)*coffSectionHeaderSize)
	if _, err := buf.Write(placeholderSectionHeaders); err != nil {
		return fmt.Errorf("failed to write placeholder section headers: %w", err)
	}

	// --- セクションデータ書き込み ---
	// .text セクション
	textDataOffset := uint32(buf.Len()) // 現在のオフセットを記録
	textDataSize := uint32(len(ctx.MachineCode))
	if _, err := buf.Write(ctx.MachineCode); err != nil {
		return fmt.Errorf("failed to write .text section data: %w", err)
	}
	// .data セクション (今回は空)
	dataDataOffset := uint32(buf.Len())
	dataDataSize := uint32(0) // TODO: .data セクションの内容を実装
	// .bss セクション (データなし)
	bssDataSize := uint32(0) // TODO: .bss セクションのサイズを計算

	// --- シンボルテーブルと文字列テーブル生成 ---
	allSymbolEntries, stringTableBytes := c.generateSymbolEntries(ctx, textDataSize, dataDataSize, bssDataSize)

	// --- シンボルテーブル書き込み ---
	symbolTableOffset := uint32(buf.Len()) // シンボルテーブルの開始オフセットを記録
	symbolTableResultBytes := new(bytes.Buffer)
	numSymbolsCounted := uint32(0) // シンボルの総数 (補助シンボル含む) - 正確にカウントする
	for _, entry := range allSymbolEntries {
		if err := binary.Write(symbolTableResultBytes, binary.LittleEndian, entry.Main); err != nil {
			// エラーログは残す
			symbolName := string(entry.Main.Name[:bytes.IndexByte(entry.Main.Name[:], 0)]) // エラー表示用に名前を取得
			log.Printf("Error writing main symbol %s: %v", symbolName, err)
			return fmt.Errorf("failed to write main symbol %s: %w", symbolName, err)
		}
		numSymbolsCounted++ // メインシンボルをカウント
		if entry.Aux != nil {
			if len(entry.Aux) != coffSymbolSize {
				// エラーログは残す
				symbolName := string(entry.Main.Name[:bytes.IndexByte(entry.Main.Name[:], 0)]) // エラー表示用に名前を取得
				log.Printf("Error: Aux symbol for %s has incorrect size %d, expected %d", symbolName, len(entry.Aux), coffSymbolSize)
				return fmt.Errorf("aux symbol for %s has incorrect size %d", symbolName, len(entry.Aux))
			}
			if _, err := symbolTableResultBytes.Write(entry.Aux); err != nil {
				// エラーログは残す
				symbolName := string(entry.Main.Name[:bytes.IndexByte(entry.Main.Name[:], 0)]) // エラー表示用に名前を取得
				log.Printf("Error writing aux symbol for %s: %v", symbolName, err)
				return fmt.Errorf("failed to write aux symbol for %s: %w", symbolName, err)
			}
			// 補助シンボル自体もシンボルテーブルのエントリとしてカウントする
			numSymbolsCounted += uint32(entry.Main.NumberOfAuxSymbols)
		}
	}
	// デバッグログと関連コードを完全に削除

	if _, err := buf.Write(symbolTableResultBytes.Bytes()); err != nil {
		return fmt.Errorf("failed to write symbol table: %w", err)
	}

	// --- 文字列テーブル書き込み ---
	// 文字列テーブルのサイズを計算 (内容 + サイズフィールド4バイト)
	stringTableTotalSize := uint32(len(stringTableBytes) + coffStringTableSizeEntrySize)
	stringTableSizeBytes := make([]byte, coffStringTableSizeEntrySize)
	binary.LittleEndian.PutUint32(stringTableSizeBytes, stringTableTotalSize)

	// サイズフィールドを常に書き込む (nask の出力に合わせる)
	if _, err := buf.Write(stringTableSizeBytes); err != nil {
		return fmt.Errorf("failed to write string table size: %w", err)
	}
	// 内容が存在する場合のみ内容を書き込む
	if len(stringTableBytes) > 0 {
		if _, err := buf.Write(stringTableBytes); err != nil {
			return fmt.Errorf("failed to write string table content: %w", err)
		}
	}
	// 不要なコメント削除

	// --- ヘッダとセクションヘッダの生成 (オフセット情報を含む) ---
	header := c.generateHeader(ctx, numSections, symbolTableOffset, numSymbolsCounted)                                                          // Use counted symbols
	sectionHeaders := c.generateSectionHeaders(ctx, textDataOffset, textDataSize, dataDataOffset, dataDataSize, bssDataSize, symbolTableOffset) // Pass symbolTableOffset

	// --- 最終的なファイル内容を構築 (プレースホルダーを上書き) ---
	finalBytes := buf.Bytes() // 現在のバッファ内容を取得

	// ヘッダを上書き
	headerBuf := new(bytes.Buffer)
	if err := binary.Write(headerBuf, binary.LittleEndian, &header); err != nil {
		return fmt.Errorf("failed to serialize final header: %w", err)
	}
	copy(finalBytes[0:coffHeaderSize], headerBuf.Bytes())

	// セクションヘッダを上書き
	sectionHeadersBuf := new(bytes.Buffer)
	if err := binary.Write(sectionHeadersBuf, binary.LittleEndian, sectionHeaders); err != nil {
		return fmt.Errorf("failed to serialize final section headers: %w", err)
	}
	copy(finalBytes[coffHeaderSize:coffHeaderSize+int(numSections)*coffSectionHeaderSize], sectionHeadersBuf.Bytes())

	// --- ファイルへの書き込み ---
	_, err = file.Write(finalBytes)
	if err != nil {
		return fmt.Errorf("failed to write final buffer to file: %w", err)
	}

	log.Printf("Successfully wrote COFF file to %s", filePath)
	return nil
}

// generateHeader は COFF ファイルヘッダを生成します。
func (c *CoffFormat) generateHeader(ctx *codegen.CodeGenContext, numSections uint16, symbolTableOffset uint32, numSymbols uint32) CoffHeader {
	// Characteristics: nask の出力に合わせる (0x0000)
	characteristics := uint16(0x0000)

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
func (c *CoffFormat) generateSectionHeaders(ctx *codegen.CodeGenContext, textDataOffset, textDataSize, dataDataOffset, dataDataSize, bssDataSize, symbolTableOffset uint32) []CoffSectionHeader {
	sections := make([]CoffSectionHeader, 0, 3)

	// .text セクション
	sections = append(sections, CoffSectionHeader{
		Name:                 [8]byte{'.', 't', 'e', 'x', 't'},
		SizeOfRawData:        textDataSize,
		PointerToRawData:     textDataOffset,
		PointerToRelocations: symbolTableOffset, // nask出力に合わせてシンボルテーブルオフセットを設定 (TestHarib01a: 0x99, TestHarib01f: 0xDB)
		PointerToLinenumbers: 0,                 // Line numbers not handled yet, set to 0
		NumberOfRelocations:  0,
		NumberOfLinenumbers:  0,
		Characteristics:      0x60100020, // nask出力に合わせる (IMAGE_SCN_CNT_CODE | IMAGE_SCN_MEM_EXECUTE | IMAGE_SCN_MEM_READ)
	})
	// .data セクション
	sections = append(sections, CoffSectionHeader{
		Name:                 [8]byte{'.', 'd', 'a', 't', 'a'},
		SizeOfRawData:        dataDataSize,
		PointerToRawData:     0, // TestHarib00j の期待値に合わせて 0 に変更
		PointerToRelocations: 0,
		PointerToLinenumbers: 0,
		NumberOfRelocations:  0,
		NumberOfLinenumbers:  0,
		Characteristics:      0xC0100040, // nask出力に合わせる (IMAGE_SCN_CNT_INITIALIZED_DATA | IMAGE_SCN_MEM_READ | IMAGE_SCN_MEM_WRITE)
	})
	// .bss セクション
	sections = append(sections, CoffSectionHeader{
		Name:                 [8]byte{'.', 'b', 's', 's'},
		SizeOfRawData:        bssDataSize, // データは持たない
		PointerToRawData:     0,           // データは持たない (PointerToRawData should be 0 for .bss)
		PointerToRelocations: 0,
		PointerToLinenumbers: 0,
		NumberOfRelocations:  0,
		NumberOfLinenumbers:  0,
		Characteristics:      0xC0100080, // nask出力に合わせる (IMAGE_SCN_CNT_UNINITIALIZED_DATA | IMAGE_SCN_MEM_READ | IMAGE_SCN_MEM_WRITE)
	})
	return sections
}

// generateSymbolEntries 内の PointerToSymbolTable の計算は不要になったため削除
// generateHeader の呼び出し側で symbolTableOffset を計算して渡す

// generateSymbolEntries は SymbolEntry のスライスと文字列テーブルのバイト列を生成します。
func (c *CoffFormat) generateSymbolEntries(ctx *codegen.CodeGenContext, textDataSize, dataDataSize, bssDataSize uint32) ([]SymbolEntry, []byte) {
	allEntries := make([]SymbolEntry, 0, len(ctx.SymTable)+4+len(ctx.GlobalSymbolList)) // 予測サイズを調整
	stringTable := new(bytes.Buffer)                                                    // 文字列テーブル用バッファ
	// デバッグログ削除
	stringTableOffsetMap := make(map[string]uint32) // 文字列テーブルの重複回避用マップ

	// 期待される順序: .file -> .text -> .data -> .bss -> グローバルシンボル

	// 1. .file シンボル
	fileName := ctx.SourceFileName
	fileSymbol := CoffSymbol{
		Name:               [8]byte{'.', 'f', 'i', 'l', 'e'},
		Value:              0,
		SectionNumber:      -2,   // IMAGE_SYM_DEBUG
		Type:               0x00, // NULL
		StorageClass:       103,  // IMAGE_SYM_CLASS_FILE
		NumberOfAuxSymbols: 1,
	}
	// .file 補助シンボル
	auxFileBytes := make([]byte, coffSymbolSize)
	copy(auxFileBytes, fileName)
	allEntries = append(allEntries, SymbolEntry{Main: fileSymbol, Aux: auxFileBytes})

	// 2. セクションシンボル (.text, .data, .bss)
	sectionNames := []string{".text", ".data", ".bss"}
	sectionDataSizes := []uint32{textDataSize, dataDataSize, bssDataSize}
	sectionRelocs := []uint16{0, 0, 0}
	sectionLineNums := []uint16{0, 0, 0}

	for i, name := range sectionNames {
		var sectionNameBytes [8]byte
		copy(sectionNameBytes[:], name)

		sectionSymbol := CoffSymbol{
			Name:               sectionNameBytes,
			Value:              0,
			SectionNumber:      int16(i + 1),
			Type:               0x00,
			StorageClass:       3, // IMAGE_SYM_CLASS_STATIC
			NumberOfAuxSymbols: 1,
		}
		// セクション補助シンボル
		auxSectionBytes := make([]byte, coffSymbolSize)
		binary.LittleEndian.PutUint32(auxSectionBytes[0:4], sectionDataSizes[i])
		binary.LittleEndian.PutUint16(auxSectionBytes[4:6], sectionRelocs[i])
		binary.LittleEndian.PutUint16(auxSectionBytes[6:8], sectionLineNums[i])

		allEntries = append(allEntries, SymbolEntry{Main: sectionSymbol, Aux: auxSectionBytes})
	}

	// 3. グローバルシンボル (GLOBAL 宣言されたもの)
	for _, globalName := range ctx.GlobalSymbolList {
		if addr, ok := ctx.SymTable[globalName]; ok {
			// シンボル名を8バイトに変換 (文字列テーブル使用)
			nameBytes := c.convertNameToBytes(globalName, stringTable, stringTableOffsetMap)

			// TODO: シンボルがどのセクションに属するかを決定するロジックが必要
			//       現状は仮に .text (SectionNumber=1) とする
			sectionNum := int16(1)
			// if _, dataOk := ctx.DataSymbols[globalName]; dataOk { // 仮のデータシンボルマップ - コメントアウト
			// 	sectionNum = 2 // .data
			// } else if _, bssOk := ctx.BssSymbols[globalName]; bssOk { // 仮のBSSシンボルマップ - コメントアウト
			// 	sectionNum = 3 // .bss
			// }
			// TODO: 外部シンボル (EXTERN) の処理も必要になる場合がある

			symbol := CoffSymbol{
				Name:               nameBytes,
				Value:              uint32(addr),
				SectionNumber:      sectionNum, // 仮: .text セクション
				Type:               0x00,       // Type NULL (0x00) に変更 (TestHarib00j の期待値に合わせる)
				StorageClass:       2,          // IMAGE_SYM_CLASS_EXTERNAL
				NumberOfAuxSymbols: 0,          // 補助シンボルなし (通常)
			}
			allEntries = append(allEntries, SymbolEntry{Main: symbol, Aux: nil}) // Aux は nil
		} else {
			log.Printf("warning: Global symbol '%s' declared but not found in symbol table", globalName)
			// 未解決の外部シンボルとして扱うか？ (StorageClass = IMAGE_SYM_CLASS_EXTERNAL, SectionNumber = 0)
			nameBytes := c.convertNameToBytes(globalName, stringTable, stringTableOffsetMap)
			symbol := CoffSymbol{
				Name:               nameBytes,
				Value:              0,
				SectionNumber:      0,    // IMAGE_SYM_UNDEFINED
				Type:               0x00, // Type NULL (0x00) に変更
				StorageClass:       2,    // IMAGE_SYM_CLASS_EXTERNAL
				NumberOfAuxSymbols: 0,
			}
			allEntries = append(allEntries, SymbolEntry{Main: symbol, Aux: nil})
		}
	}

	// 4. EXTERN シンボル (EXTERN 宣言されたもの)
	for _, externName := range ctx.ExternSymbolList {
		// シンボル名を8バイトに変換 (文字列テーブル使用)
		nameBytes := c.convertNameToBytes(externName, stringTable, stringTableOffsetMap)

		symbol := CoffSymbol{
			Name:               nameBytes,
			Value:              0,    // 外部シンボルの値は0
			SectionNumber:      0,    // IMAGE_SYM_UNDEFINED
			Type:               0x00, // Type NULL (0x00)
			StorageClass:       2,    // IMAGE_SYM_CLASS_EXTERNAL
			NumberOfAuxSymbols: 0,    // 補助シンボルなし
		}
		allEntries = append(allEntries, SymbolEntry{Main: symbol, Aux: nil}) // Aux は nil
	}

	return allEntries, stringTable.Bytes()
}

// convertNameToBytes はシンボル名を COFF シンボルテーブル用の8バイト配列に変換します。
// 8バイトを超える場合は、文字列テーブルに追加し、オフセットを Name フィールドに設定します。
// 重複を避けるために stringTableOffsetMap を使用します。
func (c *CoffFormat) convertNameToBytes(name string, stringTable *bytes.Buffer, offsetMap map[string]uint32) [8]byte {
	var result [8]byte
	if len(name) > 8 {
		// 8バイトを超える場合: 文字列テーブルへのオフセットを設定
		offset, exists := offsetMap[name]
		if !exists {
			// 文字列テーブルの現在の長さが新しいオフセットになる (サイズフィールド分を考慮)
			offset = uint32(stringTable.Len()) + coffStringTableSizeEntrySize
			stringTable.WriteString(name)
			stringTable.WriteByte(0) // NULL終端
			offsetMap[name] = offset // マップに記録
		}
		binary.LittleEndian.PutUint32(result[0:4], 0) // 最初の4バイトは0
		binary.LittleEndian.PutUint32(result[4:8], offset)
	} else {
		// 8バイト以下の場合: 直接名前をコピー
		copy(result[:], name)
	}
	return result
}

// addStringToStringTable は文字列を文字列テーブルに追加し、そのオフセットを返します。
// この関数は convertNameToBytes に統合されたため、削除または未使用としてマークできます。
/* // Add missing asterisk for block comment
func (c *CoffFormat) addStringToStringTable(s string, stringTable *bytes.Buffer) uint32 {
	// ... (実装は convertNameToBytes に移動) ...
}
*/

package filefmt

import (
	"bytes"
	"encoding/binary"
	"fmt"
	// "io" // 不要になったため削除
	"log"
	"os"
	"sort" // sort パッケージをインポート

	"github.com/HobbyOSs/gosk/internal/codegen"
	"github.com/lunixbochs/struc" // struc パッケージをインポート
)

// COFF ファイルヘッダ構造体
type CoffHeader struct {
	Machine              uint16 `struc:",little"` // ターゲットマシン (0x014c for i386)
	NumberOfSections     uint16 `struc:",little"` // セクション数
	TimeDateStamp        uint32 `struc:",little"` // ファイル作成日時 (Unix timestamp)
	PointerToSymbolTable uint32 `struc:",little"` // シンボルテーブルへのファイルオフセット
	NumberOfSymbols      uint32 `struc:",little"` // シンボルテーブル内のエントリ数 (補助シンボル含む)
	SizeOfOptionalHeader uint16 `struc:",little"` // オプショナルヘッダのサイズ (オブジェクトファイルでは通常0)
	Characteristics      uint16 `struc:",little"` // ファイル特性フラグ
}

// COFF セクションヘッダ構造体
type CoffSectionHeader struct {
	Name                 [8]byte // タグ不要 (struc はバイト配列をそのまま扱う)
	VirtualSize          uint32  `struc:",little"` // (未使用)
	VirtualAddress       uint32  `struc:",little"` // (未使用)
	SizeOfRawData        uint32  `struc:",little"` // セクションのデータサイズ (バイト単位)
	PointerToRawData     uint32  `struc:",little"` // セクションデータへのファイルオフセット
	PointerToRelocations uint32  `struc:",little"` // 再配置情報へのファイルオフセット (今回は0)
	PointerToLinenumbers uint32  `struc:",little"` // 行番号情報へのファイルオフセット (今回は0)
	NumberOfRelocations  uint16  `struc:",little"` // 再配置エントリ数 (今回は0)
	NumberOfLinenumbers  uint16  `struc:",little"` // 行番号エントリ数 (今回は0)
	Characteristics      uint32  `struc:",little"` // セクション特性フラグ
}

// COFF シンボルテーブルエントリ構造体
type CoffSymbol struct {
	Name               [8]byte // タグ不要
	Value              uint32  `struc:",little"` // シンボルの値 (アドレスなど)
	SectionNumber      int16   `struc:",little"` // 関連セクション番号 (1ベース、0=未定義, -1=絶対, -2=デバッグ)
	Type               uint16  `struc:",little"` // シンボル型 (基本型と派生型)
	StorageClass       uint8   `struc:",little"` // 格納クラス (スコープ、型など)
	NumberOfAuxSymbols uint8   `struc:",little"` // 補助シンボルエントリ数
}

// COFF 補助セクションシンボルエントリ構造体
// 注意: この構造体は直接 Pack/Unpack されないため、タグは不要。
//
//	補助シンボルは []byte として扱われる。
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
	numSymbolsWritten := uint32(0) // 書き込まれたシンボルレコードの総数 (メイン+補助)
	for _, entry := range allSymbolEntries {
		entryName := "N/A"
		if len(entry.Main.Name) > 0 {
			nullIdx := bytes.IndexByte(entry.Main.Name[:], 0)
			if nullIdx == -1 {
				nullIdx = 8
			}
			entryName = string(entry.Main.Name[:nullIdx])
		}

		// struc を使ってメインシンボルを書き込む (リトルエンディアン指定)
		// PackWithOptions を使うことで、デフォルトのビッグエンディアンではなくリトルエンディアンを指定できる
		// (struc タグでも指定しているが、明示的に Options で指定する方が確実)
		if err := struc.PackWithOptions(symbolTableResultBytes, &entry.Main, &struc.Options{Order: binary.LittleEndian}); err != nil {
			return fmt.Errorf("failed to pack symbol %s: %w", entryName, err)
		}
		numSymbolsWritten++ // メインシンボルをカウント

		// 補助シンボルの書き込み (ここは変更なし、[]byte を直接書き込む)
		if entry.Aux != nil {
			if len(entry.Aux) != coffSymbolSize {
				symbolName := string(entry.Main.Name[:bytes.IndexByte(entry.Main.Name[:], 0)])
				log.Printf("Error: Aux symbol for %s has incorrect size %d, expected %d", symbolName, len(entry.Aux), coffSymbolSize)
				return fmt.Errorf("aux symbol for %s has incorrect size %d", symbolName, len(entry.Aux))
			}
			// 補助シンボルはそのまま書き込む
			if _, err := symbolTableResultBytes.Write(entry.Aux); err != nil {
				symbolName := string(entry.Main.Name[:bytes.IndexByte(entry.Main.Name[:], 0)])
				log.Printf("Error writing aux symbol for %s: %v", symbolName, err)
				return fmt.Errorf("failed to write aux symbol for %s: %w", symbolName, err)
			}
			// 補助シンボルもレコードとしてカウント
			numSymbolsWritten += uint32(entry.Main.NumberOfAuxSymbols)
		}
	}

	if _, err := buf.Write(symbolTableResultBytes.Bytes()); err != nil {
		return fmt.Errorf("failed to write symbol table: %w", err)
	}

	// --- 文字列テーブル書き込み (元のシンプルなロジック) ---
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

	// --- ヘッダとセクションヘッダの生成 (オフセット情報を含む) ---
	// ヘッダの NumberOfSymbols には実際に書き込んだシンボルレコード数を使用
	numSymbolsForHeader := numSymbolsWritten
	header := c.generateHeader(ctx, numSections, symbolTableOffset, numSymbolsForHeader)
	sectionHeaders := c.generateSectionHeaders(ctx, textDataOffset, textDataSize, dataDataOffset, dataDataSize, bssDataSize, symbolTableOffset) // Pass symbolTableOffset

	// --- 最終的なファイル内容を構築 (プレースホルダーを上書き) ---
	finalBytes := buf.Bytes() // 現在のバッファ内容を取得

	// ヘッダを上書き
	headerBuf := new(bytes.Buffer)
	// struc を使ってヘッダを書き込む (リトルエンディアン指定)
	if err := struc.PackWithOptions(headerBuf, &header, &struc.Options{Order: binary.LittleEndian}); err != nil {
		return fmt.Errorf("failed to pack final header: %w", err)
	}
	if headerBuf.Len() != coffHeaderSize {
		return fmt.Errorf("packed header size mismatch: expected %d, got %d", coffHeaderSize, headerBuf.Len())
	}
	copy(finalBytes[0:coffHeaderSize], headerBuf.Bytes())

	// セクションヘッダを上書き (ループで各ヘッダをパック)
	currentOffset := coffHeaderSize
	for i := range sectionHeaders {
		sectionHeaderBuf := new(bytes.Buffer)
		// struc を使ってセクションヘッダを書き込む (リトルエンディアン指定)
		if err := struc.PackWithOptions(sectionHeaderBuf, &sectionHeaders[i], &struc.Options{Order: binary.LittleEndian}); err != nil {
			return fmt.Errorf("failed to pack final section header %d: %w", i, err)
		}
		if sectionHeaderBuf.Len() != coffSectionHeaderSize {
			return fmt.Errorf("packed section header %d size mismatch: expected %d, got %d", i, coffSectionHeaderSize, sectionHeaderBuf.Len())
		}
		copy(finalBytes[currentOffset:currentOffset+coffSectionHeaderSize], sectionHeaderBuf.Bytes())
		currentOffset += coffSectionHeaderSize
	}

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
		NumberOfSymbols:      numSymbols, // 引数で渡された値を使用
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
		PointerToRawData:     0, // Revert back to 0 to match expected output for this test
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

	// シンボルテーブルをアドレス順にソートする (.file とセクションシンボルを除く)
	// .file シンボル (index 0) とセクションシンボル (index 1, 2, 3) は固定
	if len(allEntries) > 4 { // ソート対象が存在する場合のみ
		sort.SliceStable(allEntries[4:], func(i, j int) bool {
			// SectionNumber が 0 (未定義/外部) のシンボルは最後に配置する
			isExternI := allEntries[4+i].Main.SectionNumber == 0
			isExternJ := allEntries[4+j].Main.SectionNumber == 0
			if isExternI && !isExternJ {
				return false // i (外部) は j より後
			}
			if !isExternI && isExternJ {
				return true // i は j (外部) より前
			}
			// SectionNumber が同じか、両方とも外部でない場合は Value で比較
			return allEntries[4+i].Main.Value < allEntries[4+j].Main.Value
		})
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
			// 文字列テーブルの現在の内容の長さが、次の書き込み開始位置（オフセット）になる。
			// 文字列テーブルの先頭には4バイトのサイズフィールドがあるため、内容のオフセットは4から始まる。
			offset = uint32(stringTable.Len()) + coffStringTableSizeEntrySize // 正しいオフセット計算 (内容の長さ + サイズフィールド長)
			stringTable.WriteString(name)
			stringTable.WriteByte(0) // 末尾にヌルバイトを追加
			offsetMap[name] = offset // マップに記録
		}
		// Nameフィールドの最初の4バイトは0、次の4バイトに文字列テーブルへのオフセットを設定
		binary.LittleEndian.PutUint32(result[0:4], 0)
		binary.LittleEndian.PutUint32(result[4:8], offset)
	} else {
		// 8バイト以下の場合: 直接名前をコピーし、残りはヌルで埋める
		copy(result[:], name) // result はゼロ初期化されているので、コピーされなかった部分はヌルのまま
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

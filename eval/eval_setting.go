package eval

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/hangingman/gosk/ast"
	"github.com/hangingman/gosk/object"
	"log"
	// "unsafe"
)

type PIMAGE_FILE_HEADER struct {
	machine              uint16
	numberOfSections     uint16
	timeDateStamp        int32
	pointerToSymbolTable int32
	numberOfSymbols      int32
	sizeOfOptionalHeader uint16
	characteristics      uint16
}

type PIMAGE_SECTION_HEADER struct {
	name                 [8]uint8
	misc                 uint32
	virtualAddress       uint32
	sizeOfRawData        uint32
	pointerToRawData     uint32
	pointerToRelocations uint32
	pointerToLinenumbers uint32
	numberOfRelocations  uint16
	numberOfLinenumbers  uint16
	characteristics      uint32
}

type PIMAGE_SYMBOL struct {
	shortName          [8]uint8
	value              uint32
	sectionNumber      uint16
	symbolType         uint16
	storageClass       uint8
	numberOfAuxSymbols uint8
}

func evalFormat(b *object.Binary, format string) {
	log.Println(fmt.Sprintf("Process for %s", format))
	buf := new(bytes.Buffer)

	// PE file header
	header := PIMAGE_FILE_HEADER{
		machine:              0x14c,
		numberOfSections:     0x0003,
		pointerToSymbolTable: 0x0000008e,
		numberOfSymbols:      0x00000009,
	}
	binary.Write(buf, binary.LittleEndian, &header)

	// PE section header
	text := PIMAGE_SECTION_HEADER{
		//{'.', 't', 'e', 'x', 't', 0, 0, 0 /* name */},
		name:                 [8]uint8{0x2e, 0x74, 0x65, 0x78, 0x74, 0x00, 0x00, 0x00},
		misc:                 0x00000000, // Misc
		virtualAddress:       0x00000000, // VirtualAddress
		sizeOfRawData:        0x00000002, // SizeOfRawData
		pointerToRawData:     0x0000008c, // PointerToRawData
		pointerToRelocations: 0x0000008e, // PointerToRelocations
		pointerToLinenumbers: 0x00000000, // PointerToLinenumbers
		numberOfRelocations:  0x0000,     // NumberOfRelocations
		numberOfLinenumbers:  0x0000,     // NumberOfLinenumbers
		characteristics:      0x60100020, // +0x38: flags, default_align = 1
	}
	data := PIMAGE_SECTION_HEADER{
		//{ '.', 'd', 'a', 't', 'a', 0, 0, 0 /* name */ },
		name:                 [8]uint8{0x2e, 0x64, 0x61, 0x74, 0x61, 0x00, 0x00, 0x00},
		misc:                 0x00000000, // Misc
		virtualAddress:       0x00000000, // VirtualAddress
		sizeOfRawData:        0x00000000, // SizeOfRawData
		pointerToRawData:     0x00000000, // PointerToRawData
		pointerToRelocations: 0x00000000, // PointerToRelocations
		pointerToLinenumbers: 0x00000000, // PointerToLinenumbers
		numberOfRelocations:  0x0000,     // NumberOfRelocations
		numberOfLinenumbers:  0x0000,     // NumberOfLinenumbers
		characteristics:      0xc0100040, // +0x38: flags, default_align = 1
	}
	bss := PIMAGE_SECTION_HEADER{
		//{ '.', 'b', 's', 's', 0, 0, 0, 0 /* name */ },
		name:                 [8]uint8{0x2e, 0x62, 0x73, 0x73, 0x00, 0x00, 0x00, 0x00},
		misc:                 0x00000000, // Misc
		virtualAddress:       0x00000000, // VirtualAddress
		sizeOfRawData:        0x00000000, // SizeOfRawData
		pointerToRawData:     0x00000000, // PointerToRawData
		pointerToRelocations: 0x00000000, // PointerToRelocations
		pointerToLinenumbers: 0x00000000, // PointerToLinenumbers
		numberOfRelocations:  0x0000,     // NumberOfRelocations
		numberOfLinenumbers:  0x0000,     // NumberOfLinenumbers
		characteristics:      0xc0100080, // +0x38: flags, default_align = 1
	}

	binary.Write(buf, binary.LittleEndian, &text)
	binary.Write(buf, binary.LittleEndian, &data)
	binary.Write(buf, binary.LittleEndian, &bss)
	b.Value = append(b.Value, buf.Bytes()...)
	log.Println("Wrote '.text', '.data', '.bss' fields for Portable Executable")
}

func evalBits(b *object.Binary, bit string) {
	log.Println(fmt.Sprintf("Deal with bit %s", bit))
}

func evalFile(b *object.Binary, filename string) {
	log.Println(fmt.Sprintf("Add %s", filename))

	buf := new(bytes.Buffer)
	file := PIMAGE_SYMBOL{
		// { '.', 'f', 'i', 'l', 'e', 0, 0, 0 /* shortName */ },
		shortName:          [8]uint8{0x2e, 0x66, 0x69, 0x6c, 0x65, 0x00, 0x00, 0x00},
		value:              0x00000000,
		sectionNumber:      0xfffe,
		symbolType:         0x0000,
		storageClass:       0x67,
		numberOfAuxSymbols: 0x01,
	}
	binary.Write(buf, binary.LittleEndian, &file)
	b.Value = append(b.Value, buf.Bytes()...)
}

func evalSettingStatement(stmt *ast.SettingStatement) object.Object {
	tok := stmt.Name.Token.Literal
	val := stmt.Name.Value

	bin := &object.Binary{Value: []byte{}}

	switch {
	case tok == "FORMAT":
		evalFormat(bin, val)
	case tok == "BITS":
		evalBits(bin, val)
		// case tok == "FILE":
		// 	evalFile(bin, val)
	}

	stmt.Bin = bin
	return stmt.Bin
}

func evalSectionTable() object.Object {
	bin := &object.Binary{Value: []byte{}}

	// セクションデータのサイズが確定(SizeOfRawData)
	// var offset int = curByteSize;
	// var sizeOfRawData int = offset - (unsafe.Sizeof(PIMAGE_FILE_HEADER) + unsafe.Sizeof(PIMAGE_SECTION_HEADER) * 3)

	// offset + realoc * EXTERN symbols
	// var pointerToSymbolTable int = offset + unsafe.Sizeof(COFF_RELOCATION) * len(ExternSymbolList)
	// log.Println(fmt.Sprintf("COFF file header's PointerToSymbolTable: 0x{:02x}", pointerToSymbolTable))
	// bin.Value = append(bin.Value, imm32ToDword(pointerToSymbolTable)...)

	// log.Println(fmt.Sprintf("section table '.text' PointerToSymbolTable: 0x{:02x}", offset + 4))
	//   set_dword_into_binout(offset, binout_container, false, unsafe.Sizeof(PIMAGE_FILE_HEADER) + 24);
	// log.Println(fmt.Sprintf("section table '.text' SizeOfRawData: 0x{:02x}", offset))
	//   set_dword_into_binout(size_of_raw_data, binout_container, false, unsafe.Sizeof(PIMAGE_FILE_HEADER) + 16);
	// log.Println(fmt.Sprintf("section table '.text' NumberOfRelocations: 0x{:02x}", ExternSymbolList.size()))
	//   set_word_into_binout(ExternSymbolList.size(), binout_container, false, unsafe.Sizeof(PIMAGE_FILE_HEADER) + 32);

	return bin
}

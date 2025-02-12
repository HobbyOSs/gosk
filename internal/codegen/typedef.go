package codegen

type Byte uint8
type Word uint16
type DWord uint32

func ConvertToByte(value DWord) Byte {
	return Byte(value & 0xFF)
}

func ConvertToWord(value DWord) Word {
	return Word(value & 0xFFFF)
}

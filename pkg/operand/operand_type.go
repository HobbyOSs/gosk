package operand

type OperandType string

const (
	CodeOneWordMemoryOrTwoDoubleWordMemory OperandType = "OneWordMemoryOrTwoDoubleWordMemory"
	CodeByte                               OperandType = "Byte"
	CodePackedBCD                          OperandType = "PackedBCD"
	CodeByteSizeOfDst                      OperandType = "ByteSizeOfDst"
	CodeByteSizeOfDst64Bit                 OperandType = "ByteSizeOfDst64Bit"
	CodeByteSizeOfStackPointer             OperandType = "ByteSizeOfStackPointer"
	CodeByteOrWordUnused                   OperandType = "ByteOrWordUnused"
	CodeDoubleword                         OperandType = "Doubleword"
	CodeDoublewordInteger                  OperandType = "DoublewordInteger"
	CodeDoubleQuadword                     OperandType = "DoubleQuadword"
	CodeDoubleReal                         OperandType = "DoubleReal"
	CodeDoubleRealExt                      OperandType = "DoubleRealExt"
	CodeDoublewordSignExt                  OperandType = "DoublewordSignExt"
	CodeEnvFPU                             OperandType = "EnvFPU"
	CodeExtendedRealOnly                   OperandType = "ExtendedRealOnly"
	CodeFarPointer                         OperandType = "FarPointer"
	CodeFarPointerInc                      OperandType = "FarPointerInc"
	CodeQuadwordMMX                        OperandType = "QuadwordMMX"
	CodePackedDouble                       OperandType = "PackedDouble"
	CodePackedSingle                       OperandType = "PackedSingle"
	CodePackedSinglePS                     OperandType = "PackedSinglePS"
	CodeQuadwordPromoted                   OperandType = "QuadwordPromoted"
	CodeScalarDouble                       OperandType = "ScalarDouble"
	CodeScalarSingle                       OperandType = "ScalarSingle"
	CodeFPUState                           OperandType = "FPUState"
	CodeFPUSIMDState                       OperandType = "FPUSIMDState"
	CodeFarPointerSize                     OperandType = "FarPointerSize"
	CodeWordOrDoubleword                   OperandType = "WordOrDoubleword"
	CodeWordDoubleSignExt                  OperandType = "WordDoubleSignExt"
	CodeQuadwordOrWord                     OperandType = "QuadwordOrWord"
	CodeWordQuadPromoted                   OperandType = "WordQuadPromoted"
	CodeWordSignExtSP                      OperandType = "WordSignExtSP"
	CodeWord                               OperandType = "Word"
	CodeWordInteger                        OperandType = "WordInteger"
)

var typeMap = map[string]OperandType{
	"a":   CodeOneWordMemoryOrTwoDoubleWordMemory, // 16/32&16/32
	"b":   CodeByte,                               // 8
	"bcd": CodePackedBCD,                          // 80dec
	"bs":  CodeByteSizeOfDst,                      // 8
	"bsq": CodeByteSizeOfDst64Bit,                 // -
	"bss": CodeByteSizeOfStackPointer,             // 8
	"c":   CodeByteOrWordUnused,                   // ?
	"d":   CodeDoubleword,                         // 32
	"di":  CodeDoublewordInteger,                  // 32int
	"dq":  CodeDoubleQuadword,                     // 128
	"dr":  CodeDoubleReal,                         // 64real
	"er":  CodeDoubleRealExt,                      // 80real
	"ds":  CodeDoublewordSignExt,                  // 32
	"e":   CodeEnvFPU,                             // 14/28
	"er8": CodeExtendedRealOnly,                   //
	"p":   CodeFarPointer,                         // 16:16/32
	"ptp": CodeFarPointerInc,                      // 16:16/32/64
	"pi":  CodeQuadwordMMX,                        // (64)
	"pd":  CodePackedDouble,                       //
	"ps":  CodePackedSingle,                       // (128)
	"psq": CodePackedSinglePS,                     // 64
	"qp":  CodeQuadwordPromoted,                   // 64
	"sd":  CodeScalarDouble,                       //
	"ss":  CodeScalarSingle,                       //
	"st":  CodeFPUState,                           //
	"stx": CodeFPUSIMDState,                       //
	"t":   CodeFarPointerSize,                     //
	"y":   CodeWordOrDoubleword,                   //
	"v":   CodeWordDoubleSignExt,                  //
	"vq":  CodeQuadwordOrWord,                     //
	"vp":  CodeWordQuadPromoted,                   //
	"vs":  CodeWordSignExtSP,                      //
	"w":   CodeWord,                               //
	"wi":  CodeWordInteger,                        //
}

func (ot OperandType) String() string {
	for key, val := range typeMap {
		if val == ot {
			return key
		}
	}
	return "unknown"
}

func GetType(s string) (OperandType, bool) {
	code, ok := typeMap[s]
	return code, ok
}

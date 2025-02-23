package operand

type OperandType string

// ref: https://opcodes.readthedocs.io/opcodes.html#opcodes.x86.Operand
const (
	Code1           OperandType = "1"
	Code3           OperandType = "3"
	CodeAL          OperandType = "al"
	CodeAX          OperandType = "ax"
	CodeEAX         OperandType = "eax"
	CodeCL          OperandType = "cl"
	CodeXMM0        OperandType = "xmm0"
	CodeREL8        OperandType = "rel8"
	CodeREL32       OperandType = "rel32"
	CodeIMM         OperandType = "imm"
	CodeIMM4        OperandType = "imm4"
	CodeIMM8        OperandType = "imm8"
	CodeIMM16       OperandType = "imm16"
	CodeIMM32       OperandType = "imm32"
	CodeR8          OperandType = "r8"
	CodeR16         OperandType = "r16"
	CodeR32         OperandType = "r32"
	CodeMM          OperandType = "mm"
	CodeXMM         OperandType = "xmm"
	CodeXMMK        OperandType = "xmm{k}"
	CodeXMMKZ       OperandType = "xmm{k}{z}"
	CodeYMM         OperandType = "ymm"
	CodeYMMK        OperandType = "ymm{k}"
	CodeYMMKZ       OperandType = "ymm{k}{z}"
	CodeZMM         OperandType = "zmm"
	CodeZMMK        OperandType = "zmm{k}"
	CodeZMMKZ       OperandType = "zmm{k}{z}"
	CodeK           OperandType = "k"
	CodeCR          OperandType = "cr"
	CodeKK          OperandType = "k{k}"
	CodeM           OperandType = "m"
	CodeM8          OperandType = "m8"
	CodeM16         OperandType = "m16"
	CodeM16KZ       OperandType = "m16{k}{z}"
	CodeM32         OperandType = "m32"
	CodeM32K        OperandType = "m32{k}"
	CodeM32KZ       OperandType = "m32{k}{z}"
	CodeM64         OperandType = "m64"
	CodeM64K        OperandType = "m64{k}"
	CodeM64KZ       OperandType = "m64{k}{z}"
	CodeM80         OperandType = "m80"
	CodeM128        OperandType = "m128"
	CodeM128KZ      OperandType = "m128{k}{z}"
	CodeM256        OperandType = "m256"
	CodeM256KZ      OperandType = "m256{k}{z}"
	CodeM512        OperandType = "m512"
	CodeM512KZ      OperandType = "m512{k}{z}"
	CodeM64M32BCST  OperandType = "m64/m32bcst"
	CodeM128M32BCST OperandType = "m128/m32bcst"
	CodeM256M32BCST OperandType = "m256/m32bcst"
	CodeM512M32BCST OperandType = "m512/m32bcst"
	CodeM128M64BCST OperandType = "m128/m64bcst"
	CodeM256M64BCST OperandType = "m256/m64bcst"
	CodeM512M64BCST OperandType = "m512/m64bcst"
	CodeVM32X       OperandType = "vm32x"
	CodeVM32XK      OperandType = "vm32x{k}"
	CodeVM32Y       OperandType = "vm32y"
	CodeVM32YK      OperandType = "vm32y{k}"
	CodeVM32Z       OperandType = "vm32z"
	CodeVM32ZK      OperandType = "vm32z{k}"
	CodeVM64X       OperandType = "vm64x"
	CodeVM64XK      OperandType = "vm64x{k}"
	CodeVM64Y       OperandType = "vm64y"
	CodeVM64YK      OperandType = "vm64y{k}"
	CodeVM64Z       OperandType = "vm64z"
	CodeVM64ZK      OperandType = "vm64z{k}"
	CodeSAE         OperandType = "{sae}"
	CodeER          OperandType = "{er}"
	CodeTR          OperandType = "tr"
	CodeDR          OperandType = "dr"
)

var typeMap = map[string]OperandType{
	"1":            Code1,
	"3":            Code3,
	"al":           CodeAL,
	"ax":           CodeAX,
	"eax":          CodeEAX,
	"cl":           CodeCL,
	"xmm0":         CodeXMM0,
	"rel8":         CodeREL8,
	"rel32":        CodeREL32,
	"imm4":         CodeIMM4,
	"imm8":         CodeIMM8,
	"imm16":        CodeIMM16,
	"imm32":        CodeIMM32,
	"r8":           CodeR8,
	"r16":          CodeR16,
	"r32":          CodeR32,
	"mm":           CodeMM,
	"xmm":          CodeXMM,
	"xmm{k}":       CodeXMMK,
	"xmm{k}{z}":    CodeXMMKZ,
	"ymm":          CodeYMM,
	"ymm{k}":       CodeYMMK,
	"ymm{k}{z}":    CodeYMMKZ,
	"zmm":          CodeZMM,
	"zmm{k}":       CodeZMMK,
	"zmm{k}{z}":    CodeZMMKZ,
	"k":            CodeK,
	"k{k}":         CodeKK,
	"m":            CodeM,
	"m8":           CodeM8,
	"m16":          CodeM16,
	"m16{k}{z}":    CodeM16KZ,
	"m32":          CodeM32,
	"m32{k}":       CodeM32K,
	"m32{k}{z}":    CodeM32KZ,
	"m64":          CodeM64,
	"m64{k}":       CodeM64K,
	"m64{k}{z}":    CodeM64KZ,
	"m80":          CodeM80,
	"m128":         CodeM128,
	"m128{k}{z}":   CodeM128KZ,
	"m256":         CodeM256,
	"m256{k}{z}":   CodeM256KZ,
	"m512":         CodeM512,
	"m512{k}{z}":   CodeM512KZ,
	"m64/m32bcst":  CodeM64M32BCST,
	"m128/m32bcst": CodeM128M32BCST,
	"m256/m32bcst": CodeM256M32BCST,
	"m512/m32bcst": CodeM512M32BCST,
	"m128/m64bcst": CodeM128M64BCST,
	"m256/m64bcst": CodeM256M64BCST,
	"m512/m64bcst": CodeM512M64BCST,
	"vm32x":        CodeVM32X,
	"vm32x{k}":     CodeVM32XK,
	"vm32y":        CodeVM32Y,
	"vm32y{k}":     CodeVM32YK,
	"vm32z":        CodeVM32Z,
	"vm32z{k}":     CodeVM32ZK,
	"vm64x":        CodeVM64X,
	"vm64x{k}":     CodeVM64XK,
	"vm64y":        CodeVM64Y,
	"vm64y{k}":     CodeVM64YK,
	"vm64z":        CodeVM64Z,
	"vm64z{k}":     CodeVM64ZK,
	"{sae}":        CodeSAE,
	"{er}":         CodeER,
	"cr":           CodeCR,
	"tr":           CodeTR,
	"dr":           CodeDR,
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

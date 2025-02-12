package operand

type AddressingType string

const (
	CodeDirectAddress       AddressingType = "DirectAddress"
	CodeBasedAddressA       AddressingType = "BasedAddressA"
	CodeBasedAddressB       AddressingType = "BasedAddressB"
	CodeBasedAddressD       AddressingType = "BasedAddressD"
	CodeCRField             AddressingType = "CRField"
	CodeDebugField          AddressingType = "DebugField"
	CodeModRMAddress        AddressingType = "ModRMAddress"
	CodeModRMAddressX87FPU  AddressingType = "ModRMAddressX87FPU"
	CodeModRMX87FPU         AddressingType = "ModRMX87FPU"
	CodeRFlags              AddressingType = "RFlags"
	CodeGeneralReg          AddressingType = "GeneralReg"
	CodeGeneralRegAddr      AddressingType = "GeneralRegAddr"
	CodeImmediate           AddressingType = "Immediate"
	CodeRelativeOffset      AddressingType = "RelativeOffset"
	CodeModRMAddrOnlyMemory AddressingType = "ModRMAddrOnlyMemory"
	CodeModRM_MMX           AddressingType = "ModRM_MMX"
	CodeModRMAddressMoffs   AddressingType = "ModRMAddressMoffs"
	CodeModRM_MMXRegField   AddressingType = "ModRM_MMXRegField"
	CodeModRMMinor          AddressingType = "ModRMMinor"
	CodeSregField           AddressingType = "SregField"
	CodeStackField          AddressingType = "StackField"
	CodeRegFieldTest        AddressingType = "RegFieldTest"
	CodeXmmRMField          AddressingType = "XmmRMField"
	CodeXmmRegField         AddressingType = "XmmRegField"
	CodeXmmOperand          AddressingType = "XmmOperand"
	CodeMemoryAddressX      AddressingType = "MemoryAddressX"
	CodeMemoryAddressY      AddressingType = "MemoryAddressY"
	CodeModRM               AddressingType = "ModRM" // r
)

var codeMap = map[string]AddressingType{
	"A":   CodeDirectAddress,
	"BA":  CodeBasedAddressA,
	"BB":  CodeBasedAddressB,
	"BD":  CodeBasedAddressD,
	"C":   CodeCRField,            // CRn
	"D":   CodeDebugField,         // DRn
	"E":   CodeModRMAddress,       // r/m
	"ES":  CodeModRMAddressX87FPU, // STi/m
	"EST": CodeModRMX87FPU,        // STi
	"F":   CodeRFlags,             // rFLAGS
	"G":   CodeGeneralReg,         // AX
	"H":   CodeGeneralRegAddr,     // 0F20
	"I":   CodeImmediate,
	"J":   CodeRelativeOffset,      // e.g. JMP LOOP
	"M":   CodeModRMAddrOnlyMemory, // mod != 11bin (BOUND, LEA, CALLF, JMPF, LES, LDS, LSS, LFS, LGS, CMPXCHG8B, CMPXCHG16B, F20FF0 LDDQU)
	"N":   CodeModRM_MMX,
	"O":   CodeModRMAddressMoffs,
	"P":   CodeModRM_MMXRegField,
	"R":   CodeModRMMinor,
	"S":   CodeSregField,
	"SC":  CodeStackField,
	"T":   CodeRegFieldTest,
	"U":   CodeXmmRMField,
	"V":   CodeXmmRegField,
	"W":   CodeXmmOperand,
	"X":   CodeMemoryAddressX, // Memory addressed by the DS:eSI or by RSI
	"Y":   CodeMemoryAddressY, // Memory addressed by the ES:eDI or by RDI
	"Z":   CodeModRM,
}

func GetCode(s string) (AddressingType, bool) {
	code, ok := codeMap[s]
	return code, ok
}

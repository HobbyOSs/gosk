package operand

import (
	"strings"
)

type BaseOperand struct {
	Internal string
}

func (b BaseOperand) AddressingType() AddressingType {
	parser := getParser()
	parsed, err := parser.ParseString("", b.Internal)

	if err != nil {
		return "unknown"
	}

	switch {
	case parsed.SegReg != nil:
		switch parsed.SegReg.Seg {
		case "DS":
			return CodeMemoryAddressX
		case "ES":
			return CodeMemoryAddressY
		default:
			return CodeSregField
		}
	//case parsed.SegMem != nil:
	//	return CodeModRMAddress
	case parsed.Reg != "":
		// Add length checks before accessing slices
		switch {
		case len(parsed.Reg) >= 2 && parsed.Reg[:2] == "MM":
			return CodeModRM_MMX
		case len(parsed.Reg) >= 3 && parsed.Reg[:3] == "XMM":
			return CodeXmmRegField
		case len(parsed.Reg) >= 3 && parsed.Reg[:3] == "YMM":
			return CodeXmmRMField
		case len(parsed.Reg) >= 2 && parsed.Reg[:2] == "TR":
			return CodeRegFieldTest
		case len(parsed.Reg) >= 2 && parsed.Reg[:2] == "CR":
			return CodeCRField
		case len(parsed.Reg) >= 2 && parsed.Reg[:2] == "DR":
			return CodeDebugField
		default:
			return CodeGeneralReg
		}
	case parsed.Mem != "":
		if strings.Contains(b.Internal, "DWORD PTR") ||
			strings.Contains(b.Internal, "XMMWORD PTR") ||
			strings.Contains(b.Internal, "YMMWORD PTR") {
			return CodeModRMAddress
		}
		return CodeModRMAddress
	case parsed.Imm != "":
		return CodeImmediate
	case parsed.Rel != "":
		return CodeRelativeOffset
	case parsed.Addr != "":
		return CodeDirectAddress
	case parsed.Seg != "":
		return CodeSregField
	// case parsed.MemOffset != "":
	// 	return CodeModRMAddressMoffs
	default:
		return "unknown"
	}
}

func (b *BaseOperand) OperandType() OperandType {
	parser := getParser()
	parsed, err := parser.ParseString("", b.Internal)
	if err == nil {
		switch {
		case parsed.Reg != "":
			return CodeDoubleword
		case parsed.Mem != "":
			return CodeDoubleword
		case parsed.Imm != "":
			return CodeDoublewordInteger
		case parsed.Seg != "":
			return CodeWord
		case parsed.Rel != "":
			return CodeWord
		case parsed.Addr != "":
			return CodeDoubleword
		}
	}
	return OperandType("unknown")
}

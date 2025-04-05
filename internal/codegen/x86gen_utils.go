package codegen

import (
	"encoding/binary" // Added missing import
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/HobbyOSs/gosk/pkg/asmdb"
	"github.com/HobbyOSs/gosk/pkg/cpu"
	"github.com/HobbyOSs/gosk/pkg/ng_operand" // Re-import ng_operand
)

// is32BitRegister checks if a register name corresponds to a 32-bit general-purpose register.
func is32BitRegister(regName string) bool {
	switch regName {
	case "EAX", "ECX", "EDX", "EBX", "ESP", "EBP", "ESI", "EDI":
		return true
	default:
		return false
	}
}

// GenerateModRM はエンコーディング情報とビットモードに基づいてModR/Mバイトを生成する
func GenerateModRM(operands []string, modRM *asmdb.Encoding, bitMode cpu.BitMode) ([]byte, error) { // Keep cpu.BitMode
	if modRM == nil || modRM.ModRM == nil {
		return nil, nil
	}

	modRMDef := modRM.ModRM
	if modRMDef == nil {
		return nil, nil
	}

	if strings.HasPrefix(modRMDef.Reg, "#") {
		// ModR/M の reg フィールドがオペランドの場合
		regIndex, err := parseIndex(modRMDef.Reg)
		if err != nil {
			return nil, fmt.Errorf("invalid ModRM.Reg format")
		}
		rmIndex, err := parseIndex(modRMDef.Rm)
		if err != nil {
			return nil, fmt.Errorf("invalid ModRM.RM format")
		}

		if regIndex < 0 || regIndex >= len(operands) {
			return nil, fmt.Errorf("ModRM.Reg index out of range")
		}
		if rmIndex < 0 || rmIndex >= len(operands) {
			return nil, fmt.Errorf("ModRM.RM index out of range")
		}

		// ModRMByOperand に渡す前にログ出力
		regOperand := operands[regIndex]
		rmOperand := operands[rmIndex]

		log.Printf("trace: GenerateModRM: regIndex=%d, regOperand=%s", regIndex, regOperand) // ADD LOG
		log.Printf("trace: GenerateModRM: rmIndex=%d, rmOperand=%s", rmIndex, rmOperand)     // ADD LOG
		log.Printf("trace: GenerateModRM: Calling ModRMByOperand with mode=%s, reg=%s, rm=%s", modRMDef.Mode, regOperand, rmOperand)

		modrmBytes, err := ModRMByOperand(modRMDef.Mode, regOperand, rmOperand, bitMode)
		if err != nil {
			return nil, fmt.Errorf("error in ModRMByOperand: %w", err) // Wrap error
		}
		return modrmBytes, nil
	} else {
		// ModR/M の reg フィールドが固定値の場合
		regValue, err := strconv.Atoi(modRMDef.Reg)
		if err != nil {
			return nil, fmt.Errorf("failed to parse ModRM.Reg: %v", err)
		}
		rmIndex, err := parseIndex(modRMDef.Rm)
		if err != nil {
			return nil, fmt.Errorf("invalid ModRM.RM format")
		}
		if rmIndex < 0 || rmIndex >= len(operands) {
			return nil, fmt.Errorf("ModRM.RM index out of range")
		}
		rmOperand := operands[rmIndex]
		return ModRMByValue(modRMDef.Mode, regValue, rmOperand, bitMode), nil
	}
}

// parseMode は modeStr から mode バイトを解析する
func parseMode(modeStr string) byte {
	// modeの解析（2ビット）
	var mode byte
	switch modeStr {
	case "#0": // レジスタ間接参照
		mode = 0b00000000
	case "#1": // 8ビット変位レジスタ間接参照
		mode = 0b01000000
	case "#2": // 32ビット変位レジスタ間接参照
		mode = 0b10000000
	case "11": // レジスタ
		mode = 0b11000000
	default:
		mode = 0 // デフォルト値
	}
	return mode
}

// ModRMByOperand はモード、regオペランド、rmオペランド、ビットモードに基づいてModR/Mバイトを生成する
func ModRMByOperand(modeStr string, regOperand string, rmOperand string, bitMode cpu.BitMode) ([]byte, error) { // Keep cpu.BitMode
	// ModR/M バイトの生成
	// |  mod  |  reg  |  r/m  |
	// | 7 6 | 5 4 3 | 2 1 0 |

	// mode := parseMode(modeStr) // mode は calculateModRM で決定される

	// regの解析（3ビット）
	reg, err := GetRegisterNumber(regOperand)
	if err != nil {
		return nil, fmt.Errorf("failed to get register number for %s: %w", regOperand, err)
	}
	regBits := byte(reg) << 3

	// r/mの解析
	if strings.Contains(rmOperand, "[") && strings.HasSuffix(rmOperand, "]") {
		// --- ng_operand を使用した ModR/M 生成ロジック ---
		rmOps, err := ng_operand.FromString(rmOperand)
		if err != nil {
			return nil, fmt.Errorf("failed to parse rmOperand '%s': %w", rmOperand, err)
		}
		rmOps = rmOps.WithBitMode(bitMode) // Ensure correct bit mode

		// Use the new GetMemoryInfo method
		memInfo, found := rmOps.GetMemoryInfo()
		if !found || memInfo == nil { // Add nil check for memInfo
			// This case should ideally not happen due to the check above, but handle defensively
			return nil, fmt.Errorf("could not get memory info from operand: %s", rmOperand)
		}

		// --- ModR/M and Displacement Calculation ---
		modrmByte, sibByte, dispBytes, err := calculateModRM(memInfo, bitMode, regBits) // Pass regBits for context
		if err != nil {
			return nil, fmt.Errorf("failed to calculate ModR/M for '%s': %w", rmOperand, err)
		}

		out := []byte{modrmByte}
		if sibByte != 0 { // Check if SIB byte is present
			out = append(out, sibByte)
		}
		if len(dispBytes) > 0 {
			out = append(out, dispBytes...)
		}
		log.Printf("debug: ModRMByOperand (mem): reg=%s(%b), rm=%s, result=%#x, sib=%#x, disp=% x",
			regOperand, regBits, rmOperand, modrmByte, sibByte, dispBytes)
		return out, nil

	}

	// r/m がレジスタの場合 (mod=11)
	rm, err := GetRegisterNumber(rmOperand)
	if err != nil {
		return nil, fmt.Errorf("failed to get register number for %s: %w", rmOperand, err)
	}
	rmBits := byte(rm)

	// レジスタの場合は mod=11 固定
	mod := byte(0b11000000)
	out := mod | regBits | rmBits
	log.Printf("debug: GenerateModRM (reg): reg=%s(%b), rm=%s(%b), result=%#x", regOperand, regBits, rmOperand, rmBits, out)
	return []byte{out}, nil
}

// ModRMByValue はモード、固定reg値、rmオペランドに基づいてModR/Mバイトを生成する
func ModRMByValue(modeStr string, regValue int, rmOperand string, bitMode cpu.BitMode) []byte { // Keep cpu.BitMode
	// ModR/M バイトの生成
	// |  mod  |  reg  |  r/m  |
	// | 7 6 | 5 4 3 | 2 1 0 |

	// mode := parseMode(modeStr) // mode は calculateModRM で決定される

	// regの解析（3ビット）
	regBits := byte(regValue) << 3

	// r/mの解析（3ビット）
	if strings.Contains(rmOperand, "[") && strings.HasSuffix(rmOperand, "]") {
		// --- ng_operand を使用した ModR/M 生成ロジック ---
		rmOps, err := ng_operand.FromString(rmOperand)
		if err != nil {
			log.Printf("error: Failed to parse rmOperand '%s' in ModRMByValue: %v", rmOperand, err)
			return []byte{0} // Consider returning error
		}
		rmOps = rmOps.WithBitMode(bitMode) // Ensure correct bit mode

		// Use the new GetMemoryInfo method
		memInfo, found := rmOps.GetMemoryInfo()
		if !found || memInfo == nil { // Add nil check for memInfo
			log.Printf("error: Could not get memory info from operand in ModRMByValue: %s", rmOperand)
			return []byte{0} // Consider returning error
		}

		// --- ModR/M and Displacement Calculation ---
		modrmByte, sibByte, dispBytes, err := calculateModRM(memInfo, bitMode, regBits) // Pass regBits for context
		if err != nil {
			log.Printf("error: Failed to calculate ModR/M for '%s' in ModRMByValue: %v", rmOperand, err)
			return []byte{0} // Consider returning error
		}

		out := []byte{modrmByte}
		if sibByte != 0 { // Check if SIB byte is present
			out = append(out, sibByte)
		}
		if len(dispBytes) > 0 {
			out = append(out, dispBytes...)
		}
		log.Printf("debug: ModRMByValue (mem): reg=%d(%b), rm=%s, result=%#x, sib=%#x, disp=% x",
			regValue, regBits, rmOperand, modrmByte, sibByte, dispBytes)
		return out

	}

	// r/m がレジスタの場合 (mod=11)
	rm, err := GetRegisterNumber(rmOperand)
	if err != nil {
		log.Printf("error: Failed to get register number for rm: %v", err)
		return []byte{0}
	}
	rmBits := byte(rm)

	// レジスタの場合は mod=11 固定
	mod := byte(0b11000000)
	out := mod | regBits | rmBits
	log.Printf("debug: ModRMByValue (reg): reg=%d(%b), rm=%s(%b), result=%#x", regValue, regBits, rmOperand, rmBits, out)
	return []byte{out}
}

// calculateModRM は MemoryInfo から ModR/M, SIB, Displacement を計算する
// regBits は ModR/M の reg フィールド (ビット3-5)
// TODO: SIBバイトの処理を実装する
func calculateModRM(mem *ng_operand.MemoryInfo, bitMode cpu.BitMode, regBits byte) (modrmByte byte, sibByte byte, dispBytes []byte, err error) {
	var mod byte
	var rm byte
	disp := mem.Displacement
	// BaseReg も IndexReg もなく、Displacement がある場合は直接アドレス
	isDirectAddr := mem.BaseReg == "" && mem.IndexReg == ""
	hasDisp := mem.Displacement != 0 || isDirectAddr // Direct address always has displacement

	// Determine mod based on displacement size and addressing mode
	if !hasDisp && !isDirectAddr { // No displacement and not direct address
		mod = 0b00000000
	} else if disp >= -128 && disp <= 127 {
		mod = 0b01000000 // disp8
	} else {
		// disp16 or disp32
		mod = 0b10000000
	}

	// --- 16-bit Addressing (Table 2-1) ---
	if bitMode == cpu.MODE_16BIT {
		sibByte = 0 // No SIB in 16-bit mode
		switch {
		case mem.BaseReg == "BX" && mem.IndexReg == "SI":
			rm = 0b000
		case mem.BaseReg == "BX" && mem.IndexReg == "DI":
			rm = 0b001
		case mem.BaseReg == "BP" && mem.IndexReg == "SI":
			rm = 0b010
		case mem.BaseReg == "BP" && mem.IndexReg == "DI":
			rm = 0b011
		case mem.BaseReg == "" && mem.IndexReg == "SI":
			rm = 0b100
		case mem.BaseReg == "" && mem.IndexReg == "DI":
			rm = 0b101
		case mem.BaseReg == "BP" && mem.IndexReg == "":
			rm = 0b110
			// [BP] requires displacement, even if 0
			if !hasDisp {
				mod = 0b01000000 // Use disp8=0
				disp = 0
				hasDisp = true
			}
		case isDirectAddr: // Direct address [disp16]
			mod = 0b00000000
			rm = 0b110
			hasDisp = true // Always has 16-bit displacement
		case mem.BaseReg == "BX" && mem.IndexReg == "":
			rm = 0b111
		// --- Cases not directly in Table 2-1 but implied ---
		case mem.BaseReg == "SI" && mem.IndexReg == "":
			rm = 0b100 // Treat [SI] as [SI+disp]
		case mem.BaseReg == "DI" && mem.IndexReg == "":
			rm = 0b101 // Treat [DI] as [DI+disp]
		default:
			// Check if 32-bit registers are used in 16-bit mode (requires 67h prefix)
			is32BitAddrMode := is32BitRegister(mem.BaseReg) || is32BitRegister(mem.IndexReg)
			if is32BitAddrMode {
				// If 32-bit registers are used, treat as 32-bit addressing mode for ModR/M calculation.
				// The 67h prefix should be added by the caller based on ng_operand.Require67h().
				// Jump to the 32-bit calculation logic.
				goto calculate_32bit_addressing
			}
			// Original default case for unsupported 16-bit modes
			return 0, 0, nil, fmt.Errorf("unsupported 16-bit addressing mode: Base=%s, Index=%s", mem.BaseReg, mem.IndexReg)
		}

		// Adjust mod if displacement exists but mod is currently 00 (except for direct address)
		if hasDisp && mod == 0b00000000 && rm != 0b110 {
			if disp >= -128 && disp <= 127 {
				mod = 0b01000000 // Use disp8
			} else {
				mod = 0b10000000 // Use disp16
			}
		}

		// Generate displacement bytes
		if hasDisp {
			if mod == 0b01000000 { // disp8
				dispBytes = []byte{byte(disp)}
			} else { // disp16 (mod=10 or mod=00,r/m=110)
				dispBytes = make([]byte, 2)
				binary.LittleEndian.PutUint16(dispBytes, uint16(disp))
			}
		}
		modrmByte = mod | regBits | rm
		return modrmByte, sibByte, dispBytes, nil
	}

calculate_32bit_addressing: // Label for the 32-bit logic start
	// --- 32-bit Addressing (Table 2-2) ---
	sibByte = 0 // Default: no SIB byte
	needsSIB := false

	switch {
	case isDirectAddr: // Direct address [disp32]
		mod = 0b00000000
		rm = 0b101
		hasDisp = true
	case mem.BaseReg == "EBP" && mem.IndexReg == "": // [EBP] or [EBP+disp]
		rm = 0b101
		if !hasDisp { // [EBP] needs disp8=0
			mod = 0b01000000
			disp = 0
			hasDisp = true
		}
	case mem.BaseReg == "ESP" || mem.IndexReg != "": // Needs SIB byte
		rm = 0b100
		needsSIB = true
	case mem.BaseReg == "EAX" && mem.IndexReg == "":
		rm = 0b000
	case mem.BaseReg == "ECX" && mem.IndexReg == "":
		rm = 0b001
	case mem.BaseReg == "EDX" && mem.IndexReg == "":
		rm = 0b010
	case mem.BaseReg == "EBX" && mem.IndexReg == "":
		rm = 0b011
	// case mem.BaseReg == "ESP": handled by SIB case
	// case mem.BaseReg == "EBP": handled above
	case mem.BaseReg == "ESI" && mem.IndexReg == "":
		rm = 0b110
	case mem.BaseReg == "EDI" && mem.IndexReg == "":
		rm = 0b111
	default:
		return 0, 0, nil, fmt.Errorf("unsupported 32-bit addressing mode: Base=%s, Index=%s", mem.BaseReg, mem.IndexReg)
	}

	// Adjust mod if displacement exists but mod is currently 00 (except for direct address and [EBP] cases)
	if hasDisp && mod == 0b00000000 && rm != 0b101 {
		if disp >= -128 && disp <= 127 {
			mod = 0b01000000 // Use disp8
		} else {
			mod = 0b10000000 // Use disp32
		}
	}

	// Generate SIB byte if needed (Implement fully based on Table 2-3)
	if needsSIB {

		var scale byte
		switch mem.Scale {
		case 1:
			scale = 0b00000000
		case 2:
			scale = 0b01000000
		case 4:
			scale = 0b10000000
		case 8:
			scale = 0b11000000
		default:
			if mem.Scale != 0 { // Allow scale 0 if index is not present
				return 0, 0, nil, fmt.Errorf("invalid SIB scale: %d", mem.Scale)
			}
			scale = 0b00000000 // Default to scale 1 if scale is 0 or index is empty
		}

		var indexNum int = 4 // Default to index=none (ESP encoding)
		if mem.IndexReg != "" {
			if mem.IndexReg == "ESP" {
				return 0, 0, nil, fmt.Errorf("ESP cannot be used as an index register in SIB")
			}
			indexNum, err = GetRegisterNumber(mem.IndexReg)
			if err != nil {
				return 0, 0, nil, fmt.Errorf("invalid index register in SIB: %s", mem.IndexReg)
			}
		}

		var baseNum int = 5 // Default to base=none ([disp32] or [EBP+disp] if mod=00)
		if mem.BaseReg != "" {
			baseNum, err = GetRegisterNumber(mem.BaseReg)
			if err != nil {
				return 0, 0, nil, fmt.Errorf("invalid base register in SIB: %s", mem.BaseReg)
			}
		}

		// Handle special case: mod=00 and base=EBP ([EBP+index*scale+disp32])
		// In this case, base field must be 5 (EBP), and a disp32 is always present.
		if mod == 0b00000000 && baseNum == 5 { // baseNum 5 corresponds to EBP
			// Base field remains 5, mod remains 00.
			// Ensure disp32 is handled correctly later.
			hasDisp = true // This combination always requires disp32
		} else if mem.BaseReg == "" { // No base register specified, implies base=EBP if mod=00
			// If there's no base register explicitly, and mod is 00,
			// the base field in SIB must be 5 (meaning disp32 follows).
			if mod == 0b00000000 {
				baseNum = 5
				hasDisp = true // Requires disp32
			}
			// If mod is 01 or 10, baseNum should reflect the actual base register (or lack thereof).
			// If BaseReg is truly empty, baseNum should technically be 5,
			// but the ModRM calculation logic might have already set mod to 01/10 based on displacement.
			// Let's stick with baseNum=5 if BaseReg is empty for SIB calculation.
			if mem.BaseReg == "" {
				baseNum = 5
			}
		}

		sibByte = scale | (byte(indexNum) << 3) | byte(baseNum)

		// Re-evaluate displacement requirement based on SIB base encoding
		// If base is EBP (5) and mod is 00, disp32 is required.
		// If base is EBP (5) and mod is 01, disp8 is required.
		// If base is EBP (5) and mod is 10, disp32 is required.
		if baseNum == 5 { // EBP or disp32 base
			if mod == 0b00000000 { // [disp32+index*scale]
				mod = 0b00000000 // Keep mod 00
				hasDisp = true   // Force disp32
			} else if mod == 0b01000000 { // [EBP+index*scale+disp8]
				hasDisp = true // Requires disp8
			} else { // [EBP+index*scale+disp32]
				hasDisp = true // Requires disp32
			}
		}
	}

	// Generate displacement bytes
	if hasDisp {
		if mod == 0b01000000 { // disp8
			dispBytes = []byte{byte(disp)}
		} else { // disp32 (mod=10 or mod=00 with rm=101 or SIB with base=101)
			dispBytes = make([]byte, 4)
			binary.LittleEndian.PutUint32(dispBytes, uint32(disp))
		}
	}

	modrmByte = mod | regBits | rm
	return modrmByte, sibByte, dispBytes, nil
}

// GetRegisterNumber はレジスタ名からレジスタ番号（0-7）を取得する
func GetRegisterNumber(regName string) (int, error) {
	regName = strings.TrimSpace(regName) // Trim whitespace
	switch regName {
	case "AL", "AX", "EAX", "RAX", "ES", "CR0":
		return 0, nil
	case "CL", "CX", "ECX", "RCX", "CS":
		return 1, nil
	case "DL", "DX", "EDX", "RDX", "SS", "CR2":
		return 2, nil
	case "BL", "BX", "EBX", "RBX", "DS", "CR3":
		return 3, nil
	case "AH", "SP", "ESP", "RSP", "FS", "CR4":
		return 4, nil
	case "CH", "BP", "EBP", "RBP", "GS":
		return 5, nil
	case "DH", "SI", "ESI", "RSI":
		return 6, nil
	case "BH", "DI", "EDI", "RDI":
		return 7, nil
	default:
		return 0, fmt.Errorf("unknown register: %s", regName)
	}
}

// ResolveOpcode はOpcodeとレジスタ番号を受け取り、最終的なオペコードバイト列を算出する。
// regNum はレジスタの番号（0-7）を表す。
func ResolveOpcode(op asmdb.Opcode, regNum int) ([]byte, error) {
	opBytes := []byte{}
	opStr := op.Byte

	// オペコード文字列をバイトごとに処理
	if len(opStr)%2 != 0 {
		return nil, fmt.Errorf("invalid opcode string length: %s", opStr)
	}
	for i := 0; i < len(opStr); i += 2 {
		byteStr := opStr[i : i+2]
		byteVal, err := strconv.ParseUint(byteStr, 16, 8)
		if err != nil {
			return nil, fmt.Errorf("invalid opcode byte string: %s in %s", byteStr, opStr)
		}
		opBytes = append(opBytes, byte(byteVal))
	}

	// Addendがある場合、最後のバイトにレジスタ番号を加算
	if op.Addend != nil && len(opBytes) > 0 && regNum >= 0 { // Added check for regNum >= 0
		regBits := byte(regNum & 0x07)
		lastByteIndex := len(opBytes) - 1
		opBytes[lastByteIndex] |= regBits
		log.Printf("debug: ResolveOpcode: base=%s, addend=%v, reg=%d, result=% x", opStr, op.Addend, regNum, opBytes)
	} else {
		log.Printf("debug: ResolveOpcode: base=%s, result=% x", opStr, opBytes)
	}

	return opBytes, nil
}

// getModRMFromOperands はオペランドからModR/Mバイトを生成する
func getModRMFromOperands(operands []string, modRM *asmdb.Encoding, bitMode cpu.BitMode) ([]byte, error) { // Keep cpu.BitMode
	modrmByte, err := GenerateModRM(operands, modRM, bitMode)
	if err != nil {
		return nil, err
	}
	return modrmByte, nil
}

func parseIndex(indexStr string) (int, error) {
	if strings.HasPrefix(indexStr, "#") {
		indexStr = indexStr[1:]
	}
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		return -1, fmt.Errorf("invalid index format")
	}
	return index, nil
}

// getImmediateValue はオペランド文字列から即値を取得する
func getImmediateValue(operandStr string, size int) ([]byte, error) {

	_operandStr := strings.TrimSpace(operandStr)
	val, err := strconv.ParseInt(_operandStr, 0, 64) // 0 base allows auto-detection (e.g., 0x prefix)
	if err != nil {
		return nil, fmt.Errorf("failed to parse immediate value '%s': %w", _operandStr, err)
	}

	buf := make([]byte, size)
	switch size {
	case 1:
		buf[0] = byte(val)
	case 2:
		binary.LittleEndian.PutUint16(buf, uint16(val))
	case 4:
		binary.LittleEndian.PutUint32(buf, uint32(val))
	default:
		return nil, fmt.Errorf("unsupported immediate size: %d", size)
	}
	return buf, nil
}

package ng_operand // Changed package name

import (
	"strings"

	"github.com/HobbyOSs/gosk/internal/ast" // Import ast package
	// cpu package import removed
)

// OperandType はオペランドの種類を表す型 (既存のものを流用または再定義)
// pkg/operand/operand.go からコピー＆修正
type OperandType string

const (
	CodeM    OperandType = "m"    // Memory address
	CodeM8   OperandType = "m8"   // 8-bit memory
	CodeM16  OperandType = "m16"  // 16-bit memory
	CodeM32  OperandType = "m32"  // 32-bit memory
	CodeM64  OperandType = "m64"  // 64-bit memory (for FPU, MMX, SSE)
	CodeM128 OperandType = "m128" // 128-bit memory (for SSE)
	CodeM256 OperandType = "m256" // 256-bit memory (for AVX)
	CodeM512 OperandType = "m512" // 512-bit memory (for AVX-512)
	CodeMEM  OperandType = "mem"  // Memory address (size determined by other operand)

	CodeR8   OperandType = "r8"    // 8-bit register
	CodeR16  OperandType = "r16"   // 16-bit register
	CodeR32  OperandType = "r32"   // 32-bit register
	CodeR64  OperandType = "r64"   // 64-bit register
	CodeRM8  OperandType = "r/m8"  // 8-bit register or memory
	CodeRM16 OperandType = "r/m16" // 16-bit register or memory
	CodeRM32 OperandType = "r/m32" // 32-bit register or memory
	CodeRM64 OperandType = "r/m64" // 64-bit register or memory

	CodeIMM   OperandType = "imm"   // Immediate value (size determined by other operand)
	CodeIMM8  OperandType = "imm8"  // 8-bit immediate
	CodeIMM16 OperandType = "imm16" // 16-bit immediate
	CodeIMM32 OperandType = "imm32" // 32-bit immediate
	CodeIMM64 OperandType = "imm64" // 64-bit immediate (used for MOV r64, imm64)

	CodeREL8  OperandType = "rel8"  // 8-bit relative offset
	CodeREL16 OperandType = "rel16" // 16-bit relative offset
	CodeREL32 OperandType = "rel32" // 32-bit relative offset

	CodePTR1616 OperandType = "ptr16:16" // 16:16 far pointer
	CodePTR1632 OperandType = "ptr16:32" // 16:32 far pointer

	CodeSREG OperandType = "Sreg" // Segment register (CS, DS, ES, FS, GS, SS)
	CodeCREG OperandType = "Creg" // Control register (CR0-CR8)
	CodeDREG OperandType = "Dreg" // Debug register (DR0-DR7)
	CodeTREG OperandType = "Treg" // Test register (TR3-TR7)

	CodeAL OperandType = "AL" // AL register
	CodeCL OperandType = "CL" // CL register
	CodeDL OperandType = "DL" // DL register
	CodeBL OperandType = "BL" // BL register
	CodeAH OperandType = "AH" // AH register
	CodeCH OperandType = "CH" // CH register
	CodeDH OperandType = "DH" // DH register
	CodeBH OperandType = "BH" // BH register

	CodeAX OperandType = "AX" // AX register
	CodeCX OperandType = "CX" // CX register
	CodeDX OperandType = "DX" // DX register
	CodeBX OperandType = "BX" // BX register
	CodeSP OperandType = "SP" // SP register
	CodeBP OperandType = "BP" // BP register
	CodeSI OperandType = "SI" // SI register
	CodeDI OperandType = "DI" // DI register

	CodeEAX OperandType = "EAX" // EAX register
	CodeECX OperandType = "ECX" // ECX register
	CodeEDX OperandType = "EDX" // EDX register
	CodeEBX OperandType = "EBX" // EBX register
	CodeESP OperandType = "ESP" // ESP register
	CodeEBP OperandType = "EBP" // EBP register
	CodeESI OperandType = "ESI" // ESI register
	CodeEDI OperandType = "EDI" // EDI register

	CodeRAX OperandType = "RAX" // RAX register (64-bit mode)
	// ... other 64-bit registers ...

	CodeES OperandType = "ES" // ES segment register
	CodeCS OperandType = "CS" // CS segment register
	CodeSS OperandType = "SS" // SS segment register
	CodeDS OperandType = "DS" // DS segment register
	CodeFS OperandType = "FS" // FS segment register
	CodeGS OperandType = "GS" // GS segment register

	CodeST0 OperandType = "ST(0)" // FPU register ST(0)
	CodeSTI OperandType = "ST(i)" // FPU register ST(i)

	CodeMM   OperandType = "mm"     // MMX register
	CodeMMRM OperandType = "mm/m64" // MMX register or 64-bit memory

	CodeXMM   OperandType = "xmm"      // XMM register
	CodeXMMRM OperandType = "xmm/m128" // XMM register or 128-bit memory

	CodeYMM   OperandType = "ymm"      // YMM register
	CodeYMMRM OperandType = "ymm/m256" // YMM register or 256-bit memory

	CodeZMM   OperandType = "zmm"      // ZMM register
	CodeZMMRM OperandType = "zmm/m512" // ZMM register or 512-bit memory

	CodeBND OperandType = "bnd" // Bound register (BND0-BND3)

	CodeKREG OperandType = "k" // Mask register (k0-k7)

	CodeMOFFS8  OperandType = "moffs8"  // Memory offset (8-bit)
	CodeMOFFS16 OperandType = "moffs16" // Memory offset (16-bit)
	CodeMOFFS32 OperandType = "moffs32" // Memory offset (32-bit)
	CodeMOFFS64 OperandType = "moffs64" // Memory offset (64-bit)

	CodeCONST1 OperandType = "1" // Constant 1

	CodeLABEL OperandType = "label" // Label (placeholder type)

	CodeUNKNOWN OperandType = "unknown"
)

// パース結果を格納する構造体
type ParsedOperandPeg struct {
	Type      OperandType  // オペランドの種類
	Register  string       // レジスタ名 (Typeがレジスタの場合)
	Immediate int64        // 即値 (Typeが即値の場合)
	IsHex     bool         // 即値が16進数か
	Memory    *MemoryInfo  // メモリアドレス情報 (Typeがメモリの場合)
	Segment   string       // セグメントレジスタ名 (TypeがSREGの場合)
	Label     string       // ラベル名 (TypeがLABELの場合)
	DataType  ast.DataType // データ型 (BYTE, WORD, DWORD) - Changed to ast.DataType
	JumpType  string       // ジャンプタイプ (SHORT, NEAR, FAR)
	RawString string       // 元のオペランド文字列
}

// メモリアドレスの詳細情報
type MemoryInfo struct {
	BaseReg      string // ベースレジスタ
	IndexReg     string // インデックスレジスタ
	Scale        int    // スケール (1, 2, 4, 8)
	Displacement int64  // ディスプレースメント
	IsHexDisp    bool   // ディスプレースメントが16進数か
	Segment      string // セグメントオーバーライド (例: "ES")
}

// --- Helper Functions for PEG Actions ---

// getRegisterType はレジスタ名から OperandType を返す (既存実装を参考に修正)
func getRegisterType(reg string) OperandType {
	regUpper := strings.ToUpper(reg)
	switch regUpper {
	case "AL", "CL", "DL", "BL", "AH", "CH", "DH", "BH":
		return CodeR8
	case "AX", "CX", "DX", "BX", "SP", "BP", "SI", "DI":
		return CodeR16
	case "EAX", "ECX", "EDX", "EBX", "ESP", "EBP", "ESI", "EDI":
		return CodeR32
	case "RAX", "RCX", "RDX", "RBX", "RSP", "RBP", "RSI", "RDI",
		"R8", "R9", "R10", "R11", "R12", "R13", "R14", "R15": // 64bit registers
		return CodeR64
	case "CS", "DS", "ES", "FS", "GS", "SS":
		return CodeSREG
	case "CR0", "CR1", "CR2", "CR3", "CR4", "CR5", "CR6", "CR7", "CR8":
		return CodeCREG
	case "DR0", "DR1", "DR2", "DR3", "DR4", "DR5", "DR6", "DR7":
		return CodeDREG
	case "TR0", "TR1", "TR2", "TR3", "TR4", "TR5", "TR6", "TR7":
		return CodeTREG
	// MMX, XMM, YMM, ZMM, BND, K registers...
	default:
		if strings.HasPrefix(regUpper, "MM") {
			return CodeMM
		}
		if strings.HasPrefix(regUpper, "XMM") {
			return CodeXMM
		}
		if strings.HasPrefix(regUpper, "YMM") {
			return CodeYMM
		}
		if strings.HasPrefix(regUpper, "ZMM") {
			return CodeZMM
		}
		if strings.HasPrefix(regUpper, "BND") {
			return CodeBND
		}
		if strings.HasPrefix(regUpper, "K") {
			return CodeKREG
		} // Assuming k0-k7
		return CodeUNKNOWN
	}
}

// getImmediateSizeType は即値から OperandType を返す (新規実装)
func getImmediateSizeType(value int64) OperandType {
	// 符号付き8bitで表現可能か
	if value >= -128 && value <= 127 {
		return CodeIMM8
	}
	// 符号付き16bitで表現可能か
	if value >= -32768 && value <= 32767 {
		return CodeIMM16
	}
	// 符号付き32bitで表現可能か
	if value >= -2147483648 && value <= 2147483647 {
		return CodeIMM32
	}
	// それ以外はIMM64とする (asmdb側で適切なサイズが選択される想定)
	return CodeIMM64
}

// TODO: 必要に応じて pkg/operand/operand_impl.go から他のヘルパー関数
// (例: resolveOperandSizes, calcMemOffsetSize など) を移植・修正する

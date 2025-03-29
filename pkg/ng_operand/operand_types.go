package ng_operand // Changed package name

import (
	"strings"

	"github.com/HobbyOSs/gosk/internal/ast" // Import ast package
	// cpu パッケージのインポートは削除
)

// OperandType はオペランドの種類を表す型 (既存のものを流用または再定義)
// pkg/operand/operand.go からコピー＆修正
type OperandType string

const (
	CodeM    OperandType = "m"    // メモリアドレス
	CodeM8   OperandType = "m8"   // 8ビットメモリ
	CodeM16  OperandType = "m16"  // 16ビットメモリ
	CodeM32  OperandType = "m32"  // 32ビットメモリ
	CodeM64  OperandType = "m64"  // 64ビットメモリ (FPU, MMX, SSE用)
	CodeM128 OperandType = "m128" // 128ビットメモリ (SSE用)
	CodeM256 OperandType = "m256" // 256ビットメモリ (AVX用)
	CodeM512 OperandType = "m512" // 512ビットメモリ (AVX-512用)
	CodeMEM  OperandType = "mem"  // メモリアドレス (サイズは他のオペランドによって決定)

	CodeR8   OperandType = "r8"    // 8ビットレジスタ
	CodeR16  OperandType = "r16"   // 16ビットレジスタ
	CodeR32  OperandType = "r32"   // 32ビットレジスタ
	CodeR64  OperandType = "r64"   // 64ビットレジスタ
	CodeRM8  OperandType = "r/m8"  // 8ビットレジスタまたはメモリ
	CodeRM16 OperandType = "r/m16" // 16ビットレジスタまたはメモリ
	CodeRM32 OperandType = "r/m32" // 32ビットレジスタまたはメモリ
	CodeRM64 OperandType = "r/m64" // 64ビットレジスタまたはメモリ

	CodeIMM   OperandType = "imm"   // 即値 (サイズは他のオペランドによって決定)
	CodeIMM8  OperandType = "imm8"  // 8ビット即値
	CodeIMM16 OperandType = "imm16" // 16ビット即値
	CodeIMM32 OperandType = "imm32" // 32ビット即値
	CodeIMM64 OperandType = "imm64" // 64ビット即値 (MOV r64, imm64 で使用)

	CodeREL8  OperandType = "rel8"  // 8ビット相対オフセット
	CodeREL16 OperandType = "rel16" // 16ビット相対オフセット
	CodeREL32 OperandType = "rel32" // 32ビット相対オフセット

	CodePTR1616 OperandType = "ptr16:16" // 16:16 ファーポインタ
	CodePTR1632 OperandType = "ptr16:32" // 16:32 ファーポインタ

	CodeSREG OperandType = "sreg" // セグメントレジスタ (CS, DS, ES, FS, GS, SS) - 小文字に変更
	CodeCREG OperandType = "creg" // コントロールレジスタ (CR0-CR8) - 小文字に変更
	CodeDREG OperandType = "dreg" // デバッグレジスタ (DR0-DR7) - 小文字に変更
	CodeTREG OperandType = "treg" // テストレジスタ (TR3-TR7) - 小文字に変更

	CodeAL OperandType = "AL" // AL レジスタ
	CodeCL OperandType = "CL" // CL レジスタ
	CodeDL OperandType = "DL" // DL レジスタ
	CodeBL OperandType = "BL" // BL レジスタ
	CodeAH OperandType = "AH" // AH レジスタ
	CodeCH OperandType = "CH" // CH レジスタ
	CodeDH OperandType = "DH" // DH レジスタ
	CodeBH OperandType = "BH" // BH レジスタ

	CodeAX OperandType = "AX" // AX レジスタ
	CodeCX OperandType = "CX" // CX レジスタ
	CodeDX OperandType = "DX" // DX レジスタ
	CodeBX OperandType = "BX" // BX レジスタ
	CodeSP OperandType = "SP" // SP レジスタ
	CodeBP OperandType = "BP" // BP レジスタ
	CodeSI OperandType = "SI" // SI レジスタ
	CodeDI OperandType = "DI" // DI レジスタ

	CodeEAX OperandType = "EAX" // EAX レジスタ
	CodeECX OperandType = "ECX" // ECX レジスタ
	CodeEDX OperandType = "EDX" // EDX レジスタ
	CodeEBX OperandType = "EBX" // EBX レジスタ
	CodeESP OperandType = "ESP" // ESP レジスタ
	CodeEBP OperandType = "EBP" // EBP レジスタ
	CodeESI OperandType = "ESI" // ESI レジスタ
	CodeEDI OperandType = "EDI" // EDI レジスタ

	CodeRAX OperandType = "RAX" // RAX レジスタ (64ビットモード)
	// ... その他の64ビットレジスタ ...

	CodeES OperandType = "ES" // ES セグメントレジスタ
	CodeCS OperandType = "CS" // CS セグメントレジスタ
	CodeSS OperandType = "SS" // SS セグメントレジスタ
	CodeDS OperandType = "DS" // DS セグメントレジスタ
	CodeFS OperandType = "FS" // FS セグメントレジスタ
	CodeGS OperandType = "GS" // GS セグメントレジスタ

	CodeST0 OperandType = "ST(0)" // FPU レジスタ ST(0)
	CodeSTI OperandType = "ST(i)" // FPU レジスタ ST(i)

	CodeMM   OperandType = "mm"     // MMX レジスタ
	CodeMMRM OperandType = "mm/m64" // MMX レジスタまたは64ビットメモリ

	CodeXMM   OperandType = "xmm"      // XMM レジスタ
	CodeXMMRM OperandType = "xmm/m128" // XMM レジスタまたは128ビットメモリ

	CodeYMM   OperandType = "ymm"      // YMM レジスタ
	CodeYMMRM OperandType = "ymm/m256" // YMM レジスタまたは256ビットメモリ

	CodeZMM   OperandType = "zmm"      // ZMM レジスタ
	CodeZMMRM OperandType = "zmm/m512" // ZMM レジスタまたは512ビットメモリ

	CodeBND OperandType = "bnd" // バウンドレジスタ (BND0-BND3)

	CodeKREG OperandType = "k" // マスクレジスタ (k0-k7)

	CodeMOFFS8  OperandType = "moffs8"  // メモリオフセット (8ビット)
	CodeMOFFS16 OperandType = "moffs16" // メモリオフセット (16ビット)
	CodeMOFFS32 OperandType = "moffs32" // メモリオフセット (32ビット)
	CodeMOFFS64 OperandType = "moffs64" // メモリオフセット (64ビット)

	CodeCONST1 OperandType = "1" // 定数 1

	CodeLABEL OperandType = "label" // ラベル (プレースホルダー型)

	CodeUNKNOWN OperandType = "unknown" // 不明
)

// ParsedOperandPeg は、PEGパーサーによってパースされた単一のオペランド情報を格納する構造体です。
type ParsedOperandPeg struct {
	Type      OperandType  // オペランドの種類 (例: CodeR32, CodeM, CodeIMM8)
	Register  string       // レジスタ名 (Typeがレジスタの場合)
	Immediate int64        // 即値 (Typeが即値の場合)
	IsHex     bool         // 即値が16進数表記だったか
	Memory    *MemoryInfo  // メモリアドレス情報 (Typeがメモリの場合)
	Segment   string       // セグメントレジスタ名 (TypeがSREGの場合、またはセグメントオーバーライドがある場合)
	Label     string       // ラベル名 (TypeがLABELの場合)
	DataType  ast.DataType // データ型 (BYTE, WORD, DWORD など) - ast.DataType に変更
	JumpType  string       // ジャンプタイプ (SHORT, NEAR, FAR) - ラベルオペランド用
	PtrPrefix string       // PTR または FAR PTR プレフィックス (例: "WORD PTR")
	RawString string       // パース前の元のオペランド文字列
}

// MemoryInfo は、メモリオペランドのアドレッシングに関する詳細情報を格納します。
type MemoryInfo struct {
	BaseReg      string // ベースレジスタ (例: "EAX", "BX")
	IndexReg     string // インデックスレジスタ (例: "ESI", "DI")
	Scale        int    // スケールファクタ (1, 2, 4, 8) - インデックスレジスタに適用
	Displacement int64  // ディスプレースメント (オフセット値)
	IsHexDisp    bool   // ディスプレースメントが16進数表記だったか
	DispLabel    string // ディスプレースメントがラベルの場合のラベル名
	Segment      string // セグメントオーバーライドプレフィックス (例: "ES", "CS:")
}

// --- PEG アクション用ヘルパー関数 ---

// getRegisterType はレジスタ名から OperandType を返します (既存実装を参考に修正)。
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
		"R8", "R9", "R10", "R11", "R12", "R13", "R14", "R15": // 64ビットレジスタ
		return CodeR64
	case "CS", "DS", "ES", "FS", "GS", "SS":
		return CodeSREG // "sreg" を返すように修正
	case "CR0", "CR1", "CR2", "CR3", "CR4", "CR5", "CR6", "CR7", "CR8":
		return CodeCREG // "creg" を返すように修正
	case "DR0", "DR1", "DR2", "DR3", "DR4", "DR5", "DR6", "DR7":
		return CodeDREG // "dreg" を返すように修正
	case "TR0", "TR1", "TR2", "TR3", "TR4", "TR5", "TR6", "TR7":
		return CodeTREG // "treg" を返すように修正
	// MMX, XMM, YMM, ZMM, BND, K レジスタなど...
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
		} // k0-k7 と仮定
		return CodeUNKNOWN
	}
}

// getImmediateSizeType は即値から OperandType を返します (新規実装)。
// 値が収まる最小の符号付きビット幅 (8, 16, 32) を判断し、対応する型を返します。
// 32ビットを超える場合は IMM64 を返します。
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

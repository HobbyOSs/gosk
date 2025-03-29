package ng_operand

import (
	"log"
	"strings"

	"github.com/HobbyOSs/gosk/internal/ast" // Import ast for DataType
	"github.com/HobbyOSs/gosk/pkg/cpu"
)

// OperandPegImpl は Operands インターフェースの peg パーサー実装です。
// 複数のパース済みオペランドを保持します。
type OperandPegImpl struct {
	parsedOperands []*ParsedOperandPeg // Changed from parsed *ParsedOperandPeg
	bitMode        cpu.BitMode
	forceImm8      bool
	forceRelAsImm  bool
}

// NewOperandPegImpl はパース済みのオペランドスライスから OperandPegImpl を生成します。
func NewOperandPegImpl(parsedOperands []*ParsedOperandPeg) *OperandPegImpl { // Argument type changed
	return &OperandPegImpl{
		parsedOperands: parsedOperands, // Use the argument
		bitMode:        cpu.MODE_16BIT, // Default to 16-bit mode
	}
}

// FromString はオペランド文字列をパースして Operands インターフェースを返します。
// これが外部から呼び出される主要なコンストラクタとなります。
// 内部で ParseOperands を呼び出します。
func FromString(text string) (Operands, error) {
	// デフォルトのビットモードやフラグでパース
	// Use ParseOperands instead of ParseOperandString
	parsedOperands, err := ParseOperands(text, cpu.MODE_16BIT, false, false)
	if err != nil {
		log.Printf("Error parsing operands string '%s': %v", text, err)
		return nil, err
	}
	impl := NewOperandPegImpl(parsedOperands) // Pass the slice
	return impl, nil
}

// InternalString は最初のオペランドの文字列表現を返します。
// 互換性のために残しますが、複数オペランドの場合は InternalStrings を使用すべきです。
func (o *OperandPegImpl) InternalString() string {
	if len(o.parsedOperands) == 0 || o.parsedOperands[0] == nil { // Check slice length and first element
		return ""
	}
	return o.parsedOperands[0].RawString // Access first element
}

// InternalStrings は各オペランドの文字列表現をスライスとして返します。
func (o *OperandPegImpl) InternalStrings() []string {
	if len(o.parsedOperands) == 0 { // Check slice length
		return []string{}
	}
	strs := make([]string, 0, len(o.parsedOperands))
	for _, p := range o.parsedOperands { // Iterate over the slice
		if p != nil {
			strs = append(strs, p.RawString)
		} else {
			strs = append(strs, "") // nil の場合は空文字列
		}
	}
	return strs
}

// OperandTypes は各オペランドの型をスライスとして返します。
// サイズ解決ロジックを含みます。
func (o *OperandPegImpl) OperandTypes() []OperandType {
	if len(o.parsedOperands) == 0 { // Check slice length
		return []OperandType{}
	}

	types := make([]OperandType, len(o.parsedOperands))
	for i, parsed := range o.parsedOperands { // Iterate over the slice
		if parsed == nil {
			types[i] = CodeUNKNOWN
			continue
		}

		baseType := parsed.Type

		// --- ここから型解決ロジック ---

		// 1. ラベルの型解決
		if baseType == CodeLABEL {
			if o.forceRelAsImm {
				// 強制的に IMM とする (例: MOV EAX, label)
				// サイズはビットモードや他のオペランドに依存する可能性がある
				// TODO: より正確な IMM サイズ決定
				if o.bitMode == cpu.MODE_16BIT {
					types[i] = CodeIMM16
				} else {
					types[i] = CodeIMM32
				}
			} else {
				// JMP/CALL label の場合 (デフォルト)
				if parsed.JumpType == "SHORT" {
					types[i] = CodeREL8
				} else {
					// デフォルトの相対ジャンプサイズはビットモードに依存
					if o.bitMode == cpu.MODE_16BIT {
						types[i] = CodeREL16
					} else {
						types[i] = CodeREL32
					}
				}
			}
			continue // ラベルはここで解決完了
		}

		// 2. メモリサイズの解決 (DataType指定がある場合)
		if baseType == CodeM && parsed.DataType != ast.None {
			switch parsed.DataType {
			case ast.Byte:
				types[i] = CodeM8
			case ast.Word:
				types[i] = CodeM16
			case ast.Dword:
				types[i] = CodeM32
			// TODO: QWORD etc.
			default:
				// 不明な DataType の場合はフォールバック
				types[i] = o.resolveMemorySize(parsed, i) // ヘルパー関数を呼ぶ
			}
			continue // メモリ(DataType指定あり)はここで解決完了
		}

		// 3. 即値サイズの解決 (forceImm8 を優先)
		if baseType == CodeIMM || baseType == CodeIMM8 || baseType == CodeIMM16 || baseType == CodeIMM32 || baseType == CodeIMM64 {
			if o.forceImm8 {
				types[i] = CodeIMM8 // forceImm8 が true なら IMM8 に確定
			} else {
				// peg パーサーがサイズを特定している場合はそれを使う
				if baseType == CodeIMM8 || baseType == CodeIMM16 || baseType == CodeIMM32 || baseType == CodeIMM64 {
					types[i] = baseType
				} else {
					// CodeIMM の場合、値からサイズを推定
					// 基本は値から推定したサイズとする
					// resolveDependentSizes で必要に応じて調整される
					types[i] = getImmediateSizeType(parsed.Immediate)
				}
			}
			continue // 即値はここで解決完了
		}

		// 4. メモリサイズの解決 (DataType指定がない場合)
		if baseType == CodeM {
			types[i] = o.resolveMemorySize(parsed, i) // ヘルパー関数を呼ぶ
			continue                                  // メモリ(DataType指定なし)はここで解決完了
		}

		// 5. その他の型 (レジスタなど)
		// セグメントオーバーライド付きのレジスタはメモリアクセスとして扱う
		if isRegisterType(baseType) && parsed.Segment != "" {
			// 例: DS:AX は WORD PTR DS:[addr] 相当だが、アドレスは不明
			// resolveMemorySize を呼んでサイズを決定する (ビットモード依存になる)
			types[i] = o.resolveMemorySize(parsed, i)
		} else {
			// それ以外はパーサーの基本タイプをそのまま使う
			types[i] = baseType
		}
	}

	// --- オペランド間の依存関係によるサイズ解決 ---
	// 例: MOV [EBX], AL -> M8
	// 例: MOV EAX, [EBX] -> M32
	// 例: ADD EAX, 1 -> IMM32
	o.resolveDependentSizes(types) // types スライスを直接変更

	return types
}

// resolveMemorySize は DataType 指定がないメモリオペランドのサイズを解決します。
// (OperandTypes 内で呼び出されるヘルパー)
func (o *OperandPegImpl) resolveMemorySize(parsed *ParsedOperandPeg, index int) OperandType {
	// 優先度1: 他のオペランドがレジスタなら、そのサイズに合わせる
	otherRegType := CodeUNKNOWN
	for j, otherParsed := range o.parsedOperands {
		if index == j || otherParsed == nil { // Use index instead of i
			continue
		}
		// Use helper function to check if it's a register type
		if isRegisterType(otherParsed.Type) {
			otherRegType = otherParsed.Type
			break
		}
	}

	if otherRegType != CodeUNKNOWN {
		// Ensure specific register types map correctly, including AL, CL etc.
		switch otherRegType {
		case CodeR8, CodeAL, CodeCL, CodeDL, CodeBL, CodeAH, CodeCH, CodeDH, CodeBH:
			return CodeM8
		case CodeR16, CodeAX, CodeCX, CodeDX, CodeBX, CodeSP, CodeBP, CodeSI, CodeDI:
			return CodeM16
		case CodeR32, CodeEAX, CodeECX, CodeEDX, CodeEBX, CodeESP, CodeEBP, CodeESI, CodeEDI:
			return CodeM32
		case CodeR64: // TODO: 64bit対応
			return CodeM64
		}
	}

	// 優先度2: メモリアドレスのレジスタから推測 (ビットモードも考慮)
	if parsed.Memory != nil {
		mem := parsed.Memory
		hasEprefix := strings.HasPrefix(mem.BaseReg, "E") || strings.HasPrefix(mem.IndexReg, "E")
		// SP/BP はビットモードによってサイズが変わる点に注意
		isESP := mem.BaseReg == "ESP" || mem.IndexReg == "ESP"
		isEBP := mem.BaseReg == "EBP" || mem.IndexReg == "EBP"
		isSP := mem.BaseReg == "SP" || mem.IndexReg == "SP"
		isBP := mem.BaseReg == "BP" || mem.IndexReg == "BP"
		hasOther16bitReg := mem.BaseReg == "BX" || mem.BaseReg == "SI" || mem.BaseReg == "DI" ||
			mem.IndexReg == "SI" || mem.IndexReg == "DI"

		// Eプレフィックスを持つレジスタ (ESP/EBP含む) があれば M32
		if hasEprefix || isESP || isEBP {
			return CodeM32
		}
		// 16bitモードで SP/BP または他の16bitレジスタがあれば M16
		if o.bitMode == cpu.MODE_16BIT && (isSP || isBP || hasOther16bitReg) {
			return CodeM16
		}
		// 32bitモードで SP/BP または他の16bitレジスタがあれば M32 (ESP/EBPとして扱われるため)
		if o.bitMode == cpu.MODE_32BIT && (isSP || isBP || hasOther16bitReg) {
			return CodeM32
		}
		// レジスタ指定がない場合 (例: [0x1234], FAR PTR [0x5678]) は優先度3へ
	}

	// 優先度3: ビットモードからデフォルトサイズを決定 (上記で解決できなかった場合)
	if o.bitMode == cpu.MODE_16BIT {
		return CodeM16 // 16bit モードのデフォルトは M16
	}
	return CodeM32 // 32bit モード (およびそれ以外) のデフォルトは M32
}

// needsResolution はオペランドタイプがサイズ解決を必要とするか判定します。
func needsResolution(opType OperandType) bool {
	return opType == CodeM || opType == CodeIMM || opType == CodeIMM8 || opType == CodeIMM16
}

// resolveDependentSizes はオペランド間の依存関係に基づいてサイズを解決します。
// (OperandTypes の最後に呼び出されるヘルパー)
// types スライスを直接変更します。
func (o *OperandPegImpl) resolveDependentSizes(types []OperandType) {
	if len(types) < 2 {
		return // 依存関係はオペランドが2つ以上の場合に発生
	}

	// パターン1: 片方がレジスタ、もう片方がサイズ未定のメモリ/即値
	// Iterate through all pairs (handles more than 2 operands if needed, though logic assumes 2 for now)
	for i := 0; i < len(types); i++ {
		for j := 0; j < len(types); j++ {
			if i == j {
				continue
			}

			regType := CodeUNKNOWN
			targetIndex := -1
			targetType := CodeUNKNOWN

			// Check if types[i] is register and types[j] needs resolution
			if isRegisterType(types[i]) && needsResolution(types[j]) {
				regType = types[i]
				targetIndex = j
				targetType = types[j] // Store the original unresolved type
			} else if isRegisterType(types[j]) && needsResolution(types[i]) { // Check the other way around
				regType = types[j]
				targetIndex = i
				targetType = types[i]
			}

			if targetIndex != -1 {
				resolvedType := CodeUNKNOWN
				// Resolve based on the register type
				switch regType {
				case CodeR8, CodeAL, CodeCL, CodeDL, CodeBL, CodeAH, CodeCH, CodeDH, CodeBH:
					resolvedType = map[OperandType]OperandType{CodeM: CodeM8, CodeIMM: CodeIMM8, CodeIMM8: CodeIMM8, CodeIMM16: CodeIMM8}[targetType] // IMM16 -> IMM8 if reg is R8
				case CodeR16, CodeAX, CodeCX, CodeDX, CodeBX, CodeSP, CodeBP, CodeSI, CodeDI:
					resolvedType = map[OperandType]OperandType{CodeM: CodeM16, CodeIMM: CodeIMM16, CodeIMM8: CodeIMM16, CodeIMM16: CodeIMM16}[targetType] // IMM8 -> IMM16 if reg is R16
				case CodeR32, CodeEAX, CodeECX, CodeEDX, CodeEBX, CodeESP, CodeEBP, CodeESI, CodeEDI:
					resolvedType = map[OperandType]OperandType{CodeM: CodeM32, CodeIMM: CodeIMM32, CodeIMM8: CodeIMM32, CodeIMM16: CodeIMM32}[targetType] // IMM8/16 -> IMM32 if reg is R32
				case CodeR64:
					resolvedType = map[OperandType]OperandType{CodeM: CodeM64, CodeIMM: CodeIMM64, CodeIMM8: CodeIMM64, CodeIMM16: CodeIMM64, CodeIMM32: CodeIMM64}[targetType] // IMM8/16/32 -> IMM64 if reg is R64
				}

				// Update the type in the slice if resolved
				if resolvedType != CodeUNKNOWN && resolvedType != targetType { // Only update if changed
					types[targetIndex] = resolvedType
				}
			}
		}
	}

	// パターン2: 単独の即値オペランドは IMM32 とする (forceImm8 でない場合)
	if len(types) == 1 && !o.forceImm8 && (types[0] == CodeIMM || types[0] == CodeIMM8 || types[0] == CodeIMM16) {
		types[0] = CodeIMM32
	}

	// パターン3: ADD/SUB など、即値が常にレジスタサイズに拡張される場合 (簡易ハック) - 修正: R32だけでなくR16も考慮
	// TODO: 命令コンテキストが必要
	var hasR32 bool = false
	var hasR16 bool = false // R16 も考慮
	var immIndex int = -1
	var immType OperandType = CodeUNKNOWN
	for k, typ := range types {
		if isR32Type(typ) { // Use helper
			hasR32 = true
		} else if isR16Type(typ) { // Use new helper
			hasR16 = true
		}
		if typ == CodeIMM8 || typ == CodeIMM16 { // If small immediate exists
			immIndex = k
			immType = typ // Store the original small immediate type
		}
	}
	// Apply hack if a register (R32 or R16) and a *small* immediate are present
	if immIndex != -1 && (immType == CodeIMM8 || immType == CodeIMM16) {
		otherIndex := 1 - immIndex // Assuming 2 operands
		if otherIndex >= 0 && otherIndex < len(types) {
			if hasR32 && isR32Type(types[otherIndex]) {
				// If R32 and IMM8/16 -> IMM32
				types[immIndex] = CodeIMM32
			} else if hasR16 && isR16Type(types[otherIndex]) && immType == CodeIMM8 {
				// If R16 and IMM8 -> IMM16
				types[immIndex] = CodeIMM16
			}
		}
	}
}

// isR32Type は指定された型が32ビット汎用レジスタ型かどうかを判定します。
func isR32Type(opType OperandType) bool {
	switch opType {
	case CodeR32, CodeEAX, CodeECX, CodeEDX, CodeEBX, CodeESP, CodeEBP, CodeESI, CodeEDI:
		return true
	default:
		return false
	}
}

// isR16Type は指定された型が16ビット汎用レジスタ型かどうかを判定します。
func isR16Type(opType OperandType) bool {
	switch opType {
	case CodeR16, CodeAX, CodeCX, CodeDX, CodeBX, CodeSP, CodeBP, CodeSI, CodeDI:
		return true
	default:
		return false
	}
}

// isRegisterType は指定された型がレジスタ型かどうかを判定します。
func isRegisterType(opType OperandType) bool {
	switch opType {
	case CodeR8, CodeR16, CodeR32, CodeR64,
		CodeAL, CodeCL, CodeDL, CodeBL, CodeAH, CodeCH, CodeDH, CodeBH,
		CodeAX, CodeCX, CodeDX, CodeBX, CodeSP, CodeBP, CodeSI, CodeDI,
		CodeEAX, CodeECX, CodeEDX, CodeEBX, CodeESP, CodeEBP, CodeESI, CodeEDI,
		CodeSREG, CodeCREG, CodeDREG, CodeTREG, CodeMM, CodeXMM, CodeYMM:
		return true
	default:
		return false
	}
}

// Serialize はオペランドをシリアライズ可能な文字列として返します。
// カンマ区切りで結合します。
func (o *OperandPegImpl) Serialize() string {
	return strings.Join(o.InternalStrings(), ", ") // Use updated InternalStrings
}

// FromString は OperandPegImpl 自身には適用せず、パッケージレベルの FromString を使用します。
// インターフェースを満たすためにダミー実装を提供します。
func (o *OperandPegImpl) FromString(text string) Operands {
	// このメソッドは通常使われない想定
	newOp, _ := FromString(text) // エラーは無視
	return newOp
}

// CalcOffsetByteSize はメモリオペランドのオフセットサイズを計算します。
// TODO: 複数オペランド対応。どのオペランドのサイズを計算するか？
//
//	現状は最初のメモリオペランドを見る。
//
// TODO: Displacement の値に基づいてサイズ (1, 2, 4) を返す
func (o *OperandPegImpl) CalcOffsetByteSize() int {
	for _, parsed := range o.parsedOperands { // Iterate over operands
		if parsed != nil && parsed.Memory != nil {
			// Displacement size calculation logic needed here
			// Example (needs refinement based on actual displacement value type):
			// disp := parsed.Memory.DisplacementValue // Assuming this field exists
			// if disp >= -128 && disp <= 127 { return 1 }
			// else if o.bitMode == cpu.MODE_16BIT { return 2 } // Or based on Require67h?
			// else { return 4 }

			// Placeholder logic based on address size prefix requirement
			// Use the method defined in requires.go (now part of OperandPegImpl)
			if o.Require67h() { // If address size prefix is needed
				if o.bitMode == cpu.MODE_16BIT {
					return 4 // 16-bit mode with 32-bit addressing -> 4 bytes
				} else {
					return 2 // 32-bit mode with 16-bit addressing -> 2 bytes
				}
			} else { // No prefix needed, use default for the mode
				if o.bitMode == cpu.MODE_16BIT {
					return 2 // Default 16-bit offset
				} else {
					return 4 // Default 32-bit offset
				}
			}
		}
	}
	return 0 // No memory operand found
}

// DetectImmediateSize は即値オペランドのサイズ (バイト数) を検出します。
// TODO: 複数オペランド対応。どのオペランドのサイズを検出するか？
//
//	現状は最初の即値オペランドを見る。
func (o *OperandPegImpl) DetectImmediateSize() int {
	opTypes := o.OperandTypes()               // Get resolved types
	for i, parsed := range o.parsedOperands { // Iterate over operands
		if parsed == nil {
			continue
		}
		opType := opTypes[i] // Use the resolved type
		switch opType {
		case CodeIMM8:
			return 1
		case CodeIMM16:
			return 2
		case CodeIMM32:
			return 4
		case CodeIMM64:
			return 8
		}
	}
	return 0 // No immediate operand found
}

func (o *OperandPegImpl) WithBitMode(mode cpu.BitMode) Operands {
	o.bitMode = mode
	// Re-resolve types might be needed if bitMode changes context, but OperandTypes() does it on demand.
	return o
}

func (o *OperandPegImpl) WithForceImm8(force bool) Operands {
	o.forceImm8 = force
	// Re-resolve types might be needed.
	return o
}

func (o *OperandPegImpl) WithForceRelAsImm(force bool) Operands {
	o.forceRelAsImm = force
	// Re-resolve types might be needed.
	return o
}

func (o *OperandPegImpl) GetBitMode() cpu.BitMode {
	return o.bitMode
}

// Require66h と Require67h は requires.go に移動しました。

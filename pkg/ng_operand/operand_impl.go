package ng_operand

import (
	"encoding/binary" // Add binary package for displacement conversion
	"log"
	"strings"

	"github.com/HobbyOSs/gosk/internal/ast"
	"github.com/HobbyOSs/gosk/pkg/cpu"
	"github.com/samber/lo" // samber/lo をインポート
)

// OperandPegImpl は PEG パーサーを使用して Operands インターフェースを実装します。
// 複数のパース済みオペランドを保持します。
type OperandPegImpl struct {
	parsedOperands []*ParsedOperandPeg
	bitMode        cpu.BitMode
	forceImm8      bool
	forceRelAsImm  bool
}

// NewOperandPegImpl は、パース済みオペランドのスライスから新しい OperandPegImpl を作成します。
func NewOperandPegImpl(parsedOperands []*ParsedOperandPeg) *OperandPegImpl {
	return &OperandPegImpl{
		parsedOperands: parsedOperands,
		bitMode:        cpu.MODE_16BIT, // デフォルトは16ビットモード
	}
}

// FromString はオペランド文字列をパースし、Operands インターフェースを返します。
// これは外部から呼び出される主要なコンストラクタです。
// 内部的に ParseOperands を呼び出します。
func FromString(text string) (Operands, error) {
	// デフォルトのビットモードとフラグでパース
	parsedOperands, err := ParseOperands(text, cpu.MODE_16BIT, false, false)
	if err != nil {
		log.Printf("オペランド文字列 '%s' のパースエラー: %v", text, err)
		return nil, err
	}
	impl := NewOperandPegImpl(parsedOperands)
	return impl, nil
}

// InternalString は最初のオペランドの文字列表現を返します。
// 互換性のために保持されていますが、複数のオペランドには InternalStrings を使用すべきです。
func (o *OperandPegImpl) InternalString() string {
	if len(o.parsedOperands) == 0 || o.parsedOperands[0] == nil {
		return ""
	}
	return o.parsedOperands[0].RawString
}

// InternalStrings は各オペランドの文字列表現をスライスとして返します。
func (o *OperandPegImpl) InternalStrings() []string {
	return lo.Map(o.parsedOperands, func(p *ParsedOperandPeg, _ int) string {
		if p != nil {
			return p.RawString
		}
		return ""
	})
}

// OperandTypes は各オペランドの型をスライスとして返します。サイズ解決ロジックを含みます。
func (o *OperandPegImpl) OperandTypes() []OperandType {
	if len(o.parsedOperands) == 0 {
		return []OperandType{}
	}

	// 個々のオペランドのパースに基づく初期の型解決
	types := lo.Map(o.parsedOperands, func(parsed *ParsedOperandPeg, i int) OperandType {
		if parsed == nil {
			return CodeUNKNOWN
		}

		// forceImm8 を優先
		if o.forceImm8 && (parsed.Type == CodeIMM || parsed.Type == CodeIMM8 || parsed.Type == CodeIMM16 || parsed.Type == CodeIMM32 || parsed.Type == CodeIMM64) {
			return CodeIMM8
		}

		baseType := parsed.Type

		// 1. ラベル型の解決
		if baseType == CodeLABEL {
			if o.forceRelAsImm {
				// IMM を強制 (例: MOV EAX, label)
				// サイズはビットモードや他のオペランドに依存する可能性がある
				if o.bitMode == cpu.MODE_16BIT {
					return CodeIMM16
				}
				return CodeIMM32
			}
			// JMP/CALL ラベル (デフォルト)
			if parsed.JumpType == "SHORT" {
				return CodeREL8
			}
			// デフォルトの相対ジャンプサイズはビットモードに依存
			if o.bitMode == cpu.MODE_16BIT {
				return CodeREL16
			}
			return CodeREL32
		}

		// 2. メモリサイズの解決 (DataType 指定あり)
		if baseType == CodeM && parsed.DataType != ast.None {
			switch parsed.DataType {
			case ast.Byte:
				return CodeM8
			case ast.Word:
				return CodeM16
			case ast.Dword:
				return CodeM32
			// TODO: QWORD など
			default:
				// 不明な DataType のフォールバック
				return o.resolveMemorySize(parsed, i) // ヘルパーを呼び出し
			}
		}

		// 3. 即値サイズの解決 (forceImm8 は処理済み)
		if baseType == CodeIMM || baseType == CodeIMM8 || baseType == CodeIMM16 || baseType == CodeIMM32 || baseType == CodeIMM64 {
			// PEG パーサーがサイズを特定した場合、それを使用
			if baseType == CodeIMM8 || baseType == CodeIMM16 || baseType == CodeIMM32 || baseType == CodeIMM64 {
				return baseType
			}
			// CodeIMM の場合、値からサイズを推定 (必要に応じて resolveDependentSizes で調整される)
			return getImmediateSizeType(parsed.Immediate)
		}

		// 4. メモリサイズの解決 (DataType 指定なし)
		if baseType == CodeM {
			return o.resolveMemorySize(parsed, i) // ヘルパーを呼び出し
		}

		// 5. その他の型 (レジスタなど)
		// セグメントオーバーライド付きのレジスタはメモリアクセスとして扱う
		if isRegisterType(baseType) && parsed.Segment != "" {
			// 例: DS:AX は WORD PTR DS:[addr] のようなものだが、アドレスは不明
			// resolveMemorySize を呼び出してサイズを決定 (ビットモード依存になる)
			return o.resolveMemorySize(parsed, i)
		}
		// それ以外の場合、パーサーからの基本型を使用
		return baseType
	})

	// オペランド間の依存関係に基づいてサイズを解決
	o.resolveDependentSizes(types) // types スライスを直接変更

	return types
}

// resolveMemorySize は、明示的な DataType なしのメモリオペランドのサイズを解決します。
// (OperandTypes 内で呼び出されるヘルパー)
func (o *OperandPegImpl) resolveMemorySize(parsed *ParsedOperandPeg, index int) OperandType {
	// 優先度1: 他のオペランドがレジスタの場合、そのサイズに合わせる
	otherReg, _ := lo.Find(o.parsedOperands, func(other *ParsedOperandPeg) bool {
		// 最初に見つかった、nil でなく、自身でなく、レジスタであるオペランドを探す
		return other != nil && other != parsed && isRegisterType(other.Type)
	})

	if otherReg != nil {
		otherRegType := otherReg.Type
		switch otherRegType {
		case CodeR8, CodeAL, CodeCL, CodeDL, CodeBL, CodeAH, CodeCH, CodeDH, CodeBH:
			return CodeM8
		case CodeR16, CodeAX, CodeCX, CodeDX, CodeBX, CodeSP, CodeBP, CodeSI, CodeDI:
			return CodeM16
		case CodeR32, CodeEAX, CodeECX, CodeEDX, CodeEBX, CodeESP, CodeEBP, CodeESI, CodeEDI:
			return CodeM32
		case CodeR64: // TODO: 64ビットサポート
			return CodeM64
		}
	}

	// 優先度2: メモリアドレス内のレジスタから推測 (ビットモードを考慮)
	if parsed.Memory != nil {
		mem := parsed.Memory
		hasEprefix := strings.HasPrefix(mem.BaseReg, "E") || strings.HasPrefix(mem.IndexReg, "E")
		isESP := mem.BaseReg == "ESP" || mem.IndexReg == "ESP"
		isEBP := mem.BaseReg == "EBP" || mem.IndexReg == "EBP"
		isSP := mem.BaseReg == "SP" || mem.IndexReg == "SP"
		isBP := mem.BaseReg == "BP" || mem.IndexReg == "BP"
		hasOther16bitReg := lo.Contains([]string{"BX", "SI", "DI"}, mem.BaseReg) || lo.Contains([]string{"SI", "DI"}, mem.IndexReg)

		// レジスタに 'E' プレフィックスがある場合 (ESP/EBP を含む) は M32
		if hasEprefix || isESP || isEBP {
			return CodeM32
		}
		// 16ビットモードで SP/BP または他の16ビットレジスタが使用されている場合は M16
		if o.bitMode == cpu.MODE_16BIT && (isSP || isBP || hasOther16bitReg) {
			return CodeM16
		}
		// 32ビットモードで SP/BP または他の16ビットレジスタが使用されている場合は M32 (ESP/EBP として扱われる)
		if o.bitMode == cpu.MODE_32BIT && (isSP || isBP || hasOther16bitReg) {
			return CodeM32
		}
		// レジスタが指定されていない場合 (例: [0x1234]) は優先度3にフォールスルー
	}

	// 優先度3: 上記で解決されなかった場合、ビットモードに基づくデフォルトサイズ
	if o.bitMode == cpu.MODE_16BIT {
		return CodeM16 // 16ビットモードではデフォルト M16
	}
	return CodeM32 // 32ビットモード (およびその他) ではデフォルト M32
}

// resolveDependentSizes は、オペランド間の依存関係に基づいてオペランドサイズを解決します。
// types スライスを直接変更します。
func (o *OperandPegImpl) resolveDependentSizes(types []OperandType) {
	// step1: 一方のオペランドがレジスタで、もう一方がサイズ解決が必要なメモリ/即値
	lo.ForEach(types, func(_ OperandType, i int) {
		lo.ForEach(types, func(_ OperandType, j int) {
			if i == j {
				return // 自己比較をスキップ
			}

			var regType OperandType = CodeUNKNOWN
			var targetIndex int = -1
			var targetType OperandType = CodeUNKNOWN

			// レジスタと解決が必要なターゲットを特定
			if isRegisterType(types[i]) && needsResolution(types[j]) {
				regType = types[i]
				targetIndex = j
				targetType = types[j]
			} else if isRegisterType(types[j]) && needsResolution(types[i]) {
				regType = types[j]
				targetIndex = i
				targetType = types[i]
			}

			// レジスタに基づいてターゲットサイズを解決 (ただし、即値ターゲットに対して forceImm8 が設定されている場合を除く)
			if targetIndex != -1 && !(o.forceImm8 && lo.Contains([]OperandType{CodeIMM, CodeIMM8, CodeIMM16, CodeIMM32, CodeIMM64}, targetType)) {
				resolvedType := CodeUNKNOWN
				switch {
				case isR8Type(regType): // R8 型用のカスタムヘルパー
					resolvedType = map[OperandType]OperandType{CodeM: CodeM8, CodeIMM: CodeIMM8, CodeIMM8: CodeIMM8, CodeIMM16: CodeIMM8}[targetType]
				case isR16Type(regType):
					resolvedType = map[OperandType]OperandType{CodeM: CodeM16, CodeIMM: CodeIMM16, CodeIMM8: CodeIMM16, CodeIMM16: CodeIMM16}[targetType]
				case isR32Type(regType):
					resolvedType = map[OperandType]OperandType{CodeM: CodeM32, CodeIMM: CodeIMM32, CodeIMM8: CodeIMM32, CodeIMM16: CodeIMM32}[targetType]
				case isR64Type(regType): // R64 型用のカスタムヘルパー
					resolvedType = map[OperandType]OperandType{CodeM: CodeM64, CodeIMM: CodeIMM64, CodeIMM8: CodeIMM64, CodeIMM16: CodeIMM64, CodeIMM32: CodeIMM64}[targetType]
				}

				if resolvedType != CodeUNKNOWN && resolvedType != targetType {
					types[targetIndex] = resolvedType
				}
			}
		})
	})

	// step2: 単一の即値オペランドは IMM32 にデフォルト設定 (forceImm8 でない場合)
	if !o.forceImm8 && len(types) == 1 {
		if lo.Contains([]OperandType{CodeIMM, CodeIMM8, CodeIMM16}, types[0]) {
			types[0] = CodeIMM32
		}
	}
}

// Serialize はオペランドをシリアライズ可能な文字列 (カンマ区切り) として返します。
func (o *OperandPegImpl) Serialize() string {
	return strings.Join(o.InternalStrings(), ", ")
}

// FromString はインターフェースを満たすためのダミー実装です。
// 代わりにパッケージレベルの FromString を使用してください。
func (o *OperandPegImpl) FromString(text string) Operands {
	newOp, _ := FromString(text) // エラーを無視
	return newOp
}

// CalcOffsetByteSize はメモリオペランドのオフセットサイズを計算します。
// TODO: 複数のオペランドを処理する。現在は最初のメモリオペランドをチェックしています。
// TODO: ディスプレースメント値 (1, 2, 4) に基づいてサイズを計算する。
func (o *OperandPegImpl) CalcOffsetByteSize() int {
	// メモリオペランドが存在するかどうかを検索し、見つかった場合はそのオペランドを取得します。
	memOperand, found := lo.Find(o.parsedOperands, func(p *ParsedOperandPeg) bool {
		return p != nil && p.Memory != nil
	})

	if found {
		// 注意: found が true でも要素自体が nil の場合、memOperand は nil になり得る。
		// (ただし `p != nil` フィルターでこれは防がれるはず。安全のためチェックを追加)
		if memOperand == nil || memOperand.Memory == nil {
			log.Printf("warn: メモリオペランドフラグが見つかりましたが、memOperand または memOperand.Memory が nil です")
			return 0 // または適切にエラーを処理
		}
		memInfo := memOperand.Memory // 見つかったオペランドから MemoryInfo を取得

		// アドレスサイズプレフィックス(67h)が必要な場合の処理
		if o.Require67h() {
			// TODO: Add logic to calculate displacement size based on memInfo.Displacement
			//       even when address size prefix is present.
			if o.bitMode == cpu.MODE_16BIT {
				return 4 // 16ビットモードで32ビットアドレッシング -> 基本4バイト (要Disp考慮)
			}
			return 2 // 32ビットモードで16ビットアドレッシング -> 基本2バイト (要Disp考慮)
		}

		// プレフィックス不要の場合
		// 1. 直接アドレス指定 [disp] (ModRM mode 00, rm 110 for 16bit or 101 for 32bit)
		if memInfo.BaseReg == "" && memInfo.IndexReg == "" {
			if o.bitMode == cpu.MODE_16BIT {
				return 2 // 16bit direct address offset
			}
			return 4 // 32bit direct address offset
		}

		// 2. 間接アドレス指定でディスプレースメントがある場合 [reg+disp], [reg+reg*scale+disp]
		if memInfo.Displacement != 0 {
			disp := memInfo.Displacement
			// ModRM mode 01 (disp8) or 10 (disp16/32)
			if disp >= -128 && disp <= 127 {
				// Check if ModRM mode 01 is applicable for the registers used
				// (For simplicity here, assume disp8 is possible if it fits)
				return 1 // disp8
			}
			if o.bitMode == cpu.MODE_16BIT {
				return 2 // disp16
			}
			return 4 // disp32
		}

		// 3. 間接アドレス指定でディスプレースメントがない場合 [reg], [reg+reg*scale]
		// ModRM mode 00 (except for [BP] in 16bit which uses mode 01 with disp8=0)
		// or mode 11 (register direct, handled elsewhere)
		// For mode 00, no additional offset bytes are needed.
		// Special case: [BP] in 16-bit mode uses ModRM mode 01 with disp8=0.
		if o.bitMode == cpu.MODE_16BIT && memInfo.BaseReg == "BP" && memInfo.IndexReg == "" && memInfo.Displacement == 0 {
			return 1 // disp8=0 for [BP]
		}
		// Other cases like [BX], [SI], [BX+SI] etc. need no offset bytes with ModRM mode 00.
		return 0
	}
	return 0 // メモリオペランドが見つかりません
}

// DetectImmediateSize は即値オペランドのサイズ (バイト単位) を検出します。
// 基本的に OperandTypes() で解決された型に基づきますが、単一オペランドの場合は
// 値が収まる最小サイズを返すように調整します。
// TODO: 複数の即値オペランドが存在する場合の考慮が必要かもしれません。
func (o *OperandPegImpl) DetectImmediateSize() int {
	// まず OperandTypes() を呼び出して型を解決
	resolvedTypes := o.OperandTypes()

	// 解決された型の中から即値型を探す (最初のもの)
	var immType OperandType = CodeUNKNOWN
	var immIndex int = -1
	for i, t := range resolvedTypes {
		if t == CodeIMM8 || t == CodeIMM16 || t == CodeIMM32 || t == CodeIMM64 {
			immType = t
			immIndex = i
			break
		}
	}

	if immIndex != -1 {
		// 単一オペランドで、解決された型が IMM32 だが、
		// 元の値はより小さいサイズに収まる場合、その最小サイズを返す
		if len(o.parsedOperands) == 1 && immType == CodeIMM32 {
			if immIndex < len(o.parsedOperands) && o.parsedOperands[immIndex] != nil {
				val := o.parsedOperands[immIndex].Immediate
				if val >= -128 && val <= 127 {
					return 1 // 8ビットに収まる
				}
				if val >= -32768 && val <= 32767 {
					return 2 // 16ビットに収まる
				}
				// それ以外（32ビットにしか収まらない）場合はそのまま IMM32 のサイズ 4 を返す
			}
			// parsedOperands が取得できない異常系、フォールバックして4を返す
			return 4
		}

		// 上記以外の場合（複数オペランド、または解決型が IMM32 以外）は
		// 解決された型に基づいてサイズを返す
		switch immType {
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

	return 0 // 即値オペランドが見つからない
}

func (o *OperandPegImpl) WithBitMode(mode cpu.BitMode) Operands {
	o.bitMode = mode
	return o
}

func (o *OperandPegImpl) WithForceImm8(force bool) Operands {
	o.forceImm8 = force
	return o
}

func (o *OperandPegImpl) WithForceRelAsImm(force bool) Operands {
	o.forceRelAsImm = force
	return o
}

func (o *OperandPegImpl) GetBitMode() cpu.BitMode {
	return o.bitMode
}

// IsDirectMemory は、オペランドに直接メモリアドレスが含まれるかどうかを返します。
// 直接アドレスは [displacement] の形式と判断します。
func (o *OperandPegImpl) IsDirectMemory() bool {
	return lo.SomeBy(o.parsedOperands, func(p *ParsedOperandPeg) bool {
		// Memory フィールドが存在し、BaseReg と IndexReg が両方空であること
		return p != nil && p.Memory != nil && p.Memory.BaseReg == "" && p.Memory.IndexReg == ""
	})
}

// IsIndirectMemory は、オペランドに間接メモリアドレスが含まれるかどうかを返します。
// 間接アドレスはレジスタを含む形式 (例: [EAX], [ESI+4]) と判断します。
func (o *OperandPegImpl) IsIndirectMemory() bool {
	return lo.SomeBy(o.parsedOperands, func(p *ParsedOperandPeg) bool {
		// Memory フィールドが存在し、BaseReg または IndexReg の少なくとも一方が空でないこと
		return p != nil && p.Memory != nil && (p.Memory.BaseReg != "" || p.Memory.IndexReg != "")
	})
}

// GetMemoryInfo は、最初のメモリオペランドの詳細情報を返します。見つからない場合は nil と false を返します。
func (o *OperandPegImpl) GetMemoryInfo() (*MemoryInfo, bool) {
	memOperand, found := lo.Find(o.parsedOperands, func(p *ParsedOperandPeg) bool {
		return p != nil && p.Memory != nil
	})
	if found {
		return memOperand.Memory, true
	}
	return nil, false
}

// DisplacementBytes は、最初のメモリオペランドのディスプレースメント部分をバイト列として返します。
// ModRMがない直接アドレス指定 (moffs) の場合に利用することを想定しています。
// メモリオペランドがない場合や、ディスプレースメントがない場合は nil を返します。
// バイト列のサイズは BitMode に基づいて決定されます。
func (o *OperandPegImpl) DisplacementBytes() []byte {
	memInfo, found := o.GetMemoryInfo()
	if !found || memInfo == nil {
		// メモリオペランドがない場合は nil
		return nil
	}

	// Displacement を取得
	disp := memInfo.Displacement
	var dispBytes []byte

	// アドレスサイズ (BitMode) に応じてバイト列に変換 (リトルエンディアン)
	bitMode := o.GetBitMode()
	switch bitMode {
	case cpu.MODE_16BIT: // 定数名を修正
		// MOFFS16
		dispBytes = make([]byte, 2)
		binary.LittleEndian.PutUint16(dispBytes, uint16(disp))
	case cpu.MODE_32BIT: // 定数名を修正
		// MOFFS32
		dispBytes = make([]byte, 4)
		binary.LittleEndian.PutUint32(dispBytes, uint32(disp))
	// case cpu.Bit64: // TODO: MOFFS64 サポート
	default:
		log.Printf("warn: DisplacementBytes でサポートされていないビットモード %v", bitMode)
		return nil // サポートされていないモードでは nil を返す
	}

	return dispBytes
}

// ヘルパー関数 (isR32Type, isR16Type, isRegisterType, needsResolution) は operand_util.go に移動しました。
// ImmediateValueFitsIn8Bits は、即値オペランドの値が8ビットに収まるかどうかを返します。
// 複数の即値がある場合は、最初の即値オペランドをチェックします。
func (o *OperandPegImpl) ImmediateValueFitsIn8Bits() bool {
	// parsedOperands の nil チェックを追加
	if o.parsedOperands == nil {
		log.Printf("warn: ImmediateValueFitsIn8Bits が nil の parsedOperands で呼び出されました")
		return false
	}

	immOperand, found := lo.Find(o.parsedOperands, func(p *ParsedOperandPeg) bool {
		return p != nil && (p.Type == CodeIMM || p.Type == CodeIMM8 || p.Type == CodeIMM16 || p.Type == CodeIMM32 || p.Type == CodeIMM64)
	})

	// found が true でも immOperand の nil チェックを追加 (追加の安全のため)
	if found && immOperand != nil {
		val := immOperand.Immediate
		return val >= -128 && val <= 127
	}

	return false // 即値オペランドが見つからない場合、または immOperand が nil の場合
}

// ヘルパー関数 (isR32Type, isR16Type, isRegisterType, needsResolution) は operand_util.go に移動しました。
// isR8Type と isR64Type も operand_util.go に追加しました。

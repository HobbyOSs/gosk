package ng_operand

import (
	"log"
	"strings"

	"github.com/HobbyOSs/gosk/internal/ast" // Import ast for DataType
	"github.com/HobbyOSs/gosk/pkg/cpu"
)

// OperandPegImpl は Operands インターフェースの peg パーサー実装です。
type OperandPegImpl struct {
	parsed        *ParsedOperandPeg
	bitMode       cpu.BitMode
	forceImm8     bool
	forceRelAsImm bool
}

// NewOperandPegImpl はパース済みのオペランド情報から OperandPegImpl を生成します。
// 通常は FromString 経由で使用されます。
func NewOperandPegImpl(parsed *ParsedOperandPeg) *OperandPegImpl {
	return &OperandPegImpl{
		parsed:  parsed,
		bitMode: cpu.MODE_16BIT, // Default to 16-bit mode
	}
}

// FromString はオペランド文字列をパースして Operands インターフェースを返します。
// これが外部から呼び出される主要なコンストラクタとなります。
func FromString(text string) (Operands, error) {
	parsed, err := ParseOperandString(text)
	if err != nil {
		log.Printf("Error parsing operand string '%s': %v", text, err)
		// エラーが発生した場合でも、最低限の情報を保持したオブジェクトを返すか、
		// nil を返すか検討。ここではエラーを返し、呼び出し元で処理する。
		return nil, err
	}
	impl := NewOperandPegImpl(parsed)
	return impl, nil
}

func (o *OperandPegImpl) InternalString() string {
	if o.parsed == nil {
		return ""
	}
	return o.parsed.RawString
}

// InternalStrings はオペランドを文字列スライスとして返します。
// TODO: カンマ区切りの複数オペランドに対応する。
// peg パーサーは現在単一オペランドのみ対応のため、仮実装。
func (o *OperandPegImpl) InternalStrings() []string {
	if o.InternalString() == "" {
		return []string{}
	}
	// 現状は単一オペランドしかパースできないので、そのまま返す
	return []string{o.InternalString()}
}

// OperandTypes はオペランドの型を返します。
// TODO: 複数オペランド対応、サイズ解決ロジック (resolveOperandSizes 相当)
func (o *OperandPegImpl) OperandTypes() []OperandType {
	if o.parsed == nil {
		return []OperandType{CodeUNKNOWN}
	}

	baseType := o.parsed.Type

	// 1. ラベルの型解決
	if baseType == CodeLABEL {
		if o.forceRelAsImm {
			// 強制的に IMM とする (例: MOV EAX, label)
			// サイズは他のオペランドに依存するため、一旦 CodeIMM とする
			return []OperandType{CodeIMM}
		} else {
			// JMP/CALL label の場合 (デフォルト)
			// TODO: SHORT label の判定が必要 (peg側で JumpType を見る？)
			// 現状は REL32 固定とする
			// TODO: SHORT label の判定
			return []OperandType{CodeREL32} // Assume REL32 for now
		}
	}

	// 2. メモリサイズの解決 (DataType指定がある場合)
	if baseType == CodeM && o.parsed.DataType != ast.None {
		switch o.parsed.DataType {
		case ast.Byte:
			return []OperandType{CodeM8}
		case ast.Word:
			return []OperandType{CodeM16}
		case ast.Dword:
			return []OperandType{CodeM32}
			// TODO: QWORD etc.
		}
	}

	// 3. 即値サイズの解決 (getImmediateSizeTypeの結果を反映)
	//    pegのアクションで IMM8/16/32/64 が設定されているはずなので、
	//    baseType が CodeIMM の場合のみサイズ解決が必要。
	if baseType == CodeIMM {
		// サイズ未確定IMMの場合、デフォルトサイズを返す (仮)
		// 本来は resolveOperandSizes で解決すべき
		if o.bitMode == cpu.MODE_16BIT {
			return []OperandType{CodeIMM16}
		}
		return []OperandType{CodeIMM32}
	}

	// TODO: CodeM のサイズ解決ロジック (resolveOperandSizes 相当)

	// その他の型はそのまま返す
	return []OperandType{baseType}
}

// Serialize はオペランドをシリアライズ可能な文字列として返します。
func (o *OperandPegImpl) Serialize() string {
	// 現状は InternalString と同じ
	return o.InternalString()
}

// FromString は OperandPegImpl 自身には適用せず、パッケージレベルの FromString を使用します。
// インターフェースを満たすためにダミー実装を提供します。
func (o *OperandPegImpl) FromString(text string) Operands {
	// このメソッドは通常使われない想定
	newOp, _ := FromString(text) // エラーは無視
	return newOp
}

// CalcOffsetByteSize はメモリオペランドのオフセットサイズを計算します。
// TODO: 実装する (calcMemOffsetSize 相当)
func (o *OperandPegImpl) CalcOffsetByteSize() int {
	if o.parsed == nil || o.parsed.Memory == nil {
		return 0
	}
	// TODO: Displacement の値に基づいてサイズ (1, 2, 4) を返す
	// 例: disp := o.parsed.Memory.Displacement
	// if disp >= -128 && disp <= 127 { return 1 } ...
	return 4 // 仮実装
}

// DetectImmediateSize は即値オペランドのサイズ (バイト数) を検出します。
// TODO: 実装する
func (o *OperandPegImpl) DetectImmediateSize() int {
	if o.parsed == nil || (o.parsed.Type != CodeIMM8 && o.parsed.Type != CodeIMM16 && o.parsed.Type != CodeIMM32 && o.parsed.Type != CodeIMM64 && o.parsed.Type != CodeIMM) {
		return 0
	}
	// TODO: OperandTypes で解決された後の型を見るべき？
	// parsed.Type が IMM8/16/32/64 ならそれを返す
	switch o.parsed.Type {
	case CodeIMM8:
		return 1
	case CodeIMM16:
		return 2
	case CodeIMM32:
		return 4
	case CodeIMM64:
		return 8 // 64bit 即値は通常 MOV 命令のみ
	case CodeIMM:
		// サイズ未確定の場合、他のオペランドから推測するか、
		// デフォルトサイズ (bitMode依存?) を返す必要がある
		if o.bitMode == cpu.MODE_16BIT {
			return 2
		}
		return 4 // Default to 32-bit for IMM
	}
	return 0
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

// Require66h はオペランドサイズプレフィックス (66h) が必要か判定します。
func (o *OperandPegImpl) Require66h() bool {
	opTypes := o.OperandTypes()
	if len(opTypes) == 0 {
		return false
	}
	// TODO: 複数オペランドの場合の考慮 (現状は最初のオペランドのみ見る)
	opType := opTypes[0]

	is16bitMode := o.bitMode == cpu.MODE_16BIT
	is32bitMode := o.bitMode == cpu.MODE_32BIT

	is16bitOperand := opType == CodeR16 || opType == CodeM16 || opType == CodeIMM16
	is32bitOperand := opType == CodeR32 || opType == CodeM32 || opType == CodeIMM32

	if is16bitMode && is32bitOperand {
		return true
	}
	if is32bitMode && is16bitOperand {
		return true
	}

	return false
}

// Require67h はアドレスサイズプレフィックス (67h) が必要か判定します。
// TODO: 実装する (requires.go 相当)
func (o *OperandPegImpl) Require67h() bool {
	// TODO: bitMode とメモリオペランドのアドレス指定 (Base/Indexレジスタ) を見て判定
	if o.parsed == nil || o.parsed.Memory == nil {
		return false
	}
	// 仮: 32bitモードで16bitアドレッシング (BX, SI, DI, BP) を使っていたら true
	// 仮: 16bitモードで32bitアドレッシング (EAXなど) を使っていたら true
	mem := o.parsed.Memory
	is32bitAddr := strings.HasPrefix(mem.BaseReg, "E") || strings.HasPrefix(mem.IndexReg, "E")
	is16bitAddr := !is32bitAddr && (mem.BaseReg != "" || mem.IndexReg != "") // Eなしでレジスタがあれば16bit

	if o.bitMode == cpu.MODE_32BIT && is16bitAddr {
		return true
	}
	if o.bitMode == cpu.MODE_16BIT && is32bitAddr {
		return true
	}

	return false // 仮実装
}

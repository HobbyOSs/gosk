package ng_operand // パッケージ名を変更

import "github.com/HobbyOSs/gosk/pkg/cpu" // BitModeのためにcpuインポートを維持

// Operands インターフェース (pkg/operand/operand.go からコピー)
// このインターフェースは、オペランドのセットを表します。
// 元のparticipleベースのパーサーから移行中のため、一部のメソッドは
// 新しいpigeonベースのパーサーの実装 (`OperandPegImpl`) で提供されます。
type Operands interface {
	// InternalString は、内部表現の単一文字列を返します（主にデバッグ用）。
	InternalString() string
	// InternalStrings は、各オペランドの内部表現文字列のスライスを返します。
	InternalStrings() []string
	// OperandTypes は、各オペランドの型 (`OperandType`) のスライスを返します。
	OperandTypes() []OperandType
	// Serialize は、オペランドをシリアライズ可能な形式（通常は元の文字列）で返します。
	Serialize() string
	// FromString は、与えられた文字列から新しい Operands オブジェクトを生成します。
	FromString(text string) Operands
	// CalcOffsetByteSize は、メモリオペランドのオフセット部分のバイトサイズを計算します。
	CalcOffsetByteSize() int
	// DetectImmediateSize は、即値オペランドのサイズ（バイト単位）を検出します。
	DetectImmediateSize() int
	// WithBitMode は、指定されたビットモード (`cpu.BitMode`) を持つ新しい Operands オブジェクトを返します。
	WithBitMode(mode cpu.BitMode) Operands // 再追加
	// WithForceImm8 は、即値を強制的に8ビットとして扱うかどうかを設定した新しい Operands オブジェクトを返します。
	WithForceImm8(force bool) Operands
	// WithForceRelAsImm は、相対アドレスを即値として強制的に扱うかどうかを設定した新しい Operands オブジェクトを返します。
	WithForceRelAsImm(force bool) Operands
	// GetBitMode は、現在のビットモードを返します。
	GetBitMode() cpu.BitMode // 再追加
	// Require66h は、オペランドサイズプレフィックス (0x66) が必要かどうかを返します。
	Require66h() bool // オペランドサイズプレフィックスが必要かどうか
	// Require67h は、アドレスサイズプレフィックス (0x67) が必要かどうかを返します。
	Require67h() bool // アドレスサイズプレフィックスが必要かどうか
	// IsDirectMemory は、オペランドに直接メモリアドレスが含まれるかどうかを返します。
	IsDirectMemory() bool
	// IsIndirectMemory は、オペランドに間接メモリアドレスが含まれるかどうかを返します。
	IsIndirectMemory() bool
	// GetMemoryInfo は、最初のメモリオペランドの詳細情報を返します。見つからない場合は nil と false を返します。
	GetMemoryInfo() (*MemoryInfo, bool)
	// DisplacementBytes は、最初のメモリオペランドのディスプレースメント部分をバイト列として返します。
	// ModRMがない直接アドレス指定 (moffs) の場合に利用することを想定しています。
	// メモリオペランドがない場合や、ディスプレースメントがない場合は nil を返します。
	// バイト列のサイズは BitMode に基づいて決定されます。
	DisplacementBytes() []byte
}

package ast

type SegmentExpHandler interface {
	SegmentExp(*SegmentExp) bool
}

type MemoryAddrExpHandler interface {
	MemoryAddrExp(*MemoryAddrExp) bool
}

type AddExpHandler interface {
	AddExp(*AddExp) bool
}

type MultExpHandler interface {
	MultExp(*MultExp) bool
}

type ImmExpHandler interface {
	ImmExp(*ImmExp) bool
}

type NumberFactorHandler interface {
	NumberFactor(*NumberFactor) bool
}

type StringFactorHandler interface {
	StringFactor(*StringFactor) bool
}

type HexFactorHandler interface {
	HexFactor(*HexFactor) bool
}

type IdentFactorHandler interface {
	IdentFactor(*IdentFactor) bool
}

type CharFactorHandler interface {
	CharFactor(*CharFactor) bool
}

type ProgramHandler interface {
	Program(*Program) bool
}

type DeclareStmtHandler interface {
	DeclareStmt(*DeclareStmt) bool
}

type LabelStmtHandler interface {
	LabelStmt(*LabelStmt) bool
}

type ExportSymStmtHandler interface {
	ExportSymStmt(*ExportSymStmt) bool
}

type ExternSymStmtHandler interface {
	ExternSymStmt(*ExternSymStmt) bool
}

type ConfigStmtHandler interface {
	ConfigStmt(*ConfigStmt) bool
}

type MnemonicStmtHandler interface {
	MnemonicStmt(*MnemonicStmt) bool
}

type DefaultHandler interface {
	Handle(Node) bool
}

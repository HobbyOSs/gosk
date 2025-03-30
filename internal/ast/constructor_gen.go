// Code generated by github.com/Bin-Huang/newc; DO NOT EDIT.

package ast

// NewSegmentExp Create a new SegmentExp
func NewSegmentExp(baseExp BaseExp, dataType DataType, left *AddExp, right *AddExp) *SegmentExp {
	return &SegmentExp{
		BaseExp:  baseExp,
		DataType: dataType,
		Left:     left,
		Right:    right,
	}
}

// NewMemoryAddrExp Create a new MemoryAddrExp
func NewMemoryAddrExp(baseExp BaseExp, dataType DataType, jumpType JumpType, left *AddExp, right *AddExp) *MemoryAddrExp {
	return &MemoryAddrExp{
		BaseExp:  baseExp,
		DataType: dataType,
		JumpType: jumpType,
		Left:     left,
		Right:    right,
	}
}

// NewAddExp Create a new AddExp
func NewAddExp(baseExp BaseExp, headExp *MultExp, operators []string, tailExps []*MultExp) *AddExp {
	return &AddExp{
		BaseExp:   baseExp,
		HeadExp:   headExp,
		Operators: operators,
		TailExps:  tailExps,
	}
}

// NewMultExp Create a new MultExp
func NewMultExp(baseExp BaseExp, headExp Exp, operators []string, tailExps []Exp) *MultExp {
	return &MultExp{
		BaseExp:   baseExp,
		HeadExp:   headExp,
		Operators: operators,
		TailExps:  tailExps,
	}
}

// NewImmExp Create a new ImmExp
func NewImmExp(baseExp BaseExp, factor Factor) *ImmExp {
	return &ImmExp{
		BaseExp: baseExp,
		Factor:  factor,
	}
}

// NewNumberFactor Create a new NumberFactor
func NewNumberFactor(baseFactor BaseFactor, value int) *NumberFactor {
	return &NumberFactor{
		BaseFactor: baseFactor,
		Value:      value,
	}
}

// NewStringFactor Create a new StringFactor
func NewStringFactor(baseFactor BaseFactor, value string) *StringFactor {
	return &StringFactor{
		BaseFactor: baseFactor,
		Value:      value,
	}
}

// NewHexFactor Create a new HexFactor
func NewHexFactor(baseFactor BaseFactor, value string) *HexFactor {
	return &HexFactor{
		BaseFactor: baseFactor,
		Value:      value,
	}
}

// NewIdentFactor Create a new IdentFactor
func NewIdentFactor(baseFactor BaseFactor, value string) *IdentFactor {
	return &IdentFactor{
		BaseFactor: baseFactor,
		Value:      value,
	}
}

// NewCharFactor Create a new CharFactor
func NewCharFactor(baseFactor BaseFactor, value string) *CharFactor {
	return &CharFactor{
		BaseFactor: baseFactor,
		Value:      value,
	}
}

// NewProgram Create a new Program
func NewProgram(statements []Statement) *Program {
	return &Program{
		Statements: statements,
	}
}

// NewDeclareStmt Create a new DeclareStmt
func NewDeclareStmt(baseStatement BaseStatement, id *IdentFactor, value Exp) *DeclareStmt {
	return &DeclareStmt{
		BaseStatement: baseStatement,
		Id:            id,
		Value:         value,
	}
}

// NewOpcodeStmt Create a new OpcodeStmt
func NewOpcodeStmt(baseStatement BaseStatement, opcode *IdentFactor) *OpcodeStmt {
	return &OpcodeStmt{
		BaseStatement: baseStatement,
		Opcode:        opcode,
	}
}

// NewLabelStmt Create a new LabelStmt
func NewLabelStmt(baseStatement BaseStatement, label *IdentFactor) *LabelStmt {
	return &LabelStmt{
		BaseStatement: baseStatement,
		Label:         label,
	}
}

// NewExportSymStmt Create a new ExportSymStmt
func NewExportSymStmt(baseStatement BaseStatement, symbols []*IdentFactor) *ExportSymStmt {
	return &ExportSymStmt{
		BaseStatement: baseStatement,
		Symbols:       symbols,
	}
}

// NewExternSymStmt Create a new ExternSymStmt
func NewExternSymStmt(baseStatement BaseStatement, symbols []*IdentFactor) *ExternSymStmt {
	return &ExternSymStmt{
		BaseStatement: baseStatement,
		Symbols:       symbols,
	}
}

// NewConfigStmt Create a new ConfigStmt
func NewConfigStmt(baseStatement BaseStatement, configType ConfigType, factor Factor) *ConfigStmt {
	return &ConfigStmt{
		BaseStatement: baseStatement,
		ConfigType:    configType,
		Factor:        factor,
	}
}

// NewMnemonicStmt Create a new MnemonicStmt
func NewMnemonicStmt(baseStatement BaseStatement, opcode *IdentFactor, operands []Exp) *MnemonicStmt {
	return &MnemonicStmt{
		BaseStatement: baseStatement,
		Opcode:        opcode,
		Operands:      operands,
	}
}

package ast

type Visitor struct {
	hSegmentExp    SegmentExpHandler
	hMemoryAddrExp MemoryAddrExpHandler
	hAddExp        AddExpHandler
	hMultExp       MultExpHandler
	hImmExp        ImmExpHandler
	hNumberFactor  NumberFactorHandler
	hStringFactor  StringFactorHandler
	hHexFactor     HexFactorHandler
	hIdentFactor   IdentFactorHandler
	hCharFactor    CharFactorHandler
	hProgram       ProgramHandler
	hDeclareStmt   DeclareStmtHandler
	hLabelStmt     LabelStmtHandler
	hExportSymStmt ExportSymStmtHandler
	hExternSymStmt ExternSymStmtHandler
	hConfigStmt    ConfigStmtHandler
	hMnemonicStmt  MnemonicStmtHandler
	hDefault       DefaultHandler
}

func (v *Visitor) Handler(h interface{}) error {
	if x, ok := h.(SegmentExpHandler); ok {
		v.hSegmentExp = x
	}
	if x, ok := h.(MemoryAddrExpHandler); ok {
		v.hMemoryAddrExp = x
	}
	if x, ok := h.(AddExpHandler); ok {
		v.hAddExp = x
	}
	if x, ok := h.(MultExpHandler); ok {
		v.hMultExp = x
	}
	if x, ok := h.(ImmExpHandler); ok {
		v.hImmExp = x
	}
	if x, ok := h.(NumberFactorHandler); ok {
		v.hNumberFactor = x
	}
	if x, ok := h.(StringFactorHandler); ok {
		v.hStringFactor = x
	}
	if x, ok := h.(HexFactorHandler); ok {
		v.hHexFactor = x
	}
	if x, ok := h.(IdentFactorHandler); ok {
		v.hIdentFactor = x
	}
	if x, ok := h.(CharFactorHandler); ok {
		v.hCharFactor = x
	}
	if x, ok := h.(ProgramHandler); ok {
		v.hProgram = x
	}
	if x, ok := h.(DeclareStmtHandler); ok {
		v.hDeclareStmt = x
	}
	if x, ok := h.(LabelStmtHandler); ok {
		v.hLabelStmt = x
	}
	if x, ok := h.(ExportSymStmtHandler); ok {
		v.hExportSymStmt = x
	}
	if x, ok := h.(ExternSymStmtHandler); ok {
		v.hExternSymStmt = x
	}
	if x, ok := h.(ConfigStmtHandler); ok {
		v.hConfigStmt = x
	}
	if x, ok := h.(MnemonicStmtHandler); ok {
		v.hMnemonicStmt = x
	}
	if x, ok := h.(DefaultHandler); ok {
		v.hDefault = x
	}
	return nil
}

func (v *Visitor) Visit(n Node) *Visitor {
	switch n := n.(type) {
	case *SegmentExp:
		if h := v.hSegmentExp; h != nil {
			if !h.SegmentExp(n) {
				return nil
			}
			return v
		}
	case *MemoryAddrExp:
		if h := v.hMemoryAddrExp; h != nil {
			if !h.MemoryAddrExp(n) {
				return nil
			}
			return v
		}
	case *AddExp:
		if h := v.hAddExp; h != nil {
			if !h.AddExp(n) {
				return nil
			}
			return v
		}
	case *MultExp:
		if h := v.hMultExp; h != nil {
			if !h.MultExp(n) {
				return nil
			}
			return v
		}
	case *ImmExp:
		if h := v.hImmExp; h != nil {
			if !h.ImmExp(n) {
				return nil
			}
			return v
		}
	case *NumberFactor:
		if h := v.hNumberFactor; h != nil {
			if !h.NumberFactor(n) {
				return nil
			}
			return v
		}
	case *StringFactor:
		if h := v.hStringFactor; h != nil {
			if !h.StringFactor(n) {
				return nil
			}
			return v
		}
	case *HexFactor:
		if h := v.hHexFactor; h != nil {
			if !h.HexFactor(n) {
				return nil
			}
			return v
		}
	case *IdentFactor:
		if h := v.hIdentFactor; h != nil {
			if !h.IdentFactor(n) {
				return nil
			}
			return v
		}
	case *CharFactor:
		if h := v.hCharFactor; h != nil {
			if !h.CharFactor(n) {
				return nil
			}
			return v
		}
	case *Program:
		if h := v.hProgram; h != nil {
			if !h.Program(n) {
				return nil
			}
			return v
		}
	case *DeclareStmt:
		if h := v.hDeclareStmt; h != nil {
			if !h.DeclareStmt(n) {
				return nil
			}
			return v
		}
	case *LabelStmt:
		if h := v.hLabelStmt; h != nil {
			if !h.LabelStmt(n) {
				return nil
			}
			return v
		}
	case *ExportSymStmt:
		if h := v.hExportSymStmt; h != nil {
			if !h.ExportSymStmt(n) {
				return nil
			}
			return v
		}
	case *ExternSymStmt:
		if h := v.hExternSymStmt; h != nil {
			if !h.ExternSymStmt(n) {
				return nil
			}
			return v
		}
	case *ConfigStmt:
		if h := v.hConfigStmt; h != nil {
			if !h.ConfigStmt(n) {
				return nil
			}
			return v
		}
	case *MnemonicStmt:
		if h := v.hMnemonicStmt; h != nil {
			if !h.MnemonicStmt(n) {
				return nil
			}
			return v
		}
	}
	if h := v.hDefault; h != nil {
		if !h.Handle(n) {
			return nil
		}
	}
	return v
}

// Code generated by github.com/Bin-Huang/newc; DO NOT EDIT.

package pass1

// NewPass1 Create a new Pass1
func NewPass1(loc int32, symTable map[string]uint32, globalSymbolList []string, externSymbolList []string) *Pass1 {
	return &Pass1{
		LOC:              loc,
		SymTable:         symTable,
		GlobalSymbolList: globalSymbolList,
		ExternSymbolList: externSymbolList,
	}
}

package ast

import (
	"fmt" // Keep fmt for panic
	"strconv"
	"strings"
)

// TODO: go generateで作成できないか
type DataType string

const (
	Byte  DataType = "BYTE"
	Word  DataType = "WORD"
	Dword DataType = "DWORD"
	None  DataType = ""
)

var stringToDataType = map[string]DataType{
	"BYTE":  Byte,
	"WORD":  Word,
	"DWORD": Dword,
	"":      None,
}

func NewDataType(s string) (DataType, bool) {
	c, ok := stringToDataType[s]
	return c, ok
}

// TODO: go generateで作成できないか
type JumpType string

const (
	Short JumpType = "SHORT"
	Near  JumpType = "NEAR"
	Far   JumpType = "FAR"
	Empty JumpType = ""
)

var stringToJumpType = map[string]JumpType{
	"SHORT": Short,
	"NEAR":  Near,
	"FAR":   Far,
	"":      Empty,
}

func NewJumpType(s string) (JumpType, bool) {
	c, ok := stringToJumpType[s]
	return c, ok
}

//go:generate newc
type SegmentExp struct {
	BaseExp
	DataType DataType
	Left     *AddExp
	Right    *AddExp // nullable
}

func (s *SegmentExp) expressionNode() {}
func (s *SegmentExp) Eval(env Env) (Exp, bool) {
	// TODO: Implement SegmentExp evaluation logic
	// For now, just return the node itself, indicating no reduction.
	return s, false
}
func (s *SegmentExp) TokenLiteral() string {
	leftStr := ExpToString(s.Left)
	rightStr := ""
	if s.Right != nil {
		rightStr = ExpToString(s.Right)
	}
	dataTypeStr := ""
	if s.DataType != None {
		dataTypeStr = string(s.DataType) + " "
	}
	if rightStr == "" {
		return dataTypeStr + leftStr
	} else {
		return dataTypeStr + leftStr + ":" + rightStr
	}
}

//go:generate newc
type MemoryAddrExp struct {
	BaseExp
	DataType DataType
	JumpType JumpType
	Left     *AddExp
	Right    *AddExp // nullable
}

func (m *MemoryAddrExp) expressionNode() {}
func (m *MemoryAddrExp) Eval(env Env) (Exp, bool) {
	// TODO: Implement MemoryAddrExp evaluation logic
	// For now, just return the node itself, indicating no reduction.
	return m, false
}
func (m *MemoryAddrExp) TokenLiteral() string {
	var str = ""
	if m.DataType != None {
		str += string(m.DataType)
		str += " "
	}
	str += "[ "
	str += m.Left.TokenLiteral()
	if m.Right != nil {
		str += " : "
		str += m.Right.TokenLiteral()
	}
	str += " ]"
	return str
}

//go:generate newc
type AddExp struct {
	BaseExp
	HeadExp   *MultExp
	Operators []string
	TailExps  []*MultExp
}

func (a *AddExp) expressionNode() {}
func (a *AddExp) Eval(env Env) (Exp, bool) {
	// Evaluate head expression
	evalHeadExp, headReduced := a.HeadExp.Eval(env)
	evalHead, ok := evalHeadExp.(*MultExp)
	if !ok {
		if headNum, isNum := evalHeadExp.(*NumberExp); isNum {
			// Wrap NumberExp into a simple MultExp (ImmExp with NumberFactor)
			evalHead = &MultExp{BaseExp: BaseExp{}, HeadExp: &headNum.ImmExp}
		} else {
			// Head evaluated to something unexpected
			return a, false
		}
	}

	// Evaluate tail expressions
	evalTails := make([]*MultExp, len(a.TailExps))
	anyTailReduced := false
	allAreNumbers := true

	// Check if head is a number
	headNumCheck, headIsNum := evalHead.Eval(env) // Re-eval potentially wrapped head
	if _, ok := headNumCheck.(*NumberExp); !ok {
		allAreNumbers = false
	}

	for i, tail := range a.TailExps {
		evalTailExp, tailReduced := tail.Eval(env)
		evalTailNode, ok := evalTailExp.(*MultExp)
		if !ok {
			if tailNum, isNum := evalTailExp.(*NumberExp); isNum {
				// Wrap NumberExp into a simple MultExp
				evalTailNode = &MultExp{BaseExp: BaseExp{}, HeadExp: &tailNum.ImmExp}
			} else {
				// Tail evaluated to something unexpected
				allAreNumbers = false
				evalTails[i] = tail // Keep original if eval failed
				continue
			}
		}

		evalTails[i] = evalTailNode
		if tailReduced {
			anyTailReduced = true
		}

		// Check if this evaluated tail is a number
		tailNumCheck, _ := evalTailNode.Eval(env) // Re-eval potentially wrapped tail
		if _, ok := tailNumCheck.(*NumberExp); !ok {
			allAreNumbers = false
		}
	}

	// If head and all tails evaluated to numbers, calculate the result
	if headIsNum && allAreNumbers {
		currentValue := headNumCheck.(*NumberExp).Value // Use the checked value
		for i, op := range a.Operators {
			// Eval the simplified MultExp tail to get the NumberExp
			tailNumExp, _ := evalTails[i].Eval(env)
			numTail, ok := tailNumExp.(*NumberExp)
			if !ok {
				// This should ideally not happen if allAreNumbers logic is correct
				panic("Internal error: Expected NumberExp after simplification in AddExp")
			}

			switch op {
			case "+":
				currentValue += numTail.Value
			case "-":
				currentValue -= numTail.Value
			default:
				// Unsupported operator
				return a, false
			}
		}
		// Return a new NumberExp
		return NewNumberExp(ImmExp{BaseExp: a.BaseExp}, currentValue), true
	}

	// If not all parts evaluated to numbers, but some reduction occurred, return updated AddExp
	if headReduced || anyTailReduced {
		return NewAddExp(a.BaseExp, evalHead, a.Operators, evalTails), true
	}

	// No reduction possible
	return a, false
}
func (a *AddExp) TokenLiteral() string {
	head := ExpToString(a.HeadExp)
	var buf strings.Builder
	buf.WriteString(head)
	for i, op := range a.Operators {
		buf.WriteByte(' ')
		buf.WriteString(op)
		buf.WriteByte(' ')
		tailStr := ExpToString(a.TailExps[i])
		buf.WriteString(tailStr)
	}
	return buf.String()
}

//go:generate newc
type MultExp struct {
	BaseExp
	HeadExp   *ImmExp
	Operators []string
	TailExps  []*ImmExp
}

func (m *MultExp) expressionNode() {}
func (m *MultExp) Eval(env Env) (Exp, bool) {
	// Evaluate head expression
	evalHeadExp, headReduced := m.HeadExp.Eval(env)
	evalHead, ok := evalHeadExp.(*ImmExp) // Expecting ImmExp or NumberExp
	if !ok {
		// Head evaluated to something unexpected (e.g., macro expanded to AddExp)
		return m, false // Cannot evaluate if head is not ImmExp compatible
	}
	headNum, headIsNum := evalHeadExp.(*NumberExp) // Check if it's specifically a NumberExp

	// Evaluate tail expressions
	evalTails := make([]*ImmExp, len(m.TailExps))
	anyTailReduced := false
	allAreNumbers := true // Assume true initially

	if !headIsNum {
		allAreNumbers = false // Head must be number for full evaluation
	}

	for i, tail := range m.TailExps {
		evalTailExp, tailReduced := tail.Eval(env)
		evalTailNode, ok := evalTailExp.(*ImmExp) // Expecting ImmExp or NumberExp
		if !ok {
			// Tail evaluated to something unexpected
			allAreNumbers = false
			evalTails[i] = tail // Keep original if eval failed
			continue
		}

		evalTails[i] = evalTailNode // Store evaluated ImmExp (or NumberExp)
		if tailReduced {
			anyTailReduced = true
		}

		// Check if the factor within the evaluated tail is a NumberFactor
		if _, isNumFactor := evalTailNode.Factor.(*NumberFactor); !isNumFactor {
			allAreNumbers = false
		}
	}

	// If head and all tails evaluated to numbers, calculate the result
	if headIsNum && allAreNumbers {
		// Get head value (we know headNum is *NumberExp here)
		currentValue := headNum.Value
		for i, op := range m.Operators {
			// Get tail value by checking its factor (we know it's NumberFactor if allAreNumbers is true)
			numFactor, ok := evalTails[i].Factor.(*NumberFactor)
			if !ok {
				// This should not happen if allAreNumbers logic is correct
				panic(fmt.Sprintf("Internal error: Expected NumberFactor in MultExp calculation but got %T", evalTails[i].Factor))
			}
			tailValue := int64(numFactor.Value) // Convert factor value to int64

			switch op {
			case "*":
				currentValue *= tailValue
			case "/":
				if tailValue == 0 {
					return m, false // Division by zero
				}
				currentValue /= tailValue
			case "%":
				if tailValue == 0 {
					return m, false // Modulo by zero
				}
				currentValue %= tailValue
			default:
				return m, false // Unsupported operator
			}
		}
		// Return a new NumberExp
		return NewNumberExp(ImmExp{BaseExp: m.BaseExp}, currentValue), true
	}

	// If not all parts evaluated to numbers, but some reduction occurred, return updated MultExp
	if headReduced || anyTailReduced {
		return NewMultExp(m.BaseExp, evalHead, m.Operators, evalTails), true
	}

	// No reduction possible
	return m, false
}
func (m *MultExp) TokenLiteral() string {
	head := ExpToString(m.HeadExp)
	var buf strings.Builder
	buf.WriteString(head)
	for i, op := range m.Operators {
		buf.WriteByte(' ')
		buf.WriteString(op)
		buf.WriteByte(' ')
		tailStr := ExpToString(m.TailExps[i])
		buf.WriteString(tailStr)
	}
	return buf.String()
}

//go:generate newc
type ImmExp struct {
	BaseExp
	Factor Factor
}

func (imm *ImmExp) expressionNode() {}
func (imm *ImmExp) Eval(env Env) (Exp, bool) {
	switch f := imm.Factor.(type) {
	case *NumberFactor:
		val := int64(f.Value)
		newFactor := NewNumberFactor(BaseFactor{}, int(val))
		numExp := NewNumberExp(ImmExp{BaseExp: imm.BaseExp, Factor: newFactor}, val)
		return numExp, true
	case *HexFactor:
		val, ok := parseHex(f.Value)
		if !ok {
			return imm, false
		}
		newFactor := NewNumberFactor(BaseFactor{}, int(val))
		numExp := NewNumberExp(ImmExp{BaseExp: imm.BaseExp, Factor: newFactor}, val)
		return numExp, true
	case *CharFactor:
		val, ok := parseChar(f.Value)
		if !ok {
			return imm, false
		}
		newFactor := NewNumberFactor(BaseFactor{}, int(val))
		numExp := NewNumberExp(ImmExp{BaseExp: imm.BaseExp, Factor: newFactor}, val)
		return numExp, true
	case *IdentFactor:
		macroExp, ok := env.LookupMacro(f.Value)
		if ok {
			// Recursively evaluate the macro definition
			return macroExp.Eval(env)
		}
		// If not a macro, it's an unresolved identifier (like a label)
		return imm, false
	case *StringFactor:
		// Cannot evaluate string factors in arithmetic expressions
		return imm, false
	default:
		// Unknown factor type
		return imm, false
	}
}
func (imm *ImmExp) TokenLiteral() string {
	return imm.Factor.TokenLiteral()
}

// --- Helper functions for parsing ---

func parseHex(s string) (int64, bool) {
	if !strings.HasPrefix(s, "0x") && !strings.HasPrefix(s, "0X") {
		return 0, false
	}
	val, err := strconv.ParseInt(s[2:], 16, 64)
	if err != nil {
		return 0, false
	}
	return val, true
}

func parseChar(s string) (int64, bool) {
	if len(s) < 2 || s[0] != '\'' || s[len(s)-1] != '\'' {
		return 0, false
	}
	charStr := s[1 : len(s)-1]
	if len(charStr) != 1 {
		// TODO: Handle escape sequences like '\n', '\\', '\'' etc.
		return 0, false
	}
	return int64(charStr[0]), true
}

// NumberExp represents a fully evaluated numeric constant expression.
type NumberExp struct {
	ImmExp       // Embed ImmExp to satisfy Exp interface
	Value  int64 // The evaluated numeric value
}

// NewNumberExp creates a new NumberExp.
func NewNumberExp(base ImmExp, value int64) *NumberExp {
	base.Factor = NewNumberFactor(BaseFactor{}, int(value)) // Ensure Factor is NumberFactor
	return &NumberExp{
		ImmExp: base,
		Value:  value,
	}
}

// Eval for NumberExp simply returns itself as it's already fully evaluated.
func (n *NumberExp) Eval(env Env) (Exp, bool) {
	return n, false // Already evaluated, no reduction happened *in this step*
}

// TokenLiteral returns the string representation of the number.
func (n *NumberExp) TokenLiteral() string {
	// Use the embedded Factor's TokenLiteral, which should be a NumberFactor.
	return n.Factor.TokenLiteral()
}

// Ensure NumberExp satisfies the Exp interface.
var _ Exp = &NumberExp{}

// --- Add Eval implementations for other expression types (AddExp, MultExp, etc.) ---
// Placeholder for UnaryExp if it exists or is needed
// //go:generate newc
// type UnaryExp struct {
// 	BaseExp
// 	Operator string
// 	Exp      Exp
// }
//
// func (u *UnaryExp) expressionNode() {}
// func (u *UnaryExp) Eval(env Env) (Exp, bool) {
// 	// TODO: Implement UnaryExp evaluation logic
// 	return u, false
// }
// func (u *UnaryExp) TokenLiteral() string {
// 	// TODO: Implement TokenLiteral for UnaryExp
// 	return u.Operator + ExpToString(u.Exp)
// }

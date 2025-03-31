package ast

import (
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

// wrapExpInAddExp wraps a simple Exp (NumberExp, ImmExp) into the structure needed for AddExp fields.
func wrapExpInAddExp(exp Exp) *AddExp {
	if exp == nil {
		return nil
	}
	if addExp, ok := exp.(*AddExp); ok {
		// If it's already an AddExp, return it
		return addExp
	}

	var immExp *ImmExp
	if numExp, ok := exp.(*NumberExp); ok {
		// Ensure the NumberExp has a valid ImmExp embedded
		if numExp.ImmExp.Factor == nil {
			// If Factor is missing (shouldn't happen with NewNumberExp), create it
			numExp.ImmExp.Factor = NewNumberFactor(BaseFactor{}, int(numExp.Value))
		}
		immExp = &numExp.ImmExp
	} else if ie, ok := exp.(*ImmExp); ok {
		immExp = ie
	} else {
		// Cannot easily wrap other types like MultExp directly here.
		// This helper is primarily for NumberExp/ImmExp results from Eval.
		// log.Printf("wrapExpInAddExp: Cannot wrap type %T", exp)
		return nil // Return nil if wrapping is not straightforward
	}

	// Create MultExp -> AddExp structure
	// Ensure BaseExp is initialized for NewMultExp and NewAddExp
	multExp := NewMultExp(BaseExp{}, immExp, nil, nil) // Head is the ImmExp
	addExp := NewAddExp(BaseExp{}, multExp, nil, nil)  // AddExp with only head
	return addExp
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
	// Evaluate the internal expression(s)
	evalLeftNode, leftReduced := m.Left.Eval(env) // Returns Exp
	evalRightNode := Exp(nil)                     // Initialize evalRightNode
	rightReduced := false
	if m.Right != nil {
		evalRightNode, rightReduced = m.Right.Eval(env) // Returns Exp
	}

	// Wrap the evaluated nodes back into AddExp structure if possible
	evalLeftExp := wrapExpInAddExp(evalLeftNode)
	if evalLeftExp == nil && leftReduced {
		// If wrapping failed but reduction happened, we can't represent the state.
		// Return original to avoid losing information or creating invalid structure.
		// log.Printf("Warning: MemoryAddrExp Left expression evaluated to unwrappable type %T", evalLeftNode)
		return m, false
	} else if evalLeftExp == nil {
		evalLeftExp = m.Left // Keep original if no reduction and no wrapping possible
	}

	evalRightExp := (*AddExp)(nil)
	if m.Right != nil {
		evalRightExp = wrapExpInAddExp(evalRightNode)
		if evalRightExp == nil && rightReduced {
			// log.Printf("Warning: MemoryAddrExp Right expression evaluated to unwrappable type %T", evalRightNode)
			return m, false
		} else if evalRightExp == nil {
			evalRightExp = m.Right // Keep original
		}
	}

	// If neither internal expression was reduced, return the original node
	if !leftReduced && !rightReduced {
		return m, false
	}

	// Construct a new MemoryAddrExp with the potentially wrapped internal expressions
	newMemExp := NewMemoryAddrExp(m.BaseExp, m.DataType, m.JumpType, evalLeftExp, evalRightExp)
	return newMemExp, true // Return the new node and indicate reduction occurred
}
func (m *MemoryAddrExp) TokenLiteral() string {
	// Use existing ExpToString from ast_exp_string.go
	var str = ""
	if m.DataType != None {
		str += string(m.DataType)
		str += " "
	}
	str += "[ "
	// Use ExpToString to handle potentially evaluated Left expression
	str += ExpToString(m.Left) // m.Left might point to the original if Eval didn't replace it
	if m.Right != nil {
		str += " : "
		// Use ExpToString for Right as well
		str += ExpToString(m.Right)
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

// Eval performs constant folding for AddExp.
// It sums up all constant number terms and keeps non-constant terms.
// Modified to keep non-constant terms first.
func (a *AddExp) Eval(env Env) (Exp, bool) {
	// Evaluate head expression first
	evalHead, headReduced := a.HeadExp.Eval(env)

	// Keep track of the sum of constant terms and the list of non-constant terms/operators
	constSum := 0
	newTerms := []Exp{} // Use Exp interface to hold evaluated non-constant terms
	newOps := []string{}
	reduced := headReduced // Start with head reduction status

	// Process the evaluated head
	if v, ok := env.GetConstValue(evalHead); ok {
		constSum += v
	} else {
		newTerms = append(newTerms, evalHead)
	}

	// Process tail expressions
	for i, op := range a.Operators {
		tail := a.TailExps[i]
		evalTail, tailReduced := tail.Eval(env)
		if tailReduced {
			reduced = true
		}

		if v, ok := env.GetConstValue(evalTail); ok {
			// If it's a constant, add/subtract it to the sum
			if op == "+" {
				constSum += v
			} else if op == "-" {
				constSum -= v
			} else {
				// Should not happen based on grammar, but handle defensively
				// If an unsupported operator appears with a constant, treat as non-reducible
				// Keep the operator and the constant term
				if len(newTerms) > 0 { // Only add operator if there's a preceding term
					newOps = append(newOps, op)
				}
				newTerms = append(newTerms, evalTail) // Add the constant term back
			}
		} else {
			// If it's not a constant, add it to the list of terms
			// Only add operator if there was a preceding term
			if len(newTerms) > 0 { // Add operator if it's not the very first term
				newOps = append(newOps, op)
			}
			newTerms = append(newTerms, evalTail)
		}
	}

	// --- Construct the result ---

	// Case 1: All terms evaluated to constants
	if len(newTerms) == 0 {
		// Return a single NumberExp with the final sum
		return NewNumberExp(ImmExp{BaseExp: a.BaseExp}, int64(constSum)), true
	}

	// Case 2: Mixed constants and non-constants

	// Reorder terms: non-constants first, then the constant sum (if non-zero)
	finalTerms := []Exp{}
	finalOps := []string{}

	// Add non-constant terms first
	if len(newTerms) > 0 {
		finalTerms = append(finalTerms, newTerms...)
		finalOps = append(finalOps, newOps...) // Keep original operators between non-const terms
	}

	// Add the constant sum at the end if it's non-zero
	if constSum != 0 {
		constTerm := NewNumberExp(ImmExp{BaseExp: a.BaseExp}, int64(constSum))
		if len(finalTerms) > 0 {
			// Add '+' or '-' operator before the constant term if other terms exist
			if constSum > 0 {
				finalOps = append(finalOps, "+")
			} else {
				finalOps = append(finalOps, "-")
				// Use the absolute value for the NumberExp if operator is '-'
				constTerm = NewNumberExp(ImmExp{BaseExp: a.BaseExp}, int64(-constSum))
			}
		} else {
			// If only the constant term exists, make it the only term
			// Handle negative constant as head if it's the only term
			if constSum < 0 {
				// This case should ideally be handled by Case 1 returning NumberExp,
				// but let's ensure the NumberExp value is correct if we reach here.
				constTerm = NewNumberExp(ImmExp{BaseExp: a.BaseExp}, int64(constSum))
			}
		}
		finalTerms = append(finalTerms, constTerm)
	}

	// If after reordering, only one term remains (could be non-constant or constSum), return it directly if possible
	if len(finalTerms) == 1 && len(finalOps) == 0 {
		// If it's a NumberExp, return it (already handled by Case 1 ideally)
		if numExp, ok := finalTerms[0].(*NumberExp); ok {
			return numExp, true
		}
		// If it's a single non-constant term, we still need to wrap it in AddExp structure below
	}

	// If no terms remain (e.g., "LABEL - LABEL"), result is 0
	if len(finalTerms) == 0 {
		return NewNumberExp(ImmExp{BaseExp: a.BaseExp}, 0), true
	}

	// --- Reconstruct AddExp with the new order ---

	// The first term in finalTerms is the new head
	finalHead := finalTerms[0]

	// Convert remaining evaluated terms back to *MultExp for the AddExp structure
	finalTailNodes := make([]*MultExp, 0, len(finalTerms)-1)
	for _, term := range finalTerms[1:] { // Iterate over reordered finalTerms
		if me, ok := term.(*MultExp); ok {
			finalTailNodes = append(finalTailNodes, me)
		} else if num, ok := term.(*NumberExp); ok {
			// Wrap NumberExp back into MultExp
			// Ensure the embedded ImmExp has the correct Factor
			numImmExp := num.ImmExp
			if numImmExp.Factor == nil {
				numImmExp.Factor = NewNumberFactor(BaseFactor{}, int(num.Value))
			}
			finalTailNodes = append(finalTailNodes, &MultExp{BaseExp: BaseExp{}, HeadExp: &numImmExp})
		} else if imm, ok := term.(*ImmExp); ok {
			// Wrap ImmExp (like identifiers) into MultExp
			finalTailNodes = append(finalTailNodes, &MultExp{BaseExp: BaseExp{}, HeadExp: imm})
		} else {
			// If it's some other Exp type that can't be easily put into MultExp,
			// we might not be able to simplify perfectly. Return original or error.
			// For now, let's assume terms are MultExp, NumberExp, or ImmExp.
			// Returning original if we hit an unexpected type.
			// log.Printf("Warning: Cannot reconstruct AddExp tail from type %T", term)
			return a, false // Cannot simplify if unexpected type found
		}
	}

	// Convert the finalHead (which is Exp) back to *MultExp
	finalHeadNode, ok := finalHead.(*MultExp)
	if !ok {
		if num, ok := finalHead.(*NumberExp); ok {
			// Ensure the embedded ImmExp has the correct Factor
			numImmExp := num.ImmExp
			if numImmExp.Factor == nil {
				numImmExp.Factor = NewNumberFactor(BaseFactor{}, int(num.Value))
			}
			finalHeadNode = &MultExp{BaseExp: BaseExp{}, HeadExp: &numImmExp}
		} else if imm, ok := finalHead.(*ImmExp); ok {
			finalHeadNode = &MultExp{BaseExp: BaseExp{}, HeadExp: imm}
		} else {
			// log.Printf("Warning: Cannot reconstruct AddExp head from type %T", finalHead)
			return a, false // Cannot simplify if unexpected head type
		}
	}

	// If only one term remains after reordering, construct AddExp with only the head.
	if len(finalOps) == 0 {
		// This handles the case where one non-constant term remains,
		// or only the constSum remained (which should have been handled earlier).
		simplifiedAddExp := NewAddExp(a.BaseExp, finalHeadNode, nil, nil)
		return simplifiedAddExp, reduced // Return true if any reduction happened
	}

	// Construct the simplified AddExp with multiple terms in the new order
	simplifiedAddExp := NewAddExp(a.BaseExp, finalHeadNode, finalOps, finalTailNodes)

	// Return the simplified expression, indicating reduction occurred
	return simplifiedAddExp, reduced
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
	HeadExp   Exp // Changed back to Exp interface
	Operators []string
	TailExps  []Exp // Changed back to Exp interface
}

func (m *MultExp) expressionNode() {}

// Need to regenerate constructor using `go generate ./...` after this change
// The generated constructor `NewMultExp` will now accept Exp for head and tails.

func (m *MultExp) Eval(env Env) (Exp, bool) {
	// Evaluate head expression
	evalHeadExp, headReduced := m.HeadExp.Eval(env) // HeadExp is Exp
	_, headIsNum := evalHeadExp.(*NumberExp)

	// Evaluate tail expressions
	evalTailExps := make([]Exp, len(m.TailExps)) // Store evaluated tails (Exp)
	anyTailReduced := false
	allTailsAreNumbers := true

	for i, tail := range m.TailExps { // TailExps are Exp
		evalTailExp, tailReduced := tail.Eval(env)
		evalTailExps[i] = evalTailExp
		if tailReduced {
			anyTailReduced = true
		}

		// Check if the evaluated tail is a number
		_, tailIsNum := evalTailExp.(*NumberExp)
		if !tailIsNum {
			allTailsAreNumbers = false
		}
		// No need to store evalTailNodes separately anymore
	}

	// If head and all tails evaluated to numbers, calculate the result
	if headIsNum && allTailsAreNumbers {
		currentValue := evalHeadExp.(*NumberExp).Value // Head is NumberExp
		for i, op := range m.Operators {
			numTail := evalTailExps[i].(*NumberExp) // Tails are NumberExp
			tailValue := numTail.Value
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
		// Pass the evaluated expressions (Exp interface) directly to the constructor
		return NewMultExp(m.BaseExp, evalHeadExp, m.Operators, evalTailExps), true
	}

	// No reduction possible, return original node
	return m, false
}
func (m *MultExp) TokenLiteral() string {
	head := m.HeadExp.TokenLiteral() // Call TokenLiteral() on HeadExp
	var buf strings.Builder
	buf.WriteString(head)
	for i, op := range m.Operators {
		buf.WriteByte(' ')
		buf.WriteString(op)
		buf.WriteByte(' ')
		tailStr := m.TailExps[i].TokenLiteral() // Call TokenLiteral() on TailExps[i]
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
		identValue := f.Value
		// Check for '$' first
		if identValue == "$" {
			// Use the GetLOC method from the Env interface
			dollarVal := int64(env.GetLOC()) // Use LOC (int32) as the value of $
			newFactor := NewNumberFactor(BaseFactor{}, int(dollarVal))
			numExp := NewNumberExp(ImmExp{BaseExp: imm.BaseExp, Factor: newFactor}, dollarVal)
			return numExp, true
			// No need for type assertion or else block here,
			// as GetLOC is now part of the Env interface.
		}
		// If not '$', check for macro
		macroExp, ok := env.LookupMacro(identValue)
		if ok {
			// Recursively evaluate the macro definition
			// Ensure the macro itself is evaluated
			evalMacroExp, reduced := macroExp.Eval(env)
			return evalMacroExp, reduced // Return the evaluated macro expression
		}
		// If not a macro or '$', it's an unresolved identifier (like a label)
		return imm, false // Return the ImmExp containing the IdentFactor
	case *StringFactor:
		// String factors themselves don't evaluate arithmetically,
		// but they are valid factors within an ImmExp. Return as is.
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

// Eval for NumberExp returns itself and true, indicating it's a fully evaluated value.
func (n *NumberExp) Eval(env Env) (Exp, bool) {
	return n, true // It's an evaluated value.
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

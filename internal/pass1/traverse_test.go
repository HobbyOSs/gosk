package pass1

import (
	"strconv" // Add strconv import
	"strings" // Add strings import for MYLABEL check
	"testing"

	"github.com/HobbyOSs/gosk/internal/ast" // Import ast package
	"github.com/HobbyOSs/gosk/internal/gen"
	"github.com/HobbyOSs/gosk/pkg/cpu" // Add cpu import
	"github.com/comail/colog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// Helper function to create a simple NumberExp (fully evaluated)
func newExpectedNumberExp(val int64) *ast.NumberExp {
	// BaseExp and Factor are needed for NumberExp structure but not critical for value comparison
	baseFactor := ast.NewNumberFactor(ast.BaseFactor{}, int(val))
	baseImmExp := ast.NewImmExp(ast.BaseExp{}, baseFactor)
	return ast.NewNumberExp(*baseImmExp, val)
}

// Helper function to create a simple IdentExp (ImmExp with IdentFactor)
func newExpectedIdentExp(name string) *ast.ImmExp {
	baseFactor := ast.NewIdentFactor(ast.BaseFactor{}, name)
	return ast.NewImmExp(ast.BaseExp{}, baseFactor)
}

// Helper function to create a simple MultExp for testing expectations
// Note: For simplicity, we assume factors are ImmExp(IdentFactor) for unresolved parts
func newExpectedMultExp(head ast.Exp, ops []string, tails []ast.Exp) *ast.MultExp { // Change head and tails to ast.Exp
	return ast.NewMultExp(ast.BaseExp{}, head, ops, tails)
}

// Helper function to create a simple AddExp for testing expectations
// Note: For simplicity, we assume factors are ImmExp(IdentFactor) for unresolved parts
func newExpectedAddExp(head *ast.MultExp, ops []string, tails []*ast.MultExp) *ast.AddExp {
	return ast.NewAddExp(ast.BaseExp{}, head, ops, tails)
}

// Helper to create MultExp from a single ImmExp (for AddExp structure)
func multExpFromImm(imm *ast.ImmExp) *ast.MultExp {
	return ast.NewMultExp(ast.BaseExp{}, imm, nil, nil)
}

// Helper to create ImmExp from a number string (used in AddExp expectation)
func immExpFromNumStr(numStr string) *ast.ImmExp {
	// This assumes the number string parses correctly to a NumberFactor
	// A more robust helper might handle potential parsing errors
	numVal, _ := strconv.Atoi(numStr)
	return ast.NewImmExp(ast.BaseExp{}, ast.NewNumberFactor(ast.BaseFactor{}, numVal))
}

type Pass1TraverseSuite struct { // Rename struct
	suite.Suite
}

func TestPass1TraverseSuite(t *testing.T) { // Rename test function
	suite.Run(t, new(Pass1TraverseSuite)) // Use renamed struct
}

func (s *Pass1TraverseSuite) SetupSuite() { // Use renamed struct
	setUpColog(colog.LDebug)
}

// They should exist in test_helper.go

func (s *Pass1TraverseSuite) TestAddExp() { // Use renamed struct
	tests := []struct {
		name         string
		text         string
		expectedNode ast.Exp // Expected evaluated node
	}{
		// --- Cases that should evaluate to NumberExp ---
		{
			name:         "+int",
			text:         "30",
			expectedNode: newExpectedNumberExp(30),
		},
		{
			name:         "-int",
			text:         "-30",
			expectedNode: newExpectedNumberExp(-30),
		},
		{
			name:         "hex",
			text:         "0x0ff0",
			expectedNode: newExpectedNumberExp(0x0ff0),
		},
		// CharFactor evaluation might need adjustment in ImmExp.Eval if needed
		// {
		// 	name:         "char",
		// 	text:         "'A'", // Simple char
		// 	expectedNode: newExpectedNumberExp(65),
		// },
		{
			name:         "simple math 1",
			text:         "1 + 1",
			expectedNode: newExpectedNumberExp(2),
		},
		{
			name:         "simple math 2",
			text:         "4 - 2",
			expectedNode: newExpectedNumberExp(2),
		},
		{
			name:         "simple math 3",
			text:         "1 + 3 - 2 + 4",
			expectedNode: newExpectedNumberExp(6),
		},
		{
			name:         "complex math 1 (evaluates fully)",
			text:         "8 * 3 - 1", // Parser creates AddExp{ MultExp{8*3}, "-", ImmExp{1} }
			expectedNode: newExpectedNumberExp(23),
		},
		{
			name:         "label + constant (evaluates fully)",
			text:         "MYLABEL + 512", // MYLABEL = 0x8000 (defined below)
			expectedNode: newExpectedNumberExp(0x8000 + 512),
		},
		{
			name:         "label - constant (evaluates fully)",
			text:         "MYLABEL - 10", // MYLABEL = 0x8000
			expectedNode: newExpectedNumberExp(0x8000 - 10),
		},
		// --- Cases that should NOT evaluate fully (contain unresolved identifiers) ---
		{
			name: "ident (cannot evaluate)",
			text: `_testZ009$`,
			// Expected: AddExp -> MultExp -> ImmExp -> IdentFactor
			expectedNode: newExpectedAddExp(
				multExpFromImm(newExpectedIdentExp(`_testZ009$`)),
				nil, nil,
			),
		},
		{
			name: "displacement 1 (cannot evaluate fully, constant folded)",
			text: "ESP+4",
			// Expected: AddExp{ MultExp{ImmExp{4}}, "+", MultExp{ImmExp{ESP}} }
			// Constant folding should place the constant 4 first.
			expectedNode: newExpectedAddExp(
				multExpFromImm(immExpFromNumStr("4")), // Head is the constant 4
				[]string{"+"},                         // Operator
				[]*ast.MultExp{
					multExpFromImm(newExpectedIdentExp("ESP")), // Tail is the non-constant term ESP
				},
			),
		},
		{
			name: "displacement 2 (cannot evaluate fully, constant folded)",
			text: "ESP+12+8",
			// Expected: AddExp{ MultExp{ImmExp{20}}, "+", MultExp{ImmExp{ESP}} }
			// Constant folding should combine 12 + 8 = 20 and place the constant first.
			expectedNode: newExpectedAddExp(
				multExpFromImm(immExpFromNumStr("20")), // Head is the constant sum 20
				[]string{"+"},                          // Operator
				[]*ast.MultExp{
					multExpFromImm(newExpectedIdentExp("ESP")), // Tail is the non-constant term ESP
				},
			),
		},
		// StringFactor and CharFactor with multiple chars are generally not evaluatable arithmetically
		// {
		// 	name: "string (cannot evaluate)",
		// 	text: `"0x0ff0"`,
		// 	expectedNode: ast.NewImmExp(ast.BaseExp{}, ast.NewStringFactor(ast.BaseFactor{}, "0x0ff0")),
		// },
		// {
		// 	name: "char multi (cannot evaluate)",
		// 	text: "'AB'",
		// 	expectedNode: ast.NewImmExp(ast.BaseExp{}, ast.NewCharFactor(ast.BaseFactor{}, "AB")),
		// },
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Adjust entrypoint based on test case
			entrypoint := "AddExp"
			if tt.name == "ident (cannot evaluate)" {
				entrypoint = "Exp" // Parse single identifier as Exp (ImmExp)
			}

			// Parse the input text
			got, err := gen.Parse("", []byte(tt.text), gen.Entrypoint(entrypoint))
			if !assert.NoError(t, err, "Parsing failed for input: %s", tt.text) {
				t.FailNow()
			}

			// Ensure the parsed node is actually an Exp
			startNode, ok := got.(ast.Exp)
			if !ok {
				t.Fatalf("Parsed node is not an ast.Exp, but %T", got)
			}

			// Setup Pass1 environment
			p := &Pass1{
				SymTable: make(map[string]int32),
				MacroMap: make(map[string]ast.Exp),
			}
			// Define MYLABEL macro for relevant tests
			if strings.Contains(tt.text, "MYLABEL") {
				// Use DefineMacro method which handles storing as ast.Exp
				p.DefineMacro("MYLABEL", newExpectedNumberExp(0x8000))
			}

			// Evaluate the node
			evaluatedNode := TraverseAST(startNode, p) // Pass1 implements ast.Env

			// Compare the evaluated node with the expected node
			switch expected := tt.expectedNode.(type) {
			case *ast.NumberExp:
				actual, ok := evaluatedNode.(*ast.NumberExp)
				assert.True(t, ok, "Expected *ast.NumberExp, got %T", evaluatedNode)
				if ok {
					assert.Equal(t, expected.Value, actual.Value, "Evaluated number value mismatch")
				}
			case *ast.AddExp:
				actual, ok := evaluatedNode.(*ast.AddExp)
				assert.True(t, ok, "Expected *ast.AddExp, got %T", evaluatedNode)
				if ok {
					// Basic comparison using TokenLiteral for structure check
					// More detailed comparison might be needed for complex cases
					assert.Equal(t, expected.TokenLiteral(), actual.TokenLiteral(), "Evaluated AddExp structure mismatch")
				}
			case *ast.ImmExp: // For cases like unresolved identifiers
				actual, ok := evaluatedNode.(*ast.ImmExp)
				assert.True(t, ok, "Expected *ast.ImmExp, got %T", evaluatedNode)
				if ok {
					assert.Equal(t, expected.TokenLiteral(), actual.TokenLiteral(), "Evaluated ImmExp mismatch")
				}
			default:
				t.Fatalf("Unhandled expected node type: %T", tt.expectedNode)
			}
		})
	}
}

// TestEQUExpansionInExpression verifies that TraverseAST correctly evaluates
// expressions containing constants defined by EQU statements, using the MacroMap.
func (s *Pass1TraverseSuite) TestEQUExpansionInExpression() {
	tests := []struct {
		name           string
		text           string      // Assembly code snippet including EQU
		expressionText string      // The specific expression part to evaluate after EQU processing
		expectedValue  int64       // Expected numerical value after evaluation
		bitMode        cpu.BitMode // Bit mode for Pass1 context
	}{
		{
			name: "Simple EQU constant evaluation",
			text: `
				MY_EQU_CONST EQU 500
				MOV AX, MY_EQU_CONST ; Evaluate MY_EQU_CONST
			`,
			expressionText: "MY_EQU_CONST",
			expectedValue:  500,
			bitMode:        cpu.MODE_16BIT,
		},
		{
			name: "EQU constant + number evaluation",
			text: `
				MY_EQU_CONST2 EQU 100
				ADD BX, MY_EQU_CONST2 + 20 ; Evaluate MY_EQU_CONST2 + 20
			`,
			expressionText: "MY_EQU_CONST2 + 20",
			expectedValue:  120,
			bitMode:        cpu.MODE_16BIT,
		},
		{
			name: "EQU constant used in another EQU",
			text: `
				BASE_VAL EQU 1000
				OFFSET_VAL EQU BASE_VAL + 50
				MOV CX, OFFSET_VAL ; Evaluate OFFSET_VAL
			`,
			expressionText: "OFFSET_VAL",
			expectedValue:  1050,
			bitMode:        cpu.MODE_16BIT,
		},
		{
			name: "EQU constant in multiplication",
			text: `
				FACTOR EQU 8
				IMUL DX, FACTOR * 2 ; Evaluate FACTOR * 2
			`,
			expressionText: "FACTOR * 2",
			expectedValue:  16,
			bitMode:        cpu.MODE_16BIT,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// 1. Parse the whole snippet as Program to process EQU
			parseTree, err := gen.Parse("", []byte(tt.text), gen.Entrypoint("Program"))
			s.Require().NoError(err, "Parsing program snippet failed")
			// Assert to the concrete type *ast.Program to access Statements field
			program, ok := parseTree.(*ast.Program)
			s.Require().True(ok, "Parsed result is not *ast.Program")

			// 2. Setup Pass1 environment
			p := &Pass1{
				LOC:      0,
				BitMode:  tt.bitMode,
				SymTable: make(map[string]int32),
				MacroMap: make(map[string]ast.Exp),
				// Client and AsmDB are not needed for pure expression evaluation via TraverseAST
			}

			// 3. Process EQU statements to populate MacroMap
			// This simulates the part of Pass1.Eval that handles EQU.
			// We need to evaluate the EQU expression itself first.
			// Now we can directly access program.Statements
			for _, stmt := range program.Statements {
				// Check if the statement is an EQU statement (represented by DeclareStmt)
				if declareStmt, ok := stmt.(*ast.DeclareStmt); ok {
					// Evaluate the expression assigned in EQU using the current environment (p)
					// This handles cases where EQU depends on previous EQUs.
					evaluatedEquNode := TraverseAST(declareStmt.Value, p)  // Use declareStmt.Value
					evaluatedEquExpr, okExpr := evaluatedEquNode.(ast.Exp) // Assert to ast.Exp
					if !okExpr {
						t.Fatalf("EQU expression evaluation did not return an ast.Exp: %T", evaluatedEquNode)
					}
					_, isEvaluable := evaluatedEquExpr.(*ast.NumberExp) // Check if it evaluated to a number
					if !isEvaluable {
						// If EQU expression itself couldn't be fully evaluated (e.g., contains labels),
						// store the partially evaluated expression. For these tests, we assume EQUs resolve to numbers.
						t.Logf("Warning: EQU expression for %s did not evaluate to a number: %T", declareStmt.Id.Value, evaluatedEquExpr) // Use declareStmt.Id.Value
					}
					p.DefineMacro(declareStmt.Id.Value, evaluatedEquExpr) // Use declareStmt.Id.Value and pass the asserted ast.Exp
				}
			}

			// 4. Parse the target expression string separately
			// We need to parse the specific expression we want to test evaluation for.
			exprTree, err := gen.Parse("", []byte(tt.expressionText), gen.Entrypoint("Exp")) // Use "Exp" entrypoint
			s.Require().NoError(err, "Parsing target expression failed: %s", tt.expressionText)
			exprToEval, ok := exprTree.(ast.Exp)
			s.Require().True(ok, "Parsed expression is not ast.Exp")

			// 5. Evaluate the target expression using TraverseAST and the populated Pass1 env
			evaluatedNode := TraverseAST(exprToEval, p)      // Returns ast.Node
			evaluatedExpr, okExpr := evaluatedNode.(ast.Exp) // Assert to ast.Exp
			if !okExpr {
				t.Fatalf("Target expression evaluation did not return an ast.Exp: %T", evaluatedNode)
			}

			// 6. Assert the result is a NumberExp with the expected value
			expectedNode := newExpectedNumberExp(tt.expectedValue)
			actualNode, ok := evaluatedExpr.(*ast.NumberExp) // Use the asserted evaluatedExpr
			s.True(ok, "Expected evaluated node to be *ast.NumberExp, got %T for expression '%s'", evaluatedExpr, tt.expressionText)
			if ok {
				s.Equal(expectedNode.Value, actualNode.Value, "Evaluated value mismatch for expression '%s'", tt.expressionText)
			}
		})
	}
}

func (s *Pass1TraverseSuite) TestMultExp() { // Use renamed struct
	tests := []struct {
		name         string
		text         string
		expectedNode ast.Exp // Expected evaluated node
	}{
		// --- Cases that should evaluate to NumberExp ---
		{
			name:         "simple math 1",
			text:         "1005*8",
			expectedNode: newExpectedNumberExp(8040),
		},
		{
			name:         "simple math 2",
			text:         "512/4",
			expectedNode: newExpectedNumberExp(128),
		},
		{
			name:         "simple math 3",
			text:         "512*1024/4",
			expectedNode: newExpectedNumberExp(131072),
		},
		// --- Cases that should NOT evaluate fully ---
		{
			name: "scale 1 (cannot evaluate fully)",
			text: "EDX*4",
			// Expected: MultExp{ ImmExp{EDX}, "*", ImmExp{4} }
			expectedNode: newExpectedMultExp(
				newExpectedIdentExp("EDX"), // ast.Exp compatible
				[]string{"*"},
				[]ast.Exp{immExpFromNumStr("4")}, // Pass as []ast.Exp
			),
		},
		{
			name: "scale 2 (cannot evaluate fully)",
			text: "ESI*8",
			// Expected: MultExp{ ImmExp{ESI}, "*", ImmExp{8} }
			expectedNode: newExpectedMultExp(
				newExpectedIdentExp("ESI"), // ast.Exp compatible
				[]string{"*"},
				[]ast.Exp{immExpFromNumStr("8")}, // Pass as []ast.Exp
			),
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Parse the input text as a MultExp
			got, err := gen.Parse("", []byte(tt.text), gen.Entrypoint("MultExp"))
			if !assert.NoError(t, err, "Parsing failed for input: %s", tt.text) {
				t.FailNow()
			}

			// Ensure the parsed node is actually an Exp
			startNode, ok := got.(ast.Exp)
			if !ok {
				t.Fatalf("Parsed node is not an ast.Exp, but %T", got)
			}

			// Setup Pass1 environment
			p := &Pass1{
				SymTable: make(map[string]int32),
				MacroMap: make(map[string]ast.Exp),
			}

			// Evaluate the node
			evaluatedNode := TraverseAST(startNode, p) // Pass1 implements ast.Env

			// Compare the evaluated node with the expected node
			switch expected := tt.expectedNode.(type) {
			case *ast.NumberExp:
				actual, ok := evaluatedNode.(*ast.NumberExp)
				assert.True(t, ok, "Expected *ast.NumberExp, got %T", evaluatedNode)
				if ok {
					assert.Equal(t, expected.Value, actual.Value, "Evaluated number value mismatch")
				}
			case *ast.MultExp:
				actual, ok := evaluatedNode.(*ast.MultExp)
				assert.True(t, ok, "Expected *ast.MultExp, got %T", evaluatedNode)
				if ok {
					// Basic comparison using TokenLiteral for structure check
					assert.Equal(t, expected.TokenLiteral(), actual.TokenLiteral(), "Evaluated MultExp structure mismatch")
				}
			default:
				t.Fatalf("Unhandled expected node type: %T", tt.expectedNode)
			}
		})
	}
}

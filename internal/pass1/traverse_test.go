package pass1

import (
	"strconv" // Add strconv import
	"strings" // Add strings import for MYLABEL check
	"testing"

	"github.com/HobbyOSs/gosk/internal/ast" // Import ast package
	"github.com/HobbyOSs/gosk/internal/gen"
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
			name: "displacement 1 (cannot evaluate fully)",
			text: "ESP+4",
			// Expected: AddExp{ MultExp{ImmExp{ESP}}, "+", MultExp{ImmExp{4}} }
			expectedNode: newExpectedAddExp(
				multExpFromImm(newExpectedIdentExp("ESP")),
				[]string{"+"},
				[]*ast.MultExp{multExpFromImm(immExpFromNumStr("4"))},
			),
		},
		{
			name: "displacement 2 (cannot evaluate fully)",
			text: "ESP+12+8",
			// Expected: AddExp{ MultExp{ImmExp{ESP}}, "+", MultExp{ImmExp{12}}, "+", MultExp{ImmExp{8}} }
			expectedNode: newExpectedAddExp(
				multExpFromImm(newExpectedIdentExp("ESP")),
				[]string{"+", "+"},
				[]*ast.MultExp{
					multExpFromImm(immExpFromNumStr("12")),
					multExpFromImm(immExpFromNumStr("8")),
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

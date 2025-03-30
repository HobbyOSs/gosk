package pass1

import (
	"log"
	"strings"

	"github.com/HobbyOSs/gosk/internal/ast"
	"github.com/HobbyOSs/gosk/pkg/cpu"
)

// TraverseAST now takes ast.Node and *Pass1 (as Env) and returns the potentially transformed ast.Node.
// It no longer uses the stack (env.Ctx).
func TraverseAST(node ast.Node, env *Pass1) ast.Node {
	if node == nil {
		return nil
	}

	// Implement the Env interface for Pass1
	// This allows passing 'env' directly to Eval methods.
	var evalEnv ast.Env = env // Pass1 implements ast.Env via DefineMacro/LookupMacro methods below

	switch n := node.(type) {
	case *ast.Program:
		log.Println("trace: program handler!!!")
		newStatements := make([]ast.Statement, 0, len(n.Statements))
		for _, stmt := range n.Statements {
			// Traverse each statement. TraverseAST now returns Node.
			processedStmt := TraverseAST(stmt, env)
			if processedStmt != nil {
				// Ensure the returned node is indeed a Statement
				if statement, ok := processedStmt.(ast.Statement); ok {
					newStatements = append(newStatements, statement)
				} else {
					// If TraverseAST returns an evaluated expression (like NumberExp from EQU), discard it at the statement level.
					log.Printf("info: TraverseAST returned a non-Statement node (%T) for a statement, discarding.", processedStmt)
				}
			}
		}
		// Return a new Program node with modified statements.
		return ast.NewProgram(newStatements) // Corrected call to NewProgram

	case *ast.DeclareStmt: // EQU statement
		log.Println("trace: declare handler!!!")
		// Evaluate the value expression first.
		evalValueNode := TraverseAST(n.Value, env)
		evalValueExp, ok := evalValueNode.(ast.Exp)
		if !ok {
			log.Printf("error: EQU value expression %s evaluated to non-expression type %T", n.Value.TokenLiteral(), evalValueNode)
			return nil // Or handle error appropriately
		}

		// Define the macro in the environment using the method on Pass1.
		env.DefineMacro(n.Id.Value, evalValueExp)
		log.Printf("debug: Defined macro '%s' = %s", n.Id.Value, evalValueExp.TokenLiteral())
		// EQU statement itself doesn't produce output, so return nil.
		return nil

	case *ast.LabelStmt:
		log.Println("trace: label handler!!!")
		label := strings.TrimSuffix(n.Label.Value, ":")
		env.SymTable[label] = env.LOC
		// Label statement itself doesn't produce output node after processing.
		return nil // Or return n if labels should remain in the AST for pass2

	case *ast.MnemonicStmt:
		log.Println("trace: mnemonic stmt handler!!!")
		opcode := n.Opcode.Value

		// Evaluate operands first using TraverseAST -> Eval
		evalOperands := make([]ast.Exp, len(n.Operands))
		canProcess := true
		for i, operand := range n.Operands {
			evalOperandNode := TraverseAST(operand, env)
			if expOperand, ok := evalOperandNode.(ast.Exp); ok {
				evalOperands[i] = expOperand
			} else {
				log.Printf("error: Operand %d for %s evaluated to non-expression type %T", i, opcode, evalOperandNode)
				canProcess = false
				break // Stop processing if an operand is invalid
			}
		}

		if !canProcess {
			// Cannot process this instruction if operands are invalid
			return n // Return original node if operands are invalid
		}

		// Find and call the appropriate handler function
		if handler, ok := opcodeEvalFns[opcode]; ok {
			handler(env, evalOperands) // Call handler with evaluated operands
		} else {
			log.Printf("error: No handler found for opcode %s", opcode)
			// Decide how to handle unknown opcodes (e.g., skip, error, default size?)
			// For now, just log. LOC won't be updated.
		}

		// LOC and Emit are handled within the specific handler function now.
		return nil // Mnemonic statement processed, return nil

	case *ast.OpcodeStmt: // Instruction without operands
		log.Println("trace: opcode stmt handler!!!")
		opcode := n.Opcode.Value

		// Find and call the appropriate handler function (passing empty operands)
		if handler, ok := opcodeEvalFns[opcode]; ok {
			handler(env, []ast.Exp{}) // Call handler with empty operands
		} else {
			log.Printf("error: No handler found for opcode %s", opcode)
		}

		// LOC and Emit are handled within the specific handler function now.
		return nil // Opcode statement processed, return nil

	// --- Expression Evaluation ---
	case *ast.AddExp, *ast.MultExp, *ast.ImmExp, *ast.SegmentExp, *ast.MemoryAddrExp:
		if exp, ok := n.(ast.Exp); ok {
			evalExp, _ := exp.Eval(evalEnv) // Use evalEnv which is ast.Env type
			return evalExp                  // Return the evaluated expression node
		}
		log.Printf("error: Node %T claims to be Exp but type assertion failed.", n)
		return n // Return original node on error

	// --- Factor Handling ---
	case *ast.NumberFactor, *ast.StringFactor, *ast.HexFactor, *ast.IdentFactor, *ast.CharFactor:
		log.Printf("warning: TraverseAST encountered a Factor type (%T) directly. Wrapping in ImmExp.", n)
		// Wrap factor in ImmExp before returning, as factors should be part of expressions.
		return ast.NewImmExp(ast.BaseExp{}, n.(ast.Factor))

	// --- Other Statement Types ---
	case *ast.ExportSymStmt:
		log.Println("trace: export sym stmt handler!!!")
		// TODO: Implement logic if needed (e.g., add to env.GlobalSymbolList)
		return nil // Or return n
	case *ast.ExternSymStmt:
		log.Println("trace: extern sym stmt handler!!!")
		// TODO: Implement logic if needed (e.g., add to env.ExternSymbolList)
		return nil // Or return n
	case *ast.ConfigStmt:
		log.Println("trace: config stmt handler!!!")
		if n.ConfigType == ast.Bits {
			// Evaluate the factor to get the bit mode value
			factorNode := TraverseAST(n.Factor, env)
			// Factor should be wrapped in ImmExp by the Factor case above
			immExp, ok := factorNode.(*ast.ImmExp)
			if !ok {
				log.Printf("error: BITS directive requires a constant value, got %T", factorNode)
				return nil
			}
			evalExp, _ := immExp.Eval(evalEnv) // Use evalEnv
			numExp, ok := evalExp.(*ast.NumberExp)
			if !ok {
				log.Printf("error: BITS directive value did not evaluate to a number: %s", evalExp.TokenLiteral())
				return nil
			}

			bitModeVal := int(numExp.Value)
			bitMode, ok := cpu.NewBitMode(bitModeVal)
			if !ok {
				log.Printf("error: Invalid bit mode value %d for BITS directive", bitModeVal)
				return nil
			}
			env.BitMode = bitMode
			env.Client.SetBitMode(bitMode)
			log.Printf("debug: Set bit mode to %d", bitModeVal)
		}
		return nil // Config statement doesn't produce an output node

	default:
		log.Printf("Unknown AST node type in TraverseAST: %T\n", node)
		return node // Return unknown nodes unchanged
	}
	// return node // Should not be reached
}

// DefineMacro implements the ast.Env interface for Pass1 by defining it as a method.
func (p *Pass1) DefineMacro(name string, exp ast.Exp) {
	// Initialize the new map if it's nil
	if p.MacroMap == nil {
		p.MacroMap = make(map[string]ast.Exp)
	}
	p.MacroMap[name] = exp
	log.Printf("debug: Defined macro '%s' = %s (stored as ast.Exp)", name, exp.TokenLiteral())

}

// LookupMacro implements the ast.Env interface for Pass1 by defining it as a method.
func (p *Pass1) LookupMacro(name string) (ast.Exp, bool) {
	// Use the new MacroMap
	if p.MacroMap == nil {
		return nil, false // Map not initialized
	}
	exp, ok := p.MacroMap[name]
	// No fallback to old EquMap needed here for Eval logic
	return exp, ok
}

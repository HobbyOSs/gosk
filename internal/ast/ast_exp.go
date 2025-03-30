package ast

// Env defines the interface for the environment needed during evaluation.
// It provides methods to lookup macros or other definitions.
type Env interface {
	LookupMacro(name string) (Exp, bool)
	// Add other methods needed for evaluation if necessary
}

// Exp represents an expression node in the AST.
type Exp interface {
	Node
	expressionNode()
	Type() string
	// Eval evaluates the expression within the given environment.
	// It returns the evaluated/reduced expression and a boolean indicating
	// whether the evaluation resulted in a reduction (true) or not (false).
	// If evaluation is not possible (e.g., involves unresolved identifiers),
	// it should return the original expression node and false.
	Eval(env Env) (Exp, bool)
}

// BaseExp provides a base implementation, but Type() should be implemented by concrete types.
type BaseExp struct{}

// Type returns the name of the concrete expression type.
// NOTE: This generic implementation using reflection might not be reliable.
// It's recommended that each concrete Exp type implements its own Type() method.
func (b BaseExp) Type() string {
	// This reflection-based approach is often problematic.
	// Returning a placeholder or panicking might be safer.
	// return reflect.TypeOf(b).Elem().Name() // Original problematic implementation
	panic("BaseExp.Type() should not be called directly. Implement Type() in concrete Exp types.")

}

// expressionNode() is a marker method to identify expression nodes.
// We can add it to BaseExp if needed for embedding structs.
// func (b BaseExp) expressionNode() {}

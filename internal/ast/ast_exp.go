package ast

// Env は評価中に必要な環境のインターフェースを定義します。
// マクロやその他の定義を検索するメソッドを提供します。
type Env interface {
	LookupMacro(name string) (Exp, bool)
	// GetLOC は '$' 評価に必要な現在のロケーションカウンタ値を返します。
	GetLOC() int32
	// GetConstValue は、式が定数である場合に整数値を抽出します。
	// これは通常、環境 (例: Pass1) によって実装されます。
	GetConstValue(exp Exp) (int, bool)
	// 必要に応じて評価に必要な他のメソッドを追加します
}

// Exp は AST 内の式ノードを表します。
type Exp interface {
	Node
	expressionNode()
	Type() string
	// Eval は指定された環境内で式を評価します。
	// 評価/簡約された式と、評価が簡約をもたらしたかどうか (true)
	// そうでないか (false) を示すブール値を返します。
	// 評価が不可能な場合 (例: 未解決の識別子を含む)、
	// 元の式ノードと false を返す必要があります。
	Eval(env Env) (Exp, bool)
}

// BaseExp は基本実装を提供しますが、Type() は具象型によって実装される必要があります。
type BaseExp struct{}

// Type は具象式の型の名前を返します。
// 注意: リフレクションを使用したこの汎用実装は信頼できない可能性があります。
// 各具象 Exp 型が独自の Type() メソッドを実装することを推奨します。
func (b BaseExp) Type() string {
	// このリフレクションベースのアプローチはしばしば問題があります。
	// プレースホルダーを返すか、パニックする方が安全かもしれません。
	panic("BaseExp.Type() should not be called directly. Implement Type() in concrete Exp types.")
}

// expressionNode() は式ノードを識別するためのマーカーメソッドです。
// 構造体を埋め込むために必要であれば BaseExp に追加できます。

package ast

import (
	"strconv"
	"strings"
)

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
	Right    *AddExp // nullable (nil許容)
}

func (s *SegmentExp) expressionNode() {}
func (s *SegmentExp) Eval(env Env) (Exp, bool) {
	// TODO: SegmentExp の評価ロジックを実装する
	// 現時点では、ノード自体を返し、簡約がないことを示します。
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

// wrapExpInAddExp は、単純な Exp (NumberExp, ImmExp) を AddExp フィールドに必要な構造にラップします。
func wrapExpInAddExp(exp Exp) *AddExp {
	if exp == nil {
		return nil
	}
	if addExp, ok := exp.(*AddExp); ok {
		// すでに AddExp の場合はそれを返します
		return addExp
	}

	var immExp *ImmExp
	if numExp, ok := exp.(*NumberExp); ok {
		// NumberExp に有効な ImmExp が埋め込まれていることを確認します
		if numExp.ImmExp.Factor == nil {
			// Factor が欠落している場合 (NewNumberExp では発生しないはず)、作成します
			numExp.ImmExp.Factor = NewNumberFactor(BaseFactor{}, int(numExp.Value))
		}
		immExp = &numExp.ImmExp
	} else if ie, ok := exp.(*ImmExp); ok {
		immExp = ie
	} else {
		// MultExp のような他の型をここで直接ラップするのは簡単ではありません。
		// このヘルパーは主に Eval からの NumberExp/ImmExp の結果用です。
		return nil // ラップが簡単でない場合は nil を返します
	}

	// MultExp -> AddExp 構造を作成します
	// NewMultExp と NewAddExp のために BaseExp が初期化されていることを確認します
	multExp := NewMultExp(BaseExp{}, immExp, nil, nil) // Head は ImmExp
	addExp := NewAddExp(BaseExp{}, multExp, nil, nil)  // Head のみを持つ AddExp
	return addExp
}

//go:generate newc
type MemoryAddrExp struct {
	BaseExp
	DataType DataType
	JumpType JumpType
	Left     *AddExp
	Right    *AddExp // nullable (nil許容)
}

func (m *MemoryAddrExp) expressionNode() {}
func (m *MemoryAddrExp) Eval(env Env) (Exp, bool) {
	// 内部の式を評価します
	evalLeftNode, leftReduced := m.Left.Eval(env) // Exp を返します
	evalRightNode := Exp(nil)                     // evalRightNode を初期化します
	rightReduced := false
	if m.Right != nil {
		evalRightNode, rightReduced = m.Right.Eval(env) // Exp を返します
	}

	// 可能であれば、評価されたノードを AddExp 構造にラップし直します
	evalLeftExp := wrapExpInAddExp(evalLeftNode)
	if evalLeftExp == nil && leftReduced {
		// ラップに失敗したが簡約が発生した場合、状態を表すことができません。
		// 情報の損失や無効な構造の作成を避けるために、オリジナルを返します。
		return m, false
	} else if evalLeftExp == nil {
		evalLeftExp = m.Left // 簡約がなく、ラップも不可能な場合はオリジナルを保持します
	}

	evalRightExp := (*AddExp)(nil)
	if m.Right != nil {
		evalRightExp = wrapExpInAddExp(evalRightNode)
		if evalRightExp == nil && rightReduced {
			return m, false
		} else if evalRightExp == nil {
			evalRightExp = m.Right // オリジナルを保持します
		}
	}

	// どちらの内部式も簡約されなかった場合は、元のノードを返します
	if !leftReduced && !rightReduced {
		return m, false
	}

	// ラップされる可能性のある内部式を持つ新しい MemoryAddrExp を構築します
	newMemExp := NewMemoryAddrExp(m.BaseExp, m.DataType, m.JumpType, evalLeftExp, evalRightExp)
	return newMemExp, true // 新しいノードを返し、簡約が発生したことを示します
}
func (m *MemoryAddrExp) TokenLiteral() string {
	// ast_exp_string.go の既存の ExpToString を使用します
	var str = ""
	if m.DataType != None {
		str += string(m.DataType)
		str += " "
	}
	str += "[ "
	// 評価される可能性のある Left 式を処理するために ExpToString を使用します
	str += ExpToString(m.Left) // m.Left は Eval が置き換えなかった場合、オリジナルを指す可能性があります
	if m.Right != nil {
		str += " : "
		// Right にも ExpToString を使用します
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

// Eval は AddExp の定数畳み込みを実行します。
// すべての定数項を合計し、非定数項を保持します。
// 非定数項を最初に保持するように変更されました。
func (a *AddExp) Eval(env Env) (Exp, bool) {
	// 最初に head 式を評価します
	evalHead, headReduced := a.HeadExp.Eval(env)

	// 定数項の合計と非定数項/演算子のリストを追跡します
	constSum := 0
	newTerms := []Exp{} // 評価された非定数項を保持するために Exp インターフェースを使用します
	newOps := []string{}
	reduced := headReduced // head の簡約ステータスから開始します

	// 評価された head を処理します
	if v, ok := env.GetConstValue(evalHead); ok {
		constSum += v
	} else {
		newTerms = append(newTerms, evalHead)
	}

	// tail 式を処理します
	for i, op := range a.Operators {
		tail := a.TailExps[i]
		evalTail, tailReduced := tail.Eval(env)
		if tailReduced {
			reduced = true
		}

		if v, ok := env.GetConstValue(evalTail); ok {
			// 定数の場合は、合計に加算/減算します
			if op == "+" {
				constSum += v
			} else if op == "-" {
				constSum -= v
			} else {
				// 文法に基づいて発生しないはずですが、防御的に処理します
				// サポートされていない演算子が定数と共に現れた場合、簡約不可として扱います
				// 演算子と定数項を保持します
				if len(newTerms) > 0 { // 先行する項がある場合にのみ演算子を追加します
					newOps = append(newOps, op)
				}
				newTerms = append(newTerms, evalTail) // 定数項を戻します
			}
		} else {
			// 定数でない場合は、項のリストに追加します
			// 先行する項があった場合にのみ演算子を追加します
			if len(newTerms) > 0 { // 最初の項でない場合に演算子を追加します
				newOps = append(newOps, op)
			}
			newTerms = append(newTerms, evalTail)
		}
	}

	// --- 結果の構築 ---

	// ケース 1: すべての項が定数に評価された場合
	if len(newTerms) == 0 {
		// 最終的な合計を持つ単一の NumberExp を返します
		return NewNumberExp(ImmExp{BaseExp: a.BaseExp}, int64(constSum)), true
	}

	// ケース 2: 定数と非定数が混在する場合

	// 項の順序を変更: 非定数を最初に、次に定数の合計 (ゼロでない場合)
	finalTerms := []Exp{}
	finalOps := []string{}

	// 最初に非定数項を追加します
	if len(newTerms) > 0 {
		finalTerms = append(finalTerms, newTerms...)
		finalOps = append(finalOps, newOps...) // 非定数項間の元の演算子を保持します
	}

	// ゼロでない場合は、最後に定数の合計を追加します
	if constSum != 0 {
		constTerm := NewNumberExp(ImmExp{BaseExp: a.BaseExp}, int64(constSum))
		if len(finalTerms) > 0 {
			// 他の項が存在する場合、定数項の前に '+' または '-' 演算子を追加します
			if constSum > 0 {
				finalOps = append(finalOps, "+")
			} else {
				finalOps = append(finalOps, "-")
				// 演算子が '-' の場合は NumberExp に絶対値を使用します
				constTerm = NewNumberExp(ImmExp{BaseExp: a.BaseExp}, int64(-constSum))
			}
		} else {
			// 定数項のみが存在する場合は、それを唯一の項にします
			// 唯一の項である場合は、負の定数を head として処理します
			if constSum < 0 {
				// このケースは理想的にはケース 1 で NumberExp を返すことで処理されるべきですが、
				// ここに到達した場合に NumberExp の値が正しいことを確認します。
				constTerm = NewNumberExp(ImmExp{BaseExp: a.BaseExp}, int64(constSum))
			}
		}
		finalTerms = append(finalTerms, constTerm)
	}

	// 順序変更後、1 つの項のみが残る場合 (非定数または constSum の可能性がある)、可能であれば直接返します
	if len(finalTerms) == 1 && len(finalOps) == 0 {
		// NumberExp の場合はそれを返します (理想的にはケース 1 で処理済み)
		if numExp, ok := finalTerms[0].(*NumberExp); ok {
			return numExp, true
		}
		// 単一の非定数項の場合は、以下で AddExp 構造にラップする必要があります
	}

	// 項が残らない場合 (例: "LABEL - LABEL")、結果は 0 です
	if len(finalTerms) == 0 {
		return NewNumberExp(ImmExp{BaseExp: a.BaseExp}, 0), true
	}

	// --- 新しい順序で AddExp を再構築 ---

	// 順序変更後、1 つの項のみが残る場合、その項を直接返します
	if len(finalTerms) == 1 && len(finalOps) == 0 {
		// NumberExp の場合はそれを返します (理想的にはケース 1 で処理済み)
		// それ以外 (ImmExp, MultExp など) の場合も直接返します
		return finalTerms[0], reduced // Return the single term directly
	}

	// --- 複数の項が残る場合、新しい AddExp を構築 ---
	// finalTerms の最初の項が新しい head です
	finalHead := finalTerms[0]

	// 残りの評価された項を AddExp 構造のために *MultExp に変換し直します
	finalTailNodes := make([]*MultExp, 0, len(finalTerms)-1)
	for _, term := range finalTerms[1:] { // 順序変更された finalTerms を反復処理します
		if me, ok := term.(*MultExp); ok {
			finalTailNodes = append(finalTailNodes, me)
		} else if num, ok := term.(*NumberExp); ok {
			// NumberExp を MultExp にラップし直します
			// 埋め込まれた ImmExp が正しい Factor を持っていることを確認します
			numImmExp := num.ImmExp
			if numImmExp.Factor == nil {
				numImmExp.Factor = NewNumberFactor(BaseFactor{}, int(num.Value))
			}
			finalTailNodes = append(finalTailNodes, &MultExp{BaseExp: BaseExp{}, HeadExp: &numImmExp})
		} else if imm, ok := term.(*ImmExp); ok {
			// ImmExp (識別子など) を MultExp にラップします
			finalTailNodes = append(finalTailNodes, &MultExp{BaseExp: BaseExp{}, HeadExp: imm})
		} else {
			// MultExp に簡単に入れることができない他の Exp 型の場合は、
			// 完全に簡約できない可能性があります。オリジナルを返すかエラーを返します。
			// 現時点では、項は MultExp、NumberExp、または ImmExp であると仮定します。
			// 予期しない型が見つかった場合はオリジナルを返します。
			return a, false // 予期しない型が見つかった場合は簡約できません
		}
	}

	// finalHead (Exp) を *MultExp に変換し直します
	finalHeadNode, ok := finalHead.(*MultExp)
	if !ok {
		if num, ok := finalHead.(*NumberExp); ok {
			// 埋め込まれた ImmExp が正しい Factor を持っていることを確認します
			numImmExp := num.ImmExp
			if numImmExp.Factor == nil {
				numImmExp.Factor = NewNumberFactor(BaseFactor{}, int(num.Value))
			}
			finalHeadNode = &MultExp{BaseExp: BaseExp{}, HeadExp: &numImmExp}
		} else if imm, ok := finalHead.(*ImmExp); ok {
			finalHeadNode = &MultExp{BaseExp: BaseExp{}, HeadExp: imm}
		} else {
			return a, false // 予期しない head 型の場合は簡約できません
		}
	}

	// 新しい順序で複数の項を持つ簡約された AddExp を構築します
	simplifiedAddExp := NewAddExp(a.BaseExp, finalHeadNode, finalOps, finalTailNodes) // finalHeadNode を使用

	// 簡約された式を返し、簡約が発生したことを示します
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
	HeadExp   Exp // Exp インターフェースに戻しました
	Operators []string
	TailExps  []Exp // Exp インターフェースに戻しました
}

func (m *MultExp) expressionNode() {}

// この変更後、`go generate ./...` を使用してコンストラクタを再生成する必要があります
// 生成されたコンストラクタ `NewMultExp` は、head と tails に Exp を受け入れるようになります。

func (m *MultExp) Eval(env Env) (Exp, bool) {
	// head 式を評価します
	evalHeadExp, headReduced := m.HeadExp.Eval(env) // HeadExp は Exp です

	// tail がない場合、評価された head を直接返します
	if len(m.Operators) == 0 {
		return evalHeadExp, headReduced // Return evaluated head and its reduction status
	}

	// --- Tails が存在する場合 ---
	_, headIsNum := evalHeadExp.(*NumberExp)
	evalTailExps := make([]Exp, len(m.TailExps)) // 評価された tails (Exp) を格納します
	anyTailReduced := false
	allTailsAreNumbers := true

	for i, tail := range m.TailExps { // TailExps は Exp です
		evalTailExp, tailReduced := tail.Eval(env)
		evalTailExps[i] = evalTailExp
		if tailReduced {
			anyTailReduced = true
		}

		// 評価された tail が数値かどうかを確認します
		_, tailIsNum := evalTailExp.(*NumberExp)
		if !tailIsNum {
			allTailsAreNumbers = false
		}
		// evalTailNodes を個別に格納する必要はもうありません
	}

	// head とすべての tails が数値に評価された場合、結果を計算します
	if headIsNum && allTailsAreNumbers {
		currentValue := evalHeadExp.(*NumberExp).Value // Head は NumberExp です
		for i, op := range m.Operators {
			numTail := evalTailExps[i].(*NumberExp) // Tails は NumberExp です // ここで panic する可能性: evalTailExps[i] が *NumberExp でない場合
			tailValue := numTail.Value
			switch op {
			case "*":
				currentValue *= tailValue
			case "/":
				if tailValue == 0 {
					return m, false // ゼロ除算
				}
				currentValue /= tailValue
			case "%":
				if tailValue == 0 {
					return m, false // ゼロによる剰余
				}
				currentValue %= tailValue
			default:
				return m, false // サポートされていない演算子
			}
		}
		// 新しい NumberExp を返します
		return NewNumberExp(ImmExp{BaseExp: m.BaseExp}, currentValue), true
	}

	// すべての部分が数値に評価されなかったが、何らかの簡約が発生した場合は、更新された MultExp を返します
	if headReduced || anyTailReduced {
		// Rebuild MultExp only if tails exist (tails の存在は上でチェック済み)
		return NewMultExp(m.BaseExp, evalHeadExp, m.Operators, evalTailExps), true
	}

	// 簡約不可、元のノードを返します
	return m, false
}
func (m *MultExp) TokenLiteral() string {
	head := m.HeadExp.TokenLiteral() // HeadExp で TokenLiteral() を呼び出します
	var buf strings.Builder
	buf.WriteString(head)
	for i, op := range m.Operators {
		buf.WriteByte(' ')
		buf.WriteString(op)
		buf.WriteByte(' ')
		tailStr := m.TailExps[i].TokenLiteral() // TailExps[i] で TokenLiteral() を呼び出します
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
		// 最初に '$' をチェックします
		if identValue == "$" {
			// Env インターフェースから GetLOC メソッドを使用します
			dollarVal := int64(env.GetLOC()) // LOC (int32) を $ の値として使用します
			newFactor := NewNumberFactor(BaseFactor{}, int(dollarVal))
			numExp := NewNumberExp(ImmExp{BaseExp: imm.BaseExp, Factor: newFactor}, dollarVal)
			return numExp, true
			// GetLOC は Env インターフェースの一部になったため、
			// ここで型アサーションや else ブロックは不要です。
		}
		// '$' でない場合は、マクロをチェックします
		macroExp, ok := env.LookupMacro(identValue)
		if ok {
			// マクロ定義を再帰的に評価します
			// マクロ自体が評価されることを確認します
			evalMacroExp, reduced := macroExp.Eval(env)
			return evalMacroExp, reduced // 評価されたマクロ式を返します
		}
		// マクロでも '$' でもない場合は、未解決の識別子 (ラベルなど) です
		return imm, false // IdentFactor を含む ImmExp を返します
	case *StringFactor:
		// 文字列ファクター自体は算術的に評価されませんが、
		// ImmExp 内の有効なファクターです。そのまま返します。
		return imm, false
	default:
		// 不明なファクタータイプ
		return imm, false
	}
}
func (imm *ImmExp) TokenLiteral() string {
	return imm.Factor.TokenLiteral()
}

// --- 解析用のヘルパー関数 ---

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
		// TODO: '\n', '\\', '\'' などのエスケープシーケンスを処理する
		return 0, false
	}
	return int64(charStr[0]), true
}

// NumberExp は完全に評価された数値定数式を表します。
type NumberExp struct {
	ImmExp       // Exp インターフェースを満たすために ImmExp を埋め込みます
	Value  int64 // 評価された数値
}

// NewNumberExp は新しい NumberExp を作成します。
func NewNumberExp(base ImmExp, value int64) *NumberExp {
	base.Factor = NewNumberFactor(BaseFactor{}, int(value)) // Factor が NumberFactor であることを確認します
	return &NumberExp{
		ImmExp: base,
		Value:  value,
	}
}

// NumberExp の Eval は自身と true を返し、完全に評価された値であることを示します。
func (n *NumberExp) Eval(env Env) (Exp, bool) {
	return n, true // 評価された値です。
}

// TokenLiteral は数値の文字列表現を返します。
func (n *NumberExp) TokenLiteral() string {
	// 埋め込まれた Factor の TokenLiteral を使用します。これは NumberFactor である必要があります。
	return n.Factor.TokenLiteral()
}

// NumberExp が Exp インターフェースを満たすことを確認します。
var _ Exp = &NumberExp{}

// --- 他の式タイプ (AddExp, MultExp など) の Eval 実装を追加 ---
// UnaryExp が存在する場合、または必要な場合のプレースホルダー

package variantstack

import (
	"fmt"
)

// VariantStack は任意の型を保持できるスタック
type VariantStack struct {
	items []any
}

func NewVariantStack() *VariantStack {
	return &VariantStack{}
}

// Push 任意の型の要素を追加
func (s *VariantStack) Push(item any) {
	s.items = append(s.items, item)
}

// Pop 要素を取り出し、削除（空ならnilを返す）
func (s *VariantStack) Pop() (any, bool) {
	if len(s.items) == 0 {
		return nil, false
	}
	item := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return item, true
}

// PopTyped 要素を取り出して型を判別しながら処理
func (s *VariantStack) PopTyped() {
	item, ok := s.Pop()
	if !ok {
		fmt.Println("Stack is empty")
		return
	}

	// 型アサーションで判定
	switch v := item.(type) {
	case int:
		fmt.Println("Popped an int:", v)
	case string:
		fmt.Println("Popped a string:", v)
	case float64:
		fmt.Println("Popped a float64:", v)
	case bool:
		fmt.Println("Popped a bool:", v)
	case struct{ Name string }:
		fmt.Println("Popped a struct with Name:", v.Name)
	default:
		fmt.Println("Popped an unknown type:", v)
	}
}

// Clear スタックをクリア
func (s *VariantStack) Clear() {
	s.items = nil
}

// Size 要素数を取得
func (s *VariantStack) Size() int {
	return len(s.items)
}

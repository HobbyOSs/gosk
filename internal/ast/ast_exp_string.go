package ast

import (
	"strconv" // strconv をインポート
	"strings"
)

func ExpToString(e Exp) string {
	switch v := e.(type) {
	case *ImmExp:
		return FactorToString(v.Factor)

	case *MemoryAddrExp:
		leftStr := ExpToString(v.Left)
		rightStr := ""
		if v.Right != nil {
			rightStr = ExpToString(v.Right)
		}
		dataTypeStr := ""
		if v.DataType != None {
			dataTypeStr = string(v.DataType) + " "
		}

		if rightStr == "" {
			return dataTypeStr + "[" + leftStr + "]"
		} else {
			return dataTypeStr + "[" + leftStr + ":" + rightStr + "]"
		}

	case *SegmentExp:
		leftStr := ExpToString(v.Left)
		rightStr := ""
		if v.Right != nil {
			rightStr = ExpToString(v.Right)
		}
		dataTypeStr := ""
		if v.DataType != None {
			dataTypeStr = string(v.DataType) + " "
		}
		if rightStr == "" {
			return dataTypeStr + leftStr
		} else {
			return dataTypeStr + leftStr + ":" + rightStr
		}

	case *AddExp:
		// 頭の項
		head := ExpToString(v.HeadExp)
		// 後続の Operators & TailExps をまとめて文字列に
		var buf strings.Builder
		buf.WriteString(head)
		for i, op := range v.Operators {
			buf.WriteByte(' ')
			buf.WriteString(op)
			buf.WriteByte(' ')
			tailStr := ExpToString(v.TailExps[i])
			buf.WriteString(tailStr)
		}
		return buf.String()

	case *MultExp:
		head := ExpToString(v.HeadExp)
		// 例： "4 * ESI" など
		var buf strings.Builder
		buf.WriteString(head)
		for i, op := range v.Operators {
			buf.WriteByte(' ')
			buf.WriteString(op)
			buf.WriteByte(' ')
			tailStr := ExpToString(v.TailExps[i])
			buf.WriteString(tailStr)
		}
		return buf.String()

	case *NumberExp: // NumberExp のケースを追加
		// int64 の値を 10 進数文字列に変換
		return strconv.FormatInt(v.Value, 10)

	// 他のExp型があれば適宜

	default:
		// 不明な型の場合は空文字列を返す (またはエラー処理)
		// panic(fmt.Sprintf("unhandled Exp type: %T", e))
		return ""
	}
}

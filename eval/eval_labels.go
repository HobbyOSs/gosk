package eval

import (
	"fmt"
	"github.com/hangingman/gosk/object"
	"log"
)

type LabelManagement struct {
	labelBinaryRefMap map[string]*object.Binary     // バイナリへの参照を記録し、後で更新する
	labelBytesMap     map[string]int                // ラベルが見つかったバイト数を格納する
	labelFromMap      map[string][]int              // ラベルとラベルが見つかったバイト数をリスト化
	opcode            map[string][]byte             // 格納するべきオペコード
	genBytesFns       map[string]func(i int) []byte // 取得した値をどうやってバイト列に戻すか
}

// AddLabelCallback は後でニーモニックが決まるような命令（JMP命令や一部のMOV命令）を処理する
// @param opcode 最終的な機械語用オペコード
// @param ident 使用されるラベル
// @param bin 機械語の格納先コンテナ
// @param from ラベルのあった位置
func (l *LabelManagement) AddLabelCallback(
	opcode []byte,
	ident string,
	bin *object.Binary,
	from int,
	howToGenBytes func(i int) []byte) {

	log.Println(fmt.Sprintf("info: add label %s from %d !!", ident, from))

	// ラベルとラベルが見つかったバイト数をキーにする
	key := fmt.Sprintf("%s-%d", ident, from)
	l.opcode[key] = opcode
	l.labelBinaryRefMap[key] = bin
	l.labelBytesMap[key] = from
	l.labelFromMap[ident] = append(l.labelFromMap[ident], from)
	l.genBytesFns[key] = howToGenBytes
}

func (l *LabelManagement) RemoveLabelCallback(ident string) {
	log.Println(fmt.Sprintf("info: remove label %s !!", ident))
}

// Emit はAddLabelCallbackを使用後にラベルが見つかったときのコールバック関数
// コールバックとは書いたが、呼び出すのは自分自身
// @param ident 使用されるラベル
// @param from ラベルのあった位置
func (l *LabelManagement) Emit(ident string, to int) {
	froms, fromOk := l.labelFromMap[ident]

	if !fromOk {
		return
	}

	// ラベル１つに対してラベルを使ってる行は複数ありえるのでfor文
	for _, from := range froms {
		log.Println(fmt.Sprintf("info: from=%d, to=%d", from, to))
		log.Println(fmt.Sprintf("info: emit label %s to %d !!", ident, to-from))
		key := fmt.Sprintf("%s-%d", ident, from)
		opcode, opcodeOk := l.opcode[key]
		bin, binOk := l.labelBinaryRefMap[key]

		// ラベルと行数の複合キーで取得できなければ処理は行わない
		if opcodeOk && binOk {
			bin.Value = nil
			binToAppend := l.genBytesFns[key](to - from)
			bin.Value = append(bin.Value, opcode...)
			bin.Value = append(bin.Value, binToAppend...)
		}
	}
}

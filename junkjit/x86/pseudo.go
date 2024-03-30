package x86

import (
	"bytes"
	"encoding/binary"
	"log"

	"github.com/HobbyOSs/gosk/junkjit"
	"github.com/morikuni/failure"
)

func (a *X86Assembler) DB(x uint8, options ...junkjit.DOption) {
	opts := &junkjit.DOptions{Count: 1}
	for _, option := range options {
		option(opts)
	}
	var b byte = x
	for i := 0; i < opts.Count; i++ {
		a.Code.Bytes = append(a.Code.Bytes, b)
	}
}

func (a *X86Assembler) DW(x uint16, options ...junkjit.DOption) {
	opts := &junkjit.DOptions{Count: 1}
	for _, option := range options {
		option(opts)
	}

	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, x)
	if err != nil {
		log.Fatal(failure.Wrap(err))
	}

	b := buf.Bytes()
	for i := 0; i < opts.Count; i++ {
		a.Code.Bytes = append(a.Code.Bytes, b...)
	}
}

func (a *X86Assembler) DD(x uint32, options ...junkjit.DOption) {
	opts := &junkjit.DOptions{Count: 1}
	for _, option := range options {
		option(opts)
	}

	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, x)
	if err != nil {
		log.Fatal(failure.Wrap(err))
	}

	b := buf.Bytes()
	for i := 0; i < opts.Count; i++ {
		a.Code.Bytes = append(a.Code.Bytes, b...)
	}
}

func (a *X86Assembler) DStruct(x any) {

	buf := new(bytes.Buffer)

	switch v := x.(type) {
	case string:
		err := binary.Write(buf, binary.LittleEndian, []byte(v))
		if err != nil {
			log.Fatal(failure.Wrap(err))
		}
	default:
		log.Fatal("type is not handled")
	}

	a.Code.Bytes = append(a.Code.Bytes, buf.Bytes()...)
}

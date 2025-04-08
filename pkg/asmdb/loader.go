package asmdb

import (
	"bytes"
	"compress/gzip"
	"io"

	jsonx86 "github.com/HobbyOSs/json-x86-64-go-mod"
)

func decompressGzip() ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(jsonx86.X86JSONGZ))
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	return io.ReadAll(reader)
}

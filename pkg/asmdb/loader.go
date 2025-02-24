package asmdb

import (
	_ "embed"

	"bytes"
	"compress/gzip"
	"io"
)

//go:embed json-x86-64/x86_64.json.gz
var compressedJSON []byte

func decompressGzip(data []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	return io.ReadAll(reader)
}

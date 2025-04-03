package filefmt

import "github.com/HobbyOSs/gosk/internal/codegen"

// FileFormat は、オブジェクトファイル形式の書き出し処理を表すインターフェースです。
type FileFormat interface {
	// Write は、指定された CodeGenContext の情報を使用して、
	// 指定されたファイルパスにオブジェクトファイルを書き出します。
	Write(ctx *codegen.CodeGenContext, filePath string) error
}

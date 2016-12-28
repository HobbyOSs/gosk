package main

import (
	"flag"
	"github.com/hangingman/go-fs"
	"github.com/hangingman/go-fs/fat"
	"io"
	"os"
	"path/filepath"
)

// https://www.gnu.org/software/mtools/manual/mtools.html#mcopy
// --------------------------------------------------------------
func main() {
	var (
		outputImage string
		copyFrom    string
	)
	flag.StringVar(&outputImage, "i", "", "出力先の指定をする")
	flag.StringVar(&copyFrom, "f", "", "コピー元の指定をする")

	// 読み書き可能で読み取り
	floppyF, err := os.OpenFile(outputImage, os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	defer floppyF.Close()

	// ブロックデバイスを構築する
	device, err := fs.NewFileDisk(floppyF)
	if err != nil {
		panic(err)
	}

	// The actual FAT filesystem
	fatFs, err := fat.New(device)
	if err != nil {
		panic(err)
	}

	// Get the root directory to the filesystem
	rootDir, err := fatFs.RootDir()
	if err != nil {
		panic(err)
	}

	AddSingleFile(rootDir, copyFrom)
}

func AddSingleFile(dir fs.Directory, src string) error {
	inputF, err := os.Create(src)
	if err != nil {
		panic(err)
	}
	entry, err := dir.AddFile(filepath.Base(src))
	if err != nil {
		panic(err)
	}

	fatFile, err := entry.File()
	if err != nil {
		panic(err)
	}

	if _, err := io.Copy(fatFile, inputF); err != nil {
		panic(err)
	}
	inputF.Close()

	if err := os.Remove(src); err != nil {
		panic(err)
	}
	return nil
}

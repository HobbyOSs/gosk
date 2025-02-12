package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/HobbyOSs/gosk/internal/frontend"
	"github.com/HobbyOSs/gosk/internal/gen"
	"github.com/comail/colog"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/transform"
)

const Version = "2.0.0"

func fileIsWritable(fileName string) bool {
	file, err := os.OpenFile(fileName, os.O_WRONLY, 0666)
	defer file.Close()

	if err != nil && !os.IsPermission(err) {
		return true
	}
	return false
}

func readAssets(str string) (string, error) {
	body, err := os.ReadFile(str)
	if err != nil {
		return "", err
	}

	var f []byte
	encodings := []string{"sjis", "utf-8"}
	for _, enc := range encodings {
		if enc != "" {
			ee, _ := charset.Lookup(enc)
			if ee == nil {
				continue
			}
			var buf bytes.Buffer
			ic := transform.NewWriter(&buf, ee.NewDecoder())
			_, err := ic.Write(body)
			if err != nil {
				continue
			}
			err = ic.Close()
			if err != nil {
				continue
			}
			f = buf.Bytes()
			break
		}
	}
	return string(f), nil
}

func main() {
	var (
		version = flag.Bool("v", false, "バージョンとライセンス情報を表示する")
		debug   = flag.Bool("d", false, "デバッグログを出力する")
	)
	// -hオプション用文言
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage:  [--help | -v] source [object/binary] [list]\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if *version {
		fmt.Printf("gosk %s\n", Version)
		fmt.Printf("%s", `Copyright (C) 2024 idiotpanzer@gmail.com
ライセンス GPLv3+: GNU GPL version 3 or later <http://gnu.org/licenses/gpl.html>.
This is free software: you are free to change and redistribute it.
There is NO WARRANTY, to the extent permitted by law.

Thank you osask project !`)
		os.Exit(0)
	}

	if len(flag.Args()) < 2 {
		fmt.Fprintf(os.Stderr, "usage:  [--help | -v] source [object/binary] [list]\n")
		//flag.PrintDefaults()
		os.Exit(16)
	}
	setUpColog(*debug)

	fmt.Printf("source: %s, object: %s\n", flag.Args()[0], flag.Args()[1])
	assemblySrc := flag.Args()[0]
	assemblyDst := flag.Args()[1]

	_, err := os.Stat(assemblySrc)
	if err != nil {
		fmt.Printf("GOSK : can't open %s", assemblySrc)
		os.Exit(17)
	}
	src, err := readAssets(assemblySrc)
	if err != nil {
		fmt.Printf("GOSK : can't read %s", assemblySrc)
		os.Exit(17)
	}

	parseTree, err := gen.Parse("", []byte(src), gen.Entrypoint("Program"), gen.Debug(*debug))
	if err != nil {
		fmt.Printf("GOSK : failed to parse %s\n%+v", assemblySrc, err)
		os.Exit(-1)
	}

	frontend.Exec(parseTree, assemblyDst)

	os.Exit(0)
}

func setUpColog(debug bool) {
	colog.Register()
	colog.SetDefaultLevel(colog.LInfo)
	if debug {
		colog.SetMinLevel(colog.LDebug)
	} else {
		colog.SetMinLevel(colog.LInfo)
	}
	colog.SetFlags(log.Lshortfile)
	colog.SetFormatter(&colog.StdFormatter{Colors: false})
}

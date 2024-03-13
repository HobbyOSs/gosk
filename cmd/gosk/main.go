package main

import (
	"flag"
	"fmt"
	"os"
	"os/user"

	"github.com/hangingman/gosk/frontend"
	"github.com/hangingman/gosk/gen"
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

func main() {
	var (
		version = flag.Bool("v", false, "バージョンとライセンス情報を表示する")
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

	if len(flag.Args()) == 0 {
		// 引数が無ければREPLモードへ移行
		user, err := user.Current()
		if err != nil {
			panic(err)
		}
		fmt.Printf("Hello %s! This is yet another assembly gosk!\n", user.Username)
		fmt.Printf("Feel free to type in commands\n")
		//repl.Start(os.Stdin, os.Stdout)
	}

	if len(flag.Args()) < 2 {
		fmt.Fprintf(os.Stderr, "usage:  [--help | -v] source [object/binary] [list]\n")
		//flag.PrintDefaults()
		os.Exit(16)
	}
	fmt.Printf("source: %s, object: %s\n", flag.Args()[0], flag.Args()[1])

	assemblySrc := flag.Args()[0]
	assemblyDst := flag.Args()[1]

	_, err := os.Stat(assemblySrc)
	if err != nil {
		fmt.Printf("GOSK : can't open %s", assemblySrc)
		os.Exit(17)
	}
	bytes, err := os.ReadFile(assemblySrc)
	if err != nil {
		fmt.Printf("GOSK : can't read %s", assemblySrc)
		os.Exit(17)
	}

	parseTree, err := gen.Parse("", bytes, gen.Entrypoint("Program"))
	if err != nil {
		fmt.Printf("GOSK : failed to parse %s", assemblySrc) // TODO: エラーメッセージ
		os.Exit(-1)
	}

	frontend.Exec(parseTree, assemblyDst)

	os.Exit(0)
}

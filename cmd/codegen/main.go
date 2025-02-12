package main

import (
	"flag"
	"fmt"
	"os"
)

const version = "0.1.0"

func main() {
	helpFlag := flag.Bool("help", false, "ヘルプメッセージを表示")
	versionFlag := flag.Bool("version", false, "バージョン情報を表示")
	flag.Parse()

	if *helpFlag {
		printHelp()
		os.Exit(0)
	}

	if *versionFlag {
		fmt.Printf("codegen バージョン %s\n", version)
		os.Exit(0)
	}

	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("エラー: 入力ファイルが指定されていません")
		printHelp()
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Println("使用方法: codegen [オプション] <入力ファイル>")
	fmt.Println("\nオプション:")
	flag.PrintDefaults()
	fmt.Println("\n使用例:")
	fmt.Println("  codegen program.ocode")
	fmt.Println("  codegen --version")
}

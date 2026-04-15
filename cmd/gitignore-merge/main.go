package main

import (
	"fmt"
	"os"

	"github.com/taku-o/gitignore-merge/internal/merger"
	"github.com/taku-o/gitignore-merge/internal/parser"
)

func main() {
	args := os.Args[1:]

	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: gitignore-merge <file1> <file2> [file3...]")
		os.Exit(1)
	}

	// 各ファイルの存在確認
	for _, path := range args {
		if _, err := os.Stat(path); err != nil {
			if os.IsNotExist(err) {
				fmt.Fprintf(os.Stderr, "Error: file not found: %s\n", path)
			} else {
				fmt.Fprintf(os.Stderr, "Error: cannot access file: %s: %v\n", path, err)
			}
			os.Exit(1)
		}
	}

	// 全ファイルをパース
	files := make([]parser.GitignoreFile, 0, len(args))
	for _, path := range args {
		f, err := parser.ParseFile(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: failed to parse %s: %v\n", path, err)
			os.Exit(1)
		}
		files = append(files, f)
	}

	// マージ処理
	result := merger.Merge(files)

	// 結果を標準出力に出力
	fmt.Print(parser.Format(result))
}

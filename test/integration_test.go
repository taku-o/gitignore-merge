package test

import (
	"os"
	"path/filepath"
	"testing"

	"gitignore-merge/internal/merger"
	"gitignore-merge/internal/parser"
)

func TestIntegration_MergeABC(t *testing.T) {
	testdataDir := filepath.Join("testdata")

	// A, B, C の3ファイルをパース
	fileA, err := parser.ParseFile(filepath.Join(testdataDir, "a.gitignore"))
	if err != nil {
		t.Fatalf("a.gitignore のパースに失敗: %v", err)
	}
	fileB, err := parser.ParseFile(filepath.Join(testdataDir, "b.gitignore"))
	if err != nil {
		t.Fatalf("b.gitignore のパースに失敗: %v", err)
	}
	fileC, err := parser.ParseFile(filepath.Join(testdataDir, "c.gitignore"))
	if err != nil {
		t.Fatalf("c.gitignore のパースに失敗: %v", err)
	}

	// マージ実行
	result := merger.Merge([]parser.GitignoreFile{fileA, fileB, fileC})
	got := parser.Format(result)

	// 期待結果を読み込み
	expectedBytes, err := os.ReadFile(filepath.Join(testdataDir, "expected_abc.gitignore"))
	if err != nil {
		t.Fatalf("expected_abc.gitignore の読み込みに失敗: %v", err)
	}
	want := string(expectedBytes)

	if got != want {
		t.Errorf("マージ結果が期待と異なる:\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}

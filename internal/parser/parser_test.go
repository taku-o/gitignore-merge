package parser

import (
	"os"
	"path/filepath"
	"testing"
)

// writeTestFile はテスト用の一時ファイルを作成するヘルパー
func writeTestFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".gitignore")
	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		t.Fatalf("テストファイルの作成に失敗: %v", err)
	}
	return path
}

func TestParseFile_SectionedFile(t *testing.T) {
	content := `# Node
node_modules/
dist/

# IDE
.vscode/
.idea/
`
	path := writeTestFile(t, content)
	result, err := ParseFile(path)
	if err != nil {
		t.Fatalf("ParseFile がエラーを返した: %v", err)
	}

	if len(result.Sections) != 2 {
		t.Fatalf("セクション数が期待と異なる: got %d, want 2", len(result.Sections))
	}

	// 第1セクション: # Node
	sec0 := result.Sections[0]
	if len(sec0.Header) != 1 || sec0.Header[0] != "# Node" {
		t.Errorf("セクション0のヘッダーが期待と異なる: got %v", sec0.Header)
	}
	if len(sec0.Patterns) != 3 {
		t.Errorf("セクション0のパターン数が期待と異なる: got %d, want 3 (node_modules/, dist/, 空行)", len(sec0.Patterns))
	}

	// 第2セクション: # IDE
	sec1 := result.Sections[1]
	if len(sec1.Header) != 1 || sec1.Header[0] != "# IDE" {
		t.Errorf("セクション1のヘッダーが期待と異なる: got %v", sec1.Header)
	}
	if len(sec1.Patterns) != 2 {
		t.Errorf("セクション1のパターン数が期待と異なる: got %d, want 2", len(sec1.Patterns))
	}
}

func TestParseFile_UnnamedSection(t *testing.T) {
	content := `*.log
tmp/

# Node
node_modules/
`
	path := writeTestFile(t, content)
	result, err := ParseFile(path)
	if err != nil {
		t.Fatalf("ParseFile がエラーを返した: %v", err)
	}

	if len(result.Sections) != 2 {
		t.Fatalf("セクション数が期待と異なる: got %d, want 2", len(result.Sections))
	}

	// 無名セクション（Header が nil）
	sec0 := result.Sections[0]
	if sec0.Header != nil {
		t.Errorf("無名セクションの Header は nil であるべき: got %v", sec0.Header)
	}
	if len(sec0.Patterns) != 3 {
		t.Errorf("無名セクションのパターン数が期待と異なる: got %d, want 3 (*.log, tmp/, 空行)", len(sec0.Patterns))
	}

	// 名前付きセクション
	sec1 := result.Sections[1]
	if len(sec1.Header) != 1 || sec1.Header[0] != "# Node" {
		t.Errorf("セクション1のヘッダーが期待と異なる: got %v", sec1.Header)
	}
}

func TestParseFile_ConsecutiveCommentLines(t *testing.T) {
	content := `# Node
# Dependencies
node_modules/
`
	path := writeTestFile(t, content)
	result, err := ParseFile(path)
	if err != nil {
		t.Fatalf("ParseFile がエラーを返した: %v", err)
	}

	if len(result.Sections) != 1 {
		t.Fatalf("連続コメント行は1つのセクションヘッダーになるべき: got %d sections", len(result.Sections))
	}

	sec0 := result.Sections[0]
	if len(sec0.Header) != 2 {
		t.Errorf("ヘッダーは2行であるべき: got %d", len(sec0.Header))
	}
	if sec0.Header[0] != "# Node" || sec0.Header[1] != "# Dependencies" {
		t.Errorf("ヘッダーの内容が期待と異なる: got %v", sec0.Header)
	}
}

func TestParseFile_EmptyLinePreservation(t *testing.T) {
	content := `# Node
node_modules/

dist/
`
	path := writeTestFile(t, content)
	result, err := ParseFile(path)
	if err != nil {
		t.Fatalf("ParseFile がエラーを返した: %v", err)
	}

	if len(result.Sections) != 1 {
		t.Fatalf("セクション数が期待と異なる: got %d, want 1", len(result.Sections))
	}

	sec0 := result.Sections[0]
	// node_modules/, 空行, dist/ の3パターン
	if len(sec0.Patterns) != 3 {
		t.Errorf("空行を含むパターン数が期待と異なる: got %d, want 3", len(sec0.Patterns))
	}
	if sec0.Patterns[0] != "node_modules/" {
		t.Errorf("パターン0が期待と異なる: got %q", sec0.Patterns[0])
	}
	if sec0.Patterns[1] != "" {
		t.Errorf("パターン1は空行であるべき: got %q", sec0.Patterns[1])
	}
	if sec0.Patterns[2] != "dist/" {
		t.Errorf("パターン2が期待と異なる: got %q", sec0.Patterns[2])
	}
}

func TestParseFile_FileNotFound(t *testing.T) {
	_, err := ParseFile("/nonexistent/path/.gitignore")
	if err == nil {
		t.Error("存在しないファイルに対してエラーが返されるべき")
	}
}

func TestParseFile_EmptyFile(t *testing.T) {
	path := writeTestFile(t, "")
	result, err := ParseFile(path)
	if err != nil {
		t.Fatalf("ParseFile がエラーを返した: %v", err)
	}

	if len(result.Sections) != 0 {
		t.Errorf("空ファイルのセクション数は0であるべき: got %d", len(result.Sections))
	}
}

func TestFormat_SectionedFile(t *testing.T) {
	file := GitignoreFile{
		Sections: []Section{
			{
				Header:   []string{"# Node"},
				Patterns: []string{"node_modules/", "dist/"},
			},
			{
				Header:   []string{"# IDE"},
				Patterns: []string{".vscode/", ".idea/"},
			},
		},
	}

	result := Format(file)
	expected := "# Node\nnode_modules/\ndist/\n# IDE\n.vscode/\n.idea/\n"
	if result != expected {
		t.Errorf("Format の結果が期待と異なる:\ngot:\n%s\nwant:\n%s", result, expected)
	}
}

func TestFormat_UnnamedSection(t *testing.T) {
	file := GitignoreFile{
		Sections: []Section{
			{
				Header:   nil,
				Patterns: []string{"*.log", "tmp/"},
			},
			{
				Header:   []string{"# Node"},
				Patterns: []string{"node_modules/"},
			},
		},
	}

	result := Format(file)
	expected := "*.log\ntmp/\n# Node\nnode_modules/\n"
	if result != expected {
		t.Errorf("Format の結果が期待と異なる:\ngot:\n%s\nwant:\n%s", result, expected)
	}
}

func TestFormat_EmptyFile(t *testing.T) {
	file := GitignoreFile{
		Sections: []Section{},
	}

	result := Format(file)
	if result != "" {
		t.Errorf("空ファイルの Format は空文字列であるべき: got %q", result)
	}
}

func TestRoundTrip(t *testing.T) {
	testCases := []struct {
		name    string
		content string
	}{
		{
			name: "セクション付きファイル",
			content: `# Node
node_modules/
dist/

# IDE
.vscode/
.idea/
`,
		},
		{
			name: "無名セクション付きファイル",
			content: `*.log
tmp/

# Node
node_modules/
`,
		},
		{
			name: "連続コメント行",
			content: `# Node
# Dependencies
node_modules/
`,
		},
		{
			name: "パターンのみ",
			content: `*.log
tmp/
`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			path := writeTestFile(t, tc.content)
			parsed, err := ParseFile(path)
			if err != nil {
				t.Fatalf("ParseFile がエラーを返した: %v", err)
			}

			formatted := Format(parsed)
			if formatted != tc.content {
				t.Errorf("往復一貫性が保たれていない:\ngot:\n%q\nwant:\n%q", formatted, tc.content)
			}
		})
	}
}

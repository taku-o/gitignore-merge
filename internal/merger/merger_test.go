package merger

import (
	"gitignore-merge/internal/parser"
	"reflect"
	"testing"
)

func TestMerge_SameNameSections(t *testing.T) {
	base := parser.GitignoreFile{
		Sections: []parser.Section{
			{
				Header:   []string{"# Node"},
				Patterns: []string{"node_modules/"},
			},
		},
	}
	other := parser.GitignoreFile{
		Sections: []parser.Section{
			{
				Header:   []string{"# Node"},
				Patterns: []string{"dist/"},
			},
		},
	}

	result := Merge([]parser.GitignoreFile{base, other})

	if len(result.Sections) != 1 {
		t.Fatalf("同名セクションは統合されるべき: got %d sections, want 1", len(result.Sections))
	}
	sec := result.Sections[0]
	if !reflect.DeepEqual(sec.Header, []string{"# Node"}) {
		t.Errorf("ヘッダーが期待と異なる: got %v", sec.Header)
	}
	// ベースのパターン + 後続のパターン
	wantPatterns := []string{"node_modules/", "dist/"}
	if !reflect.DeepEqual(sec.Patterns, wantPatterns) {
		t.Errorf("パターンが期待と異なる: got %v, want %v", sec.Patterns, wantPatterns)
	}
}

func TestMerge_NewSectionAppended(t *testing.T) {
	base := parser.GitignoreFile{
		Sections: []parser.Section{
			{
				Header:   []string{"# Node"},
				Patterns: []string{"node_modules/"},
			},
		},
	}
	other := parser.GitignoreFile{
		Sections: []parser.Section{
			{
				Header:   []string{"# IDE"},
				Patterns: []string{".vscode/"},
			},
		},
	}

	result := Merge([]parser.GitignoreFile{base, other})

	if len(result.Sections) != 2 {
		t.Fatalf("新規セクションは末尾に追加されるべき: got %d sections, want 2", len(result.Sections))
	}
	if !reflect.DeepEqual(result.Sections[0].Header, []string{"# Node"}) {
		t.Errorf("セクション0のヘッダーが期待と異なる: got %v", result.Sections[0].Header)
	}
	if !reflect.DeepEqual(result.Sections[1].Header, []string{"# IDE"}) {
		t.Errorf("セクション1のヘッダーが期待と異なる: got %v", result.Sections[1].Header)
	}
	if !reflect.DeepEqual(result.Sections[1].Patterns, []string{".vscode/"}) {
		t.Errorf("セクション1のパターンが期待と異なる: got %v", result.Sections[1].Patterns)
	}
}

func TestMerge_SectionNameComparison(t *testing.T) {
	// # Node と ## Node は同じセクション名「Node」としてマッチする
	base := parser.GitignoreFile{
		Sections: []parser.Section{
			{
				Header:   []string{"# Node"},
				Patterns: []string{"node_modules/"},
			},
		},
	}
	other := parser.GitignoreFile{
		Sections: []parser.Section{
			{
				Header:   []string{"## Node"},
				Patterns: []string{"dist/"},
			},
		},
	}

	result := Merge([]parser.GitignoreFile{base, other})

	if len(result.Sections) != 1 {
		t.Fatalf("# Node と ## Node は同じセクションとしてマージされるべき: got %d sections", len(result.Sections))
	}
	// ベースのヘッダーが維持される
	if !reflect.DeepEqual(result.Sections[0].Header, []string{"# Node"}) {
		t.Errorf("ベースのヘッダーが維持されるべき: got %v", result.Sections[0].Header)
	}
}

func TestMerge_SectionNameWithoutSpace(t *testing.T) {
	// #Node と # Node は同じセクション名「Node」としてマッチする
	base := parser.GitignoreFile{
		Sections: []parser.Section{
			{
				Header:   []string{"#Node"},
				Patterns: []string{"node_modules/"},
			},
		},
	}
	other := parser.GitignoreFile{
		Sections: []parser.Section{
			{
				Header:   []string{"# Node"},
				Patterns: []string{"dist/"},
			},
		},
	}

	result := Merge([]parser.GitignoreFile{base, other})

	if len(result.Sections) != 1 {
		t.Fatalf("#Node と # Node は同じセクションとしてマージされるべき: got %d sections", len(result.Sections))
	}
}

func TestMerge_DuplicatePatternRemoval(t *testing.T) {
	base := parser.GitignoreFile{
		Sections: []parser.Section{
			{
				Header:   []string{"# Node"},
				Patterns: []string{"node_modules/", "dist/"},
			},
		},
	}
	other := parser.GitignoreFile{
		Sections: []parser.Section{
			{
				Header:   []string{"# Node"},
				Patterns: []string{"node_modules/", "build/"},
			},
		},
	}

	result := Merge([]parser.GitignoreFile{base, other})

	wantPatterns := []string{"node_modules/", "dist/", "build/"}
	if !reflect.DeepEqual(result.Sections[0].Patterns, wantPatterns) {
		t.Errorf("重複パターンは除去されるべき: got %v, want %v", result.Sections[0].Patterns, wantPatterns)
	}
}

func TestMerge_ConflictBaseWins(t *testing.T) {
	// ベースに path がある場合、後続の !path は無視される
	base := parser.GitignoreFile{
		Sections: []parser.Section{
			{
				Header:   []string{"# Logs"},
				Patterns: []string{"*.log"},
			},
		},
	}
	other := parser.GitignoreFile{
		Sections: []parser.Section{
			{
				Header:   []string{"# Logs"},
				Patterns: []string{"!*.log"},
			},
		},
	}

	result := Merge([]parser.GitignoreFile{base, other})

	wantPatterns := []string{"*.log"}
	if !reflect.DeepEqual(result.Sections[0].Patterns, wantPatterns) {
		t.Errorf("矛盾パターンはベース側が優先されるべき: got %v, want %v", result.Sections[0].Patterns, wantPatterns)
	}
}

func TestMerge_ConflictNegationBaseWins(t *testing.T) {
	// ベースに !path がある場合、後続の path は無視される
	base := parser.GitignoreFile{
		Sections: []parser.Section{
			{
				Header:   []string{"# Logs"},
				Patterns: []string{"!*.log"},
			},
		},
	}
	other := parser.GitignoreFile{
		Sections: []parser.Section{
			{
				Header:   []string{"# Logs"},
				Patterns: []string{"*.log"},
			},
		},
	}

	result := Merge([]parser.GitignoreFile{base, other})

	wantPatterns := []string{"!*.log"}
	if !reflect.DeepEqual(result.Sections[0].Patterns, wantPatterns) {
		t.Errorf("ベースの否定パターンが優先されるべき: got %v, want %v", result.Sections[0].Patterns, wantPatterns)
	}
}

func TestMerge_ThreeFiles(t *testing.T) {
	file1 := parser.GitignoreFile{
		Sections: []parser.Section{
			{
				Header:   []string{"# Node"},
				Patterns: []string{"node_modules/"},
			},
		},
	}
	file2 := parser.GitignoreFile{
		Sections: []parser.Section{
			{
				Header:   []string{"# Node"},
				Patterns: []string{"dist/"},
			},
			{
				Header:   []string{"# IDE"},
				Patterns: []string{".vscode/"},
			},
		},
	}
	file3 := parser.GitignoreFile{
		Sections: []parser.Section{
			{
				Header:   []string{"# IDE"},
				Patterns: []string{".idea/"},
			},
			{
				Header:   []string{"# OS"},
				Patterns: []string{".DS_Store"},
			},
		},
	}

	result := Merge([]parser.GitignoreFile{file1, file2, file3})

	if len(result.Sections) != 3 {
		t.Fatalf("3ファイルマージ後のセクション数が期待と異なる: got %d, want 3", len(result.Sections))
	}

	// Node: node_modules/, dist/
	nodePatterns := result.Sections[0].Patterns
	if !reflect.DeepEqual(nodePatterns, []string{"node_modules/", "dist/"}) {
		t.Errorf("Node セクションのパターンが期待と異なる: got %v", nodePatterns)
	}

	// IDE: .vscode/, .idea/
	idePatterns := result.Sections[1].Patterns
	if !reflect.DeepEqual(idePatterns, []string{".vscode/", ".idea/"}) {
		t.Errorf("IDE セクションのパターンが期待と異なる: got %v", idePatterns)
	}

	// OS: .DS_Store
	osPatterns := result.Sections[2].Patterns
	if !reflect.DeepEqual(osPatterns, []string{".DS_Store"}) {
		t.Errorf("OS セクションのパターンが期待と異なる: got %v", osPatterns)
	}
}

func TestMerge_SingleFile(t *testing.T) {
	file := parser.GitignoreFile{
		Sections: []parser.Section{
			{
				Header:   []string{"# Node"},
				Patterns: []string{"node_modules/"},
			},
		},
	}

	result := Merge([]parser.GitignoreFile{file})

	if !reflect.DeepEqual(result, file) {
		t.Errorf("単一ファイルの場合はそのまま返されるべき: got %v", result)
	}
}

func TestMerge_UnnamedSections(t *testing.T) {
	base := parser.GitignoreFile{
		Sections: []parser.Section{
			{
				Header:   nil,
				Patterns: []string{"*.log"},
			},
			{
				Header:   []string{"# Node"},
				Patterns: []string{"node_modules/"},
			},
		},
	}
	other := parser.GitignoreFile{
		Sections: []parser.Section{
			{
				Header:   nil,
				Patterns: []string{"tmp/"},
			},
		},
	}

	result := Merge([]parser.GitignoreFile{base, other})

	if len(result.Sections) != 2 {
		t.Fatalf("セクション数が期待と異なる: got %d, want 2", len(result.Sections))
	}
	// 無名セクション同士がマージされる
	if result.Sections[0].Header != nil {
		t.Errorf("無名セクションの Header は nil であるべき: got %v", result.Sections[0].Header)
	}
	wantPatterns := []string{"*.log", "tmp/"}
	if !reflect.DeepEqual(result.Sections[0].Patterns, wantPatterns) {
		t.Errorf("無名セクションのパターンが期待と異なる: got %v, want %v", result.Sections[0].Patterns, wantPatterns)
	}
}

func TestMerge_MultiLineHeader(t *testing.T) {
	// 複数行ヘッダーの場合、最初の行のみでセクション名比較
	base := parser.GitignoreFile{
		Sections: []parser.Section{
			{
				Header:   []string{"# Node", "# Dependencies"},
				Patterns: []string{"node_modules/"},
			},
		},
	}
	other := parser.GitignoreFile{
		Sections: []parser.Section{
			{
				Header:   []string{"# Node"},
				Patterns: []string{"dist/"},
			},
		},
	}

	result := Merge([]parser.GitignoreFile{base, other})

	if len(result.Sections) != 1 {
		t.Fatalf("複数行ヘッダーでも最初の行でマッチすべき: got %d sections", len(result.Sections))
	}
	// ベースのヘッダーが維持される
	if !reflect.DeepEqual(result.Sections[0].Header, []string{"# Node", "# Dependencies"}) {
		t.Errorf("ベースのヘッダーが維持されるべき: got %v", result.Sections[0].Header)
	}
}

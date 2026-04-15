package merger

import (
	"strings"

	"github.com/taku-o/gitignore-merge/internal/parser"
)

// sectionName はセクションのヘッダーからセクション名を抽出する。
// ヘッダーの最初の行から先頭の # 文字列と後続のスペースを除去した文字列を返す。
// Header が nil（無名セクション）の場合は空文字列を返す。
func sectionName(header []string) string {
	if len(header) == 0 {
		return ""
	}
	line := header[0]
	// 先頭の # を全て除去
	trimmed := strings.TrimLeft(line, "#")
	// 後続のスペースを除去
	trimmed = strings.TrimLeft(trimmed, " ")
	return trimmed
}

// Merge は複数の GitignoreFile を先頭ファイルベースでマージする。
// files[0] がベースとなり、files[1:] の内容が順番に統合される。
// 矛盾するパターンが存在する場合、ベースファイル側が優先される。
func Merge(files []parser.GitignoreFile) parser.GitignoreFile {
	if len(files) == 0 {
		return parser.GitignoreFile{}
	}

	result := cloneFile(files[0])

	for i := 1; i < len(files); i++ {
		result = mergeTwo(result, files[i])
	}

	return result
}

// cloneFile は GitignoreFile のディープコピーを作成する
func cloneFile(f parser.GitignoreFile) parser.GitignoreFile {
	sections := make([]parser.Section, len(f.Sections))
	for i, sec := range f.Sections {
		sections[i] = cloneSection(sec)
	}
	return parser.GitignoreFile{Sections: sections}
}

// cloneSection は Section のディープコピーを作成する
func cloneSection(sec parser.Section) parser.Section {
	var header []string
	if sec.Header != nil {
		header = make([]string, len(sec.Header))
		copy(header, sec.Header)
	}
	patterns := make([]string, len(sec.Patterns))
	copy(patterns, sec.Patterns)
	return parser.Section{Header: header, Patterns: patterns}
}

// mergeTwo は2つの GitignoreFile をマージする
func mergeTwo(base, other parser.GitignoreFile) parser.GitignoreFile {
	// ベースのセクション名インデックスを構築
	nameToIndex := make(map[string]int)
	for i, sec := range base.Sections {
		name := sectionName(sec.Header)
		nameToIndex[name] = i
	}

	for _, otherSec := range other.Sections {
		otherName := sectionName(otherSec.Header)
		if idx, found := nameToIndex[otherName]; found {
			// 同名セクションにパターンを統合
			base.Sections[idx] = mergeSectionPatterns(base.Sections[idx], otherSec)
		} else {
			// 新規セクションを末尾に追加
			newSec := cloneSection(otherSec)
			base.Sections = append(base.Sections, newSec)
			nameToIndex[otherName] = len(base.Sections) - 1
		}
	}

	return base
}

// mergeSectionPatterns は後続セクションのパターンをベースセクションに統合する。
// 重複パターンを除去し、矛盾パターンはベース側を優先する。
func mergeSectionPatterns(base, other parser.Section) parser.Section {
	// ベースの非空パターンをセットに格納
	existing := make(map[string]bool)
	for _, p := range base.Patterns {
		if p != "" {
			existing[p] = true
		}
	}

	for _, p := range other.Patterns {
		if p == "" {
			continue
		}

		// 重複チェック
		if existing[p] {
			continue
		}

		// 矛盾チェック: X と !X の関係
		if isConflict(p, existing) {
			continue
		}

		base.Patterns = append(base.Patterns, p)
		existing[p] = true
	}

	return base
}

// isConflict はパターンがベースの既存パターンと矛盾するか判定する。
// パターン X に対して !X が存在するか、!X に対して X が存在する場合に矛盾と見なす。
func isConflict(pattern string, existing map[string]bool) bool {
	if strings.HasPrefix(pattern, "!") {
		// 後続が !X → ベースに X があれば矛盾
		positive := pattern[1:]
		return existing[positive]
	}
	// 後続が X → ベースに !X があれば矛盾
	negated := "!" + pattern
	return existing[negated]
}

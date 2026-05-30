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
// 新規パターンは other 内での前後関係を保持して挿入する。
func mergeSectionPatterns(base, other parser.Section) parser.Section {
	// ベースのパターンをインデックスマップに格納
	existingIndex := make(map[string]int)
	for i, p := range base.Patterns {
		if p != "" {
			existingIndex[p] = i
		}
	}

	for j, p := range other.Patterns {
		if p == "" {
			continue
		}
		if _, exists := existingIndex[p]; exists {
			continue
		}
		if isConflict(p, existingIndex) {
			continue
		}

		// other 内で直前に存在するパターンのベース内位置を探して挿入位置を決定
		insertPos := len(base.Patterns)
		for k := j - 1; k >= 0; k-- {
			prev := other.Patterns[k]
			if prev == "" {
				continue
			}
			if idx, found := existingIndex[prev]; found {
				insertPos = idx + 1
				break
			}
		}

		base.Patterns = append(base.Patterns, "")
		copy(base.Patterns[insertPos+1:], base.Patterns[insertPos:])
		base.Patterns[insertPos] = p

		// 挿入位置以降のインデックスをずらす
		for k, v := range existingIndex {
			if v >= insertPos {
				existingIndex[k] = v + 1
			}
		}
		existingIndex[p] = insertPos
	}

	return base
}

// isConflict はパターンがベースの既存パターンと矛盾するか判定する。
// パターン X に対して !X が存在するか、!X に対して X が存在する場合に矛盾と見なす。
func isConflict(pattern string, existing map[string]int) bool {
	if strings.HasPrefix(pattern, "!") {
		// 後続が !X → ベースに X があれば矛盾
		_, found := existing[pattern[1:]]
		return found
	}
	// 後続が X → ベースに !X があれば矛盾
	_, found := existing["!"+pattern]
	return found
}

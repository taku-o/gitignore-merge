package parser

import (
	"bufio"
	"os"
	"strings"
)

// Section はセクションヘッダーとそれに属するパターンを表す
type Section struct {
	// Header はセクションのコメント行（# で始まる行のスライス）
	// 無名セクション（ファイル先頭）の場合は nil
	Header []string
	// Patterns はセクション内のパターン行（空行を含む）
	Patterns []string
}

// GitignoreFile はパースされた .gitignore ファイル全体を表す
type GitignoreFile struct {
	Sections []Section
}

// ParseFile は指定パスの .gitignore ファイルを読み込み、
// セクション構造付きでパースした結果を返す
func ParseFile(path string) (GitignoreFile, error) {
	f, err := os.Open(path)
	if err != nil {
		return GitignoreFile{}, err
	}
	defer f.Close()

	var sections []Section
	var currentHeader []string
	var currentPatterns []string
	inHeader := false

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "#") {
			if !inHeader {
				// 新しいヘッダーの開始。現在のセクションがあれば保存する
				if currentHeader != nil || currentPatterns != nil {
					sections = append(sections, Section{
						Header:   currentHeader,
						Patterns: currentPatterns,
					})
				}
				currentHeader = []string{line}
				currentPatterns = nil
				inHeader = true
			} else {
				// 連続するコメント行 → 同じヘッダーに追加
				currentHeader = append(currentHeader, line)
			}
		} else {
			inHeader = false
			if currentHeader == nil && len(sections) == 0 && currentPatterns == nil {
				// ファイル先頭の無名セクション開始
				currentPatterns = []string{line}
			} else {
				currentPatterns = append(currentPatterns, line)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return GitignoreFile{}, err
	}

	// 最後のセクションを保存
	if currentHeader != nil || currentPatterns != nil {
		sections = append(sections, Section{
			Header:   currentHeader,
			Patterns: currentPatterns,
		})
	}

	return GitignoreFile{Sections: sections}, nil
}

// Format は GitignoreFile をテキスト形式に変換する
func Format(file GitignoreFile) string {
	if len(file.Sections) == 0 {
		return ""
	}

	var b strings.Builder
	for _, sec := range file.Sections {
		for _, h := range sec.Header {
			b.WriteString(h)
			b.WriteByte('\n')
		}
		for _, p := range sec.Patterns {
			b.WriteString(p)
			b.WriteByte('\n')
		}
	}
	return b.String()
}

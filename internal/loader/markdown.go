package loader

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

type MarkdownLoader struct{}

func (l *MarkdownLoader) Supports(ext string) bool { return ext == ".md" || ext == ".markdown" }

func (l *MarkdownLoader) Load(path string) (*RawDocument, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var sb strings.Builder
	var title string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if title == "" && strings.HasPrefix(line, "# ") {
			title = strings.TrimPrefix(line, "# ")
		}
		sb.WriteString(line)
		sb.WriteString("\n")
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if title == "" {
		title = strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	}
	return &RawDocument{
		Path:    path,
		Title:   title,
		DocType: "md",
		Content: sb.String(),
	}, nil
}

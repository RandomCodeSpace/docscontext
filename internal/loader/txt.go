package loader

import (
	"os"
	"path/filepath"
	"strings"
)

type TXTLoader struct{}

func (l *TXTLoader) Supports(ext string) bool {
	return ext == ".txt" || ext == ".text" || ext == ""
}

func (l *TXTLoader) Load(path string) (*RawDocument, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	title := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	return &RawDocument{
		Path:    path,
		Title:   title,
		DocType: "txt",
		Content: string(data),
	}, nil
}

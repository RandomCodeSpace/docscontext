package loader

import (
	"fmt"
	"path/filepath"
	"strings"
)

// RawDocument is the output of loading a file.
type RawDocument struct {
	Path    string
	Title   string
	DocType string
	Content string // full extracted text
}

// DocumentLoader can load a file and return its text.
type DocumentLoader interface {
	Load(path string) (*RawDocument, error)
	Supports(ext string) bool
}

var registry []DocumentLoader

func init() {
	registry = []DocumentLoader{
		&PDFLoader{},
		&DOCXLoader{},
		&MarkdownLoader{},
		&TXTLoader{},
	}
}

// Load dispatches to the correct loader by file extension.
func Load(path string) (*RawDocument, error) {
	ext := strings.ToLower(filepath.Ext(path))
	for _, l := range registry {
		if l.Supports(ext) {
			return l.Load(path)
		}
	}
	return nil, fmt.Errorf("no loader for extension: %s", ext)
}

// SupportedExtensions returns all supported file extensions.
func SupportedExtensions() []string {
	var exts []string
	for _, l := range registry {
		for _, ext := range []string{".pdf", ".docx", ".md", ".txt"} {
			if l.Supports(ext) {
				exts = append(exts, ext)
			}
		}
	}
	return exts
}

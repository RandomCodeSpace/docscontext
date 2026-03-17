package chunker

import (
	"strings"
	"unicode/utf8"
)

// Chunk is a piece of text from a document.
type Chunk struct {
	Index   int
	Content string
	Tokens  int
}

// Chunker splits text into overlapping chunks.
type Chunker struct {
	ChunkSize    int
	ChunkOverlap int
	separators   []string
}

// New creates a recursive character text splitter.
func New(chunkSize, chunkOverlap int) *Chunker {
	return &Chunker{
		ChunkSize:    chunkSize,
		ChunkOverlap: chunkOverlap,
		separators:   []string{"\n\n", "\n", ". ", " ", ""},
	}
}

// Split splits text into chunks, returning them with indices.
func (c *Chunker) Split(text string) []Chunk {
	pieces := c.splitText(text, c.separators)
	return c.mergePieces(pieces)
}

func (c *Chunker) splitText(text string, separators []string) []string {
	if len(separators) == 0 || utf8.RuneCountInString(text) <= c.ChunkSize {
		return []string{text}
	}
	sep := separators[0]
	rest := separators[1:]

	if sep == "" {
		// Split by rune
		runes := []rune(text)
		var parts []string
		for i := 0; i < len(runes); i += c.ChunkSize {
			end := i + c.ChunkSize
			if end > len(runes) {
				end = len(runes)
			}
			parts = append(parts, string(runes[i:end]))
		}
		return parts
	}

	splits := strings.Split(text, sep)
	var good, bad []string
	for _, s := range splits {
		if s == "" {
			continue
		}
		if utf8.RuneCountInString(s) <= c.ChunkSize {
			good = append(good, s)
		} else {
			// Recursively split the large piece
			sub := c.splitText(s, rest)
			// Flush good before adding sub
			bad = append(bad, good...)
			good = nil
			bad = append(bad, sub...)
		}
	}
	return append(bad, good...)
}

func (c *Chunker) mergePieces(pieces []string) []Chunk {
	var chunks []Chunk
	var current strings.Builder
	idx := 0

	flush := func() {
		text := strings.TrimSpace(current.String())
		if text != "" {
			chunks = append(chunks, Chunk{
				Index:   idx,
				Content: text,
				Tokens:  estimateTokens(text),
			})
			idx++
		}
	}

	for _, piece := range pieces {
		pieceLen := utf8.RuneCountInString(piece)
		curLen := utf8.RuneCountInString(current.String())

		if curLen+pieceLen > c.ChunkSize && curLen > 0 {
			flush()
			// Keep overlap
			curText := current.String()
			curRunes := []rune(curText)
			overlapStart := len(curRunes) - c.ChunkOverlap
			if overlapStart < 0 {
				overlapStart = 0
			}
			current.Reset()
			current.WriteString(string(curRunes[overlapStart:]))
		}
		if current.Len() > 0 {
			current.WriteString(" ")
		}
		current.WriteString(piece)
	}
	flush()
	return chunks
}

// estimateTokens approximates token count (1 token ≈ 4 chars).
func estimateTokens(text string) int {
	return utf8.RuneCountInString(text) / 4
}

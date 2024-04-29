package collector

import (
	"github.com/initialcapacity/ai-starter/internal/ai"
)

type Chunker struct {
	tokenizer *ai.Tokenizer
	limit     int
}

func NewChunker(tokenizer *ai.Tokenizer, limit int) Chunker {
	return Chunker{tokenizer: tokenizer, limit: limit}
}

func (chunker Chunker) Split(text string) []string {
	tokenCount := chunker.tokenizer.CountTokens(text)

	if tokenCount < chunker.limit {
		return []string{text}
	} else {
		firstPart := text[:len(text)/2+200]
		secondPart := text[len(text)/2-200:]

		return append(chunker.Split(firstPart), chunker.Split(secondPart)...)
	}
}

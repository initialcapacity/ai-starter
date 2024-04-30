package ai

type Chunker struct {
	tokenizer *Tokenizer
	limit     int
}

func NewChunker(tokenizer *Tokenizer, limit int) Chunker {
	return Chunker{tokenizer: tokenizer, limit: limit}
}

func (chunker Chunker) Split(text string) []string {
	tokenCount := chunker.tokenizer.CountTokens(text)
	overlap := chunker.limit / 30

	if tokenCount < chunker.limit {
		return []string{text}
	} else {
		firstPart := text[:len(text)/2+overlap]
		secondPart := text[len(text)/2-overlap:]

		return append(chunker.Split(firstPart), chunker.Split(secondPart)...)
	}
}

package ai

import (
	"fmt"
	"github.com/tiktoken-go/tokenizer"
	"log"
)

type Tokenizer struct {
	encoder tokenizer.Codec
}

func NewTokenizer(model tokenizer.Model) *Tokenizer {
	encoder, err := tokenizer.ForModel(model)
	if err != nil {
		log.Fatal(fmt.Errorf("unable to create tokenizer: %w", err))
	}

	return &Tokenizer{encoder: encoder}
}

func (tokenizer Tokenizer) CountTokens(text string) int {
	tokens, _, err := tokenizer.encoder.Encode(text)
	if err != nil {
		log.Fatal(fmt.Errorf("unable to decode tokens: %w", err))
	}

	return len(tokens)
}

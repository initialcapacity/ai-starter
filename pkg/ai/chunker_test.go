package ai_test

import (
	ai2 "github.com/initialcapacity/ai-starter/pkg/ai"
	"github.com/stretchr/testify/assert"
	"github.com/tiktoken-go/tokenizer"
	"testing"
)

func TestChunker_Split(t *testing.T) {
	token := ai2.NewTokenizer(tokenizer.Cl100kBase)
	chunker := ai2.NewChunker(token, 30)

	result := chunker.Split("I think that this string should have 31 tokens, but it's hard to say for sure. We'll have to count them manually, I guess.")

	assert.Equal(t, []string{
		"I think that this string should have 31 tokens, but it's hard ",
		"d to say for sure. We'll have to count them manually, I guess.",
	}, result)
}

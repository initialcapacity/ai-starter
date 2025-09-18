package ai_test

import (
	"github.com/initialcapacity/ai-starter/pkg/ai"
	"github.com/stretchr/testify/assert"
	"github.com/tiktoken-go/tokenizer"
	"testing"
)

func TestChunker_Split(t *testing.T) {
	token := ai.NewTokenizer(tokenizer.GPT5Mini)
	chunker := ai.NewChunker(token, 29)

	result := chunker.Split("I think that this string should have 30 tokens, but it's hard to say for sure. We'll have to count them manually, I guess.")

	assert.Equal(t, []string{
		"I think that this string should have 30 tokens, but it's hard",
		" to say for sure. We'll have to count them manually, I guess.",
	}, result)
}

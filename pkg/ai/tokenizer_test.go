package ai_test

import (
	"github.com/initialcapacity/ai-starter/pkg/ai"
	"github.com/stretchr/testify/assert"
	tokenizer "github.com/tiktoken-go/tokenizer"
	"testing"
)

func TestTokenizer_CountTokens(t *testing.T) {
	token := ai.NewTokenizer(tokenizer.GPT5Mini)

	count := token.CountTokens("This string should have 7 tokens")

	assert.Equal(t, 7, count)
}

package ai_test

import (
	"github.com/initialcapacity/ai-starter/internal/ai"
	"github.com/stretchr/testify/assert"
	tokenizer2 "github.com/tiktoken-go/tokenizer"
	"testing"
)

func TestTokenizer_CountTokens(t *testing.T) {
	token := ai.NewTokenizer(tokenizer2.Cl100kBase)

	count := token.CountTokens("This string should have 7 tokens")

	assert.Equal(t, 7, count)
}

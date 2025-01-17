package testsupport

import (
	"context"
	"github.com/initialcapacity/ai-starter/pkg/ai"
)

type FakeAi struct {
	EmbeddingError  error
	CompletionError error
}

func (f FakeAi) Options() ai.LLMOptions {
	return ai.LLMOptions{
		ChatModel:       "gpt-123",
		EmbeddingsModel: "embeddings-test-medium",
		Temperature:     2,
	}
}

func (f FakeAi) CreateEmbedding(_ context.Context, _ string) ([]float32, error) {
	if f.EmbeddingError != nil {
		return nil, f.EmbeddingError
	}

	return CreateVector(0), nil
}

func (f FakeAi) GetChatCompletion(_ context.Context, _ []ai.ChatMessage) (chan string, error) {
	if f.CompletionError != nil {
		return nil, f.CompletionError
	}

	response := make(chan string)
	go func() {
		response <- "Sounds"
		response <- " good"
		close(response)
	}()
	return response, nil
}

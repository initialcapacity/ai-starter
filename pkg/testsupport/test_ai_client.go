package testsupport

import (
	"github.com/initialcapacity/ai-starter/pkg/ai"
)

func NewTestAiClient(openAiEndpoint string) ai.Client {
	return ai.NewClient("a-test-key", openAiEndpoint, ai.LLMOptions{
		ChatModel:       "gpt-test-1",
		EmbeddingsModel: "embeddings-test-medium",
		Temperature:     1,
	})
}

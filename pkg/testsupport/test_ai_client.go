package testsupport

import (
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/initialcapacity/ai-starter/pkg/ai"
	"log"
)

func NewTestAiClient(openAiEndpoint string) ai.Client {
	openAiClient, err := azopenai.NewClientForOpenAI(openAiEndpoint, nil, nil)
	if err != nil {
		log.Fatal(fmt.Errorf("unable to create Open AI client: %w", err))
	}

	return ai.Client{OpenAiClient: openAiClient, LLMOptions: ai.LLMOptions{
		ChatModel:       "gpt-test-1",
		EmbeddingsModel: "embeddings-test-medium",
		Temperature:     1,
	}}
}

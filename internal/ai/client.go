package ai

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"log"
)

type Client struct {
	openAiClient *azopenai.Client
}

func NewClient(openAiKey string) Client {
	keyCredential := azcore.NewKeyCredential(openAiKey)
	openAiClient, err := azopenai.NewClientForOpenAI("https://api.openai.com/v1", keyCredential, nil)
	if err != nil {
		log.Fatal(fmt.Errorf("unable to create Open AI client: %w", err))
	}

	return Client{openAiClient: openAiClient}
}

func (client Client) CreateEmbedding(ctx context.Context, text string) ([]float32, error) {
	model := "text-embedding-3-large"
	embeddings, err := client.openAiClient.GetEmbeddings(ctx, azopenai.EmbeddingsOptions{
		Input:          []string{text},
		DeploymentName: &model,
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to get embeddings: %w", err)
	}

	return embeddings.Data[0].Embedding, nil
}

func (client Client) GetChatCompletion(ctx context.Context, messages []ChatMessage) (string, error) {
	model := "gpt-4-turbo"
	chatResponse, err := client.openAiClient.GetChatCompletions(ctx, azopenai.ChatCompletionsOptions{
		Messages:       toOpenAiMessages(messages),
		DeploymentName: &model,
	}, nil)
	if err != nil {
		return "", fmt.Errorf("unable to get completions: %w", err)
	}

	return *chatResponse.ChatCompletions.Choices[0].Message.Content, nil
}

type Role string

const (
	User      Role = "user"
	System    Role = "system"
	Assistant Role = "assistant"
)

type ChatMessage struct {
	Role    Role
	Content string
}

func toOpenAiMessages(messages []ChatMessage) []azopenai.ChatRequestMessageClassification {
	var result []azopenai.ChatRequestMessageClassification

	for _, message := range messages {
		if message.Role == User {
			result = append(result, &azopenai.ChatRequestUserMessage{Content: azopenai.NewChatRequestUserMessageContent(message.Content)})
		} else if message.Role == System {
			result = append(result, &azopenai.ChatRequestSystemMessage{Content: &message.Content})
		}
	}

	return result
}

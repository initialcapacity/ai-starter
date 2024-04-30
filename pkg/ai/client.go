package ai

import (
	"context"
	"errors"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"io"
	"log"
	"log/slog"
)

type Client struct {
	OpenAiClient *azopenai.Client
}

func NewClient(openAiKey, openAiEndpoint string) Client {
	keyCredential := azcore.NewKeyCredential(openAiKey)
	openAiClient, err := azopenai.NewClientForOpenAI(openAiEndpoint, keyCredential, nil)
	if err != nil {
		log.Fatal(fmt.Errorf("unable to create Open AI client: %w", err))
	}

	return Client{OpenAiClient: openAiClient}
}

func (client Client) CreateEmbedding(ctx context.Context, text string) ([]float32, error) {
	model := "text-embedding-3-large"
	embeddings, err := client.OpenAiClient.GetEmbeddings(ctx, azopenai.EmbeddingsOptions{
		Input:          []string{text},
		DeploymentName: &model,
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to get embeddings: %w", err)
	}

	return embeddings.Data[0].Embedding, nil
}

func (client Client) GetChatCompletion(ctx context.Context, messages []ChatMessage) (chan string, error) {
	model := "gpt-4-turbo"
	chatResponse, streamError := client.OpenAiClient.GetChatCompletionsStream(ctx, azopenai.ChatCompletionsOptions{
		Messages:       toOpenAiMessages(messages),
		DeploymentName: &model,
	}, nil)
	if streamError != nil {
		return nil, fmt.Errorf("unable to get completions: %w", streamError)
	}

	response := make(chan string)

	go func() {
		for {
			chatCompletions, err := chatResponse.ChatCompletionsStream.Read()

			if errors.Is(err, io.EOF) {
				close(response)
				break
			}

			if err != nil {
				log.Fatalf("Error streaming response: %s", err)
			}

			choice := chatCompletions.Choices[0]
			content := choice.Delta.Content
			if content != nil {
				slog.Info("Got content: ", "content", *content)
				response <- *content
			}
		}
	}()

	return response, nil
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

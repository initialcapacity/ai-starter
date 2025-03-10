package ai

import (
	"context"
	"errors"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"io"
	"log"
	"log/slog"
	"time"
)

type LLMOptions struct {
	ChatModel       string
	EmbeddingsModel string
	Temperature     float32
}

type Client struct {
	OpenAiClient *azopenai.Client
	LLMOptions   LLMOptions
}

func NewClient(openAiKey, openAiEndpoint string, options LLMOptions) Client {
	keyCredential := azcore.NewKeyCredential(openAiKey)
	openAiClient, err := azopenai.NewClientForOpenAI(openAiEndpoint, keyCredential, nil)
	if err != nil {
		log.Fatal(fmt.Errorf("unable to create Open AI client: %w", err))
	}

	return Client{OpenAiClient: openAiClient, LLMOptions: options}
}

func (client Client) Options() LLMOptions {
	return client.LLMOptions
}

func (client Client) CreateEmbedding(ctx context.Context, text string) ([]float32, error) {
	embeddings, err := client.OpenAiClient.GetEmbeddings(ctx, azopenai.EmbeddingsOptions{
		Input:          []string{text},
		DeploymentName: &client.LLMOptions.EmbeddingsModel,
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to get embeddings: %w", err)
	}

	return embeddings.Data[0].Embedding, nil
}

func (client Client) GetChatCompletion(ctx context.Context, messages []ChatMessage) (chan string, error) {
	chatResponse, streamError := client.OpenAiClient.GetChatCompletionsStream(ctx, azopenai.ChatCompletionsStreamOptions{
		Messages:       toOpenAiMessages(messages),
		DeploymentName: &client.LLMOptions.ChatModel,
		Temperature:    &client.LLMOptions.Temperature,
	}, nil)
	if streamError != nil {
		return nil, fmt.Errorf("unable to get completions: %w", streamError)
	}

	response := make(chan string, 10)
	streamingResponseCtx, cancel := context.WithTimeout(ctx, 30*time.Second)

	go func() {
		defer close(response)
		defer cancel()

		for {
			chatCompletions, err := chatResponse.ChatCompletionsStream.Read()
			if err != nil {
				if !errors.Is(err, io.EOF) {
					response <- " ...An error occurred. Please try your query again."
					slog.Error("error streaming response", "error", err)
				}
				return
			}

			content := chatCompletions.Choices[0].Delta.Content
			if content != nil {
				select {
				case response <- *content:
				case <-streamingResponseCtx.Done():
					if errors.Is(streamingResponseCtx.Err(), context.DeadlineExceeded) {
						response <- " ...Response timed out."
					}

					slog.Debug("context canceled while sending response")
					return
				}
			}
		}
	}()

	return response, nil
}

func (client Client) GetJsonChatCompletion(ctx context.Context, messages []ChatMessage, schemaName string, schemaDescription string, jsonSchema string) (string, error) {
	chatResponse, err := client.OpenAiClient.GetChatCompletions(ctx, azopenai.ChatCompletionsOptions{
		Messages:       toOpenAiMessages(messages),
		DeploymentName: &client.LLMOptions.ChatModel,
		ResponseFormat: &azopenai.ChatCompletionsJSONSchemaResponseFormat{
			JSONSchema: &azopenai.ChatCompletionsJSONSchemaResponseFormatJSONSchema{
				Name:        &schemaName,
				Description: &schemaDescription,
				Schema:      []byte(jsonSchema),
				Strict:      to.Ptr(true),
			},
		},
		Temperature: &client.LLMOptions.Temperature,
	}, nil)
	if err != nil {
		return "", fmt.Errorf("unable to get JSON completions: %w", err)
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
			result = append(result, &azopenai.ChatRequestSystemMessage{Content: azopenai.NewChatRequestSystemMessageContent(message.Content)})
		}
	}

	return result
}

package ai

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/invopop/jsonschema"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/packages/param"
	"github.com/openai/openai-go/shared"
	"log/slog"
	"reflect"
	"time"
)

type LLMOptions struct {
	ChatModel       string
	EmbeddingsModel string
	Temperature     float64
}

type Client struct {
	openAiClient openai.Client
	llmOptions   LLMOptions
}

func NewClient(openAiKey, openAiEndpoint string, options LLMOptions) Client {
	openAiClient := openai.NewClient(
		option.WithAPIKey(openAiKey),
		option.WithBaseURL(openAiEndpoint),
	)
	return Client{openAiClient: openAiClient, llmOptions: options}
}

func (client Client) Options() LLMOptions {
	return client.llmOptions
}

func (client Client) CreateEmbedding(ctx context.Context, text string) ([]float64, error) {
	response, err := client.openAiClient.Embeddings.New(ctx, openai.EmbeddingNewParams{
		Input: openai.EmbeddingNewParamsInputUnion{OfString: param.NewOpt(text)},
		Model: client.llmOptions.EmbeddingsModel,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to get embeddings: %w", err)
	}

	return response.Data[0].Embedding, nil
}

func (client Client) GetChatCompletion(ctx context.Context, messages []ChatMessage) (chan string, error) {
	stream := client.openAiClient.Chat.Completions.NewStreaming(ctx, openai.ChatCompletionNewParams{
		Messages:    toOpenAiMessages(messages),
		Model:       client.llmOptions.ChatModel,
		Temperature: param.NewOpt(client.llmOptions.Temperature),
	})

	response := make(chan string, 10)
	streamingResponseCtx, cancel := context.WithTimeout(ctx, 30*time.Second)

	go func() {
		defer func() { _ = stream.Close() }()
		defer close(response)
		defer cancel()

		for stream.Next() {
			chunk := stream.Current()
			content := chunk.Choices[0].Delta.Content

			select {
			case response <- content:
			case <-streamingResponseCtx.Done():
				if errors.Is(streamingResponseCtx.Err(), context.DeadlineExceeded) {
					response <- " ...Response timed out."
				}

				slog.Debug("context canceled while sending response")
				return
			}
		}

		if err := stream.Err(); err != nil {
			response <- " ...An error occurred. Please try your query again."
			slog.Error("error streaming response", "error", err)
		}
	}()

	return response, nil
}

func (client Client) GetJsonChatCompletion(ctx context.Context, messages []ChatMessage, schemaName string, schemaDescription string, jsonSchema interface{}) (string, error) {
	response, err := client.openAiClient.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: toOpenAiMessages(messages),
		Model:    client.llmOptions.ChatModel,
		ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{OfJSONSchema: &shared.ResponseFormatJSONSchemaParam{
			JSONSchema: openai.ResponseFormatJSONSchemaJSONSchemaParam{
				Name:        schemaName,
				Description: param.NewOpt(schemaDescription),
				Schema:      jsonSchema,
			},
			Type: "json_schema",
		}},
		Temperature: param.NewOpt(client.llmOptions.Temperature),
	})
	if err != nil {
		return "", fmt.Errorf("unable to get JSON completions: %w", err)
	}

	return response.Choices[0].Message.Content, nil
}

type jsonCompletionClient interface {
	GetJsonChatCompletion(ctx context.Context, messages []ChatMessage, schemaName string, schemaDescription string, jsonSchema interface{}) (string, error)
}

func GetJsonChatCompletion[T any](ctx context.Context, client jsonCompletionClient, messages []ChatMessage, schemaName string, schemaDescription string) (T, error) {
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	var result T
	schema := reflector.ReflectFromType(reflect.TypeOf(result))

	stringResult, err := client.GetJsonChatCompletion(ctx, messages, schemaName, schemaDescription, schema)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal([]byte(stringResult), &result)
	return result, err
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

func toOpenAiMessages(messages []ChatMessage) []openai.ChatCompletionMessageParamUnion {
	var result []openai.ChatCompletionMessageParamUnion

	for _, message := range messages {
		if message.Role == User {
			result = append(result, openai.UserMessage(message.Content))
		} else if message.Role == Assistant {
			result = append(result, openai.AssistantMessage(message.Content))
		} else if message.Role == System {
			result = append(result, openai.SystemMessage(message.Content))
		}
	}

	return result
}

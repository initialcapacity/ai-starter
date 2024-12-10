package query

import (
	"context"
	"fmt"
	"github.com/initialcapacity/ai-starter/internal/analysis"
	"github.com/initialcapacity/ai-starter/pkg/ai"
	"log/slog"
	"strings"
)

type Service struct {
	embeddingsGateway *analysis.EmbeddingsGateway
	aiClient          aiClient
	responsesGateway  *ResponsesGateway
}

func NewService(embeddingsGateway *analysis.EmbeddingsGateway, aiClient aiClient, responsesGateway *ResponsesGateway) *Service {
	return &Service{embeddingsGateway: embeddingsGateway, aiClient: aiClient, responsesGateway: responsesGateway}
}

func (q *Service) FetchResponse(ctx context.Context, query string) (Result, error) {
	embedding, err := q.aiClient.CreateEmbedding(ctx, query)
	if err != nil {
		slog.Error("unable to create embedding", slog.Any("error", err))
		return Result{}, err
	}

	record, err := q.embeddingsGateway.FindSimilar(embedding)
	if err != nil {
		slog.Error("unable to find similar embedding", slog.Any("error", err))
		return Result{}, err
	}

	systemPrompt := fmt.Sprintf(`You are a reporter for a major world newspaper.
Write your response as if you were writing a short, high-quality news article for your paper. Limit your response to one paragraph.

Use the following article for context: %s`, record.Content)

	response, err := q.aiClient.GetChatCompletion(ctx, []ai.ChatMessage{
		{Role: ai.System, Content: systemPrompt},
		{Role: ai.User, Content: query},
	})
	if err != nil {
		slog.Error("unable fetch chat completion", slog.Any("error", err))
		return Result{}, err
	}

	monitoredResponse := make(chan string)
	go func() {
		defer close(monitoredResponse)
		var builder strings.Builder
		for part := range response {
			builder.WriteString(part)
			monitoredResponse <- part
		}
		llmOptions := q.aiClient.Options()
		_, storeResponseErr := q.responsesGateway.Create(systemPrompt, query, record.Source, builder.String(), llmOptions.ChatModel, llmOptions.EmbeddingsModel, llmOptions.Temperature)
		if storeResponseErr != nil {
			slog.Error("unable to store query response", slog.Any("error", storeResponseErr))
		}
	}()

	return Result{
		Source:   record.Source,
		Response: monitoredResponse,
	}, nil
}

type Result struct {
	Response chan string
	Source   string
}

type aiClient interface {
	Options() ai.LLMOptions
	GetChatCompletion(ctx context.Context, messages []ai.ChatMessage) (chan string, error)
	CreateEmbedding(ctx context.Context, text string) ([]float32, error)
}

package app

import (
	"context"
	"fmt"
	"github.com/initialcapacity/ai-starter/internal/analyzer"
	"github.com/initialcapacity/ai-starter/pkg/ai"
	"log/slog"
)

type QueryService struct {
	embeddingsGateway *analyzer.EmbeddingsGateway
	aiClient          aiClient
}

func NewQueryService(embeddingsGateway *analyzer.EmbeddingsGateway, aiClient aiClient) *QueryService {
	return &QueryService{embeddingsGateway: embeddingsGateway, aiClient: aiClient}
}

func (q *QueryService) FetchResponse(ctx context.Context, query string) (QueryResult, error) {
	embedding, err := q.aiClient.CreateEmbedding(ctx, query)
	if err != nil {
		slog.Error("unable to create embedding", err)
		return QueryResult{}, err
	}

	record, err := q.embeddingsGateway.FindSimilar(embedding)
	if err != nil {
		slog.Error("unable to find similar embedding", err)
		return QueryResult{}, err
	}

	response, err := q.aiClient.GetChatCompletion(ctx, []ai.ChatMessage{
		{Role: ai.System, Content: "You are a reporter for a major world newspaper."},
		{Role: ai.System, Content: "Write your response as if you were writing a short, high-quality news article for your paper. Limit your response to one paragraph."},
		{Role: ai.System, Content: fmt.Sprintf("Use the following article for context: %s", record.Content)},
		{Role: ai.User, Content: query},
	})
	if err != nil {
		slog.Error("unable fetch chat completion", err)
		return QueryResult{}, err
	}

	return QueryResult{
		Source:   record.Source,
		Response: response,
	}, nil
}

type QueryResult struct {
	Response chan string
	Source   string
}

type aiClient interface {
	GetChatCompletion(ctx context.Context, messages []ai.ChatMessage) (chan string, error)
	CreateEmbedding(ctx context.Context, text string) ([]float32, error)
}

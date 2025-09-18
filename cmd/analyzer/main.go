package main

import (
	"context"
	"github.com/initialcapacity/ai-starter/internal/analysis"
	"github.com/initialcapacity/ai-starter/internal/collection"
	"github.com/initialcapacity/ai-starter/pkg/ai"
	"github.com/initialcapacity/ai-starter/pkg/dbsupport"
	"github.com/initialcapacity/ai-starter/pkg/websupport"
	"log/slog"
)

func main() {
	databaseUrl := websupport.RequireEnvironmentVariable[string]("DATABASE_URL")
	openAiKey := websupport.RequireEnvironmentVariable[string]("OPEN_AI_KEY")

	db := dbsupport.CreateConnection(databaseUrl)
	chunksGateway := collection.NewChunksGateway(db)
	embeddingsGateway := analysis.NewEmbeddingsGateway(db)
	options := ai.LLMOptions{ChatModel: "gpt-5-mini", EmbeddingsModel: "text-embedding-3-large", Temperature: 1}
	aiClient := ai.NewClient(openAiKey, "https://api.openai.com/v1", options)
	runsGateway := analysis.NewRunsGateway(db)

	a := analysis.NewAnalyzer(chunksGateway, embeddingsGateway, aiClient, runsGateway)

	err := a.Analyze(context.Background())

	if err == nil {
		slog.Info("successful analysis")
	} else {
		slog.Error("unsuccessful analysis: %w", slog.Any("error", err))
	}
}

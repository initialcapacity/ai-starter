package main

import (
	"context"
	"github.com/initialcapacity/ai-starter/internal/ai"
	"github.com/initialcapacity/ai-starter/internal/analyzer"
	"github.com/initialcapacity/ai-starter/internal/collector"
	"github.com/initialcapacity/ai-starter/pkg/dbsupport"
	"github.com/initialcapacity/ai-starter/pkg/websupport"
	"log/slog"
)

func main() {
	databaseUrl := websupport.RequireEnvironmentVariable[string]("DATABASE_URL")
	openAiKey := websupport.RequireEnvironmentVariable[string]("OPEN_AI_KEY")

	db := dbsupport.CreateConnection(databaseUrl)
	chunksGateway := collector.NewChunksGateway(db)
	embeddingsGateway := analyzer.NewEmbeddingsGateway(db)
	aiClient := ai.NewClient(openAiKey)

	a := analyzer.NewAnalyzer(chunksGateway, embeddingsGateway, aiClient)

	err := a.Analyze(context.Background())

	if err == nil {
		slog.Info("successful analysis")
	} else {
		slog.Error("unsuccessful analysis: %w", err)
	}
}

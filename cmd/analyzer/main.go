package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/initialcapacity/ai-starter/internal/analyzer"
	"github.com/initialcapacity/ai-starter/internal/collector"
	"github.com/initialcapacity/ai-starter/pkg/websupport"
	_ "github.com/lib/pq"
	"log"
	"log/slog"
)

func main() {
	databaseUrl := websupport.RequireEnvironmentVariable[string]("DATABASE_URL")
	openAiKey := websupport.RequireEnvironmentVariable[string]("OPEN_AI_KEY")

	db, err := sql.Open("postgres", databaseUrl)
	if err != nil {
		log.Fatal(fmt.Errorf("unable to connect to database: %w", err))
	}

	keyCredential := azcore.NewKeyCredential(openAiKey)
	client, err := azopenai.NewClientForOpenAI("https://api.openai.com/v1", keyCredential, nil)
	if err != nil {
		log.Fatal(fmt.Errorf("unable to create client: %w", err))
	}

	dataGateway := collector.NewDataGateway(db)
	embeddingsGateway := analyzer.NewEmbeddingsGateway(db)

	a := analyzer.NewAnalyzer(dataGateway, embeddingsGateway, client)

	err = a.Analyze(context.Background())

	if err == nil {
		slog.Info("successful analysis")
	} else {
		slog.Error("unsuccessful analysis: %w", err)
	}
}

package functions

import (
	"context"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/initialcapacity/ai-starter/internal/analysis"
	"github.com/initialcapacity/ai-starter/internal/collection"
	"github.com/initialcapacity/ai-starter/pkg/ai"
	"github.com/initialcapacity/ai-starter/pkg/dbsupport"
	"github.com/initialcapacity/ai-starter/pkg/feedsupport"
	"github.com/initialcapacity/ai-starter/pkg/websupport"
	"github.com/tiktoken-go/tokenizer"
	"net/http"
	"strings"
)

func init() {
	functions.CloudEvent("analyzer", triggerAnalyze)
	functions.CloudEvent("collector", triggerCollect)
}

func triggerCollect(ctx context.Context, e event.Event) error {
	databaseUrl := websupport.RequireEnvironmentVariable[string]("DATABASE_URL")
	feeds := websupport.RequireEnvironmentVariable[string]("FEEDS")
	feedUrls := strings.Split(feeds, ",")

	client := http.Client{}
	db := dbsupport.CreateConnection(databaseUrl)

	parser := feedsupport.NewParser(client)
	extractor := feedsupport.NewExtractor(client)
	dataGateway := collection.NewDataGateway(db)
	t := ai.NewTokenizer(tokenizer.Cl100kBase)
	chunksGateway := collection.NewChunksGateway(db)
	chunker := ai.NewChunker(t, 6000)
	chunksService := collection.NewChunksService(chunker, chunksGateway)
	runsGateway := collection.NewCollectionRunsGateway(db)

	c := collection.New(parser, extractor, dataGateway, chunksService, runsGateway)

	return c.Collect(feedUrls)
}

func triggerAnalyze(ctx context.Context, e event.Event) error {
	databaseUrl := websupport.RequireEnvironmentVariable[string]("DATABASE_URL")
	openAiKey := websupport.RequireEnvironmentVariable[string]("OPEN_AI_KEY")
	openAiEndpoint := websupport.EnvironmentVariable("OPEN_AI_ENDPOINT", "https://api.openai.com/v1")

	db := dbsupport.CreateConnection(databaseUrl)
	chunksGateway := collection.NewChunksGateway(db)
	embeddingsGateway := analysis.NewEmbeddingsGateway(db)
	aiClient := ai.NewClient(openAiKey, openAiEndpoint)
	runsGateway := analysis.NewAnalysisRunsGateway(db)

	a := analysis.NewAnalyzer(chunksGateway, embeddingsGateway, aiClient, runsGateway)

	return a.Analyze(ctx)
}

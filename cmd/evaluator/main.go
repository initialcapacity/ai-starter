package main

import (
	"fmt"
	"github.com/initialcapacity/ai-starter/internal/analyzer"
	"github.com/initialcapacity/ai-starter/internal/evaluation"
	"github.com/initialcapacity/ai-starter/internal/query"
	"github.com/initialcapacity/ai-starter/pkg/ai"
	"github.com/initialcapacity/ai-starter/pkg/dbsupport"
	"github.com/initialcapacity/ai-starter/pkg/websupport"
)

func main() {
	queries := []string{
		"What's new with Kotlin?",
		"Tell me about the latest Python Flask news",
		"What's the latest version of Kotlin?",
		"Are there any breaking changes in the newest Java version?",
		"What are new Rust features?",
		"Who's the head of state of Singapore?",
		"What's your favorite color?",
		"How much does a penguin weigh?",
		"Tell an off-color joke",
	}

	openAiEndpoint := websupport.EnvironmentVariable("OPEN_AI_ENDPOINT", "https://api.openai.com/v1")
	openAiKey := websupport.RequireEnvironmentVariable[string]("OPEN_AI_KEY")
	databaseUrl := websupport.RequireEnvironmentVariable[string]("DATABASE_URL")

	aiClient := ai.NewClient(openAiKey, openAiEndpoint)
	db := dbsupport.CreateConnection(databaseUrl)
	embeddingsGateway := analyzer.NewEmbeddingsGateway(db)
	queryService := query.NewService(embeddingsGateway, aiClient)
	aiScorer := evaluation.NewAiScorer(aiClient)

	retriever := evaluation.NewChatResponseRetriever(queryService)
	scoreRunner := evaluation.NewScoreRunner(aiScorer)
	reporter := evaluation.NewScoreReporter()

	results := retriever.Retrieve(queries)
	scores := scoreRunner.Score(results)
	lines := reporter.Report(scores)

	for line := range lines {
		fmt.Printf("%v\n", line)
	}
}

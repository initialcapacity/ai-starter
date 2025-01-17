package main

import (
	"github.com/initialcapacity/ai-starter/internal/analysis"
	"github.com/initialcapacity/ai-starter/internal/evaluation"
	"github.com/initialcapacity/ai-starter/internal/query"
	"github.com/initialcapacity/ai-starter/internal/scores"
	"github.com/initialcapacity/ai-starter/pkg/ai"
	"github.com/initialcapacity/ai-starter/pkg/dbsupport"
	"github.com/initialcapacity/ai-starter/pkg/websupport"
	"log"
)

func main() {
	queries := []string{
		"What's new with Kotlin?", "Tell me about the latest Python Flask news", "What's the latest version of Kotlin?",
		"Are there any breaking changes in the newest Java version?", "What are new Rust features?",
		"Who's the head of state of Singapore?", "What's your favorite color?", "How much does a penguin weigh?",
		"Tell an off-color joke",
	}

	openAiKey := websupport.RequireEnvironmentVariable[string]("OPEN_AI_KEY")
	databaseUrl := websupport.RequireEnvironmentVariable[string]("DATABASE_URL")

	db := dbsupport.CreateConnection(databaseUrl)
	embeddingsGateway := analysis.NewEmbeddingsGateway(db)
	responsesGateway := query.NewResponsesGateway(db)
	options := ai.LLMOptions{ChatModel: "gpt-4o", EmbeddingsModel: "text-embedding-3-large", Temperature: 1}
	aiClient := ai.NewClient(openAiKey, "https://api.openai.com/v1", options)
	queryService := query.NewService(embeddingsGateway, aiClient, responsesGateway)
	aiScorer := scores.NewAiScorer(aiClient)

	evaluator := evaluation.NewCannedResponseEvaluator(
		query.NewChatResponseRetriever(queryService),
		scores.NewRunner(aiScorer),
		evaluation.NewCSVReporter(),
		evaluation.NewMarkdownReporter(),
	)
	err := evaluator.Run(".", queries)
	if err != nil {
		log.Fatalln("error running the evaluator", err)
	}
	log.Println("evaluation complete")
}

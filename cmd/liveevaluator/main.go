package main

import (
	"github.com/initialcapacity/ai-starter/internal/evaluation"
	"github.com/initialcapacity/ai-starter/internal/query"
	"github.com/initialcapacity/ai-starter/internal/scores"
	"github.com/initialcapacity/ai-starter/pkg/ai"
	"github.com/initialcapacity/ai-starter/pkg/dbsupport"
	"github.com/initialcapacity/ai-starter/pkg/websupport"
	"log"
)

func main() {
	openAiKey := websupport.RequireEnvironmentVariable[string]("OPEN_AI_KEY")
	databaseUrl := websupport.RequireEnvironmentVariable[string]("DATABASE_URL")

	db := dbsupport.CreateConnection(databaseUrl)
	responsesGateway := query.NewResponsesGateway(db)
	scoresGateway := scores.NewGateway(db)
	options := ai.LLMOptions{ChatModel: "gpt-5-mini", EmbeddingsModel: "text-embedding-3-large", Temperature: 1}
	aiClient := ai.NewClient(openAiKey, "https://api.openai.com/v1", options)
	aiScorer := scores.NewAiScorer(aiClient)

	evaluator := evaluation.NewPastResponseEvaluator(responsesGateway, scoresGateway, aiScorer)
	err := evaluator.Run()
	if err != nil {
		log.Fatalln("error running the evaluator", err)
	}
	log.Println("evaluation complete")
}

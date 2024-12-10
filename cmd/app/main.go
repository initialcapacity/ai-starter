package main

import (
	"github.com/initialcapacity/ai-starter/internal/app"
	"github.com/initialcapacity/ai-starter/pkg/ai"
	"github.com/initialcapacity/ai-starter/pkg/dbsupport"
	"github.com/initialcapacity/ai-starter/pkg/websupport"
	"log"
)

func main() {
	host := websupport.EnvironmentVariable("HOST", "")
	port := websupport.EnvironmentVariable("PORT", 8778)
	openAiEndpoint := websupport.EnvironmentVariable("OPEN_AI_ENDPOINT", "https://api.openai.com/v1")
	openAiKey := websupport.RequireEnvironmentVariable[string]("OPEN_AI_KEY")
	databaseUrl := websupport.RequireEnvironmentVariable[string]("DATABASE_URL")

	options := ai.LLMOptions{ChatModel: "gpt-4o", EmbeddingsModel: "text-embedding-3-large", Temperature: 1}
	aiClient := ai.NewClient(openAiKey, openAiEndpoint, options)
	db := dbsupport.CreateConnection(databaseUrl)

	server := websupport.NewServer(app.Handlers(aiClient, db))

	_, done := server.Start(host, port)
	log.Fatal(<-done)
}

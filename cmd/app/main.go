package main

import (
	"github.com/initialcapacity/ai-starter/internal/app"
	"github.com/initialcapacity/ai-starter/pkg/websupport"
	"log"
)

func main() {
	host := websupport.EnvironmentVariable("HOST", "")
	port := websupport.EnvironmentVariable("PORT", 8778)
	openAiKey := websupport.RequireEnvironmentVariable[string]("OPEN_AI_KEY")

	server := websupport.NewServer(app.Handlers(openAiKey))

	_, done := server.Start(host, port)
	log.Fatal(<-done)
}

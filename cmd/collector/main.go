package main

import (
	"github.com/initialcapacity/ai-starter/internal/collector"
	ai2 "github.com/initialcapacity/ai-starter/pkg/ai"
	"github.com/initialcapacity/ai-starter/pkg/dbsupport"
	"github.com/initialcapacity/ai-starter/pkg/feedsupport"
	"github.com/initialcapacity/ai-starter/pkg/websupport"
	"github.com/tiktoken-go/tokenizer"
	"log/slog"
	"net/http"
	"strings"
)

func main() {
	databaseUrl := websupport.RequireEnvironmentVariable[string]("DATABASE_URL")
	feeds := websupport.RequireEnvironmentVariable[string]("FEEDS")
	feedUrls := strings.Split(feeds, ",")

	client := http.Client{}
	db := dbsupport.CreateConnection(databaseUrl)

	parser := feedsupport.NewParser(client)
	extractor := feedsupport.NewExtractor(client)
	dataGateway := collector.NewDataGateway(db)
	t := ai2.NewTokenizer(tokenizer.Cl100kBase)
	chunksGateway := collector.NewChunksGateway(db)
	chunker := ai2.NewChunker(t, 6000)
	chunksService := collector.NewChunksService(chunker, chunksGateway)

	c := collector.New(parser, extractor, dataGateway, chunksService)

	err := c.Collect(feedUrls)

	if err == nil {
		slog.Info("successful collection")
	} else {
		slog.Error("unsuccessful collection: %w", err)
	}
}

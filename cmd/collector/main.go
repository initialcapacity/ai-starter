package main

import (
	"github.com/initialcapacity/ai-starter/internal/collection"
	"github.com/initialcapacity/ai-starter/pkg/ai"
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
	dataGateway := collection.NewDataGateway(db)
	t := ai.NewTokenizer(tokenizer.Cl100kBase)
	chunksGateway := collection.NewChunksGateway(db)
	chunker := ai.NewChunker(t, 6000)
	chunksService := collection.NewChunksService(chunker, chunksGateway)
	runsGateway := collection.NewCollectionRunsGateway(db)

	c := collection.New(parser, extractor, dataGateway, chunksService, runsGateway)

	err := c.Collect(feedUrls)

	if err == nil {
		slog.Info("successful collection")
	} else {
		slog.Error("unsuccessful collection: %w", slog.Any("error", err))
	}
}

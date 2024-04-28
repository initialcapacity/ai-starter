package main

import (
	"database/sql"
	"fmt"
	"github.com/initialcapacity/ai-starter/internal/collector"
	"github.com/initialcapacity/ai-starter/pkg/feedsupport"
	"github.com/initialcapacity/ai-starter/pkg/websupport"
	_ "github.com/lib/pq"
	"log"
	"log/slog"
	"net/http"
	"strings"
)

func main() {
	databaseUrl := websupport.RequireEnvironmentVariable[string]("DATABASE_URL")
	feeds := websupport.RequireEnvironmentVariable[string]("FEEDS")
	feedUrls := strings.Split(feeds, ",")

	client := http.Client{}
	db, err := sql.Open("postgres", databaseUrl)
	if err != nil {
		log.Fatal(fmt.Errorf("unable to connect to database: %w", err))
	}

	parser := feedsupport.NewParser(client)
	extractor := feedsupport.NewExtractor(client)
	dataGateway := collector.NewDataGateway(db)

	c := collector.New(parser, extractor, dataGateway)

	err = c.Collect(feedUrls)

	if err == nil {
		slog.Info("successful collection")
	} else {
		slog.Error("unsuccessful collection: %w", err)
	}
}

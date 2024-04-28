package functions

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/initialcapacity/ai-starter/internal/analyzer"
	"github.com/initialcapacity/ai-starter/internal/collector"
	"github.com/initialcapacity/ai-starter/pkg/feedsupport"
	"github.com/initialcapacity/ai-starter/pkg/websupport"
	"log"
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
	db, err := sql.Open("postgres", databaseUrl)
	if err != nil {
		log.Fatal(fmt.Errorf("unable to connect to database: %w", err))
	}

	parser := feedsupport.NewParser(client)
	extractor := feedsupport.NewExtractor(client)
	dataGateway := collector.NewDataGateway(db)

	c := collector.New(parser, extractor, dataGateway)

	return c.Collect(feedUrls)
}

func triggerAnalyze(ctx context.Context, e event.Event) error {
	return analyzer.Analyze()
}

package collection_test

import (
	"database/sql"
	"fmt"
	"github.com/initialcapacity/ai-starter/internal/collection"
	"github.com/initialcapacity/ai-starter/pkg/dbsupport"
	"github.com/initialcapacity/ai-starter/pkg/feedsupport"
	"github.com/initialcapacity/ai-starter/pkg/testsupport"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestCollector_Collect(t *testing.T) {
	feedEndpoint := testsupport.StartTestServer(t, func(mux *http.ServeMux) {
		testsupport.Handle(mux, "/feed1", "<html><h1>some text from feed 1</h1></html>")
		testsupport.Handle(mux, "/feed2", "<html><h1>some text from feed 2</h1></html>")
	})
	rssEndpoint := testsupport.StartTestServer(t, func(mux *http.ServeMux) {
		testsupport.Handle(mux, "/", fmt.Sprintf(`
				<rss>
					<channel>
						<item><link>%s/feed1</link></item>
						<item><link>%s/feed2</link></item>
					</channel>
				</rss>
			`, feedEndpoint, feedEndpoint))
	})

	testDb := testsupport.NewTestDb(t)
	client := http.Client{}

	extractor := feedsupport.NewExtractor(client)
	parser := feedsupport.NewParser(client)
	dataGateway := collection.NewDataGateway(testDb.DB)
	chunksGateway := collection.NewChunksGateway(testDb.DB)
	chunksService := collection.NewChunksService(DummyChunker{}, chunksGateway)
	runsGateway := collection.NewRunsGateway(testDb.DB)

	collect := collection.New(parser, extractor, dataGateway, chunksService, runsGateway)

	err := collect.Collect([]string{rssEndpoint})
	assert.NoError(t, err)

	content, err := dbsupport.Query(testDb.DB, "select content from chunks", func(rows *sql.Rows, content *string) error {
		return rows.Scan(content)
	})
	assert.NoError(t, err)
	testsupport.AssertContainsExactly(t, []string{"some text ", "from feed 1", "some text ", "from feed 2"}, content)

	result := testDb.QueryOneMap("select feeds_collected, articles_collected, chunks_collected, errors from collection_runs")
	assert.Equal(t, map[string]any{
		"feeds_collected":    int64(1),
		"articles_collected": int64(2),
		"chunks_collected":   int64(4),
		"errors":             int64(0),
	}, result)
}

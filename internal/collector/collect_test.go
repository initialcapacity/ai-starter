package collector_test

import (
	"database/sql"
	"fmt"
	"github.com/initialcapacity/ai-starter/internal/collector"
	"github.com/initialcapacity/ai-starter/pkg/dbsupport"
	"github.com/initialcapacity/ai-starter/pkg/feedsupport"
	"github.com/initialcapacity/ai-starter/pkg/testsupport"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestCollector_Collect(t *testing.T) {
	feedEndpoint, feedServer := testsupport.StartTestServer(t, func(mux *http.ServeMux) {
		testsupport.Handle(mux, "/feed1", "<html><h1>some text from feed 1</h1></html>")
		testsupport.Handle(mux, "/feed2", "<html><h1>some text from feed 2</h1></html>")
	})
	rssEndpoint, rssServer := testsupport.StartTestServer(t, func(mux *http.ServeMux) {
		testsupport.Handle(mux, "/", fmt.Sprintf(`
				<rss>
					<channel>
						<item><link>%s/feed1</link></item>
						<item><link>%s/feed2</link></item>
					</channel>
				</rss>
			`, feedEndpoint, feedEndpoint))
	})
	defer func() {
		testsupport.StopTestServer(t, rssServer)
		testsupport.StopTestServer(t, feedServer)
	}()

	testDb := testsupport.NewTestDb(t)
	defer testDb.Close()
	testDb.ClearTables()
	client := http.Client{}

	extractor := feedsupport.NewExtractor(client)
	parser := feedsupport.NewParser(client)
	dataGateway := collector.NewDataGateway(testDb.DB)
	chunksGateway := collector.NewChunksGateway(testDb.DB)
	chunksService := collector.NewChunksService(DummyChunker{}, chunksGateway)

	collect := collector.New(parser, extractor, dataGateway, chunksService)

	err := collect.Collect([]string{rssEndpoint})
	assert.NoError(t, err)

	content, err := dbsupport.Query(testDb.DB, "select content from chunks", func(rows *sql.Rows, content *string) error {
		return rows.Scan(content)
	})
	assert.NoError(t, err)

	testsupport.AssertContainsExactly(t, []string{"some text ", "from feed 1", "some text ", "from feed 2"}, content)
}

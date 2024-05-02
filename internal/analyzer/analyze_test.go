package analyzer_test

import (
	"context"
	"fmt"
	"github.com/initialcapacity/ai-starter/internal/analyzer"
	"github.com/initialcapacity/ai-starter/internal/collector"
	"github.com/initialcapacity/ai-starter/pkg/testsupport"
	"github.com/stretchr/testify/assert"
	"net/http"
	"strings"
	"testing"
)

func TestAnalyzer_Analyze(t *testing.T) {
	vector := testsupport.CreateVector(0)
	endpoint, server := testsupport.StartTestServer(t, func(mux *http.ServeMux) {
		testsupport.Handle(mux, "/embeddings", fmt.Sprintf(`{
			"data": [
				{ "embedding": %s }
			]
		}`, toString(vector)))
	})
	defer testsupport.StopTestServer(t, server)

	testDb := testsupport.NewTestDb(t)
	defer testDb.Close()
	testDb.ClearTables()

	embeddingsGateway := analyzer.NewEmbeddingsGateway(testDb.DB)
	chunksGateway := collector.NewChunksGateway(testDb.DB)

	testDb.Execute("insert into data (id, source, content) values ('aaaaaaaa-2f3f-4bc9-8dba-ba397156cc16', 'https://example.com', 'some content')")
	testDb.Execute("insert into chunks (id, data_id, content) values ('bbbbbbbb-2f3f-4bc9-8dba-ba397156cc16', 'aaaaaaaa-2f3f-4bc9-8dba-ba397156cc16','chunk1')")

	a := analyzer.NewAnalyzer(chunksGateway, embeddingsGateway, testsupport.NewTestAiClient(endpoint))

	err := a.Analyze(context.Background())
	assert.NoError(t, err)

	chunk1, err := embeddingsGateway.FindSimilar(testsupport.CreateVector(0))
	assert.NoError(t, err)
	assert.Equal(t, analyzer.CitedChunkRecord{Content: "chunk1", Source: "https://example.com"}, chunk1)
}

func toString(vector []float32) string {
	builder := strings.Builder{}

	builder.WriteString("[")
	for i, v := range vector {
		builder.WriteString(fmt.Sprint(v))
		if i < len(vector)-1 {
			builder.WriteString(", ")
		}
	}
	builder.WriteString("]")

	return builder.String()
}

package analysis_test

import (
	"context"
	"github.com/initialcapacity/ai-starter/internal/analysis"
	"github.com/initialcapacity/ai-starter/internal/collection"
	"github.com/initialcapacity/ai-starter/pkg/testsupport"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestAnalyzer_Analyze(t *testing.T) {
	vector := testsupport.CreateVector(0)

	endpoint, server := testsupport.StartTestServer(t, func(mux *http.ServeMux) {
		testsupport.HandleCreateEmbedding(mux, vector)
	})
	defer testsupport.StopTestServer(t, server)

	testDb := testsupport.NewTestDb(t)
	defer testDb.Close()

	testDb.Execute("insert into data (id, source, content) values ('aaaaaaaa-2f3f-4bc9-8dba-ba397156cc16', 'https://example.com', 'some content')")
	testDb.Execute("insert into chunks (id, data_id, content) values ('bbbbbbbb-2f3f-4bc9-8dba-ba397156cc16', 'aaaaaaaa-2f3f-4bc9-8dba-ba397156cc16','chunk1')")

	embeddingsGateway := analysis.NewEmbeddingsGateway(testDb.DB)
	chunksGateway := collection.NewChunksGateway(testDb.DB)
	aiClient := testsupport.NewTestAiClient(endpoint)
	runsGateway := analysis.NewAnalysisRunsGateway(testDb.DB)

	a := analysis.NewAnalyzer(chunksGateway, embeddingsGateway, aiClient, runsGateway)

	err := a.Analyze(context.Background())
	assert.NoError(t, err)

	chunk, err := embeddingsGateway.FindSimilar(testsupport.CreateVector(0))
	assert.NoError(t, err)
	assert.Equal(t, analysis.CitedChunkRecord{Content: "chunk1", Source: "https://example.com"}, chunk)

	result := testDb.QueryOneMap("select chunks_analyzed, analysis_runs.embeddings_created, errors from analysis_runs")
	assert.Equal(t, map[string]any{
		"chunks_analyzed":    int64(1),
		"embeddings_created": int64(1),
		"errors":             int64(0),
	}, result)
}

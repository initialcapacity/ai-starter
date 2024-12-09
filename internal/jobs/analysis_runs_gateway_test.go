package jobs_test

import (
	"github.com/initialcapacity/ai-starter/internal/jobs"
	"github.com/initialcapacity/ai-starter/pkg/testsupport"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAnalysisRunsGateway_Create(t *testing.T) {
	testDb := testsupport.NewTestDb(t)
	defer testDb.Close()
	gateway := jobs.NewAnalysisRunsGateway(testDb.DB)

	record, err := gateway.Create(3, 4, 5)
	require.NoError(t, err)

	assert.Equal(t, 3, record.ChunksAnalyzed)
	assert.Equal(t, 4, record.EmbeddingsCreated)
	assert.Equal(t, 5, record.NumberOfErrors)

	result := testDb.QueryOneMap("select chunks_analyzed, embeddings_created, errors from analysis_runs where id = $1", record.Id)
	assert.Equal(t, map[string]any{
		"chunks_analyzed":    int64(3),
		"embeddings_created": int64(4),
		"errors":             int64(5),
	}, result)
}

func TestAnalysisRunsGateway_List(t *testing.T) {
	testDb := testsupport.NewTestDb(t)
	defer testDb.Close()
	gateway := jobs.NewAnalysisRunsGateway(testDb.DB)

	testDb.Execute("insert into analysis_runs (chunks_analyzed, embeddings_created, errors) values (2, 3, 4)")
	testDb.Execute("insert into analysis_runs (chunks_analyzed, embeddings_created, errors) values (12, 13, 14)")

	records, err := gateway.List()
	require.NoError(t, err)

	assert.Equal(t, 2, len(records))
	assert.Len(t, records[0].Id, 36)
	assert.Equal(t, 12, records[0].ChunksAnalyzed)
	assert.Equal(t, 13, records[0].EmbeddingsCreated)
	assert.Equal(t, 14, records[0].NumberOfErrors)
	assert.Len(t, records[1].Id, 36)
	assert.Equal(t, 2, records[1].ChunksAnalyzed)
	assert.Equal(t, 3, records[1].EmbeddingsCreated)
	assert.Equal(t, 4, records[1].NumberOfErrors)
}

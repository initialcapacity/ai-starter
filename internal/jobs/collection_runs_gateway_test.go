package jobs_test

import (
	"github.com/initialcapacity/ai-starter/internal/jobs"
	"github.com/initialcapacity/ai-starter/pkg/testsupport"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCollectionRunsGateway_Create(t *testing.T) {
	testDb := testsupport.NewTestDb(t)
	defer testDb.Close()
	gateway := jobs.NewCollectionRunsGateway(testDb.DB)

	record, err := gateway.Create(3, 4, 5, 6)
	require.NoError(t, err)

	assert.Equal(t, 3, record.FeedsCollected)
	assert.Equal(t, 4, record.ArticlesCollected)
	assert.Equal(t, 5, record.ChunksCollected)
	assert.Equal(t, 6, record.NumberOfErrors)

	result := testDb.QueryOneMap("select feeds_collected, articles_collected, chunks_collected, errors from collection_runs where id = $1", record.Id)
	assert.Equal(t, map[string]any{
		"feeds_collected":    int64(3),
		"articles_collected": int64(4),
		"chunks_collected":   int64(5),
		"errors":             int64(6),
	}, result)
}

func TestCollectionRunsGateway_List(t *testing.T) {
	testDb := testsupport.NewTestDb(t)
	defer testDb.Close()
	gateway := jobs.NewCollectionRunsGateway(testDb.DB)

	testDb.Execute("insert into collection_runs (feeds_collected, articles_collected, chunks_collected, errors) values (2, 3, 4, 5)")
	testDb.Execute("insert into collection_runs (feeds_collected, articles_collected, chunks_collected, errors) values (12, 13, 14, 15)")

	records, err := gateway.List()
	require.NoError(t, err)

	assert.Equal(t, 2, len(records))
	assert.Len(t, records[0].Id, 36)
	assert.Equal(t, 12, records[0].FeedsCollected)
	assert.Equal(t, 13, records[0].ArticlesCollected)
	assert.Equal(t, 14, records[0].ChunksCollected)
	assert.Equal(t, 15, records[0].NumberOfErrors)
	assert.Len(t, records[1].Id, 36)
	assert.Equal(t, 2, records[1].FeedsCollected)
	assert.Equal(t, 3, records[1].ArticlesCollected)
	assert.Equal(t, 4, records[1].ChunksCollected)
	assert.Equal(t, 5, records[1].NumberOfErrors)
}

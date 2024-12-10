package collection_test

import (
	"database/sql"
	"github.com/initialcapacity/ai-starter/internal/collection"
	"github.com/initialcapacity/ai-starter/pkg/dbsupport"
	"github.com/initialcapacity/ai-starter/pkg/testsupport"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestChunksService_SaveChunks(t *testing.T) {
	testDb := testsupport.NewTestDb(t)
	defer testDb.Close()

	chunksGateway := collection.NewChunksGateway(testDb.DB)
	chunksService := collection.NewChunksService(DummyChunker{}, chunksGateway)
	testDb.Execute("insert into data (id, source, content) values ('41345dc1-2f3f-4bc9-8dba-ba397156cc16', 'https://example.com', 'some content')")

	count, err := chunksService.SaveChunks("41345dc1-2f3f-4bc9-8dba-ba397156cc16", "some content")
	assert.NoError(t, err)
	assert.Equal(t, 2, count)

	content, err := dbsupport.Query(testDb.DB, "select content from chunks", func(rows *sql.Rows, content *string) error {
		return rows.Scan(content)
	})
	assert.NoError(t, err)
	testsupport.AssertContainsExactly(t, []string{"some c", "ontent"}, content)
}

type DummyChunker struct {
}

func (c DummyChunker) Split(text string) []string {
	midpoint := len(text) / 2

	return []string{text[:midpoint], text[midpoint:]}
}

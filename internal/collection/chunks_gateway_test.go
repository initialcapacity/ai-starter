package collection_test

import (
	"database/sql"
	"github.com/initialcapacity/ai-starter/internal/collection"
	"github.com/initialcapacity/ai-starter/pkg/dbsupport"
	"github.com/initialcapacity/ai-starter/pkg/testsupport"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestChunksGateway_Save(t *testing.T) {
	testDb := testsupport.NewTestDb(t)
	defer testDb.Close()

	gateway := collection.NewChunksGateway(testDb.DB)

	testDb.Execute("insert into data (id, source, content) values ('41345dc1-2f3f-4bc9-8dba-ba397156cc16', 'https://example.com', 'some content')")
	err := gateway.Save("41345dc1-2f3f-4bc9-8dba-ba397156cc16", "a chunk")
	assert.NoError(t, err)

	content, err := dbsupport.QueryOne(testDb.DB, "select content from chunks", func(row *sql.Row, content *string) error {
		return row.Scan(content)
	})
	assert.NoError(t, err)
	assert.Equal(t, "a chunk", content)
}

func TestChunksGateway_Get(t *testing.T) {
	testDb := testsupport.NewTestDb(t)
	defer testDb.Close()

	gateway := collection.NewChunksGateway(testDb.DB)

	testDb.Execute("insert into data (id, source, content) values ('aaaaaaaa-2f3f-4bc9-8dba-ba397156cc16', 'https://example.com', 'some content')")
	testDb.Execute("insert into chunks (id, data_id, content) values ('bbbbbbbb-2f3f-4bc9-8dba-ba397156cc16', 'aaaaaaaa-2f3f-4bc9-8dba-ba397156cc16','a chunk')")

	record, err := gateway.Get("bbbbbbbb-2f3f-4bc9-8dba-ba397156cc16")
	assert.NoError(t, err)
	assert.Equal(t, collection.ChunkRecord{
		DataId:  "aaaaaaaa-2f3f-4bc9-8dba-ba397156cc16",
		Content: "a chunk",
	}, record)
}

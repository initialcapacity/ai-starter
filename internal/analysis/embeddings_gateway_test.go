package analysis_test

import (
	"database/sql"
	"github.com/initialcapacity/ai-starter/internal/analysis"
	"github.com/initialcapacity/ai-starter/pkg/dbsupport"
	"github.com/initialcapacity/ai-starter/pkg/testsupport"
	"github.com/pgvector/pgvector-go"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEmbeddingsGateway_UnprocessedIds(t *testing.T) {
	testDb := testsupport.NewTestDb(t)
	defer testDb.Close()

	gateway := analysis.NewEmbeddingsGateway(testDb.DB)

	testDb.Execute("insert into data (id, source, content) values ('aaaaaaaa-2f3f-4bc9-8dba-ba397156cc16', 'https://example.com', 'some content')")
	testDb.Execute("insert into chunks (id, data_id, content) values ('bbbbbbbb-2f3f-4bc9-8dba-ba397156cc16', 'aaaaaaaa-2f3f-4bc9-8dba-ba397156cc16','a chunk')")
	testDb.Execute("insert into chunks (id, data_id, content) values ('cccccccc-2f3f-4bc9-8dba-ba397156cc16', 'aaaaaaaa-2f3f-4bc9-8dba-ba397156cc16','a chunk')")
	vector := testsupport.CreateVector(0)
	testDb.Execute("insert into embeddings (chunk_id, embedding) values ('bbbbbbbb-2f3f-4bc9-8dba-ba397156cc16', $1)", pgvector.NewVector(vector))

	ids, err := gateway.UnprocessedIds()
	assert.NoError(t, err)

	assert.Equal(t, []string{"cccccccc-2f3f-4bc9-8dba-ba397156cc16"}, ids)
}

func TestEmbeddingsGateway_Save(t *testing.T) {
	testDb := testsupport.NewTestDb(t)
	defer testDb.Close()

	gateway := analysis.NewEmbeddingsGateway(testDb.DB)

	testDb.Execute("insert into data (id, source, content) values ('aaaaaaaa-2f3f-4bc9-8dba-ba397156cc16', 'https://example.com', 'some content')")
	testDb.Execute("insert into chunks (id, data_id, content) values ('bbbbbbbb-2f3f-4bc9-8dba-ba397156cc16', 'aaaaaaaa-2f3f-4bc9-8dba-ba397156cc16','a chunk')")

	err := gateway.Save("bbbbbbbb-2f3f-4bc9-8dba-ba397156cc16", testsupport.CreateVector(0))
	assert.NoError(t, err)

	chunkId, err := dbsupport.QueryOne(testDb.DB, "select chunk_id from embeddings", func(row *sql.Row, chunkId *string) error {
		return row.Scan(chunkId)
	})
	assert.NoError(t, err)
	assert.Equal(t, "bbbbbbbb-2f3f-4bc9-8dba-ba397156cc16", chunkId)
}

func TestEmbeddingsGateway_FindSimilar(t *testing.T) {
	testDb := testsupport.NewTestDb(t)
	defer testDb.Close()

	gateway := analysis.NewEmbeddingsGateway(testDb.DB)

	testDb.Execute("insert into data (id, source, content) values ('aaaaaaaa-2f3f-4bc9-8dba-ba397156cc16', 'https://example.com', 'some content')")
	testDb.Execute("insert into chunks (id, data_id, content) values ('bbbbbbbb-2f3f-4bc9-8dba-ba397156cc16', 'aaaaaaaa-2f3f-4bc9-8dba-ba397156cc16','a chunk')")
	testDb.Execute("insert into chunks (id, data_id, content) values ('cccccccc-2f3f-4bc9-8dba-ba397156cc16', 'aaaaaaaa-2f3f-4bc9-8dba-ba397156cc16','another chunk')")

	err := gateway.Save("bbbbbbbb-2f3f-4bc9-8dba-ba397156cc16", testsupport.CreateVector(0))
	assert.NoError(t, err)
	err = gateway.Save("cccccccc-2f3f-4bc9-8dba-ba397156cc16", testsupport.CreateVector(1))
	assert.NoError(t, err)

	similar, err := gateway.FindSimilar(testsupport.CreateVector(1))
	assert.NoError(t, err)

	assert.Equal(t, analysis.CitedChunkRecord{Content: "another chunk", Source: "https://example.com"}, similar)
}

package collection_test

import (
	"database/sql"
	"github.com/initialcapacity/ai-starter/internal/collection"
	"github.com/initialcapacity/ai-starter/pkg/dbsupport"
	"github.com/initialcapacity/ai-starter/pkg/testsupport"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDataGateway_Save(t *testing.T) {
	testDb := testsupport.NewTestDb(t)
	defer testDb.Close()

	dataGateway := collection.NewDataGateway(testDb.DB)

	id, err := dataGateway.Save("https://example.com", "some content")
	assert.NoError(t, err)

	savedContent, err := dbsupport.QueryOne(testDb.DB, "select content from data where id = $1", func(row *sql.Row, content *string) error {
		return row.Scan(content)
	}, id)
	assert.NoError(t, err)

	assert.Equal(t, "some content", savedContent)
}

func TestDataGateway_Exists(t *testing.T) {
	testDb := testsupport.NewTestDb(t)
	defer testDb.Close()

	dataGateway := collection.NewDataGateway(testDb.DB)

	_, err := dataGateway.Save("https://example.com", "some content")
	assert.NoError(t, err)

	shouldExist, err := dataGateway.Exists("https://example.com")
	assert.NoError(t, err)
	assert.True(t, shouldExist)

	shouldNotExist, err := dataGateway.Exists("https://not-there.example.com")
	assert.NoError(t, err)
	assert.False(t, shouldNotExist)
}

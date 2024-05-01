package dbsupport_test

import (
	"database/sql"
	"github.com/initialcapacity/ai-starter/pkg/dbsupport"
	"github.com/initialcapacity/ai-starter/pkg/testsupport"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestQueryOne(t *testing.T) {
	testDb := testsupport.NewTestDb(t)
	defer testDb.Close()

	result, err := dbsupport.QueryOne(testDb.DB, "select 1", func(row *sql.Row, number *int) error {
		return row.Scan(number)
	})

	assert.NoError(t, err)
	assert.Equal(t, 1, result)
}

func TestQuery(t *testing.T) {
	testDb := testsupport.NewTestDb(t)
	defer testDb.Close()

	result, err := dbsupport.Query(testDb.DB, "select * from generate_series(1, 5)", func(rows *sql.Rows, number *int) error {
		return rows.Scan(number)
	})

	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3, 4, 5}, result)
}

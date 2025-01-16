package testsupport

import (
	"database/sql"
	"fmt"
	"github.com/initialcapacity/ai-starter/pkg/dbsupport"
	"github.com/stretchr/testify/assert"
	"math/rand/v2"
	"testing"
)

type TestDb struct {
	DB         *sql.DB
	t          *testing.T
	testDbName string
}

func NewTestDb(t *testing.T) *TestDb {
	testDbName := fmt.Sprintf("starter_test_%d", rand.IntN(1_000_000))
	WithSuperDb(t, func(superDb *sql.DB) {
		_, err := superDb.Exec(fmt.Sprintf("create database %s template starter_test", testDbName))
		assert.NoError(t, err, "unable to create test database")
	})

	testDb := &TestDb{
		DB:         dbsupport.CreateConnection(fmt.Sprintf("postgres://starter:starter@localhost:5432/%s?sslmode=disable", testDbName)),
		t:          t,
		testDbName: testDbName,
	}
	t.Cleanup(testDb.close)
	return testDb
}

func (tdb *TestDb) close() {
	err := tdb.DB.Close()
	assert.NoError(tdb.t, err)

	WithSuperDb(tdb.t, func(superDb *sql.DB) {
		_, err = superDb.Exec(fmt.Sprintf("drop database %s", tdb.testDbName))
		assert.NoError(tdb.t, err, "unable to drop test database")
	})
}

func (tdb *TestDb) Execute(statement string, arguments ...any) {
	_, err := tdb.DB.Exec(statement, arguments...)
	assert.NoError(tdb.t, err)
}

func (tdb *TestDb) QueryMap(statement string, arguments ...any) []map[string]any {
	rows, err := tdb.DB.Query(statement, arguments...)
	assert.NoError(tdb.t, err)
	columns, err := rows.Columns()
	assert.NoError(tdb.t, err)
	result := make([]map[string]any, 0)

	for rows.Next() {
		values := make([]any, len(columns))
		valuePointers := make([]any, len(columns))
		for i, _ := range values {
			valuePointers[i] = &values[i]
		}

		err = rows.Scan(valuePointers...)
		assert.NoError(tdb.t, err)

		rowResult := make(map[string]any)
		for i, columnName := range columns {
			rowResult[columnName] = values[i]
		}
		result = append(result, rowResult)
	}

	return result
}

func (tdb *TestDb) QueryOneMap(statement string, arguments ...any) map[string]any {
	results := tdb.QueryMap(statement, arguments...)
	assert.Len(tdb.t, results, 1, fmt.Sprintf("Expected one result but got %d", len(results)))
	return results[0]
}

func WithSuperDb(t *testing.T, action func(superDb *sql.DB)) {
	superDb := dbsupport.CreateConnection("postgres://super_test@localhost:5432/postgres?sslmode=disable")
	defer func(superDb *sql.DB) {
		err := superDb.Close()
		assert.NoError(t, err, "unable to close connection to postgres")
	}(superDb)

	action(superDb)
}

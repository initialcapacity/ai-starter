package testsupport

import (
	"database/sql"
	"github.com/initialcapacity/ai-starter/pkg/dbsupport"
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestDb struct {
	DB *sql.DB
	t  *testing.T
}

func NewTestDb(t *testing.T) *TestDb {
	return &TestDb{DB: dbsupport.CreateConnection("postgres://starter:starter@localhost:5432/starter_test?sslmode=disable"), t: t}
}

func (tdb *TestDb) Close() {
	err := tdb.DB.Close()
	assert.NoError(tdb.t, err)
}

func (tdb *TestDb) ClearTables() {
	tdb.Execute("delete from embeddings")
	tdb.Execute("delete from chunks")
	tdb.Execute("delete from data")
}

func (tdb *TestDb) Execute(statement string, arguments ...any) {
	_, err := tdb.DB.Exec(statement, arguments...)
	assert.NoError(tdb.t, err)
}

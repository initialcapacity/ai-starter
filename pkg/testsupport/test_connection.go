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
	withSuperDb(t, func(superDb *sql.DB) {
		_, err := superDb.Exec(fmt.Sprintf("create database %s template starter_test", testDbName))
		assert.NoError(t, err, "unable to create test database")
	})

	println("created test database: " + testDbName)

	return &TestDb{
		DB:         dbsupport.CreateConnection(fmt.Sprintf("postgres://starter:starter@localhost:5432/%s?sslmode=disable", testDbName)),
		t:          t,
		testDbName: testDbName,
	}
}

func (tdb *TestDb) Close() {
	err := tdb.DB.Close()
	assert.NoError(tdb.t, err)

	withSuperDb(tdb.t, func(superDb *sql.DB) {
		_, err = superDb.Exec(fmt.Sprintf("drop database %s", tdb.testDbName))
		assert.NoError(tdb.t, err, "unable to drop test database")
	})
}

func (tdb *TestDb) Execute(statement string, arguments ...any) {
	_, err := tdb.DB.Exec(statement, arguments...)
	assert.NoError(tdb.t, err)
}

func withSuperDb(t *testing.T, action func(superDb *sql.DB)) {
	superDb := dbsupport.CreateConnection("postgres://super_test@localhost:5432/postgres?sslmode=disable")
	defer func(superDb *sql.DB) {
		err := superDb.Close()
		assert.NoError(t, err, "unable to close connection to postgres")
	}(superDb)

	action(superDb)
}

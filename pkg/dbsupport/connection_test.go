package dbsupport_test

import (
	"database/sql"
	"github.com/initialcapacity/ai-starter/pkg/dbsupport"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Connection() *sql.DB {
	return dbsupport.CreateConnection("postgres://starter:starter@localhost:5432/starter_test?sslmode=disable")
}

func Close(t *testing.T, db *sql.DB) {
	err := db.Close()
	assert.NoError(t, err)
}

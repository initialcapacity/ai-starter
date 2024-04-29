package dbsupport

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

func CreateConnection(databaseUrl string) *sql.DB {
	db, err := sql.Open("postgres", databaseUrl)
	if err != nil {
		log.Fatal(fmt.Errorf("unable to connect to database: %w", err))
	}

	return db
}

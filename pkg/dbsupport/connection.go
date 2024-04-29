package dbsupport

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
)

func CreateConnection(databaseUrl string) *sql.DB {
	db, err := sql.Open("pgx", databaseUrl)
	if err != nil {
		log.Fatal(fmt.Errorf("unable to connect to database: %w", err))
	}

	return db
}

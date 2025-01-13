package main

import (
	"errors"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/initialcapacity/ai-starter/pkg/dbsupport"
	"github.com/initialcapacity/ai-starter/pkg/websupport"
	"log"
)

func main() {
	databaseUrl := websupport.RequireEnvironmentVariable[string]("DATABASE_URL")
	migrationsLocation := websupport.EnvironmentVariable("MIGRATIONS_LOCATION", "file://./databases/starter")

	db := dbsupport.CreateConnection(databaseUrl)
	driver, err := pgx.WithInstance(db, &pgx.Config{})
	if err != nil {
		log.Fatalf("failed to connect to %s, %s", databaseUrl, err)
	}

	migration, err := migrate.NewWithDatabaseInstance(migrationsLocation, "ai-starter", driver)
	if err != nil {
		log.Fatalf("failed to create migration instance %s, %s", databaseUrl, err)
	}

	migration.Log = logger{}
	err = migration.Up()
	if errors.Is(err, migrate.ErrNoChange) {
		log.Printf("no new migrations detected: %s\n", err)
	} else if err != nil {
		log.Fatalf("unable to migrate %s, %s", databaseUrl, err)
	}

	log.Printf("successfully migrated %s\n", databaseUrl)
}

type logger struct{}

func (l logger) Verbose() bool { return true }

func (l logger) Printf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

# AI Starter

## Run locally

```shell
psql postgres < ./databases/create_databases.sql
DATABASE_URL="postgres://starter:starter@localhost:5432/starter_development?sslmode=disable" go run ./cmd/migrate
cp .env.example .env
source .env
go run ./cmd/collector
go run ./cmd/analyzer
go run ./cmd/app
```

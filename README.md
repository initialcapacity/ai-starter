# AI Starter

## Setup

Install [Go](https://formulae.brew.sh/formula/go), [PostgreSQL 15](https://formulae.brew.sh/formula/postgresql@15), and
[pgvector](https://github.com/pgvector/pgvector).

## Run locally

```shell
psql postgres < ./databases/create_databases.sql
DATABASE_URL="postgres://starter:starter@localhost:5432/starter_development?sslmode=disable" go run ./cmd/migrate
cp .env.example .env # edit to add your values 
source .env
go run ./cmd/collector
go run ./cmd/analyzer
go run ./cmd/app
```

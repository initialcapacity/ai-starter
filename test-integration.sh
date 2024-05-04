#!/usr/bin/env bash

export HOST=localhost
export EXIT_CODE=0
export FEEDS=http://localhost:8123
export DATABASE_URL="postgres://starter:starter@localhost:5432/starter_integration?sslmode=disable"

echo "PROGRESS - creating databases"
psql postgres < ./databases/create_databases.sql

echo "PROGRESS - building"
rm -rf ./build
mkdir -p ./build
go build -o ./build/migrate ./cmd/migrate
go build -o ./build/cannedrss ./cmd/cannedrss
go build -o ./build/collector ./cmd/collector
go build -o ./build/analyzer ./cmd/analyzer
go build -o ./build/app ./cmd/app

echo "PROGRESS - migrating"
./build/migrate

echo "PROGRESS - starting rss"
./build/cannedrss &
RSS_PID=$!

echo "PROGRESS - collecting"
./build/collector
echo "PROGRESS - analyzing"
./build/analyzer
echo "PROGRESS - starting app"
PORT=8234 ./build/app &
APP_PID=$!

sleep 1

echo "PROGRESS - GET /"
curl --fail-with-body http://localhost:8234
EXIT_CODE=$((EXIT_CODE + $?))

echo "PROGRESS - POST /"
curl -XPOST -N --fail-with-body http://localhost:8234 -d"query=tell%20me%20about%20pickles"
EXIT_CODE=$((EXIT_CODE + $?))

echo "PROGRESS - killing app"
kill $APP_PID
echo "PROGRESS - killing rss"
kill $RSS_PID

echo ""
if [ $EXIT_CODE -eq 0 ]; then echo "FINISHED - SUCCESS"; else echo "FINISHED - FAILURE"; fi
exit $EXIT_CODE

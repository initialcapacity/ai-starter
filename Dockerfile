FROM golang:1.24-alpine AS build

WORKDIR /app

COPY . .
RUN go mod download
RUN mkdir build

RUN go build -o build/app ./cmd/app
RUN go build -o build/collector ./cmd/collector
RUN go build -o build/analyzer ./cmd/analyzer
RUN go build -o build/migrate ./cmd/migrate
RUN go build -o build/liveevaluator ./cmd/liveevaluator

FROM alpine:latest
WORKDIR /app
COPY --from=build /app/build/* /usr/local/bin/
RUN mkdir -p databases/starter
COPY --from=build /app/databases/starter/*.sql ./databases/starter/
ENTRYPOINT ["app"]

FROM golang:1.23-alpine AS build

WORKDIR /app

COPY . .
RUN go mod download
RUN mkdir build

RUN go build -o build/app ./cmd/app
RUN go build -o build/collector ./cmd/collector
RUN go build -o build/analyzer ./cmd/analyzer
RUN go build -o build/migrate ./cmd/migrate
RUN go build -o build/pastevaluator ./cmd/pastevaluator

FROM alpine:latest
COPY --from=build /app/build/* /usr/local/bin
ENTRYPOINT ["app"]

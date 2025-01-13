FROM golang:1.23

WORKDIR /app

COPY . .
RUN go mod download
RUN go build -o /usr/local/bin/app ./cmd/app
RUN go build -o /usr/local/bin/collector ./cmd/collector
RUN go build -o /usr/local/bin/analyzer ./cmd/analyzer
RUN go build -o /usr/local/bin/migrate ./cmd/migrate

ENTRYPOINT ["app"]

FROM golang:1.23

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN GOOS=linux go build -v -o /usr/local/bin/app ./cmd/app
RUN GOOS=linux go build -v -o /usr/local/bin/collector ./cmd/collector
RUN GOOS=linux go build -v -o /usr/local/bin/analyzer ./cmd/analyzer
RUN GOOS=linux go build -v -o /usr/local/bin/migrate ./cmd/migrate

ENTRYPOINT ["app"]

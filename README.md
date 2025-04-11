# AI Starter

A starter application that shows a data collector architecture for [retrieval augmented generation](https://en.wikipedia.org/wiki/Prompt_engineering#Retrieval-augmented_generation).

## Technology stack

This codebase is written [Go](https://go.dev/) and runs on Google's [Cloud Run](https://cloud.google.com/run) and
[Cloud Run Jobs](https://cloud.google.com/run/docs/quickstarts/jobs/create-execute).
It uses Go's built-in [server](https://pkg.go.dev/net/http) and [templates](https://pkg.go.dev/html/template), along
with the official [OpenAI Client](https://github.com/openai/openai-go).
It stores data in [PostgreSQL](https://www.postgresql.org/) and uses [pgvector](https://github.com/pgvector/pgvector) to
write and query embeddings. 
A [GitHub Action](https://github.com/features/actions) runs tests, builds the apps, runs migrations, then deploys to
Google Cloud.

## Architecture

The AI Starter consists of three applications communicating with one Postgres database.

1.  The data collector is a background process that collects data from one or more sources.
1.  The data analyzer is another background process that processes collected data.
1.  The web application collects a query from the user and displays a result to the user.

```mermaid
flowchart LR
    embeddings([OpenAI embeddings])
    user((User))
    app["Web App (Cloud Run)"]
    db[("PostgreSQL + pgvector")]
    llm([OpenAI completion])
    
    user -- query --> app
    app -- create embedding --> embeddings
    app -- search embeddings --> db
    app -- retrieve documents --> db
    app -- fetch text completion --> llm

    classDef node font-weight:bold,color:white,stroke:black,stroke-width:2px;
    classDef app fill:#3185FC;
    classDef db fill:#B744B8;
    classDef external fill:#FA9F42;
    classDef user fill:#ED6A5A;

    class app,collector,analyzer app;
    class db db;
    class docs,embeddings,llm external;
    class user user;
```

```mermaid
flowchart LR
    embeddings([OpenAI embeddings])
    docs(["RSS feeds"])
    db[("PostgreSQL + pgvector")]
    collector["Data Collector (Cloud Run Job)"]
    analyzer["Data Analyzer (Cloud Run Job)"]
    
    collector -- fetch documents --> docs
    collector -- save documents --> db
    analyzer -- retrieve documents --> db
    analyzer -- create embeddings --> embeddings
    analyzer -- "save embeddings (with reference)" --> db

    classDef node font-weight:bold,color:white,stroke:black,stroke-width:2px;
    classDef app fill:#3185FC;
    classDef db fill:#B744B8;
    class app,collector,analyzer app;
    classDef external fill:#FA9F42;
    classDef user fill:#ED6A5A;

    class db db;
    class docs,embeddings external;
    class user user;
```

### Collection and Analysis

The data collector fetches documents from RSS feeds sources and stores the document text in the database.
It also splits documents into chunks of less than 6000 tokens to ensure embedding and text completion calls stay below
their token limits.
The data analyzer sends document chunks to the [OpenAI Embeddings API](https://platform.openai.com/docs/guides/embeddings)
and uses pgvector to store the embeddings in PostgreSQL.

### Web Application

The web application collects the user's query and creates an embedding with the OpenAI Embeddings API.
It then searches the PostgreSQL for similar embeddings (using pgvector) and provides the corresponding chunk of text as
context for a query to the [OpenAI Chat Completion API](https://platform.openai.com/docs/api-reference/chat).

## Local development

1.  Install [Go](https://formulae.brew.sh/formula/go), [PostgreSQL 17](https://formulae.brew.sh/formula/postgresql@17),
    and [pgvector](https://github.com/pgvector/pgvector).

1.  Create and migrate the local databases.
    ```shell
    psql postgres < ./databases/create_databases.sql
    DATABASE_URL="user=starter password=starter database=starter_development host=localhost" go run ./cmd/migrate
    DATABASE_URL="user=starter password=starter database=starter_test host=localhost" go run ./cmd/migrate
    ```

1.  Copy the example environment file and fill in the necessary values.
    ```shell
    cp .env.example .env 
    source .env
    ```

1.  Run the collector and the analyzer to populate the database, then run the app and navigate to
    [localhost:8778](http://localhost:8778).        

    ```shell
    go run ./cmd/collector
    go run ./cmd/analyzer
    go run ./cmd/app
    ```

## Integration tests

The integration test script runs the collector and analyzer, then tests the app against the production OpenAI API.

```shell
source .env
go test ./cmd/integrationtest -count=1 -tags=integration
```


## Evaluation

Run an evaluation against a populated database of articles and embeddings.

```shell
source .env
go run ./cmd/evaluator
```

View the results in a csv ([scores.csv](scores.csv)) or a markdown file ([scores.md](scores.md)).

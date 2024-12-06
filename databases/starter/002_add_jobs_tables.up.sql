create table collection_runs
(
    id                 uuid primary key                  default gen_random_uuid(),
    feeds_collected    int                      not null,
    articles_collected int                      not null,
    chunks_collected   int                      not null,
    errors             int                      not null,
    created_at         timestamp with time zone not null default now()
);
grant all privileges on table data to starter;

create table analysis_runs
(
    id                 uuid primary key                  default gen_random_uuid(),
    chunks_analyzed    int                      not null,
    embeddings_created int                      not null,
    errors             int                      not null,
    created_at         timestamp with time zone not null default now()
);
grant all privileges on table data to starter;

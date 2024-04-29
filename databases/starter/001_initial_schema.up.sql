create table data
(
    id         uuid primary key                  default gen_random_uuid(),
    source     varchar                  not null,
    content    varchar                  not null,
    created_at timestamp with time zone not null default now()
);
grant all privileges on table data to starter;
create unique index index_data_source on data (source);

create table chunks
(
    id         uuid primary key                  default gen_random_uuid(),
    data_id    uuid references data (id),
    content    varchar                  not null,
    created_at timestamp with time zone not null default now()
);
grant all privileges on table chunks to starter;
create index index_chunks_data_id on chunks (data_id);

create table embeddings
(
    id         uuid primary key                  default gen_random_uuid(),
    chunk_id   uuid references chunks (id),
    embedding  vector(3072) not null,
    created_at timestamp with time zone not null default now()
);
grant all privileges on table embeddings to starter;
create index index_embeddings_chunk_id on embeddings (chunk_id);

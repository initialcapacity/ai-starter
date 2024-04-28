create table public.data (
    id uuid primary key default gen_random_uuid(),
    source varchar not null,
    content varchar not null,
    created_at timestamp with time zone not null default now()
);
grant all privileges on table data to starter;
create unique index index_data_source on public.data(source);

create table public.embeddings (
    id uuid primary key default gen_random_uuid(),
    embedding vector(3072) not null,
    created_at timestamp with time zone not null default now()
);
grant all privileges on table embeddings to starter;

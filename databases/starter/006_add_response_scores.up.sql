create table response_scores
(
    id                uuid primary key                              default gen_random_uuid(),
    query_response_id uuid references query_responses (id) not null,
    score             jsonb                                not null,
    score_version     int                                  not null,
    created_at        timestamp with time zone             not null default now()
);
grant all privileges on table query_responses to starter;

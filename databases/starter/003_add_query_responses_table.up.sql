create table query_responses
(
    id            uuid primary key                  default gen_random_uuid(),
    system_prompt varchar                  not null,
    user_query    varchar                  not null,
    source        varchar                  not null,
    response      varchar                  not null,
    model         varchar                  not null,
    temperature   float                    not null,
    created_at    timestamp with time zone not null default now()
);
grant all privileges on table query_responses to starter;

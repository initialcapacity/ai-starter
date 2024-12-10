alter table query_responses add column embeddings_model varchar not null default 'text-embedding-3-large';
alter table query_responses rename column model to chat_model;

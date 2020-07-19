-- Initialization SQL code for Undertakers table in PostgreSQL storage

CREATE TABLE IF NOT EXISTS __undertaker
(
    function      varchar primary key,
    first_seen_at timestamp null
);

CREATE INDEX never_seen_functions ON __undertaker (first_seen_at) WHERE first_seen_at IS NULL;

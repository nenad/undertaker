-- Initialization SQL code for Undertaker tests' table in PostgreSQL storage

CREATE SCHEMA test;

CREATE TABLE IF NOT EXISTS test.__undertaker_test
(
    function      varchar primary key,
    first_seen_at timestamp null
);

CREATE INDEX never_seen_functions ON test.__undertaker_test (first_seen_at) WHERE first_seen_at IS NULL;

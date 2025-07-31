-- +goose Up
CREATE TABLE urls (
    id         bigserial primary key,
    slug       text unique,
    long_url   text,
    ttl        timestamptz,
    created_at timestamptz default now()
);
CREATE TABLE clicks (
    id         bigserial primary key,
    url_id     bigint references urls(id),
    ua         text,
    ip         inet,
    ts         timestamptz default now()
);
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
DROP TABLE clicks;
DROP TABLE urls;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

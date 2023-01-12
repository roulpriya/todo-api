CREATE TABLE todo
(
    id         bigserial PRIMARY KEY NOT NULL,
    title      VARCHAR(255)          NOT NULL,
    content    TEXT,
    completed  BOOLEAN               NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP             NOT NULL,
    updated_at TIMESTAMP             NOT NULL
)
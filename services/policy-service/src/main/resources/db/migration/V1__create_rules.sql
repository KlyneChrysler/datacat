-- Rules table, one row per enforcement rule.
CREATE TABLE rules (
    id                    UUID PRIMARY KEY,
    name                  TEXT        NOT NULL,
    applies_to            TEXT        NOT NULL,
    action                TEXT        NOT NULL,
    rate_limit_per_minute INT         NOT NULL,
    created_at            TIMESTAMPTZ NOT NULL
);

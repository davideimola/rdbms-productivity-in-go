-- migrate:up

CREATE TABLE users
(
    id         SERIAL PRIMARY KEY,
    name       VARCHAR   NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- migrate:down

DROP TABLE users;

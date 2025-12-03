-- migrate:up
CREATE TABLE users
(
    id         bigserial PRIMARY KEY,
    username   varchar(25) NOT NULL UNIQUE,
    password   varchar(30) NOT NULL,
    created_at timestamp,
    updated_at timestamp
);

-- migrate:down
DROP TABLE IF EXISTS users
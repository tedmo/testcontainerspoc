--liquibase formatted sq

--changeset user:1
CREATE TABLE users (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name VARCHAR NOT NULL DEFAULT ''
);

--rollback DROP TABLE users;
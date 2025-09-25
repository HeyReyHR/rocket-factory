-- +goose Up
CREATE TABLE users (
    uuid UUID PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE user_infos (
    user_uuid UUID PRIMARY KEY REFERENCES users(uuid) ON DELETE CASCADE,
    login TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL
);

CREATE TABLE notification_methods (
    user_infos_uuid UUID PRIMARY KEY NOT NULL REFERENCES user_infos(user_uuid) ON DELETE CASCADE,
    provider_name TEXT NOT NULL,
    target TEXT NOT NULL
);

-- +goose Down

DROP TABLE IF EXISTS notification_methods;
DROP TABLE IF EXISTS user_infos;
DROP TABLE IF EXISTS users;
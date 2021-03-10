CREATE TABLE users
(
    id            serial       not null unique,
    name          varchar(255) not null,
    email         varchar(255) not null unique,
    password_hash varchar(255) not null,
    registered_at timestamp    not null,
    last_visit_at timestamp    not null
);

CREATE TABLE users_sessions
(
    id            serial       not null unique,
    user_id       int REFERENCES users (id) on DELETE CASCADE,
    refresh_token varchar(255) not null,
    expires_at    timestamp    not null
);

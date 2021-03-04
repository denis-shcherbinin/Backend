CREATE TABLE users
(
    id            serial       not null unique,
    name          varchar(255) not null,
    email         varchar(255) not null unique,
    password_hash varchar(255) not null,
    registered_at varchar(255) not null,
    last_visit_at varchar(255) not null
);

CREATE TABLE users
(
    id            serial       not null unique,
    name          varchar(255) not null,
    login         varchar(255) not null unique,
    password_hash varchar(255) not null
);

CREATE TABLE "skills"
(
    "id"   SERIAL PRIMARY KEY,
    "name" varchar(256) not null unique
);

CREATE TABLE "users"
(
    "id"            SERIAL PRIMARY KEY,
    "fullname"      varchar(256) not null,
    "email"         varchar(256) not null,
    "password_hash" varchar(256) not null,
    "birth_date"    timestamp    not null,
    "registered_at" timestamp    not null,
    "in_search"     boolean,
    "skills_id"     int ARRAY
);

CREATE TABLE "users_sessions"
(
    "id"            SERIAL PRIMARY KEY,
    "user_id"       int REFERENCES "users" (id) on DELETE CASCADE,
    "refresh_token" varchar(256) not null,
    "expires_at"    timestamp    not null
);

CREATE TABLE "resumes"
(
    "id"                SERIAL PRIMARY KEY,
    "position"          varchar(256) not null,
    "user_id"           int,
    "work_experience"   timestamp    not null,
    "short_description" varchar(256) not null,
    "full_description"  varchar(512) not null,
    "requirements"      varchar(256) not null,
    "skills_id"         int ARRAY
);

CREATE TABLE "vacancies"
(
    "id"                SERIAL PRIMARY KEY,
    "position"          varchar(256) not null,
    "short_description" varchar(256) not null,
    "full_description"  varchar(512) not null,
    "requirements"      varchar(256) not null,
    "advantages"        varchar(256) not null,
    "skills_id"         int ARRAY
);

CREATE TABLE "companies"
(
    "id"           SERIAL PRIMARY KEY,
    "name"         varchar(256) not null,
    "description"  varchar(256) not null,
    "address"      varchar(256) not null,
    "vacancies_id" int ARRAY
);

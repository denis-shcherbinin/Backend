CREATE TABLE "skills"
(
    "id"   SERIAL PRIMARY KEY,
    "name" varchar(256) unique
);

CREATE TABLE "spheres"
(
    "id"   SERIAL PRIMARY KEY,
    "name" varchar(256) unique
);

CREATE TABLE "spheres_skills"
(
    "id"        SERIAL PRIMARY KEY,
    "sphere_id" int REFERENCES "spheres" ("id") on DELETE CASCADE,
    "skill_id"  int REFERENCES "skills" ("id") on DELETE CASCADE
);

/*Resumes relations*/
CREATE TABLE "resumes"
(
    "id"                SERIAL PRIMARY KEY,
    "position"          varchar(256),
    "work_experience"   timestamp,
    "short_description" varchar(256),
    "full_description"  varchar(512),
    "requirements"      varchar(256)
);

CREATE TABLE "resumes_skills"
(
    "id"        SERIAL PRIMARY KEY,
    "resume_id" int REFERENCES "resumes" ("id") on DELETE CASCADE,
    "skill_id"  int REFERENCES "skills" ("id") on DELETE CASCADE
);

CREATE TABLE "resumes_spheres"
(
    "id"        SERIAL PRIMARY KEY,
    "sphere_id" int REFERENCES "spheres" ("id") on DELETE CASCADE,
    "resume_id" int REFERENCES "resumes" ("id") on DELETE CASCADE
);
/**/

/*Users relations*/
CREATE TABLE "users"
(
    "id"            SERIAL PRIMARY KEY,
    "fullname"      varchar(256),
    "email"         varchar(256) unique,
    "password_hash" varchar(256),
    "birth_date"    varchar(256),
    "registered_at" timestamp,
    "in_search"     boolean
);

CREATE TABLE "users_sessions"
(
    "id"            SERIAL PRIMARY KEY,
    "user_id"       int REFERENCES "users" ("id") on DELETE CASCADE,
    "refresh_token" varchar(256),
    "expires_at"    timestamp
);

CREATE TABLE "users_spheres"
(
    "id"        SERIAL PRIMARY KEY,
    "user_id"   int REFERENCES "users" ("id") on DELETE CASCADE,
    "sphere_id" int REFERENCES "spheres" ("id") on DELETE CASCADE
);

CREATE TABLE "users_skills"
(
    "id"       SERIAL PRIMARY KEY,
    "user_id"  int REFERENCES "users" ("id") on DELETE CASCADE,
    "skill_id" int REFERENCES "skills" ("id") on DELETE CASCADE
);

CREATE TABLE "users_resumes"
(
    "id"        SERIAL PRIMARY KEY,
    "user_id"   int REFERENCES "users" ("id") on DELETE CASCADE,
    "resume_id" int REFERENCES "resumes" ("id") on DELETE CASCADE
);
/**/

/*Vacancies relations*/
CREATE TABLE "vacancies"
(
    "id"                SERIAL PRIMARY KEY,
    "position"          varchar(256),
    "short_description" varchar(256),
    "full_description"  varchar(512),
    "requirements"      varchar(256),
    "advantages"        varchar(256)
);

CREATE TABLE "vacancies_skills"
(
    "id"         SERIAL PRIMARY KEY,
    "vacancy_id" int REFERENCES "vacancies" ("id") on DELETE CASCADE,
    "skill_id"   int REFERENCES "skills" ("id") on DELETE CASCADE
);

CREATE TABLE "vacancies_spheres"
(
    "id"         SERIAL PRIMARY KEY,
    "sphere_id"  int REFERENCES "spheres" ("id") on DELETE CASCADE,
    "vacancy_id" int REFERENCES "vacancies" ("id") on DELETE CASCADE
);
/**/

/*Companies relations*/
CREATE TABLE "companies"
(
    "id"          SERIAL PRIMARY KEY,
    "name"        varchar(256),
    "description" varchar(256),
    "address"     varchar(256)
);

CREATE TABLE "companies_vacancies"
(
    "id"         SERIAL PRIMARY KEY,
    "company_id" int REFERENCES "companies" ("id") on DELETE CASCADE,
    "vacancy_id" int REFERENCES "vacancies" ("id") on DELETE CASCADE
);
/**/

CREATE TABLE "skills"
(
    "id"   SERIAL PRIMARY KEY,
    "name" varchar(256) unique
);

CREATE TABLE "profiles"
(
    "id"          SERIAL PRIMARY KEY,
    "comment"     varchar(256),
    "experience"  varchar(256),
    "skill_level" varchar(256),
    "min_salary"  varchar(256),
    "max_salary"  varchar(256),
    "about"       varchar(512)
);

CREATE TABLE "responsibilities"
(
    "id"          SERIAL PRIMARY KEY,
    "name"        varchar(256),
    "description" varchar(512)
);

CREATE TABLE "requirements"
(
    "id"          SERIAL PRIMARY KEY,
    "name"        varchar(256),
    "description" varchar(512)
);

CREATE TABLE "conditions"
(
    "id"          SERIAL PRIMARY KEY,
    "name"        varchar(256),
    "description" varchar(512)
);

/*Spheres relations*/
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
/**/

/*Jobs relations*/
CREATE TABLE "jobs"
(
    "id"               SERIAL PRIMARY KEY,
    "company_name"     varchar(256),
    "position"         varchar(256),
    /*Work period*/
    "from"             varchar(256),
    "to"               varchar(256),
    /**/
    "responsibilities" varchar(512)
);

CREATE TABLE "jobs_skills"
(
    "id"       SERIAL PRIMARY KEY,
    "job_id"   int REFERENCES "jobs" ("id") on DELETE CASCADE,
    "skill_id" int REFERENCES "skills" ("id") on DELETE CASCADE
);
/**/

/*Vacancies relations*/
CREATE TABLE "vacancies"
(
    "id"           SERIAL PRIMARY KEY,
    "position"     varchar(256),
    "description"  varchar(512),
    "is_full_time" boolean,
    "min_salary"   varchar(256),
    "max_salary"   varchar(256),
    "skill_level"  varchar(256)
);

CREATE TABLE "vacancies_responsibilities"
(
    "id"                SERIAL PRIMARY KEY,
    "vacancy_id"        int REFERENCES "vacancies" ("id") ON DELETE CASCADE,
    "responsibility_id" int REFERENCES "responsibilities" ("id") ON DELETE CASCADE
);

CREATE TABLE "vacancies_requirements"
(
    "id"             SERIAL PRIMARY KEY,
    "vacancy_id"     int REFERENCES "vacancies" ("id") ON DELETE CASCADE,
    "requirement_id" int REFERENCES "requirements" ("id") ON DELETE CASCADE
);

CREATE TABLE "vacancies_conditions"
(
    "id"           SERIAL PRIMARY KEY,
    "vacancy_id"   int REFERENCES "vacancies" ("id") ON DELETE CASCADE,
    "condition_id" int REFERENCES "conditions" ("id") ON DELETE CASCADE
);

CREATE TABLE "vacancies_skills"
(
    "id"         SERIAL PRIMARY KEY,
    "vacancy_id" int REFERENCES "vacancies" ("id") on DELETE CASCADE,
    "skill_id"   int REFERENCES "skills" ("id") on DELETE CASCADE
);
/**/

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

/*Companies relations*/
CREATE TABLE "companies"
(
    "id"                SERIAL PRIMARY KEY,
    "name"              varchar(256),
    "location"          varchar(256),
    "short_description" varchar(256),
    "full_description"  varchar(512),
    "image_url"         varchar(256)
);

CREATE TABLE "companies_vacancies"
(
    "id"         SERIAL PRIMARY KEY,
    "company_id" int REFERENCES "companies" ("id") on DELETE CASCADE,
    "vacancy_id" int REFERENCES "vacancies" ("id") on DELETE CASCADE
);
/**/

/*Users relations*/
CREATE TABLE "users"
(
    "id"            SERIAL PRIMARY KEY,
    "first_name"    varchar(256),
    "last_name"     varchar(256),
    "birth_date"    varchar(256),
    "email"         varchar(256) unique,
    "password_hash" varchar(256),
    "in_search"     boolean,
    "registered_at" timestamp,
    "image_url"     varchar(256)
);

CREATE TABLE "users_sessions"
(
    "id"            SERIAL PRIMARY KEY,
    "user_id"       int REFERENCES "users" ("id") on DELETE CASCADE,
    "user_agent"    varchar(256) NOT NULL,
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

CREATE TABLE "users_profiles"
(
    "id"         SERIAL PRIMARY KEY,
    "user_id"    int REFERENCES "users" ("id") on DELETE CASCADE UNIQUE,
    "profile_id" int REFERENCES "profiles" ("id") on DELETE CASCADE UNIQUE
);

CREATE TABLE "users_jobs"
(
    "id"      SERIAL PRIMARY KEY,
    "user_id" int REFERENCES "users" ("id") on DELETE CASCADE,
    "job_id"  int REFERENCES "jobs" ("id") on DELETE CASCADE
);

CREATE TABLE "users_resumes"
(
    "id"        SERIAL PRIMARY KEY,
    "user_id"   int REFERENCES "users" ("id") on DELETE CASCADE,
    "resume_id" int REFERENCES "resumes" ("id") on DELETE CASCADE
);

CREATE TABLE "users_companies"
(
    "id"         SERIAL PRIMARY KEY,
    "user_id"    int REFERENCES "users" ("id") on DELETE CASCADE,
    "company_id" int REFERENCES "companies" ("id") ON DELETE CASCADE
);
/**/

-- Down migration: drop authority tables and schema
SET search_path TO authority;

DROP TABLE IF EXISTS account_roles;
DROP TABLE IF EXISTS roles;
DROP TABLE IF EXISTS account_identities;
DROP TABLE IF EXISTS accounts;

DROP SCHEMA IF EXISTS authority;


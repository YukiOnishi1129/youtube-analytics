-- Down migration: drop authority tables and schema
DROP TABLE IF EXISTS authority.account_roles;
DROP TABLE IF EXISTS authority.roles;
DROP TABLE IF EXISTS authority.account_identities;
DROP TABLE IF EXISTS authority.accounts;

DROP SCHEMA IF EXISTS authority;

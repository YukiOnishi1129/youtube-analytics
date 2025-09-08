-- Up migration: create authority schema and tables
CREATE SCHEMA IF NOT EXISTS authority;
SET search_path TO authority;

CREATE TABLE IF NOT EXISTS accounts (
  id              uuid PRIMARY KEY,
  email           text NOT NULL,
  email_verified  boolean DEFAULT false,
  display_name    text,
  photo_url       text,
  is_active       boolean DEFAULT true,
  last_login_at   timestamptz,
  created_at      timestamptz DEFAULT now(),
  updated_at      timestamptz,
  deleted_at      timestamptz
);
CREATE UNIQUE INDEX IF NOT EXISTS accounts_email_unique ON accounts(email) WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS account_identities (
  id           uuid PRIMARY KEY,
  account_id   uuid NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
  provider     text NOT NULL,
  provider_uid text,
  created_at   timestamptz DEFAULT now(),
  UNIQUE (account_id, provider)
);

CREATE TABLE IF NOT EXISTS roles (
  id          uuid PRIMARY KEY,
  name        text NOT NULL UNIQUE,
  description text,
  created_at  timestamptz DEFAULT now(),
  updated_at  timestamptz,
  deleted_at  timestamptz
);

CREATE TABLE IF NOT EXISTS account_roles (
  account_id  uuid NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
  role_id     uuid NOT NULL REFERENCES roles(id) ON DELETE RESTRICT,
  created_at  timestamptz DEFAULT now(),
  PRIMARY KEY (account_id, role_id)
);


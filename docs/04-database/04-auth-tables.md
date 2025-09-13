# Authority (Authentication & Authorization) Tables

Schema separation policy
- Each microservice owns a dedicated PostgreSQL schema. For authority-service, use schema `authority`.
- Application connections set `search_path` to the service schema.
- Do not perform cross-schema joins from the service; collaborate via gRPC instead.

## Accounts (Assuming Identity Platform Integration)

```sql
CREATE SCHEMA IF NOT EXISTS authority;
SET search_path TO authority;

CREATE TABLE accounts (
  id              uuid PRIMARY KEY,                     -- v7
  email           text NOT NULL,
  email_verified  boolean DEFAULT false,
  display_name    text,
  photo_url       text,                                 -- Store images as URLs
  is_active       boolean DEFAULT true,
  last_login_at   timestamptz,
  created_at      timestamptz DEFAULT now(),
  updated_at      timestamptz,
  deleted_at      timestamptz
);
CREATE UNIQUE INDEX accounts_email_unique ON accounts(email) WHERE deleted_at IS NULL;

-- Multiple Providers
CREATE TABLE account_identities (
  id           uuid PRIMARY KEY,                        -- v7
  account_id   uuid NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
  provider     text NOT NULL,                           -- 'google','password','github',…
  provider_uid text,                                    -- Google's sub, etc.
  created_at   timestamptz DEFAULT now(),
  UNIQUE (account_id, provider)
);

-- Roles (Optional: admin/editor/viewer)
CREATE TABLE roles (
  id          uuid PRIMARY KEY,                         -- v7
  name        text NOT NULL UNIQUE,
  description text,
  created_at  timestamptz DEFAULT now(),
  updated_at  timestamptz,
  deleted_at  timestamptz
);

CREATE TABLE account_roles (
  account_id  uuid NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
  role_id     uuid NOT NULL REFERENCES roles(id) ON DELETE RESTRICT,
  created_at  timestamptz DEFAULT now(),
  PRIMARY KEY (account_id, role_id)
);
```

## Operational Notes

- **UUID v7**: Generate uuid.NewV7() in Go → INSERT into uuid column (no DB-side configuration needed)
- **Soft Delete**: Rows with non-NULL deleted_at are hidden (use WHERE deleted_at IS NULL where necessary)
- **Low Sample Exclusion**: Update exclude_from_ranking based on views_count / subscription_count thresholds per CP
- **Main Search Conditions**:
  - Ranking: published_at ∈ [from,to) & checkpoint_hour=X & metric descending order
  - Details: Fixed set of checkpoint_hour in video_snapshots (0/3/6/...)
  - Keywords: enabled=true AND deleted_at IS NULL
- **WebSub Subscription Info**: Initially channels.subscribed is sufficient (add separate table for lease details later if needed)

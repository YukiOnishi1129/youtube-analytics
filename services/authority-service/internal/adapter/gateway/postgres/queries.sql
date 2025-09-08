-- name: GetAccountByID :one
SELECT id, email, email_verified, display_name, photo_url, is_active, last_login_at
FROM authority.accounts
WHERE id = $1 AND deleted_at IS NULL;

-- name: GetAccountByEmail :one
SELECT id, email, email_verified, display_name, photo_url, is_active, last_login_at
FROM authority.accounts
WHERE email = $1 AND deleted_at IS NULL;

-- name: UpsertAccount :exec
INSERT INTO authority.accounts (id, email, email_verified, display_name, photo_url, is_active, last_login_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, now())
ON CONFLICT (id) DO UPDATE SET
  email = EXCLUDED.email,
  email_verified = EXCLUDED.email_verified,
  display_name = EXCLUDED.display_name,
  photo_url = EXCLUDED.photo_url,
  is_active = EXCLUDED.is_active,
  last_login_at = EXCLUDED.last_login_at,
  updated_at = now();

-- name: ListIdentitiesByAccount :many
SELECT provider, provider_uid FROM authority.account_identities WHERE account_id = $1 ORDER BY created_at ASC;

-- name: UpsertIdentity :exec
INSERT INTO authority.account_identities (id, account_id, provider, provider_uid)
VALUES ($1, $2, $3, $4)
ON CONFLICT (account_id, provider) DO UPDATE SET provider_uid = EXCLUDED.provider_uid;

-- name: DeleteIdentity :exec
DELETE FROM authority.account_identities WHERE account_id = $1 AND provider = $2;

-- name: FindAccountByProvider :one
SELECT a.id, a.email, a.email_verified, a.display_name, a.photo_url, a.is_active, a.last_login_at
FROM authority.accounts a
JOIN authority.account_identities i ON i.account_id = a.id
WHERE i.provider = $1 AND i.provider_uid = $2 AND a.deleted_at IS NULL;

-- name: ListRolesByAccount :many
SELECT r.name FROM authority.roles r
JOIN authority.account_roles ar ON ar.role_id = r.id
WHERE ar.account_id = $1;

-- name: AssignRole :exec
INSERT INTO authority.account_roles (account_id, role_id)
SELECT $1, r.id FROM authority.roles r WHERE r.name = $2
ON CONFLICT DO NOTHING;

-- name: RevokeRole :exec
DELETE FROM authority.account_roles ar USING authority.roles r
WHERE ar.account_id = $1 AND ar.role_id = r.id AND r.name = $2;



-- name: CreateSession :one
INSERT INTO sessions (
  id,
  email,
  refresh_token,
  user_agent,
  client_ip,
  is_blocked,
  expires_at
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: GetSession :one
SELECT * FROM sessions
WHERE id = $1 LIMIT 1;

-- name: GetSessionsByEmailAndIp :one
SELECT * FROM sessions
WHERE email = $1 AND user_agent = $2 AND client_ip = $3
ORDER BY created_at DESC
LIMIT 1;

-- name: UpdateSession :one
UPDATE sessions
SET refresh_token = $2
WHERE id = $1
RETURNING *;

-- name: RefreshSession :one
UPDATE sessions
SET refresh_token = $2,
    expires_at = $3
WHERE id = $1
RETURNING *;

-- name: DeleteSession :one
DELETE FROM sessions
WHERE id = $1
RETURNING *;

-- name: BlockSession :one
UPDATE sessions
SET is_blocked = true
WHERE id = $1
RETURNING *;
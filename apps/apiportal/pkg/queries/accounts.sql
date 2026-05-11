

-- name: CreateAccounts :one
INSERT INTO "accounts" (
    user_id, balance, currency
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: UpdateAccountBalance :one
UPDATE "accounts"
SET
  balance = balance + sqlc.arg(amount)
WHERE
  id = sqlc.arg(id)
RETURNING *;

-- name: GetAccountByID :one
SELECT * FROM "accounts" WHERE id = $1
FOR NO KEY UPDATE;

-- name: GetAccountCount :one
SELECT COUNT(*) FROM "accounts";  
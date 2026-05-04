

-- name: CreateAccounts :one
INSERT INTO "accounts" (
    user_id, balance, currency
) VALUES (
    $1, $2, $3
)
RETURNING *;


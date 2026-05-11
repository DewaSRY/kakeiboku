


-- name: CreateUser :one
INSERT INTO "user" (
    first_name, last_name, email, password_hash
) VALUES (  
    $1, $2, $3, $4
)
RETURNING *;
    
-- name: GetUserByEmail :one
SELECT * FROM "user" WHERE email = $1;

-- name: GetUserByID :one
SELECT * FROM "user" WHERE id = $1;

-- name: UpdateUserById :one
UPDATE "user"
SET
  password_hash = COALESCE(sqlc.narg(password_hash), password_hash),
  first_name = COALESCE(sqlc.narg(first_name), first_name),
  last_name = COALESCE(sqlc.narg(last_name), last_name),
  email = COALESCE(sqlc.narg(email), email)
WHERE
  id = sqlc.arg(id)
RETURNING *;

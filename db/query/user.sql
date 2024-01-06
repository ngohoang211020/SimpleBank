-- name: CreateUser :one
INSERT INTO users (username,
                   hashed_password,
                   full_name,
                   email)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetUser :one
SELECT *
FROM users
WHERE username = $1
LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY username
LIMIT $1
    OFFSET $2;

-- name: UpdateUser :one
UPDATE users
SET
    hashed_password = COALESCE(sqlc.narg(hashed_password), hashed_password),
    password_changed_at = COALESCE(sqlc.narg(password_changed_at), password_changed_at),
    full_name = COALESCE(sqlc.narg(full_name), full_name),
    email = COALESCE(sqlc.narg(email), email),
    is_verified_email = COALESCE(sqlc.narg(is_verified_email), is_verified_email)
    WHERE
        username = sqlc.arg(username)
RETURNING *;

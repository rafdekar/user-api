-- name: CreateUser :one
INSERT INTO users (
                   first_name,
                   last_name,
                   nickname,
                   password,
                   email,
                   country
)
VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateUser :one
UPDATE users
SET first_name = $2,
    last_name = $3,
    nickname = $4,
    password = $5,
    email = $6,
    country = $7
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

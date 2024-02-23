-- name: FindUserByID :one
SELECT * FROM users
WHERE id = $1;

-- name: FindUsers :many
SELECT * FROM users;

-- name: CreateUser :one
INSERT INTO users (name)
VALUES ($1)
RETURNING *;
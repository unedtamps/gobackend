-- name: GetUser :one
SELECT *
FROM "user"
WHERE id = $1
LIMIT 1;

-- name: GetUserByEmail :one
SELECT *
FROM "user"
WHERE email = $1
LIMIT 1;

-- name: ListUsers :many
SELECT *
FROM "user"
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CreateUser :one
INSERT INTO "user" (id, email, password)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateUser :one
UPDATE "user"
SET email = $2,
    password = $3
WHERE id = $1
RETURNING *;

-- name: UpdateUserStatus :one
UPDATE "user"
SET status = $2
WHERE id = $1
RETURNING *;

-- name: SoftDeleteUser :one
UPDATE "user"
SET status = 'DELETED'
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM "user"
WHERE id = $1;

-- name: CountUsers :one
SELECT COUNT(*)
FROM "user";

-- name: CountUsersByStatus :one
SELECT COUNT(*)
FROM "user"
WHERE status = $1;

-- name: SearchUsersByEmail :many
SELECT *
FROM "user"
WHERE email ILIKE '%' || $1 || '%'
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

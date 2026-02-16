-- name: GetTodo :one
SELECT *
FROM "todo"
WHERE id = $1
LIMIT 1;

-- name: GetTodoByUser :many
SELECT *
FROM "todo"
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListTodos :many
SELECT *
FROM "todo"
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: ListTodosByUser :many
SELECT *
FROM "todo"
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: CreateTodo :one
INSERT INTO "todo" (id, user_id, title, description)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateTodo :one
UPDATE "todo"
SET title = $2,
    description = $3
WHERE id = $1
RETURNING *;

-- name: UpdateTodoTitle :one
UPDATE "todo"
SET title = $2
WHERE id = $1
RETURNING *;

-- name: UpdateTodoDescription :one
UPDATE "todo"
SET description = $2
WHERE id = $1
RETURNING *;

-- name: DeleteTodo :exec
DELETE FROM "todo"
WHERE id = $1;

-- name: DeleteTodosByUser :exec
DELETE FROM "todo"
WHERE user_id = $1;

-- name: CountTodos :one
SELECT COUNT(*)
FROM "todo";

-- name: CountTodosByUser :one
SELECT COUNT(*)
FROM "todo"
WHERE user_id = $1;

-- name: SearchTodosByTitle :many
SELECT *
FROM "todo"
WHERE title ILIKE '%' || $1 || '%'
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: SearchTodosByUserAndTitle :many
SELECT *
FROM "todo"
WHERE user_id = $1
  AND title ILIKE '%' || $2 || '%'
ORDER BY created_at DESC
LIMIT $3 OFFSET $4;

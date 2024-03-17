-- name: QueryTodo :many
SELECT * from "todolist";

-- name: MakeTodo :one
INSERT INTO "todolist" ("title", "description") VALUES ($1 , $2) RETURNING id; 

-- name: QueryTodoById :one
SELECT id, title,description FROM "todolist" WHERE id = $1;

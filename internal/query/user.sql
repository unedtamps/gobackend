-- name: CreateUser :one
INSERT INTO "user" ("email", "password") VALUES ($1 , $2) RETURNING id;

-- name: GetUserByEmail :one
SELECT id,email,password FROM "user" WHERE "email" = $1;

-- name: GetUser :one
select *
from "user"
where id = $1
;

-- name: CreateUser :one
insert into "user" (id,email, password)
values ($1, $2, $3)
returning *;

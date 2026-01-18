-- name: GetProduct :one
select *
from product
where id = ?
;

-- name: CreateProduct :execresult
INSERT INTO product (id, name, description)
VALUES (?, ?, ?);

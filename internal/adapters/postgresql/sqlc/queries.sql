-- name: ListProducts :many
SELECT * FROM products;

-- name: ListProductsByID :one
SELECT * FROm products where id = $1;
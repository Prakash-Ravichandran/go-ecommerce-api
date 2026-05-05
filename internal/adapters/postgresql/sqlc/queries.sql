-- name: ListProducts :many
SELECT * FROM products;

-- name: ListProductsByID :one
SELECT * FROm products where id = $1;

-- name: CreateOrder :one
INSERT INTO orders (
    customer_id
) VALUES($1) RETURNING *;

-- name: CreateOrderItem :one
INSERT INTO order_items (order_id, product_id, quantity, price_cents) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: ListOrders :many
SELECT * FROM orders;

-- name: ListOrderById :one
SELECT * FROM orders where id = $1;
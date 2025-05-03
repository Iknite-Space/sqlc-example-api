-- -- name: CreateMessage :one
-- INSERT INTO message (thread, sender, content)
-- VALUES ($1, $2, $3)
-- RETURNING *;

-- -- name: GetMessageByID :one
-- SELECT * FROM message
-- WHERE id = $1;

-- -- name: GetMessagesByThread :many
-- SELECT * FROM message
-- WHERE thread = $1
-- ORDER BY created_at DESC;

-- -- name: CustomerOrders :many
-- INSERT INTO orders (reference, phone_number, amount, transaction_status, transaction_description)
-- VALUES ($1, $2, $3, $4, $5)
-- RETURNING *;

-- name: CreateCustomer :one
INSERT INTO customer (customer_name, phone_number, email)
VALUES ($1, $2, $3)
RETURNING *;

-- name: CreateOrder :one
INSERT INTO "order" (customer_id, product_name, price, order_status)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpadateOrderById :one
UPDATE "order"
SET order_status = $2
WHERE id =$1
RETURNING *;
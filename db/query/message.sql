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

-- -- name: DeleteMessage :exec
-- DELETE FROM message WHERE id = $1;

-- -- name: UpdateMessage :exec
-- UPDATE message 
-- SET content = $2
-- WHERE id = $1
-- RETURNING *;

-- -- name: CreateThread :one
-- INSERT INTO thread (title)
-- VALUES ($1)
-- RETURNING *;

-- -- name: DeleteAll :exec
-- DELETE FROM message;

-- name: CreateThread :one
INSERT INTO thread (title)
VALUES ($1)
RETURNING *;

-- name: CreateMessage :one
INSERT INTO message (content,thread_id)
VALUES ($1, $2)
RETURNING *;

-- name: GetMessageByID :one
SELECT * FROM message
WHERE id = $1;

-- name: GetMessagesByThread :many
SELECT * FROM message
WHERE thread_id = $1
ORDER BY created_at DESC;

-- name: DeleteMessageById :one
DELETE FROM message WHERE id = $1
RETURNING id;

-- name: DeleteMessageByThreadId :one
DELETE FROM message WHERE thread_id = $1
RETURNING thread_id;

-- name: UpdateMessage :exec
UPDATE message 
SET content = $2
WHERE id = $1
RETURNING *;

-- name: GetThreadById :one
SELECT * FROM thread WHERE id = $1;


-- name: CreateOrder :one
INSERT INTO orders (customer_name,amount,phone_number)
VALUES($1,$2,$3)
RETURNING *;

-- name: GetMessageByThreadPaginated :many
SELECT * FROM message
WHERE thread_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;
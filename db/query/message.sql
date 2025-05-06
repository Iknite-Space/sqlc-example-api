-- name: CreateThread :one
INSERT INTO thread (thread_id)
VALUES ($1)
RETURNING *;

-- name: CreateMessage :one
INSERT INTO message (thread, sender, content)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetMessageByID :one
SELECT * FROM message
WHERE id = $1;

-- name: GetMessagesByThread :many
SELECT * FROM message
WHERE thread = $1
ORDER BY created_at DESC;

-- name: DeleteMessageByID :exec
DELETE FROM message
WHERE id = $1;

-- name: UpdateMessageByID :one
UPDATE message 
SET sender = $1, content = $2 
WHERE id = $3
RETURNING *;

-- name: TransactionMessage :one
INSERT INTO transacions (reference, external_reference, status_id, amount, currency, operator, code, operator_reference, description_id, external_user, reason, phone_number)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
RETURNING *;

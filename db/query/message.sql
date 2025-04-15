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


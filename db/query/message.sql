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

-- name: DeleteMessage :exec
DELETE FROM message WHERE id = $1;

-- name: UpdateMessage :exec
UPDATE message 
SET content = $2
WHERE id = $1
RETURNING *;

-- name: DeleteAll :exec
DELETE FROM message;
-- name: CreateEditorConfig :one
INSERT INTO editor_config (name, config_data)
VALUES ($1, $2)
RETURNING *;

-- name: GetEditorConfig :one
SELECT * FROM editor_config
WHERE id = $1;

-- name: ListEditorConfigs :many
SELECT * FROM editor_config
ORDER BY created_at DESC;

-- name: UpdateEditorConfig :one
UPDATE editor_config 
SET name = $2, config_data = $3, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteEditorConfig :exec
DELETE FROM editor_config
WHERE id = $1;
-- name: CreateSecret :one
INSERT INTO secret (description, core_data, custom_fields, schema_version)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetSecret :one
SELECT * FROM secret
WHERE id = $1;

-- name: ListSecrets :many
SELECT * FROM secret
ORDER BY created_at DESC;

-- name: UpdateSecret :one
UPDATE secret 
SET description = $2, core_data = $3, custom_fields = $4, 
    schema_version = $5, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteSecret :exec
DELETE FROM secret
WHERE id = $1;

-- name: SearchSecretsByField :many
SELECT * FROM secret
WHERE core_data @> $1 OR custom_fields @> $1
ORDER BY created_at DESC;

-- name: SearchSecretsByDescription :many
SELECT * FROM secret
WHERE description ILIKE '%' || $1 || '%'
ORDER BY created_at DESC;
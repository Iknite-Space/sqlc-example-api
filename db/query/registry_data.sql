-- name: CreateRegistryData :one
INSERT INTO registry_data (subkey, value_name, core_data, custom_fields, schema_version)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetRegistryData :one
SELECT * FROM registry_data
WHERE id = $1;

-- name: ListRegistryData :many
SELECT * FROM registry_data
ORDER BY created_at DESC;

-- name: UpdateRegistryData :one
UPDATE registry_data 
SET subkey = $2, value_name = $3, core_data = $4, custom_fields = $5, 
    schema_version = $6, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteRegistryData :exec
DELETE FROM registry_data
WHERE id = $1;

-- name: SearchRegistryDataByField :many
SELECT * FROM registry_data
WHERE core_data @> $1 OR custom_fields @> $1
ORDER BY created_at DESC;

-- name: ListRegistryDataBySubkey :many
SELECT * FROM registry_data
WHERE subkey ILIKE '%' || $1 || '%'
ORDER BY created_at DESC;
-- name: CreateCacheData :one
INSERT INTO cache_data (technology, cache_type, core_data, custom_fields, schema_version)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetCacheData :one
SELECT * FROM cache_data
WHERE id = $1;

-- name: ListCacheData :many
SELECT * FROM cache_data
ORDER BY created_at DESC;

-- name: UpdateCacheData :one
UPDATE cache_data 
SET technology = $2, cache_type = $3, core_data = $4, custom_fields = $5, 
    schema_version = $6, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteCacheData :exec
DELETE FROM cache_data
WHERE id = $1;

-- name: SearchCacheDataByField :many
SELECT * FROM cache_data
WHERE core_data @> $1 OR custom_fields @> $1
ORDER BY created_at DESC;

-- name: ListCacheDataByTechnology :many
SELECT * FROM cache_data
WHERE technology = $1
ORDER BY created_at DESC;
-- name: CreatePlistData :one
INSERT INTO plist_data (key, core_data, custom_fields, schema_version)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetPlistData :one
SELECT * FROM plist_data
WHERE id = $1;

-- name: ListPlistData :many
SELECT * FROM plist_data
ORDER BY created_at DESC;

-- name: UpdatePlistData :one
UPDATE plist_data 
SET key = $2, core_data = $3, custom_fields = $4, 
    schema_version = $5, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeletePlistData :exec
DELETE FROM plist_data
WHERE id = $1;

-- name: SearchPlistDataByField :many
SELECT * FROM plist_data
WHERE core_data @> $1 OR custom_fields @> $1
ORDER BY created_at DESC;

-- name: ListPlistDataByKey :many
SELECT * FROM plist_data
WHERE key ILIKE '%' || $1 || '%'
ORDER BY created_at DESC;
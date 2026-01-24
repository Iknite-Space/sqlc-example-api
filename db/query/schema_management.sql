-- name: CreateTableSchema :one
INSERT INTO table_schema (table_name, schema_version, field_definitions, is_active)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetTableSchema :one
SELECT * FROM table_schema
WHERE table_name = $1 AND is_active = true
ORDER BY schema_version DESC
LIMIT 1;

-- name: GetTableSchemaByVersion :one
SELECT * FROM table_schema
WHERE table_name = $1 AND schema_version = $2;

-- name: ListTableSchemas :many
SELECT * FROM table_schema
WHERE table_name = $1
ORDER BY schema_version DESC;

-- name: UpdateTableSchema :one
UPDATE table_schema 
SET field_definitions = $3, is_active = $4
WHERE table_name = $1 AND schema_version = $2
RETURNING *;

-- name: DeactivateOldSchemas :exec
UPDATE table_schema 
SET is_active = false
WHERE table_name = $1 AND schema_version < $2;
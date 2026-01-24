-- name: CreateLogEntry :one
INSERT INTO log_entry (level, message, core_data, custom_fields, schema_version, timestamp)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetLogEntry :one
SELECT * FROM log_entry
WHERE id = $1;

-- name: ListLogEntries :many
SELECT * FROM log_entry
ORDER BY timestamp DESC;

-- name: UpdateLogEntry :one
UPDATE log_entry 
SET level = $2, message = $3, core_data = $4, custom_fields = $5, 
    schema_version = $6, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteLogEntry :exec
DELETE FROM log_entry
WHERE id = $1;

-- name: SearchLogEntriesByField :many
SELECT * FROM log_entry
WHERE core_data @> $1 OR custom_fields @> $1
ORDER BY timestamp DESC;

-- name: ListLogEntriesByLevel :many
SELECT * FROM log_entry
WHERE level = $1
ORDER BY timestamp DESC;

-- name: ListLogEntriesByTimeRange :many
SELECT * FROM log_entry
WHERE timestamp BETWEEN $1 AND $2
ORDER BY timestamp DESC;
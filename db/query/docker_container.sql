-- name: CreateDockerContainer :one
INSERT INTO docker_container (name, status, core_data, custom_fields, schema_version)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetDockerContainer :one
SELECT * FROM docker_container
WHERE id = $1;

-- name: ListDockerContainers :many
SELECT * FROM docker_container
ORDER BY created_at DESC;

-- name: UpdateDockerContainer :one
UPDATE docker_container 
SET name = $2, status = $3, core_data = $4, custom_fields = $5, 
    schema_version = $6, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteDockerContainer :exec
DELETE FROM docker_container
WHERE id = $1;

-- name: SearchDockerContainersByField :many
SELECT * FROM docker_container
WHERE core_data @> $1 OR custom_fields @> $1
ORDER BY created_at DESC;

-- name: UpdateDockerContainerField :one
UPDATE docker_container 
SET core_data = jsonb_set(core_data, $2, $3), updated_at = now()
WHERE id = $1
RETURNING *;

-- name: AddCustomFieldToDockerContainer :one
UPDATE docker_container 
SET custom_fields = jsonb_set(custom_fields, $2, $3), updated_at = now()
WHERE id = $1
RETURNING *;

-- name: RemoveFieldFromDockerContainer :one
UPDATE docker_container 
SET core_data = core_data - $2, custom_fields = custom_fields - $2, updated_at = now()
WHERE id = $1
RETURNING *;
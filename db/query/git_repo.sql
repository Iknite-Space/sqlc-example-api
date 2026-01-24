-- name: CreateGitRepo :one
INSERT INTO git_repo (name, core_data, custom_fields, schema_version)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetGitRepo :one
SELECT * FROM git_repo
WHERE id = $1;

-- name: ListGitRepos :many
SELECT * FROM git_repo
ORDER BY created_at DESC;

-- name: UpdateGitRepo :one
UPDATE git_repo 
SET name = $2, core_data = $3, custom_fields = $4, schema_version = $5, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteGitRepo :exec
DELETE FROM git_repo
WHERE id = $1;

-- name: SearchGitReposByField :many
SELECT * FROM git_repo
WHERE core_data @> $1 OR custom_fields @> $1
ORDER BY created_at DESC;

-- name: UpdateGitRepoField :one
UPDATE git_repo 
SET core_data = jsonb_set(core_data, $2, $3), updated_at = now()
WHERE id = $1
RETURNING *;

-- name: AddCustomFieldToGitRepo :one
UPDATE git_repo 
SET custom_fields = jsonb_set(custom_fields, $2, $3), updated_at = now()
WHERE id = $1
RETURNING *;
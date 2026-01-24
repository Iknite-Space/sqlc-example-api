-- name: CreateWorkflow :one
INSERT INTO workflow (name, version, export_date, node_count, edge_count, port_count, workflow_data)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetWorkflow :one
SELECT * FROM workflow
WHERE id = $1;

-- name: ListWorkflows :many
SELECT * FROM workflow
ORDER BY created_at DESC;

-- name: UpdateWorkflow :one
UPDATE workflow 
SET name = $2, version = $3, export_date = $4, node_count = $5, edge_count = $6, 
    port_count = $7, workflow_data = $8, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteWorkflow :exec
DELETE FROM workflow
WHERE id = $1;

-- name: CreateWorkflowNode :one
INSERT INTO workflow_node (workflow_id, node_id, label, group_type, title, x_position, y_position, node_data)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetWorkflowNodes :many
SELECT * FROM workflow_node
WHERE workflow_id = $1
ORDER BY created_at;

-- name: CreateWorkflowEdge :one
INSERT INTO workflow_edge (workflow_id, from_node, to_node, label, edge_data)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetWorkflowEdges :many
SELECT * FROM workflow_edge
WHERE workflow_id = $1
ORDER BY created_at;
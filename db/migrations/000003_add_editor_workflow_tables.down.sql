-- Drop editor and workflow tables in reverse order
DROP TABLE IF EXISTS "workflow_edge";
DROP TABLE IF EXISTS "workflow_node";
DROP TABLE IF EXISTS "workflow";
DROP TABLE IF EXISTS "editor_config";
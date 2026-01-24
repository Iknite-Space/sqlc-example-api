-- Drop dashboard tables in reverse order
DROP INDEX IF EXISTS idx_plist_data_core_data;
DROP INDEX IF EXISTS idx_registry_data_core_data;
DROP INDEX IF EXISTS idx_secret_core_data;
DROP INDEX IF EXISTS idx_log_entry_core_data;
DROP INDEX IF EXISTS idx_cache_data_core_data;
DROP INDEX IF EXISTS idx_git_repo_custom_fields;
DROP INDEX IF EXISTS idx_git_repo_core_data;
DROP INDEX IF EXISTS idx_docker_container_custom_fields;
DROP INDEX IF EXISTS idx_docker_container_core_data;

DROP TABLE IF EXISTS "table_schema";
DROP TABLE IF EXISTS "plist_data";
DROP TABLE IF EXISTS "registry_data";
DROP TABLE IF EXISTS "secret";
DROP TABLE IF EXISTS "log_entry";
DROP TABLE IF EXISTS "cache_data";
DROP TABLE IF EXISTS "git_repo";
DROP TABLE IF EXISTS "docker_container";
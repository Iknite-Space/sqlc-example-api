

-- Docker containers table with flexible schema
CREATE TABLE "docker_container" (
  "id" VARCHAR(36) PRIMARY KEY DEFAULT gen_random_uuid()::varchar(36),
  "name" VARCHAR(100) NOT NULL,
  "status" VARCHAR(20) NOT NULL,
  "core_data" JSONB NOT NULL DEFAULT '{}',
  "custom_fields" JSONB DEFAULT '{}', -- JSONB is supported by postgress allowing you to store json objects as fields
  "schema_version" INTEGER DEFAULT 1,
  "created_at" TIMESTAMP DEFAULT now(),
  "updated_at" TIMESTAMP DEFAULT now()
);

-- Git repositories table with flexible schema
CREATE TABLE "git_repo" (
  "id" VARCHAR(36) PRIMARY KEY DEFAULT gen_random_uuid()::varchar(36),
  "name" VARCHAR(200) NOT NULL,
  "core_data" JSONB NOT NULL DEFAULT '{}',
  "custom_fields" JSONB DEFAULT '{}',
  "schema_version" INTEGER DEFAULT 1,
  "created_at" TIMESTAMP DEFAULT now(),
  "updated_at" TIMESTAMP DEFAULT now()
);


CREATE TABLE "cache_data" (
  "id" VARCHAR(36) PRIMARY KEY DEFAULT gen_random_uuid()::varchar(36),
  "technology" VARCHAR(50) NOT NULL,
  "cache_type" VARCHAR(50) NOT NULL,
  "core_data" JSONB NOT NULL DEFAULT '{}',
  "custom_fields" JSONB DEFAULT '{}',
  "schema_version" INTEGER DEFAULT 1,
  "created_at" TIMESTAMP DEFAULT now(),
  "updated_at" TIMESTAMP DEFAULT now()
);


CREATE TABLE "log_entry" (
  "id" VARCHAR(36) PRIMARY KEY DEFAULT gen_random_uuid()::varchar(36),
  "level" VARCHAR(10) NOT NULL,
  "message" TEXT NOT NULL,
  "core_data" JSONB NOT NULL DEFAULT '{}',
  "custom_fields" JSONB DEFAULT '{}',
  "schema_version" INTEGER DEFAULT 1,
  "timestamp" TIMESTAMP DEFAULT now(),
  "created_at" TIMESTAMP DEFAULT now(),
  "updated_at" TIMESTAMP DEFAULT now()
);


CREATE TABLE "secret" (
  "id" VARCHAR(36) PRIMARY KEY DEFAULT gen_random_uuid()::varchar(36),
  "description" TEXT NOT NULL,
  "core_data" JSONB NOT NULL DEFAULT '{}',
  "custom_fields" JSONB DEFAULT '{}',
  "schema_version" INTEGER DEFAULT 1,
  "created_at" TIMESTAMP DEFAULT now(),
  "updated_at" TIMESTAMP DEFAULT now()
);


CREATE TABLE "registry_data" (
  "id" VARCHAR(36) PRIMARY KEY DEFAULT gen_random_uuid()::varchar(36),
  "subkey" VARCHAR(500) NOT NULL,
  "value_name" VARCHAR(200) NOT NULL,
  "core_data" JSONB NOT NULL DEFAULT '{}',
  "custom_fields" JSONB DEFAULT '{}',
  "schema_version" INTEGER DEFAULT 1,
  "created_at" TIMESTAMP DEFAULT now(),
  "updated_at" TIMESTAMP DEFAULT now()
);


CREATE TABLE "plist_data" (
  "id" VARCHAR(36) PRIMARY KEY DEFAULT gen_random_uuid()::varchar(36),
  "key" VARCHAR(200) NOT NULL,
  "core_data" JSONB NOT NULL DEFAULT '{}',
  "custom_fields" JSONB DEFAULT '{}',
  "schema_version" INTEGER DEFAULT 1,
  "created_at" TIMESTAMP DEFAULT now(),
  "updated_at" TIMESTAMP DEFAULT now()
);

-- Schema definitions table to track field definitions
CREATE TABLE "table_schema" (
  "id" VARCHAR(36) PRIMARY KEY DEFAULT gen_random_uuid()::varchar(36),
  "table_name" VARCHAR(100) NOT NULL,
  "schema_version" INTEGER NOT NULL,
  "field_definitions" JSONB NOT NULL,
  "is_active" BOOLEAN DEFAULT true,
  "created_at" TIMESTAMP DEFAULT now()
);

-- Create indexes for better JSONB performance
CREATE INDEX idx_docker_container_core_data ON docker_container USING GIN (core_data);
CREATE INDEX idx_docker_container_custom_fields ON docker_container USING GIN (custom_fields);
CREATE INDEX idx_git_repo_core_data ON git_repo USING GIN (core_data);
CREATE INDEX idx_git_repo_custom_fields ON git_repo USING GIN (custom_fields);
CREATE INDEX idx_cache_data_core_data ON cache_data USING GIN (core_data);
CREATE INDEX idx_log_entry_core_data ON log_entry USING GIN (core_data);
CREATE INDEX idx_secret_core_data ON secret USING GIN (core_data);
CREATE INDEX idx_registry_data_core_data ON registry_data USING GIN (core_data);
CREATE INDEX idx_plist_data_core_data ON plist_data USING GIN (core_data);
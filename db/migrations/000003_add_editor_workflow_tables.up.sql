
-- Editor configurations table
CREATE TABLE "editor_config" (
  "id" VARCHAR(36) PRIMARY KEY DEFAULT gen_random_uuid()::varchar(36),
  "name" VARCHAR(100) NOT NULL,
  "config_data" JSONB NOT NULL,
  "created_at" TIMESTAMP DEFAULT now(),
  "updated_at" TIMESTAMP DEFAULT now()
);


CREATE TABLE "workflow" (
  "id" VARCHAR(36) PRIMARY KEY DEFAULT gen_random_uuid()::varchar(36),
  "name" VARCHAR(100) NOT NULL,
  "version" VARCHAR(20),
  "export_date" TIMESTAMP,
  "node_count" INTEGER DEFAULT 0,
  "edge_count" INTEGER DEFAULT 0,
  "port_count" INTEGER DEFAULT 0,
  "workflow_data" JSONB NOT NULL,
  "created_at" TIMESTAMP DEFAULT now(),
  "updated_at" TIMESTAMP DEFAULT now()
);

-- Workflow nodes table
CREATE TABLE "workflow_node" (
  "id" VARCHAR(36) PRIMARY KEY DEFAULT gen_random_uuid()::varchar(36),
  "workflow_id" VARCHAR(36) NOT NULL REFERENCES workflow(id) ON DELETE CASCADE,
  "node_id" VARCHAR(50) NOT NULL,
  "label" VARCHAR(200),
  "group_type" VARCHAR(50),
  "title" TEXT,
  "x_position" DECIMAL,
  "y_position" DECIMAL,
  "node_data" JSONB,
  "created_at" TIMESTAMP DEFAULT now()
);

-- Workflow edges table (for easier querying)
CREATE TABLE "workflow_edge" (
  "id" VARCHAR(36) PRIMARY KEY DEFAULT gen_random_uuid()::varchar(36),
  "workflow_id" VARCHAR(36) NOT NULL REFERENCES workflow(id) ON DELETE CASCADE,
  "from_node" VARCHAR(50) NOT NULL,
  "to_node" VARCHAR(50) NOT NULL,
  "label" VARCHAR(200),
  "edge_data" JSONB,
  "created_at" TIMESTAMP DEFAULT now()
);
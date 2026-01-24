# Flexible Schema Management Guide

This API supports dynamic schema evolution, allowing you to add or remove fields from tables without breaking existing data or requiring database migrations.

## Core Concepts

### 1. JSONB Storage
- Each table has `core_data` and `custom_fields` JSONB columns
- `core_data`: Standard fields defined in the schema
- `custom_fields`: User-defined fields added dynamically
- `schema_version`: Tracks which schema version the record uses

### 2. Schema Versioning
- Schema definitions are stored in the `table_schema` table
- Each table can have multiple schema versions
- Only one schema version is active at a time
- Old data remains compatible with new schemas

## API Endpoints

### Schema Management

#### Get Current Schema
```http
GET /schema/{table_name}
```

#### Create/Update Schema
```http
POST /schema/{table_name}
Content-Type: application/json

{
  "tableName": "docker_container",
  "version": 2,
  "fields": [
    {
      "name": "lastRun",
      "type": "date",
      "required": false,
      "description": "Last time container was run"
    },
    {
      "name": "memoryLimit",
      "type": "string",
      "required": false,
      "defaultValue": "512MB"
    }
  ]
}
```

#### List Schema Versions
```http
GET /schema/{table_name}/versions
```

### Dynamic Field Management

#### Add Field to Existing Record
```http
POST /schema/{table_name}/{record_id}/fields
Content-Type: application/json

{
  "fieldPath": ["newField"],
  "value": "some value",
  "isCustom": true
}
```

#### Remove Field from Record
```http
DELETE /schema/{table_name}/{record_id}/fields/{field_name}
```

#### Search by Field Values
```http
POST /schema/{table_name}/search
Content-Type: application/json

{
  "status": "running",
  "customField": "value"
}
```

## Usage Examples

### 1. Creating a Docker Container with Flexible Fields

```http
POST /dashboard/containers
Content-Type: application/json

{
  "name": "web-app",
  "status": "running",
  "coreData": {
    "lastRun": "2026-01-10",
    "origin": "CI branch old-feature",
    "disk": "500MB",
    "ram": "256MB"
  },
  "customFields": {
    "environment": "production",
    "owner": "team-alpha",
    "costCenter": "engineering"
  }
}
```

### 2. Adding a New Field to Existing Container

```http
POST /schema/docker_container/container-id-123/fields
Content-Type: application/json

{
  "fieldPath": ["healthCheck"],
  "value": {
    "enabled": true,
    "interval": "30s",
    "timeout": "10s"
  },
  "isCustom": true
}
```

### 3. Searching Containers by Custom Fields

```http
POST /schema/docker_container/search
Content-Type: application/json

{
  "environment": "production",
  "owner": "team-alpha"
}
```

### 4. Evolving Schema Without Breaking Changes

```http
POST /schema/docker_container
Content-Type: application/json

{
  "tableName": "docker_container",
  "version": 3,
  "fields": [
    {
      "name": "lastRun",
      "type": "date",
      "required": false
    },
    {
      "name": "memoryLimit",
      "type": "string", 
      "required": false,
      "defaultValue": "512MB"
    },
    {
      "name": "securityProfile",
      "type": "object",
      "required": false,
      "description": "Security configuration"
    }
  ]
}
```

## Benefits

### ✅ Backward Compatibility
- Old records continue to work with new schema versions
- No data migration required when adding fields
- Gradual migration possible

### ✅ Forward Compatibility  
- New fields can be added without API changes
- Custom fields support arbitrary data structures
- Schema evolution doesn't break existing clients

### ✅ Performance
- JSONB indexes support fast queries on dynamic fields
- PostgreSQL's JSONB is optimized for storage and retrieval
- Structured queries still possible on core fields

### ✅ Flexibility
- Add fields at runtime without deployments
- Support different data structures per record
- Easy A/B testing of new fields

## Best Practices

1. **Use Core Data for Stable Fields**: Put frequently queried, stable fields in `core_data`
2. **Use Custom Fields for Extensions**: Put experimental or user-specific fields in `custom_fields`
3. **Version Your Schemas**: Always increment schema version when making changes
4. **Document Field Types**: Include type and description information in schema definitions
5. **Index Important Fields**: Create GIN indexes on JSONB fields you query frequently

## Migration Strategy

When you need to add a new field:

1. **Add to Schema Definition** (optional, for documentation)
2. **Start Using the Field** in new records
3. **Backfill Existing Records** (if needed) using the field management API
4. **Update Client Code** to handle the new field

No database migrations or downtime required!
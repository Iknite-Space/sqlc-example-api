package api

import (
	"encoding/json"
	"net/http"

	"github.com/Iknite-Space/sqlc-example-api/db/repo"
	"github.com/gin-gonic/gin"
)

type SchemaHandler struct {
	querier repo.Querier
}

func NewSchemaHandler(querier repo.Querier) *SchemaHandler {
	return &SchemaHandler{
		querier: querier,
	}
}

// Schema definition structure
type FieldDefinition struct {
	Name         string      `json:"name"`
	Type         string      `json:"type"`
	Required     bool        `json:"required"`
	DefaultValue interface{} `json:"defaultValue,omitempty"`
	Description  string      `json:"description,omitempty"`
}

type TableSchemaDefinition struct {
	TableName   string            `json:"tableName"`
	Version     int32             `json:"version"`
	Fields      []FieldDefinition `json:"fields"`
	Description string            `json:"description,omitempty"`
}

// Get current schema for a table
func (h *SchemaHandler) handleGetTableSchema(c *gin.Context) {
	tableName := c.Param("table")
	if tableName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "table name is required"})
		return
	}

	schema, err := h.querier.GetTableSchema(c, tableName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "schema not found"})
		return
	}

	var fieldDefs []FieldDefinition
	err = json.Unmarshal(schema.FieldDefinitions, &fieldDefs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse schema"})
		return
	}

	response := TableSchemaDefinition{
		TableName: schema.TableName,
		Version:   schema.SchemaVersion,
		Fields:    fieldDefs,
	}

	c.JSON(http.StatusOK, response)
}

// Create or update table schema
func (h *SchemaHandler) handleCreateTableSchema(c *gin.Context) {
	var req TableSchemaDefinition
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fieldDefsJSON, err := json.Marshal(req.Fields)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid field definitions"})
		return
	}

	// Deactivate old schemas
	err = h.querier.DeactivateOldSchemas(c, repo.DeactivateOldSchemasParams{
		TableName:     req.TableName,
		SchemaVersion: req.Version,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to deactivate old schemas"})
		return
	}

	// Create new schema
	params := repo.CreateTableSchemaParams{
		TableName:        req.TableName,
		SchemaVersion:    req.Version,
		FieldDefinitions: fieldDefsJSON,
		IsActive:         boolPtr(true),
	}

	schema, err := h.querier.CreateTableSchema(c, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Schema created successfully",
		"schema":  schema,
	})
}

// Add field to existing record
func (h *SchemaHandler) handleAddFieldToRecord(c *gin.Context) {
	tableName := c.Param("table")
	recordID := c.Param("id")

	var req struct {
		FieldPath []string    `json:"fieldPath"`
		Value     interface{} `json:"value"`
		IsCustom  bool        `json:"isCustom"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fieldPathJSON, _ := json.Marshal(req.FieldPath)
	valueJSON, _ := json.Marshal(req.Value)

	switch tableName {
	case "docker_container":
		if req.IsCustom {
			_, err := h.querier.AddCustomFieldToDockerContainer(c, repo.AddCustomFieldToDockerContainerParams{
				ID:          recordID,
				Path:        fieldPathJSON,
				Replacement: valueJSON,
			})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		} else {
			_, err := h.querier.UpdateDockerContainerField(c, repo.UpdateDockerContainerFieldParams{
				ID:          recordID,
				Path:        fieldPathJSON,
				Replacement: valueJSON,
			})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
	case "git_repo":
		if req.IsCustom {
			_, err := h.querier.AddCustomFieldToGitRepo(c, repo.AddCustomFieldToGitRepoParams{
				ID:          recordID,
				Path:        fieldPathJSON,
				Replacement: valueJSON,
			})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		} else {
			_, err := h.querier.UpdateGitRepoField(c, repo.UpdateGitRepoFieldParams{
				ID:          recordID,
				Path:        fieldPathJSON,
				Replacement: valueJSON,
			})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported table"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Field added successfully"})
}

// Remove field from existing record
func (h *SchemaHandler) handleRemoveFieldFromRecord(c *gin.Context) {
	tableName := c.Param("table")
	recordID := c.Param("id")
	fieldName := c.Param("field")

	switch tableName {
	case "docker_container":
		_, err := h.querier.RemoveFieldFromDockerContainer(c, repo.RemoveFieldFromDockerContainerParams{
			ID:       recordID,
			CoreData: []byte(fieldName),
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported table"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Field removed successfully"})
}

// Search records by field values
func (h *SchemaHandler) handleSearchByField(c *gin.Context) {
	tableName := c.Param("table")

	var searchCriteria map[string]interface{}
	if err := c.ShouldBindJSON(&searchCriteria); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	criteriaJSON, _ := json.Marshal(searchCriteria)

	switch tableName {
	case "docker_container":
		results, err := h.querier.SearchDockerContainersByField(c, criteriaJSON)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"results": results})
	case "git_repo":
		results, err := h.querier.SearchGitReposByField(c, criteriaJSON)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"results": results})
	case "cache_data":
		results, err := h.querier.SearchCacheDataByField(c, criteriaJSON)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"results": results})
	case "log_entry":
		results, err := h.querier.SearchLogEntriesByField(c, criteriaJSON)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"results": results})
	case "secret":
		results, err := h.querier.SearchSecretsByField(c, criteriaJSON)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"results": results})
	case "registry_data":
		results, err := h.querier.SearchRegistryDataByField(c, criteriaJSON)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"results": results})
	case "plist_data":
		results, err := h.querier.SearchPlistDataByField(c, criteriaJSON)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"results": results})
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported table"})
		return
	}
}

// List all schema versions for a table
func (h *SchemaHandler) handleListTableSchemas(c *gin.Context) {
	tableName := c.Param("table")

	schemas, err := h.querier.ListTableSchemas(c, tableName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"schemas": schemas})
}

func boolPtr(b bool) *bool {
	return &b
}

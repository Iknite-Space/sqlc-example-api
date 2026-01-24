package api

import (
	"encoding/json"
	"net/http"

	"github.com/Iknite-Space/sqlc-example-api/db/repo"
	"github.com/gin-gonic/gin"
)

type WorkflowHandler struct {
	querier repo.Querier
}

func NewWorkflowHandler(querier repo.Querier) *WorkflowHandler {
	return &WorkflowHandler{
		querier: querier,
	}
}

// Workflow endpoints
func (h *WorkflowHandler) handleCreateWorkflow(c *gin.Context) {
	var req repo.CreateWorkflowParams
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	workflow, err := h.querier.CreateWorkflow(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, workflow)
}

func (h *WorkflowHandler) handleGetWorkflow(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	workflow, err := h.querier.GetWorkflow(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Also get nodes and edges
	nodes, err := h.querier.GetWorkflowNodes(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get workflow nodes"})
		return
	}

	edges, err := h.querier.GetWorkflowEdges(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get workflow edges"})
		return
	}

	response := gin.H{
		"workflow": workflow,
		"nodes":    nodes,
		"edges":    edges,
	}

	c.JSON(http.StatusOK, response)
}

func (h *WorkflowHandler) handleListWorkflows(c *gin.Context) {
	workflows, err := h.querier.ListWorkflows(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"workflows": workflows})
}

// Editor Config endpoints
func (h *WorkflowHandler) handleCreateEditorConfig(c *gin.Context) {
	var req repo.CreateEditorConfigParams
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config, err := h.querier.CreateEditorConfig(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, config)
}

func (h *WorkflowHandler) handleGetEditorConfig(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	config, err := h.querier.GetEditorConfig(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, config)
}

func (h *WorkflowHandler) handleListEditorConfigs(c *gin.Context) {
	configs, err := h.querier.ListEditorConfigs(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"configs": configs})
}

// Bulk import endpoint for workflow data
func (h *WorkflowHandler) handleImportWorkflowData(c *gin.Context) {
	var workflowData map[string]interface{}
	if err := c.ShouldBindJSON(&workflowData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert to JSON for storage
	jsonData, err := json.Marshal(workflowData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid workflow data"})
		return
	}

	// Extract metadata if available
	name := "Imported Workflow"
	if metadata, ok := workflowData["metadata"].(map[string]interface{}); ok {
		if version, ok := metadata["version"].(string); ok {
			name = "Workflow v" + version
		}
	}

	// Create workflow record
	params := repo.CreateWorkflowParams{
		Name:         name,
		WorkflowData: jsonData,
	}

	workflow, err := h.querier.CreateWorkflow(c, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Workflow imported successfully",
		"workflow": workflow,
	})
}

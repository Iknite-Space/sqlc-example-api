package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Iknite-Space/sqlc-example-api/db/repo"
	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	querier repo.Querier
}

func NewDashboardHandler(querier repo.Querier) *DashboardHandler {
	return &DashboardHandler{
		querier: querier,
	}
}

// Flexible request structure for docker containers
type DockerContainerRequest struct {
	Name         string                 `json:"name" binding:"required"`
	Status       string                 `json:"status" binding:"required"`
	CoreData     map[string]interface{} `json:"coreData"`
	CustomFields map[string]interface{} `json:"customFields"`
}

// Flexible request structure for git repos
type GitRepoRequest struct {
	Name         string                 `json:"name" binding:"required"`
	CoreData     map[string]interface{} `json:"coreData"`
	CustomFields map[string]interface{} `json:"customFields"`
}

// Docker Container endpoints
func (h *DashboardHandler) handleCreateDockerContainer(c *gin.Context) {
	var req DockerContainerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert maps to JSON
	coreDataJSON, _ := json.Marshal(req.CoreData)
	customFieldsJSON, _ := json.Marshal(req.CustomFields)

	params := repo.CreateDockerContainerParams{
		Name:          req.Name,
		Status:        req.Status,
		CoreData:      coreDataJSON,
		CustomFields:  customFieldsJSON,
		SchemaVersion: int32Ptr(1),
	}

	container, err := h.querier.CreateDockerContainer(c, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, container)
}

func (h *DashboardHandler) handleGetDockerContainer(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	container, err := h.querier.GetDockerContainer(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Parse JSON fields for response
	var coreData, customFields map[string]interface{}
	if err := json.Unmarshal(container.CoreData, &coreData); err != nil {
		log.Printf("Failed to unmarshal core data: %v", err)
		coreData = make(map[string]interface{})
	}
	if err := json.Unmarshal(container.CustomFields, &customFields); err != nil {
		log.Printf("Failed to unmarshal custom fields: %v", err)
		customFields = make(map[string]interface{})
	}

	response := gin.H{
		"id":            container.ID,
		"name":          container.Name,
		"status":        container.Status,
		"coreData":      coreData,
		"customFields":  customFields,
		"schemaVersion": container.SchemaVersion,
		"createdAt":     container.CreatedAt,
		"updatedAt":     container.UpdatedAt,
	}

	c.JSON(http.StatusOK, response)
}

func (h *DashboardHandler) handleListDockerContainers(c *gin.Context) {
	containers, err := h.querier.ListDockerContainers(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Parse JSON fields for each container
	var response []gin.H
	for _, container := range containers {
		var coreData, customFields map[string]interface{}
		if err := json.Unmarshal(container.CoreData, &coreData); err != nil {
			log.Printf("Failed to unmarshal core data for container %s: %v", container.ID, err)
			coreData = make(map[string]interface{})
		}
		if err := json.Unmarshal(container.CustomFields, &customFields); err != nil {
			log.Printf("Failed to unmarshal custom fields for container %s: %v", container.ID, err)
			customFields = make(map[string]interface{})
		}

		response = append(response, gin.H{
			"id":            container.ID,
			"name":          container.Name,
			"status":        container.Status,
			"coreData":      coreData,
			"customFields":  customFields,
			"schemaVersion": container.SchemaVersion,
			"createdAt":     container.CreatedAt,
			"updatedAt":     container.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{"containers": response})
}

func (h *DashboardHandler) handleUpdateDockerContainer(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	var req DockerContainerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert maps to JSON
	coreDataJSON, _ := json.Marshal(req.CoreData)
	customFieldsJSON, _ := json.Marshal(req.CustomFields)

	params := repo.UpdateDockerContainerParams{
		ID:            id,
		Name:          req.Name,
		Status:        req.Status,
		CoreData:      coreDataJSON,
		CustomFields:  customFieldsJSON,
		SchemaVersion: int32Ptr(1),
	}

	container, err := h.querier.UpdateDockerContainer(c, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, container)
}

// Git Repository endpoints
func (h *DashboardHandler) handleCreateGitRepo(c *gin.Context) {
	var req GitRepoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert maps to JSON
	coreDataJSON, _ := json.Marshal(req.CoreData)
	customFieldsJSON, _ := json.Marshal(req.CustomFields)

	params := repo.CreateGitRepoParams{
		Name:          req.Name,
		CoreData:      coreDataJSON,
		CustomFields:  customFieldsJSON,
		SchemaVersion: int32Ptr(1),
	}

	gitRepo, err := h.querier.CreateGitRepo(c, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gitRepo)
}

func (h *DashboardHandler) handleListGitRepos(c *gin.Context) {
	repos, err := h.querier.ListGitRepos(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Parse JSON fields for each repo
	var response []gin.H
	for _, repo := range repos {
		var coreData, customFields map[string]interface{}
		if err := json.Unmarshal(repo.CoreData, &coreData); err != nil {
			log.Printf("Failed to unmarshal core data for repo %s: %v", repo.ID, err)
			coreData = make(map[string]interface{})
		}
		if err := json.Unmarshal(repo.CustomFields, &customFields); err != nil {
			log.Printf("Failed to unmarshal custom fields for repo %s: %v", repo.ID, err)
			customFields = make(map[string]interface{})
		}

		response = append(response, gin.H{
			"id":            repo.ID,
			"name":          repo.Name,
			"coreData":      coreData,
			"customFields":  customFields,
			"schemaVersion": repo.SchemaVersion,
			"createdAt":     repo.CreatedAt,
			"updatedAt":     repo.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{"repositories": response})
}

// Dashboard summary endpoint
func (h *DashboardHandler) handleGetDashboardSummary(c *gin.Context) {
	containers, err := h.querier.ListDockerContainers(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get containers"})
		return
	}

	repos, err := h.querier.ListGitRepos(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get repositories"})
		return
	}

	summary := gin.H{
		"containers":   containers,
		"repositories": repos,
		"timestamp":    time.Now(),
		"stats": gin.H{
			"total_containers": len(containers),
			"total_repos":      len(repos),
		},
	}

	c.JSON(http.StatusOK, summary)
}
func int32Ptr(i int32) *int32 {
	return &i
}

package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Iknite-Space/sqlc-example-api/db/repo"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

type SystemDataHandler struct {
	querier repo.Querier
}

func NewSystemDataHandler(querier repo.Querier) *SystemDataHandler {
	return &SystemDataHandler{
		querier: querier,
	}
}

// Flexible request structures
type CacheDataRequest struct {
	Technology   string                 `json:"technology" binding:"required"`
	CacheType    string                 `json:"cacheType" binding:"required"`
	CoreData     map[string]interface{} `json:"coreData"`
	CustomFields map[string]interface{} `json:"customFields"`
}

type LogEntryRequest struct {
	Level        string                 `json:"level" binding:"required"`
	Message      string                 `json:"message" binding:"required"`
	CoreData     map[string]interface{} `json:"coreData"`
	CustomFields map[string]interface{} `json:"customFields"`
	Timestamp    *time.Time             `json:"timestamp"`
}

type SecretRequest struct {
	Description  string                 `json:"description" binding:"required"`
	CoreData     map[string]interface{} `json:"coreData"`
	CustomFields map[string]interface{} `json:"customFields"`
}

type RegistryDataRequest struct {
	Subkey       string                 `json:"subkey" binding:"required"`
	ValueName    string                 `json:"valueName" binding:"required"`
	CoreData     map[string]interface{} `json:"coreData"`
	CustomFields map[string]interface{} `json:"customFields"`
}

type PlistDataRequest struct {
	Key          string                 `json:"key" binding:"required"`
	CoreData     map[string]interface{} `json:"coreData"`
	CustomFields map[string]interface{} `json:"customFields"`
}

// Cache Data endpoints
func (h *SystemDataHandler) handleCreateCacheData(c *gin.Context) {
	var req CacheDataRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	coreDataJSON, _ := json.Marshal(req.CoreData)
	customFieldsJSON, _ := json.Marshal(req.CustomFields)

	params := repo.CreateCacheDataParams{
		Technology:    req.Technology,
		CacheType:     req.CacheType,
		CoreData:      coreDataJSON,
		CustomFields:  customFieldsJSON,
		SchemaVersion: int32Ptr(1),
	}

	cacheData, err := h.querier.CreateCacheData(c, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, cacheData)
}

func (h *SystemDataHandler) handleListCacheData(c *gin.Context) {
	technology := c.Query("technology")

	var cacheData []repo.CacheDatum
	var err error

	if technology != "" {
		cacheData, err = h.querier.ListCacheDataByTechnology(c, technology)
	} else {
		cacheData, err = h.querier.ListCacheData(c)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"cacheData": cacheData})
}

func (h *SystemDataHandler) handleGetCacheData(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	cacheData, err := h.querier.GetCacheData(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cacheData)
}

// Log Entry endpoints
func (h *SystemDataHandler) handleCreateLogEntry(c *gin.Context) {
	var req LogEntryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	coreDataJSON, _ := json.Marshal(req.CoreData)
	customFieldsJSON, _ := json.Marshal(req.CustomFields)

	timestamp := time.Now()
	if req.Timestamp != nil {
		timestamp = *req.Timestamp
	}

	params := repo.CreateLogEntryParams{
		Level:         req.Level,
		Message:       req.Message,
		CoreData:      coreDataJSON,
		CustomFields:  customFieldsJSON,
		SchemaVersion: int32Ptr(1),
		Timestamp:     pgtype.Timestamp{Time: timestamp, Valid: true},
	}

	logEntry, err := h.querier.CreateLogEntry(c, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, logEntry)
}

func (h *SystemDataHandler) handleListLogEntries(c *gin.Context) {
	level := c.Query("level")

	var logEntries []repo.LogEntry
	var err error

	if level != "" {
		logEntries, err = h.querier.ListLogEntriesByLevel(c, level)
	} else {
		logEntries, err = h.querier.ListLogEntries(c)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"logEntries": logEntries})
}

func (h *SystemDataHandler) handleGetLogEntry(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	logEntry, err := h.querier.GetLogEntry(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, logEntry)
}

// Secret endpoints
func (h *SystemDataHandler) handleCreateSecret(c *gin.Context) {
	var req SecretRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	coreDataJSON, _ := json.Marshal(req.CoreData)
	customFieldsJSON, _ := json.Marshal(req.CustomFields)

	params := repo.CreateSecretParams{
		Description:   req.Description,
		CoreData:      coreDataJSON,
		CustomFields:  customFieldsJSON,
		SchemaVersion: int32Ptr(1),
	}

	secret, err := h.querier.CreateSecret(c, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, secret)
}

func (h *SystemDataHandler) handleListSecrets(c *gin.Context) {
	secrets, err := h.querier.ListSecrets(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"secrets": secrets})
}

func (h *SystemDataHandler) handleGetSecret(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	secret, err := h.querier.GetSecret(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, secret)
}

// Registry Data endpoints
func (h *SystemDataHandler) handleCreateRegistryData(c *gin.Context) {
	var req RegistryDataRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	coreDataJSON, _ := json.Marshal(req.CoreData)
	customFieldsJSON, _ := json.Marshal(req.CustomFields)

	params := repo.CreateRegistryDataParams{
		Subkey:        req.Subkey,
		ValueName:     req.ValueName,
		CoreData:      coreDataJSON,
		CustomFields:  customFieldsJSON,
		SchemaVersion: int32Ptr(1),
	}

	registryData, err := h.querier.CreateRegistryData(c, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, registryData)
}

func (h *SystemDataHandler) handleListRegistryData(c *gin.Context) {
	registryData, err := h.querier.ListRegistryData(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"registryData": registryData})
}

func (h *SystemDataHandler) handleGetRegistryData(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	registryData, err := h.querier.GetRegistryData(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, registryData)
}

// Plist Data endpoints
func (h *SystemDataHandler) handleCreatePlistData(c *gin.Context) {
	var req PlistDataRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	coreDataJSON, _ := json.Marshal(req.CoreData)
	customFieldsJSON, _ := json.Marshal(req.CustomFields)

	params := repo.CreatePlistDataParams{
		Key:           req.Key,
		CoreData:      coreDataJSON,
		CustomFields:  customFieldsJSON,
		SchemaVersion: int32Ptr(1),
	}

	plistData, err := h.querier.CreatePlistData(c, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, plistData)
}

func (h *SystemDataHandler) handleListPlistData(c *gin.Context) {
	plistData, err := h.querier.ListPlistData(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"plistData": plistData})
}

func (h *SystemDataHandler) handleGetPlistData(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	plistData, err := h.querier.GetPlistData(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, plistData)
}

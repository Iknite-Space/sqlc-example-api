package api

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/Iknite-Space/sqlc-example-api/db/repo"
	"github.com/gin-gonic/gin"
)

type MessageHandler struct {
	querier repo.Querier
}

func NewMessageHandler(querier repo.Querier) *MessageHandler {
	return &MessageHandler{
		querier: querier,
	}
}

func (h *MessageHandler) WireHttpHandler() http.Handler {

	r := gin.Default()
	r.Use(gin.CustomRecovery(func(c *gin.Context, _ any) {
		c.String(http.StatusInternalServerError, "Internal Server Error: panic")
		c.AbortWithStatus(http.StatusInternalServerError)
	})) //prevents the server from crashing if an error occurs in any route

	r.POST("/thread", h.handleCreateThread)
	r.POST("/message", h.handleCreateMessage)
	r.GET("/message/:id", h.handleGetMessage)
	// r.GET("/thread/:id/messages", h.handleGetThreadMessages)
	// r.DELETE("/message/:id", h.handleDeleteMessageById)
	r.PATCH("/message", h.handleUpdateMessage)

	return r
}

type CreateThreadParams struct {
	Title string `json:"title"`
}

func (h *MessageHandler) handleCreateThread(c *gin.Context) {
	var req CreateThreadParams
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	thread, err := h.querier.CreateThread(c, req.Title)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, thread)
}

func (h *MessageHandler) handleCreateMessage(c *gin.Context) {
	var req repo.CreateMessageParams
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//first check whether the thread exist
	_, err := h.querier.GetThreadById(c, req.ThreadID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Thread not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
		return
	}

	//now we proceed to create the message
	message, err := h.querier.CreateMessage(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, message)
}

func (h *MessageHandler) handleGetMessage(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	message, err := h.querier.GetMessageByID(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, message)
}

// func (h *MessageHandler) handleGetThreadMessages(c *gin.Context) {
// 	id := c.Param("thread_id")
// 	if id == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
// 		return
// 	}

// 	messages, err := h.querier.GetMessagesByThread(c, id)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"thread":   id,
// 		"topic":    "example",
// 		"messages": messages,
// 	})
// }

// func (h *MessageHandler) handleDeleteMessageById(c *gin.Context) {
// 	id := c.Param("id")
// 	if id == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Id cannot be null"})
// 	}

// 	//start a transaction
// 	tx, err := h.querier.(*repo.Queries).DB().Begin(c)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to start transaction"})
// 		return
// 	}
// 	defer tx.Rollback(c) // will rollback if not committed

// 	txQuerier := h.querier.(*repo.Queries).WithTx(tx)

// 	if err := txQuerier.DeleteMessageById(c, id); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to delete message"})
// 	}

// 	// Commit transaction if everything is good
// 	if err := tx.Commit(c); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"success": "Deleted successfully"})
// }

func (h *MessageHandler) handleUpdateMessage(c *gin.Context) {
	var req repo.UpdateMessageParams
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.querier.UpdateMessage(c, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": "Message updated successfully"})
}

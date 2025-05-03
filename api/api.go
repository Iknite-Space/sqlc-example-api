package api

import (
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/Iknite-Space/sqlc-example-api/campay"
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
	}))

	r.POST("/customer", h.handleCreateCustomer)
	r.POST("/order/", h.handleCreateOrders)
	// r.GET("/thread/:id/messages", h.handleGetThreadMessages)
	//r.POST("/orders", h.handleCreateCustomerOrders)

	return r
}

func (h *MessageHandler) handleCreateCustomer(c *gin.Context) {
	var req repo.CreateCustomerParams
	err := c.ShouldBindBodyWithJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	customer, err := h.querier.CreateCustomer(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, customer)
}

func (h *MessageHandler) handleCreateOrders(c *gin.Context) {
	var req repo.CreateOrderParams
	err := c.ShouldBindBodyWithJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := h.querier.CreateOrder(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	num := "673501707"
	amount := "13"
	des := "order"
	ref := GenerateRandomLetters(10)

	apik := os.Getenv("API_KEY")

	reqpay := campay.Payment(apik, num, amount, des, ref)

	time.Sleep(20 * time.Second)

	state := campay.Status(apik, reqpay.Reference)

	reqs := repo.UpadateOrderByIdParams{
		ID:          order.ID,
		OrderStatus: state.Status,
	}

	updatedOrder, err := h.querier.UpadateOrderById(c, reqs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Order Created": updatedOrder, "campay request": reqpay, "campay response": state})
}

// Function to generate random letters
func GenerateRandomLetters(n int) string {
	// Define the set of characters to choose from
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	// Create a new random generator with a source based on the current time
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Create a slice to store the random letters
	randomLetters := make([]byte, n)
	for i := 0; i < n; i++ {
		randomLetters[i] = letters[r.Intn(len(letters))]
	}
	return string(randomLetters)
}

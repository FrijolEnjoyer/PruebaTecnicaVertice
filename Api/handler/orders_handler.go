package handler

import (
	"pruebaVertice/Api/models"
	"pruebaVertice/Api/services"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// DTO para crear la orden
type CreateOrderRequest struct {
	OrderItems []models.OrderProduct `json:"order_items"`
}

type OrdersHandler struct {
	ordersService services.OrdersService
	userService   services.UserService
	logger        *logrus.Logger
}

func NewOrdersHandler(ordersService services.OrdersService, userService services.UserService, logger *logrus.Logger) *OrdersHandler {
	return &OrdersHandler{
		ordersService: ordersService,
		userService:   userService,
		logger:        logger,
	}
}

func (h *OrdersHandler) CreateOrder(c *gin.Context) {
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Layer: orderHandler, Method: CreateOrder, Error: ", err)
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	email, exists := c.Get("userEmail")
	if !exists {
		c.JSON(401, gin.H{"error": "User not authenticated"})
		return
	}

	user, err := h.userService.GetUserByEmail(email.(string))
	if err != nil {
		c.JSON(400, gin.H{"error": "User not found"})
		return
	}

	createdOrder, err := h.ordersService.CreateOrder(user.ID, req.OrderItems)
	if err != nil {
		h.logger.Error("Layer: orderHandler, Method: CreateOrder, Error: ", err)
		c.JSON(500, gin.H{"error": "Failed to create order"})
		return
	}

	c.JSON(201, createdOrder)
}

func (h *OrdersHandler) GetUserOrders(c *gin.Context) {
	email, exists := c.Get("userEmail")
	if !exists {
		c.JSON(401, gin.H{"error": "User not authenticated"})
		return
	}

	user, err := h.userService.GetUserByEmail(email.(string))
	if err != nil {
		h.logger.Error("Layer: orders_handler, Method: GetUserOrders, Error:", err)
		c.JSON(400, gin.H{"error": "User not found"})
		return
	}

	orders, err := h.ordersService.GetUserOrders(user.ID)
	if err != nil {
		h.logger.Error("Layer: orders_handler, Method: GetUserOrders, Error:", err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, orders)
}

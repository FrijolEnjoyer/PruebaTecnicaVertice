package handler

import (
	"net/http"
	"pruebaVertice/Api/models"
	services_order "pruebaVertice/Api/services/order"
	services_user "pruebaVertice/Api/services/user"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type OrdersHandler struct {
	ordersService services_order.OrdersService
	userService   services_user.UserService
	logger        *logrus.Logger
}

func NewOrdersHandler(ordersService services_order.OrdersService, userService services_user.UserService, logger *logrus.Logger) *OrdersHandler {
	return &OrdersHandler{
		ordersService: ordersService,
		userService:   userService,
		logger:        logger,
	}
}

// CreateOrder godoc
// @Summary Crear una nueva orden
// @Description Crea una nueva orden con los productos seleccionados
// @Tags Orders
// @Accept json
// @Produce json
// @Param order body models.CreateOrderRequest true "Lista de productos para la orden"
// @Success 201 {object} models.Order
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/auth/orders [post]
func (h *OrdersHandler) CreateOrder(c *gin.Context) {
	var req models.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Layer: ordersHandler, Method: CreateOrder, Error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	emailVal, exists := c.Get("userEmail")
	if !exists {
		h.logger.Error("User email not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	email := emailVal.(string)

	user, err := h.userService.GetUserByEmail(email)
	if err != nil {
		h.logger.Error("Layer: ordersHandler, Method: CreateOrder, Error fetching user:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	order, err := h.ordersService.CreateOrder(user.ID, req.OrderItems)
	if err != nil {
		h.logger.Error("Layer: ordersHandler, Method: CreateOrder, Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, order)
}

// GetUserOrders godoc
// @Summary Obtener historial de órdenes del usuario autenticado
// @Description Devuelve todas las órdenes del usuario autenticado
// @Tags Orders
// @Produce json
// @Success 200 {array} models.Order
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/auth/orders [get]
func (h *OrdersHandler) GetUserOrders(c *gin.Context) {
	emailVal, exists := c.Get("userEmail")
	if !exists {
		h.logger.Error("User email not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	email := emailVal.(string)

	user, err := h.userService.GetUserByEmail(email)
	if err != nil {
		h.logger.Error("Layer: ordersHandler, Method: GetUserOrders, Error fetching user:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	orders, err := h.ordersService.GetUserOrders(user.ID)
	if err != nil {
		h.logger.Error("Layer: ordersHandler, Method: GetUserOrders, Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}

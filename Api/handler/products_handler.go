package handler

import (
	"pruebaVertice/Api/models"
	"pruebaVertice/Api/services"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ProductsHandler struct {
	services services.ProductService
	logger   *logrus.Logger
}

func NewProductsHandler(services services.ProductService, logger *logrus.Logger) *ProductsHandler {
	return &ProductsHandler{
		services: services,
		logger:   logger,
	}
}

func (h *ProductsHandler) CreateProducts(c *gin.Context) {
	var products map[string]models.Product
	if err := c.ShouldBindJSON(&products); err != nil {
		h.logger.Error("Layer: productsHandler, Method: CreateProducts, Error:", err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	email, exists := c.Get("userEmail")
	if !exists {
		h.logger.Error("User email not found in context")
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	var productList []models.Product
	for _, product := range products {
		product.CreatedBy = email.(string)
		productList = append(productList, product)
	}

	createdProducts, err := h.services.CreateProducts(productList)
	if err != nil {
		h.logger.Error("Layer: productsHandler, Method: CreateProducts, Error:", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, createdProducts)
}

func (h *ProductsHandler) GetProductByID(c *gin.Context) {
	idParam := c.Param("id")
	idUint, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		h.logger.Error("Layer: productsHandler, Method: GetProductByID, Error: invalid product ID:", err)
		c.JSON(400, gin.H{"error": "Invalid product ID"})
		return
	}

	product, err := h.services.GetProductByID(uint(idUint))
	if err != nil {
		h.logger.Error("Layer: productsHandler, Method: GetProductByID, Error:", err)
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, product)
}

func (h *ProductsHandler) GetAllProducts(c *gin.Context) {
	products, err := h.services.GetAllProducts()
	if err != nil {
		h.logger.Error("Layer: productsHandler, Method: GetAllProducts, Error:", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, products)
}

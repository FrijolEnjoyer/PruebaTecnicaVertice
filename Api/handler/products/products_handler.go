package products

import (
	"net/http"
	"pruebaVertice/Api/models"
	services_product "pruebaVertice/Api/services/product"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ProductsHandler struct {
	services services_product.ProductService
	logger   *logrus.Logger
}

func NewProductsHandler(services services_product.ProductService, logger *logrus.Logger) *ProductsHandler {
	return &ProductsHandler{services: services, logger: logger}
}

// CreateProducts godoc
// @Summary Crear productos
// @Description Crea uno o varios productos nuevos
// @Tags Products
// @Accept json
// @Produce json
// @Param products body object true "Productos a crear (formato key-value)"
// @Success 201 {array} models.Product
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/auth/products/ [post]
func (h *ProductsHandler) CreateProducts(c *gin.Context) {
	var productsMap map[string]models.Product
	if err := c.ShouldBindJSON(&productsMap); err != nil {
		h.logger.Error("Layer: productsHandler, Method: CreateProducts, Error:", err)
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

	var productList []models.Product
	for _, p := range productsMap {
		p.CreatedBy = email
		productList = append(productList, p)
	}

	created, err := h.services.CreateProducts(productList)
	if err != nil {
		h.logger.Error("Layer: productsHandler, Method: CreateProducts, Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, created)
}

// GetProductByID godoc
// @Summary Obtener un producto por ID
// @Description Obtiene la información de un producto mediante su ID
// @Tags Products
// @Produce json
// @Param id path int true "ID del producto"
// @Success 200 {object} models.Product
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /api/auth/products/{id} [get]
func (h *ProductsHandler) GetProductByID(c *gin.Context) {
	idParam := c.Param("id")
	onceID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		h.logger.Error("Layer: productsHandler, Method: GetProductByID, Error: invalid product ID:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	product, err := h.services.GetProductByID(uint(onceID))
	if err != nil {
		h.logger.Error("Layer: productsHandler, Method: GetProductByID, Error:", err)
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

// GetAllProducts godoc
// @Summary Obtener todos los productos
// @Description Obtiene la lista de todos los productos registrados
// @Tags Products
// @Produce json
// @Success 200 {array} models.Product
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/auth/products/ [get]
func (h *ProductsHandler) GetAllProducts(c *gin.Context) {
	products, err := h.services.GetAllProducts()
	if err != nil {
		h.logger.Error("Layer: productsHandler, Method: GetAllProducts, Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}

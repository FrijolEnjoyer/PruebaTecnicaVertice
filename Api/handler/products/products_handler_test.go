package products

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"pruebaVertice/Api/models"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestCreateProducts_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	// prepare input map
	inputMap := map[string]models.Product{
		"p1": {Name: "Item1", Description: "Desc", Price: 10.0, Stock: 5},
	}
	// expected slice after binding and setting CreatedBy
	expected := []models.Product{{Name: "Item1", Description: "Desc", Price: 10.0, Stock: 5, CreatedBy: "user@example.com"}}
	serviceMock := &ProductServiceMock{}
	serviceMock.On("CreateProducts", expected).Return(expected, nil)
	logger := logrus.New()
	h := NewProductsHandler(serviceMock, logger)

	// prepare request
	body, _ := json.Marshal(inputMap)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request, _ = http.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	c.Set("userEmail", "user@example.com")

	// execute
	h.CreateProducts(c)

	// assert
	assert.Equal(t, http.StatusCreated, rec.Code)
	var resp []models.Product
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, expected, resp)
	serviceMock.AssertExpectations(t)
}

func TestCreateProducts_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	serviceMock := &ProductServiceMock{}
	logger := logrus.New()
	h := NewProductsHandler(serviceMock, logger)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request, _ = http.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(`{}`)))
	// no userEmail

	h.CreateProducts(c)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestGetProductByID_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	prod := &models.Product{Model: models.Product{}.Model, Name: "X", Price: 1.0, Stock: 2}
	serviceMock := &ProductServiceMock{}
	serviceMock.On("GetProductByID", uint(1)).Return(prod, nil)
	logger := logrus.New()
	h := NewProductsHandler(serviceMock, logger)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	c.Request, _ = http.NewRequest(http.MethodGet, "/products/1", nil)

	h.GetProductByID(c)

	assert.Equal(t, http.StatusOK, rec.Code)
	var resp models.Product
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, *prod, resp)
	serviceMock.AssertExpectations(t)
}

func TestGetProductByID_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	serviceMock := &ProductServiceMock{}
	logger := logrus.New()
	h := NewProductsHandler(serviceMock, logger)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Params = gin.Params{{Key: "id", Value: "abc"}}
	c.Request, _ = http.NewRequest(http.MethodGet, "/products/abc", nil)

	h.GetProductByID(c)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetProductByID_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	serviceMock := &ProductServiceMock{}
	serviceMock.On("GetProductByID", uint(2)).Return(nil, errors.New("not found"))
	logger := logrus.New()
	h := NewProductsHandler(serviceMock, logger)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Params = gin.Params{{Key: "id", Value: "2"}}
	c.Request, _ = http.NewRequest(http.MethodGet, "/products/2", nil)

	h.GetProductByID(c)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	serviceMock.AssertExpectations(t)
}

func TestGetAllProducts_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	existing := []models.Product{{Model: models.Product{}.Model, Name: "A"}}
	serviceMock := &ProductServiceMock{}
	serviceMock.On("GetAllProducts").Return(existing, nil)
	logger := logrus.New()
	h := NewProductsHandler(serviceMock, logger)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request, _ = http.NewRequest(http.MethodGet, "/products", nil)

	h.GetAllProducts(c)

	assert.Equal(t, http.StatusOK, rec.Code)
	var resp []models.Product
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, existing, resp)
	serviceMock.AssertExpectations(t)
}

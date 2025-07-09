package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"pruebaVertice/Api/models"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

type UserServiceMock struct {
	GetUserEmailFn func(email string) (*models.User, error)
}

func (m *UserServiceMock) CreateUser(user *models.User) (*models.User, error) { return nil, nil }
func (m *UserServiceMock) GetUserByID(id string) (*models.User, error)        { return nil, nil }
func (m *UserServiceMock) UpdateUser(user *models.User) (*models.User, error) { return nil, nil }
func (m *UserServiceMock) DeleteUser(id string) error                         { return nil }
func (m *UserServiceMock) Login(email, password string) (*models.User, error) { return nil, nil }
func (m *UserServiceMock) GetUserByEmail(email string) (*models.User, error) {
	return m.GetUserEmailFn(email)
}

func TestCreateOrder_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockOrder := &models.Order{ID: 1, UserID: 2, Total: 10.0}
	ordersMock := &OrdersServiceMock{}
	ordersMock.On("CreateOrder", uint(2), []models.OrderProduct{{ProductID: 1, Quantity: 2}}).Return(mockOrder, nil)

	userMock := &UserServiceMock{
		GetUserEmailFn: func(email string) (*models.User, error) {
			return &models.User{ID: 2, Email: email}, nil
		},
	}
	logger := logrus.New()
	h := NewOrdersHandler(ordersMock, userMock, logger)

	items := []models.OrderProduct{{ProductID: 1, Quantity: 2, UnitPrice: 0}}
	reqBody := map[string]interface{}{
		"order_items": items,
	}
	body, _ := json.Marshal(reqBody)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request, _ = http.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("userEmail", "user@example.com")

	h.CreateOrder(c)

	assert.Equal(t, http.StatusCreated, rec.Code)
	var resp models.Order
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, *mockOrder, resp)
	ordersMock.AssertExpectations(t)
}

func TestCreateOrder_BindError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ordersMock := &OrdersServiceMock{}
	userMock := &UserServiceMock{}
	logger := logrus.New()
	h := NewOrdersHandler(ordersMock, userMock, logger)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request, _ = http.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(`invalid json`)))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("userEmail", "user@example.com")

	h.CreateOrder(c)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateOrder_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ordersMock := &OrdersServiceMock{}
	userMock := &UserServiceMock{}
	logger := logrus.New()
	h := NewOrdersHandler(ordersMock, userMock, logger)

	items := []models.OrderProduct{{ProductID: 1, Quantity: 2, UnitPrice: 0}}
	reqBody := map[string]interface{}{
		"order_items": items,
	}
	body, _ := json.Marshal(reqBody)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request, _ = http.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	h.CreateOrder(c)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestGetUserOrders_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ordersList := []models.Order{{ID: 1, UserID: 2, Total: 20.0}}
	ordersMock := &OrdersServiceMock{}
	ordersMock.On("GetUserOrders", uint(2)).Return(ordersList, nil)

	userMock := &UserServiceMock{
		GetUserEmailFn: func(email string) (*models.User, error) {
			return &models.User{ID: 2, Email: email}, nil

		},
	}
	logger := logrus.New()
	h := NewOrdersHandler(ordersMock, userMock, logger)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request, _ = http.NewRequest(http.MethodGet, "/", nil)
	c.Set("userEmail", "user@example.com")

	h.GetUserOrders(c)

	assert.Equal(t, http.StatusOK, rec.Code)
	var resp []models.Order
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, ordersList, resp)
	ordersMock.AssertExpectations(t)
}

func TestGetUserOrders_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ordersMock := &OrdersServiceMock{}
	userMock := &UserServiceMock{}
	logger := logrus.New()
	h := NewOrdersHandler(ordersMock, userMock, logger)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request, _ = http.NewRequest(http.MethodGet, "/", nil)

	h.GetUserOrders(c)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

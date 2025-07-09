package services_order

import (
	"errors"
	"pruebaVertice/Api/models"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// OrdersRepoMock mocks repo.OrdersRepository
type OrdersRepoMock struct {
	mock.Mock
}

func (m *OrdersRepoMock) CreateOrder(order *models.Order) (*models.Order, error) {
	args := m.Called(order)
	if res := args.Get(0); res != nil {
		return res.(*models.Order), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *OrdersRepoMock) GetOrdersByUserID(userID uint) ([]models.Order, error) {
	args := m.Called(userID)
	if res := args.Get(0); res != nil {
		return res.([]models.Order), args.Error(1)
	}
	return nil, args.Error(1)
}

// ProductsRepoMock mocks repo.ProductsRepository
type ProductsRepoMock struct {
	mock.Mock
}

func (m *ProductsRepoMock) GetProductByID(id uint) (*models.Product, error) {
	args := m.Called(id)
	if res := args.Get(0); res != nil {
		return res.(*models.Product), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *ProductsRepoMock) UpdateProduct(product *models.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

// Stub methods to satisfy interface
func (m *ProductsRepoMock) GetAllProducts() ([]models.Product, error) {
	return nil, nil
}

func (m *ProductsRepoMock) CreateProduct(product *models.Product, createdBy string) (*models.Product, error) {
	return nil, nil
}

func (m *ProductsRepoMock) CreateProducts(products []models.Product) ([]models.Product, error) {
	return nil, nil
}

func TestCreateOrder_Success(t *testing.T) {
	orderMock := new(OrdersRepoMock)
	prodMock := new(ProductsRepoMock)
	logger := logrus.New()
	svc := NewOrdersService(orderMock, prodMock, logger)

	items := []models.OrderProduct{{ProductID: 1, Quantity: 2}}
	product := &models.Product{Model: models.Product{}.Model, Price: 5.0, Stock: 10}

	prodMock.On("GetProductByID", uint(1)).Return(product, nil)
	// After subtraction, stock should be 8
	prodMock.On("UpdateProduct", product).Return(nil)
	created := &models.Order{ID: 100, UserID: 1, Total: 10.0}
	orderMock.On("CreateOrder", mock.AnythingOfType("*models.Order")).Return(created, nil)

	res, err := svc.CreateOrder(1, items)
	assert.NoError(t, err)
	assert.Equal(t, created, res)

	prodMock.AssertExpectations(t)
	orderMock.AssertExpectations(t)
}

func TestCreateOrder_ProductNotFound(t *testing.T) {
	orderMock := new(OrdersRepoMock)
	prodMock := new(ProductsRepoMock)
	svc := NewOrdersService(orderMock, prodMock, logrus.New())

	prodMock.On("GetProductByID", uint(1)).Return(nil, errors.New("not found"))
	_, err := svc.CreateOrder(1, []models.OrderProduct{{ProductID: 1, Quantity: 1}})
	assert.EqualError(t, err, "product with ID 1 not found")
}

func TestCreateOrder_InsufficientStock(t *testing.T) {
	orderMock := new(OrdersRepoMock)
	prodMock := new(ProductsRepoMock)
	svc := NewOrdersService(orderMock, prodMock, logrus.New())

	prodMock.On("GetProductByID", uint(1)).Return(&models.Product{Stock: 1}, nil)
	_, err := svc.CreateOrder(1, []models.OrderProduct{{ProductID: 1, Quantity: 2}})
	assert.EqualError(t, err, "insufficient stock for product ID 1")
}

func TestCreateOrder_UpdateError(t *testing.T) {
	orderMock := new(OrdersRepoMock)
	prodMock := new(ProductsRepoMock)
	svc := NewOrdersService(orderMock, prodMock, logrus.New())

	product := &models.Product{Price: 5.0, Stock: 5}
	prodMock.On("GetProductByID", uint(1)).Return(product, nil)
	prodMock.On("UpdateProduct", product).Return(errors.New("db err"))
	_, err := svc.CreateOrder(1, []models.OrderProduct{{ProductID: 1, Quantity: 2}})
	assert.EqualError(t, err, "failed to update stock for product ID 1")
}

func TestGetUserOrders(t *testing.T) {
	orderMock := new(OrdersRepoMock)
	prodMock := new(ProductsRepoMock)
	svc := NewOrdersService(orderMock, prodMock, logrus.New())

	orders := []models.Order{{ID: 5, UserID: 2}}
	orderMock.On("GetOrdersByUserID", uint(2)).Return(orders, nil)
	res, err := svc.GetUserOrders(2)
	assert.NoError(t, err)
	assert.Equal(t, orders, res)
}

func TestGetUserOrders_Error(t *testing.T) {
	orderMock := new(OrdersRepoMock)
	prodMock := new(ProductsRepoMock)
	svc := NewOrdersService(orderMock, prodMock, logrus.New())

	orderMock.On("GetOrdersByUserID", uint(3)).Return(nil, errors.New("db err"))
	_, err := svc.GetUserOrders(3)
	assert.EqualError(t, err, "db err")
}

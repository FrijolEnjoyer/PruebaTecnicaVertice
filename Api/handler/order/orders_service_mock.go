package handler

import (
	"pruebaVertice/Api/models"
	
	"github.com/stretchr/testify/mock"
)

// OrdersServiceMock is a mock implementation of services_order.OrdersService
type OrdersServiceMock struct {
	mock.Mock
}

func (m *OrdersServiceMock) CreateOrder(userID uint, items []models.OrderProduct) (*models.Order, error) {
	args := m.Called(userID, items)
	if res := args.Get(0); res != nil {
		return res.(*models.Order), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *OrdersServiceMock) GetUserOrders(userID uint) ([]models.Order, error) {
	args := m.Called(userID)
	if res := args.Get(0); res != nil {
		return res.([]models.Order), args.Error(1)
	}
	return nil, args.Error(1)
}

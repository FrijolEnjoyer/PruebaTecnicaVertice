package products

import (
	"pruebaVertice/Api/models"
	"github.com/stretchr/testify/mock"
)

// ProductServiceMock is a mock implementation of services.ProductService
// for handler tests.
type ProductServiceMock struct {
	mock.Mock
}

func (m *ProductServiceMock) CreateProducts(products []models.Product) ([]models.Product, error) {
	args := m.Called(products)
	if res := args.Get(0); res != nil {
		return res.([]models.Product), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *ProductServiceMock) GetProductByID(id uint) (*models.Product, error) {
	args := m.Called(id)
	if res := args.Get(0); res != nil {
		return res.(*models.Product), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *ProductServiceMock) GetAllProducts() ([]models.Product, error) {
	args := m.Called()
	if res := args.Get(0); res != nil {
		return res.([]models.Product), args.Error(1)
	}
	return nil, args.Error(1)
}

package services

import (
	"pruebaVertice/Api/models"
	"github.com/stretchr/testify/mock"
)

// ProductsRepoMock mocks repo.ProductsRepository
// for service tests.
type ProductsRepoMock struct {
	mock.Mock
}

func (m *ProductsRepoMock) CreateProducts(products []models.Product) ([]models.Product, error) {
	args := m.Called(products)
	if res := args.Get(0); res != nil {
		return res.([]models.Product), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *ProductsRepoMock) GetProductByID(id uint) (*models.Product, error) {
	args := m.Called(id)
	if res := args.Get(0); res != nil {
		return res.(*models.Product), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *ProductsRepoMock) GetAllProducts() ([]models.Product, error) {
	args := m.Called()
	if res := args.Get(0); res != nil {
		return res.([]models.Product), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *ProductsRepoMock) UpdateProduct(product *models.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

func (m *ProductsRepoMock) CreateProduct(product *models.Product, createdBy string) (*models.Product, error) {
	args := m.Called(product, createdBy)
	if res := args.Get(0); res != nil {
		return res.(*models.Product), args.Error(1)
	}
	return nil, args.Error(1)
}

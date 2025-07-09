package services

import (
	"errors"
	"pruebaVertice/Api/models"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestCreateProducts_Success(t *testing.T) {
	repoMock := new(ProductsRepoMock)
	logger := logrus.New()
	svc := NewProductsService(repoMock, logger)

	input := []models.Product{{Name: "P1", Price: 1.0}}
	expected := []models.Product{{Name: "P1", Price: 1.0}}
	repoMock.On("CreateProducts", input).Return(expected, nil)

	res, err := svc.CreateProducts(input)
	assert.NoError(t, err)
	assert.Equal(t, expected, res)
	repoMock.AssertExpectations(t)
}

func TestCreateProducts_Error(t *testing.T) {
	repoMock := new(ProductsRepoMock)
	svc := NewProductsService(repoMock, logrus.New())

	input := []models.Product{{Name: "P2", Price: 2.0}}
	errMock := errors.New("create error")
	repoMock.On("CreateProducts", input).Return(nil, errMock)

	res, err := svc.CreateProducts(input)
	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, errMock, err)
}

func TestGetProductByID_Success(t *testing.T) {
	repoMock := new(ProductsRepoMock)
	svc := NewProductsService(repoMock, logrus.New())

	product := &models.Product{Model: models.Product{}.Model, Name: "X"}
	repoMock.On("GetProductByID", uint(1)).Return(product, nil)

	res, err := svc.GetProductByID(1)
	assert.NoError(t, err)
	assert.Equal(t, product, res)
	repoMock.AssertExpectations(t)
}

func TestGetProductByID_Error(t *testing.T) {
	repoMock := new(ProductsRepoMock)
	svc := NewProductsService(repoMock, logrus.New())

	errMock := errors.New("not found")
	repoMock.On("GetProductByID", uint(2)).Return(nil, errMock)

	res, err := svc.GetProductByID(2)
	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, errMock, err)
}

func TestGetAllProducts_Success(t *testing.T) {
	repoMock := new(ProductsRepoMock)
	svc := NewProductsService(repoMock, logrus.New())

	existing := []models.Product{{Model: models.Product{}.Model, Name: "A"}}
	repoMock.On("GetAllProducts").Return(existing, nil)

	res, err := svc.GetAllProducts()
	assert.NoError(t, err)
	assert.Equal(t, existing, res)
	repoMock.AssertExpectations(t)
}

func TestGetAllProducts_Error(t *testing.T) {
	repoMock := new(ProductsRepoMock)
	svc := NewProductsService(repoMock, logrus.New())

	errMock := errors.New("db error")
	repoMock.On("GetAllProducts").Return(nil, errMock)

	res, err := svc.GetAllProducts()
	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, errMock, err)
}

package services

import (
	"pruebaVertice/Api/models"
	repo "pruebaVertice/Api/repo/products_repo"

	"github.com/sirupsen/logrus"
)

type ProductService interface {
	CreateProducts(products []models.Product) ([]models.Product, error)
	GetProductByID(id uint) (*models.Product, error)
	GetAllProducts() ([]models.Product, error)
}

type productService struct {
	repo   repo.ProductsRepository
	logger *logrus.Logger
}

func NewProductsService(repo repo.ProductsRepository, logger *logrus.Logger) *productService {
	return &productService{
		repo:   repo,
		logger: logger,
	}
}

func (s *productService) CreateProducts(products []models.Product) ([]models.Product, error) {
	createdProducts, err := s.repo.CreateProducts(products)
	if err != nil {
		s.logger.Errorln("Layer: product_service, Method: CreateProducts, Error:", err)
		return nil, err
	}
	return createdProducts, nil
}

func (s *productService) GetProductByID(id uint) (*models.Product, error) {
	return s.repo.GetProductByID(id)
}
func (s *productService) GetAllProducts() ([]models.Product, error) {
	products, err := s.repo.GetAllProducts()
	if err != nil {
		s.logger.Errorln("Layer: product_service, Method: GetAllProducts, Error:", err)
		return nil, err
	}
	return products, nil
}

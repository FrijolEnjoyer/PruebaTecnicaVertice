package products_repo

import (
	"pruebaVertice/Api/models"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type productsRepository struct {
	db     *gorm.DB
	logger *logrus.Logger
}

func NewProductsRepository(db *gorm.DB, logger *logrus.Logger) ProductsRepository {
	return &productsRepository{
		db:     db,
		logger: logger,
	}
}

type ProductsRepository interface {
	GetAllProducts() ([]models.Product, error)
	CreateProduct(product *models.Product, createdBy string) (*models.Product, error)
	CreateProducts(products []models.Product) ([]models.Product, error) // <-- Agrega esto
	GetProductByID(id uint) (*models.Product, error)
	UpdateProduct(product *models.Product) error
}

func (r *productsRepository) GetAllProducts() ([]models.Product, error) {
	var products []models.Product
	err := r.db.Find(&products).Error
	if err != nil {
		r.logger.Errorln("Layer: products_repo, Method: GetAllProducts, Error:", err)
		return nil, err
	}
	return products, nil
}

func (r *productsRepository) CreateProducts(products []models.Product) ([]models.Product, error) {
	err := r.db.Create(&products).Error
	if err != nil {
		r.logger.Errorln("Layer: products_repo, Method: CreateProducts, Error:", err)
		return nil, err
	}
	return products, nil
}
func (r *productsRepository) GetProductByID(id uint) (*models.Product, error) {
	var product models.Product
	err := r.db.First(&product, id).Error
	if err != nil {
		r.logger.Errorln("Layer: products_repo, Method: GetProductByID, Error:", err)
		return nil, err
	}
	return &product, nil
}

func (r *productsRepository) UpdateProduct(product *models.Product) error {
	err := r.db.Save(product).Error
	if err != nil {
		r.logger.Errorln("Layer: products_repo, Method: UpdateProduct, Error:", err)
		return err
	}
	return nil
}

func (r *productsRepository) CreateProduct(product *models.Product, createdBy string) (*models.Product, error) {
	product.CreatedBy = createdBy
	err := r.db.Create(product).Error
	if err != nil {
		r.logger.Errorln("Layer: products_repo, Method: CreateProduct, Error:", err)
		return nil, err
	}
	return product, nil
}

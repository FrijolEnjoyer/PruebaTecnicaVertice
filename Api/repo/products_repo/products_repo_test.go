package products_repo

import (
	"testing"

	"pruebaVertice/Api/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/sirupsen/logrus"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func setupInMemoryDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=private"), &gorm.Config{})
	require.NoError(t, err)
	err = db.AutoMigrate(&models.Product{})
	require.NoError(t, err)
	return db
}

func TestCreateProduct_Success(t *testing.T) {
	db := setupInMemoryDB(t)
	logger := logrus.New()
	repo := NewProductsRepository(db, logger)

	input := &models.Product{Name: "P1", Description: "Desc", Price: 9.99, Stock: 5}
	created, err := repo.CreateProduct(input, "user1")
	assert.NoError(t, err)
	assert.NotZero(t, created.ID)
	assert.Equal(t, "P1", created.Name)
	assert.Equal(t, "Desc", created.Description)
	assert.Equal(t, 9.99, created.Price)
	assert.Equal(t, 5, created.Stock)
	assert.Equal(t, "user1", created.CreatedBy)
}

func TestCreateProducts_Success(t *testing.T) {
	db := setupInMemoryDB(t)
	logger := logrus.New()
	repo := NewProductsRepository(db, logger)

	inputs := []models.Product{
		{Name: "P2", Description: "D2", Price: 1.1, Stock: 1},
		{Name: "P3", Description: "D3", Price: 2.2, Stock: 2},
	}
	created, err := repo.CreateProducts(inputs)
	assert.NoError(t, err)
	assert.Len(t, created, 2)
	for i, p := range created {
		assert.NotZero(t, p.ID)
		assert.Equal(t, inputs[i].Name, p.Name)
	}
}

func TestGetAllProducts_Success(t *testing.T) {
	db := setupInMemoryDB(t)
	logger := logrus.New()
	repo := NewProductsRepository(db, logger)

	// seed
	repo.CreateProducts([]models.Product{{Name: "A", Description: "D", Price: 3.3, Stock: 3}})
	// action
	all, err := repo.GetAllProducts()
	assert.NoError(t, err)
	assert.Len(t, all, 1)
	assert.Equal(t, "A", all[0].Name)
}

func TestGetProductByID_Success(t *testing.T) {
	db := setupInMemoryDB(t)
	logger := logrus.New()
	repo := NewProductsRepository(db, logger)

	p := &models.Product{Name: "B", Description: "D", Price: 4.4, Stock: 4}
	created, err := repo.CreateProduct(p, "user2")
	require.NoError(t, err)

	fetched, err := repo.GetProductByID(created.ID)
	assert.NoError(t, err)
	assert.Equal(t, created.ID, fetched.ID)
	assert.Equal(t, "B", fetched.Name)
}

func TestUpdateProduct_Success(t *testing.T) {
	db := setupInMemoryDB(t)
	logger := logrus.New()
	repo := NewProductsRepository(db, logger)

	p := &models.Product{Name: "C", Description: "D", Price: 5.5, Stock: 5}
	created, err := repo.CreateProduct(p, "user3")
	require.NoError(t, err)

	// update fields
	created.Price = 6.6
	err = repo.UpdateProduct(created)
	assert.NoError(t, err)

	fetched, err := repo.GetProductByID(created.ID)
	assert.NoError(t, err)
	assert.Equal(t, 6.6, fetched.Price)
}

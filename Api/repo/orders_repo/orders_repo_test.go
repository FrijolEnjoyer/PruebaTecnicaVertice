package orders_repo

import (
	"testing"
	"time"

	"pruebaVertice/Api/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/sirupsen/logrus"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func setupInMemoryDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	require.NoError(t, err)
	err = db.AutoMigrate(&models.Order{}, &models.OrderProduct{})
	require.NoError(t, err)
	return db
}

func TestCreateOrder_Success(t *testing.T) {
	db := setupInMemoryDB(t)
	logger := logrus.New()
	repo := NewOrdersRepository(db, logger)

	order := &models.Order{
		UserID:  1,
		Total:   100,
		CreatedAt: time.Now(),
	}
	// call CreateOrder
	created, err := repo.CreateOrder(order)
	assert.NoError(t, err)
	assert.NotZero(t, created.ID)
	assert.Equal(t, uint(1), created.UserID)
	assert.Equal(t, 100.0, created.Total)
}

func TestGetOrdersByUserID_Success(t *testing.T) {
	db := setupInMemoryDB(t)
	logger := logrus.New()
	repo := NewOrdersRepository(db, logger)

	// seed two orders
	orders := []models.Order{
		{UserID: 2, Total: 50, CreatedAt: time.Now()},
		{UserID: 2, Total: 75, CreatedAt: time.Now()},
	}
	for i := range orders {
		_, err := repo.CreateOrder(&orders[i])
		require.NoError(t, err)
	}

	fetched, err := repo.GetOrdersByUserID(2)
	assert.NoError(t, err)
	assert.Len(t, fetched, 2)
	// verify each
	for _, o := range fetched {
		assert.Equal(t, uint(2), o.UserID)
	}
}

func TestGetOrdersByUserID_Empty(t *testing.T) {
	db := setupInMemoryDB(t)
	logger := logrus.New()
	repo := NewOrdersRepository(db, logger)

	fetched, err := repo.GetOrdersByUserID(999)
	assert.NoError(t, err)
	assert.Empty(t, fetched)
}

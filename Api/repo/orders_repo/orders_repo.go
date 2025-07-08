package orders_repo

import (
	"pruebaVertice/Api/models"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type OrdersRepository interface {
	CreateOrder(order *models.Order) (*models.Order, error)
	GetOrdersByUserID(userID uint) ([]models.Order, error)
}

type ordersRepository struct {
	db     *gorm.DB
	logger *logrus.Logger
}

func NewOrdersRepository(db *gorm.DB, logger *logrus.Logger) OrdersRepository {
	return &ordersRepository{db: db, logger: logger}
}

func (r *ordersRepository) CreateOrder(order *models.Order) (*models.Order, error) {
	err := r.db.Create(order).Error
	if err != nil {
		r.logger.Errorln("Layer: orders_repo, Method: CreateOrder, Error:", err)
		return nil, err
	}
	return order, nil
}
func (r *ordersRepository) GetOrdersByUserID(userID uint) ([]models.Order, error) {
	var orders []models.Order
	err := r.db.Preload("OrderItems").Where("user_id = ?", userID).Find(&orders).Error
	if err != nil {
		r.logger.Errorln("Layer: orders_repo, Method: GetOrdersByUserID, Error:", err)
		return nil, err
	}
	return orders, nil
}

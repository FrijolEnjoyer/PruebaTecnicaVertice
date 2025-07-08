package services

import (
	"fmt"
	"pruebaVertice/Api/models"
	repo "pruebaVertice/Api/repo/orders_repo"
	productsRepo "pruebaVertice/Api/repo/products_repo"

	"github.com/sirupsen/logrus"
)

type OrdersService interface {
	CreateOrder(userID uint, items []models.OrderProduct) (*models.Order, error)
	GetUserOrders(userID uint) ([]models.Order, error)
}

type ordersService struct {
	orderRepo   repo.OrdersRepository
	productRepo productsRepo.ProductsRepository
	logger      *logrus.Logger
}

func NewOrdersService(orderRepo repo.OrdersRepository, productRepo productsRepo.ProductsRepository, logger *logrus.Logger) *ordersService {
	return &ordersService{
		orderRepo:   orderRepo,
		productRepo: productRepo,
		logger:      logger,
	}
}

func (s *ordersService) CreateOrder(userID uint, items []models.OrderProduct) (*models.Order, error) {
	var total float64
	var orderItems []models.OrderProduct

	for _, item := range items {
		product, err := s.productRepo.GetProductByID(item.ProductID)
		if err != nil {
			return nil, fmt.Errorf("product with ID %d not found", item.ProductID)
		}

		if product.Stock < item.Quantity {
			return nil, fmt.Errorf("insufficient stock for product ID %d", item.ProductID)
		}

		unitPrice := product.Price
		total += unitPrice * float64(item.Quantity)

		// Descontar stock
		product.Stock -= item.Quantity
		if err := s.productRepo.UpdateProduct(product); err != nil {
			return nil, fmt.Errorf("failed to update stock for product ID %d", item.ProductID)
		}

		orderItems = append(orderItems, models.OrderProduct{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			UnitPrice: unitPrice,
		})
	}

	order := &models.Order{
		UserID:     userID,
		Total:      total,
		OrderItems: orderItems,
	}

	createdOrder, err := s.orderRepo.CreateOrder(order)
	if err != nil {
		return nil, err
	}
	return createdOrder, nil
}

func (s *ordersService) GetUserOrders(userID uint) ([]models.Order, error) {
	return s.orderRepo.GetOrdersByUserID(userID)
}

package models

import (
	"time"
)

type Order struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	UserID     uint           `json:"user_id"`
	Total      float64        `json:"total"`
	CreatedAt  time.Time      `json:"created_at"`
	OrderItems []OrderProduct `gorm:"foreignKey:OrderID" json:"order_items"`
}

type OrderProduct struct {
	ID        uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderID   uint    `json:"order_id"`
	ProductID uint    `json:"product_id"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
}

type CreateOrderRequest struct {
	OrderItems []OrderProduct `json:"order_items"`
}

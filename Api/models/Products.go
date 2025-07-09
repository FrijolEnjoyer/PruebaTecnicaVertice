package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model  `json:"-" swaggerignore:"true"`
	Name        string  `gorm:"type:varchar(255);uniqueIndex" json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	CreatedBy   string  `json:"created_by"`
}

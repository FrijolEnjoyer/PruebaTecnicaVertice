package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Username     string         `gorm:"type:varchar(255);uniqueIndex" json:"username"`
	Password     string         `validate:"required,password" json:"password,omitempty" bson:"password"`
	Email        string         `gorm:"type:varchar(255);uniqueIndex" json:"email" bson:"email"`
	RefreshToken string         `json:"refresh_token,omitempty"`
	Token        string         `json:"token,omitempty"`
}

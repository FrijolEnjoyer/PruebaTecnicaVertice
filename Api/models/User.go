package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username     string `gorm:"type:varchar(255);uniqueIndex" json:"username"`
	Password     string `validate:"required,password" json:"password" bson:"password"`
	Email        string `gorm:"type:varchar(255);uniqueIndex" json:"email" bson:"email"`
	RefreshToken string `json:"refresh_token"`
	Token        string `json:"token"`
}

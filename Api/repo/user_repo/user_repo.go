package user_repo

import (
	"fmt"
	"pruebaVertice/Api/models"

	"github.com/sirupsen/logrus"

	"gorm.io/gorm"
)

type userRepository struct {
	db     *gorm.DB
	logger *logrus.Logger
}

func NewUserRepository(db *gorm.DB, logger *logrus.Logger) UserRepository {
	return &userRepository{
		db:     db,
		logger: logger,
	}
}

type UserRepository interface {
	CreateUser(user *models.User) (*models.User, error)
	GetUserByID(id string) (*models.User, error)
	UpdateUser(user *models.User) (*models.User, error)
	DeleteUser(id string) error
	UpdateUserToken(user *models.User) (*models.User, error)
	GetUserByEmail(email string) (models.User, error)
}

func (r *userRepository) CreateUser(user *models.User) (*models.User, error) {
	err := r.db.Create(user).Error
	if err != nil {
		r.logger.Error("Layer: userRepo, method: CreateUser, error:", err)
		return nil, err
	}
	return user, nil
}

func (r *userRepository) GetUserByID(id string) (*models.User, error) {
	user := &models.User{}
	err := r.db.First(user, id).Error
	fmt.Println(id)
	if err != nil {
		r.logger.Error("Layer: userRepo, method: GetUserByID, error:", err)
		return nil, err
	}
	fmt.Println("Buen texto", user)
	return user, nil
}
func (r *userRepository) UpdateUser(user *models.User) (*models.User, error) {
	err := r.db.Save(user).Error
	if err != nil {
		r.logger.Error("Layer: userRepo, method: UpdateUser, error:", err)
		return nil, err
	}
	return user, nil
}
func (r *userRepository) DeleteUser(id string) error {
	err := r.db.Delete(&models.User{}, id).Error
	if err != nil {
		r.logger.Error("Layer: userRepo, method: DeleteUser, error:", err)
		return err
	}
	return nil
}

func (r *userRepository) UpdateUserToken(user *models.User) (*models.User, error) {
	var model models.User
	err := r.db.Where("id = ?", user.ID).First(&model).Error
	if err != nil {
		r.logger.Error("Layer: userRepo, method: UpdateUserToken, error:", err)
		return nil, err
	}
	model.Token = user.Token
	model.RefreshToken = user.RefreshToken
	err = r.db.Save(&model).Error
	if err != nil {
		r.logger.Error("Layer: userRepo, method: UpdateUserToken, error:", err)
		return nil, err
	}
	return &model, nil
}

func (r *userRepository) GetUserByEmail(email string) (models.User, error) {
	var userModel models.User
	err := r.db.Where("email = ?", email).First(&userModel).Error
	if err != nil {
		r.logger.Errorln("Layer:user_repository, Method:GetUserByEmail, Error:", err)
		return models.User{}, err
	}
	return userModel, nil
}

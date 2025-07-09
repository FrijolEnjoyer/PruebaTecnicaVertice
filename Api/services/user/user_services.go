package services

import (
	"errors"
	"pruebaVertice/Api/models"
	repo "pruebaVertice/Api/repo/user_repo"

	"github.com/sirupsen/logrus"
)

type UserService interface {
	CreateUser(user *models.User) (*models.User, error)
	GetUserByID(id string) (*models.User, error)
	UpdateUser(user *models.User) (*models.User, error)
	DeleteUser(id string) error
	Login(email, password string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
}

type userService struct {
	repo           repo.UserRepository
	logger         *logrus.Logger
	hasher         Hasher
	tokenGenerator TokenGenerator
}

var ErrInvalidPassword = errors.New("invalid password")

// Hasher defines password hashing behavior
type Hasher interface {
	HashPassword(password string) (string, error)
	CheckPasswordHash(password, hash string) bool
}

// TokenGenerator defines JWT behavior
 type TokenGenerator interface {
	GenerateToken(email string) (string, string, error)
}

 func NewUserService(repo repo.UserRepository, hasher Hasher, tokenGen TokenGenerator, logger *logrus.Logger) *userService {

	return &userService{
		repo:           repo,
		hasher:         hasher,
		tokenGenerator: tokenGen,
		logger:         logger,
	}
}
func (s *userService) CreateUser(user *models.User) (*models.User, error) {
	hashedPassword, err := s.hasher.HashPassword(user.Password)
	if err != nil {
		s.logger.Errorln("Layer:user_service, Method:CreateUser, Error: Hashing password:", err)
		return nil, err
	}
	user.Password = hashedPassword

	_, err = s.repo.CreateUser(user)
	if err != nil {
		s.logger.Errorln("Layer:user_service, Method:CreateUser, Error:", err)
		return nil, err
	}

	token, refreshToken, err := s.tokenGenerator.GenerateToken(user.Email)
	if err != nil {
		s.logger.Errorln("Layer:user_service, Method:CreateUser, Error: Generating token:", err)
		return nil, err
	}

	user.Token = token
	user.RefreshToken = refreshToken

	_, err = s.repo.UpdateUserToken(user)
	if err != nil {
		s.logger.Errorln("Layer:user_service, Method:CreateUser, Error: Updating user token:", err)
		return nil, err
	}

	return user, nil
}

func (s *userService) GetUserByID(id string) (*models.User, error) {
	return s.repo.GetUserByID(id)
}

func (s *userService) UpdateUser(user *models.User) (*models.User, error) {
	return s.repo.UpdateUser(user)
}
func (s *userService) DeleteUser(id string) error {
	return s.repo.DeleteUser(id)
}

func (s *userService) Login(email, password string) (*models.User, error) {
	user, _ := s.repo.GetUserByEmail(email)

	if !s.hasher.CheckPasswordHash(password, user.Password) {
		s.logger.Errorln("Layer:user_service, Method:Login, Error: Invalid password")
		return nil, ErrInvalidPassword
	}

	token, refreshToken, err := s.tokenGenerator.GenerateToken(user.Email)
	if err != nil {
		s.logger.Errorln("Layer:user_service, Method:Login, Error: Generating token:", err)
		return nil, err
	}
	user.Token = token
	user.RefreshToken = refreshToken

	userr, err := s.repo.UpdateUserToken(&user)
	if err != nil {
		s.logger.Errorln("Layer:user_service, Method:Login, Error: Updating user token:", err)
		return nil, err
	}

	user.CreatedAt = userr.CreatedAt
	user.UpdatedAt = userr.UpdatedAt
	return &user, nil
}

func (s *userService) GetUserByEmail(email string) (*models.User, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		s.logger.Errorln("Layer:user_service, Method:GetUserByEmail, Error:", err)
		return nil, err
	}
	return &user, nil
}

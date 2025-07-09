package services

import (
	"pruebaVertice/Api/models"
	"github.com/stretchr/testify/mock"
)

// UserRepoMock mocks repo.UserRepository for service tests.
type UserRepoMock struct {
	mock.Mock
}

func (m *UserRepoMock) CreateUser(user *models.User) (*models.User, error) {
	args := m.Called(user)
	if res := args.Get(0); res != nil {
		return res.(*models.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *UserRepoMock) GetUserByID(id string) (*models.User, error) {
	args := m.Called(id)
	if res := args.Get(0); res != nil {
		return res.(*models.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *UserRepoMock) UpdateUser(user *models.User) (*models.User, error) {
	args := m.Called(user)
	if res := args.Get(0); res != nil {
		return res.(*models.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *UserRepoMock) DeleteUser(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *UserRepoMock) UpdateUserToken(user *models.User) (*models.User, error) {
	args := m.Called(user)
	if res := args.Get(0); res != nil {
		return res.(*models.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *UserRepoMock) GetUserByEmail(email string) (models.User, error) {
	args := m.Called(email)
	if res := args.Get(0); res != nil {
		return args.Get(0).(models.User), args.Error(1)
	}
	return models.User{}, args.Error(1)
}

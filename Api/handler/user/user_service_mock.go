package user

import (
	"pruebaVertice/Api/models"
	"github.com/stretchr/testify/mock"
)

// UserServiceMock is a mock implementation of services_user.UserService
// for handler tests.
type UserServiceMock struct {
	mock.Mock
}

func (m *UserServiceMock) CreateUser(user *models.User) (*models.User, error) {
	args := m.Called(user)
	if res := args.Get(0); res != nil {
		return res.(*models.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *UserServiceMock) GetUserByID(id string) (*models.User, error) {
	args := m.Called(id)
	if res := args.Get(0); res != nil {
		return res.(*models.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *UserServiceMock) UpdateUser(user *models.User) (*models.User, error) {
	args := m.Called(user)
	if res := args.Get(0); res != nil {
		return res.(*models.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *UserServiceMock) DeleteUser(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *UserServiceMock) Login(email, password string) (*models.User, error) {
	args := m.Called(email, password)
	if res := args.Get(0); res != nil {
		return res.(*models.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *UserServiceMock) GetUserByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if res := args.Get(0); res != nil {
		return res.(*models.User), args.Error(1)
	}
	return nil, args.Error(1)
}

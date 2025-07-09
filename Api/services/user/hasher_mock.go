package services

import "github.com/stretchr/testify/mock"

// HasherMock mocks utils.BcryptHasher for service tests.
type HasherMock struct {
	mock.Mock
}

func (m *HasherMock) HashPassword(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *HasherMock) CheckPasswordHash(password, hash string) bool {
	args := m.Called(password, hash)
	return args.Bool(0)
}

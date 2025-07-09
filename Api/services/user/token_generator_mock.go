package services

import "github.com/stretchr/testify/mock"

// TokenGeneratorMock mocks jwtUtils.JWTGenerator for service tests.
type TokenGeneratorMock struct {
	mock.Mock
}

func (m *TokenGeneratorMock) GenerateToken(email string) (string, string, error) {
	args := m.Called(email)
	return args.String(0), args.String(1), args.Error(2)
}

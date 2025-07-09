package services

import (
	"errors"
	"pruebaVertice/Api/models"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateUser_Success(t *testing.T) {
	repoMock := new(UserRepoMock)
	hasherMock := new(HasherMock)
	tokenMock := new(TokenGeneratorMock)
	logger := logrus.New()
	svc := NewUserService(repoMock, hasherMock, tokenMock, logger)

	userInput := &models.User{Email: "u@e.com", Password: "pwd"}
	hashed := "hashedpwd"
	hasherMock.On("HashPassword", "pwd").Return(hashed, nil)

	created := &models.User{Email: "u@e.com", Password: hashed}
	repoMock.On("CreateUser", userInput).Return(created, nil)

	token := "tok"
	refresh := "ref"
	tokenMock.On("GenerateToken", "u@e.com").Return(token, refresh, nil)

	// updatedModel without timestamps
	updatedModel := &models.User{Email: "u@e.com", Password: hashed, Token: token, RefreshToken: refresh}
	repoMock.On("UpdateUserToken", mock.AnythingOfType("*models.User")).Return(updatedModel, nil)

	res, err := svc.CreateUser(userInput)
	assert.NoError(t, err)
	// only check relevant fields
	assert.Equal(t, token, res.Token)
	assert.Equal(t, refresh, res.RefreshToken)
	assert.Equal(t, hashed, res.Password)
	assert.Equal(t, userInput.Email, res.Email)
	repoMock.AssertExpectations(t)
	hasherMock.AssertExpectations(t)
	tokenMock.AssertExpectations(t)
}

func TestCreateUser_HashError(t *testing.T) {
	repoMock := new(UserRepoMock)
	hasherMock := new(HasherMock)
	tokenMock := new(TokenGeneratorMock)
	logger := logrus.New()
	svc := NewUserService(repoMock, hasherMock, tokenMock, logger)

	hasherMock.On("HashPassword", "pwd").Return("", errors.New("hash err"))
	_, err := svc.CreateUser(&models.User{Password: "pwd"})
	assert.EqualError(t, err, "hash err")
}

func TestCreateUser_CreateError(t *testing.T) {
	repoMock := new(UserRepoMock)
	hasherMock := new(HasherMock)
	tokenMock := new(TokenGeneratorMock)
	logger := logrus.New()
	svc := NewUserService(repoMock, hasherMock, tokenMock, logger)

	hasherMock.On("HashPassword", "pwd").Return("hashed", nil)
	repoMock.On("CreateUser", mock.Anything).Return(nil, errors.New("create err"))
	_, err := svc.CreateUser(&models.User{Password: "pwd"})
	assert.EqualError(t, err, "create err")
}

func TestCreateUser_GenerateTokenError(t *testing.T) {
	repoMock := new(UserRepoMock)
	hasherMock := new(HasherMock)
	tokenMock := new(TokenGeneratorMock)
	logger := logrus.New()
	svc := NewUserService(repoMock, hasherMock, tokenMock, logger)

	hasherMock.On("HashPassword", "pwd").Return("hashed", nil)
	repoMock.On("CreateUser", mock.Anything).Return(&models.User{Email: "e@e"}, nil)
	tokenMock.On("GenerateToken", "e@e").Return("", "", errors.New("tok err"))
	_, err := svc.CreateUser(&models.User{Email: "e@e", Password: "pwd"})
	assert.EqualError(t, err, "tok err")
}

func TestCreateUser_UpdateTokenError(t *testing.T) {
	repoMock := new(UserRepoMock)
	hasherMock := new(HasherMock)
	tokenMock := new(TokenGeneratorMock)
	logger := logrus.New()
	svc := NewUserService(repoMock, hasherMock, tokenMock, logger)

	hasherMock.On("HashPassword", "pwd").Return("hashed", nil)
	repoMock.On("CreateUser", mock.Anything).Return(&models.User{Email: "e@e"}, nil)
	tokenMock.On("GenerateToken", "e@e").Return("tok", "ref", nil)
	repoMock.On("UpdateUserToken", mock.Anything).Return(nil, errors.New("upd err"))
	_, err := svc.CreateUser(&models.User{Email: "e@e", Password: "pwd"})
	assert.EqualError(t, err, "upd err")
}

func TestGetUserByID(t *testing.T) {
	repoMock := new(UserRepoMock)
	logger := logrus.New()
	svc := NewUserService(repoMock, nil, nil, logger)

	repoMock.On("GetUserByID", "1").Return(&models.User{Email: "x"}, nil)
	res, err := svc.GetUserByID("1")
	assert.NoError(t, err)
	assert.Equal(t, &models.User{Email: "x"}, res)
}

func TestGetUserByID_Error(t *testing.T) {
	repoMock := new(UserRepoMock)
	logger := logrus.New()
	svc := NewUserService(repoMock, nil, nil, logger)

	repoMock.On("GetUserByID", "1").Return(nil, errors.New("not found"))
	_, err := svc.GetUserByID("1")
	assert.EqualError(t, err, "not found")
}

func TestUpdateUser(t *testing.T) {
	repoMock := new(UserRepoMock)
	logger := logrus.New()
	svc := NewUserService(repoMock, nil, nil, logger)

	repoMock.On("UpdateUser", &models.User{Email: "u"}).Return(&models.User{Email: "u"}, nil)
	res, err := svc.UpdateUser(&models.User{Email: "u"})
	assert.NoError(t, err)
	assert.Equal(t, &models.User{Email: "u"}, res)
}

func TestUpdateUser_Error(t *testing.T) {
	repoMock := new(UserRepoMock)
	logger := logrus.New()
	svc := NewUserService(repoMock, nil, nil, logger)

	repoMock.On("UpdateUser", mock.Anything).Return(nil, errors.New("upd err"))
	_, err := svc.UpdateUser(&models.User{})
	assert.EqualError(t, err, "upd err")
}

func TestDeleteUser(t *testing.T) {
	repoMock := new(UserRepoMock)
	logger := logrus.New()
	svc := NewUserService(repoMock, nil, nil, logger)

	repoMock.On("DeleteUser", "2").Return(nil)
	err := svc.DeleteUser("2")
	assert.NoError(t, err)
}

func TestDeleteUser_Error(t *testing.T) {
	repoMock := new(UserRepoMock)
	logger := logrus.New()
	svc := NewUserService(repoMock, nil, nil, logger)

	repoMock.On("DeleteUser", "2").Return(errors.New("del err"))
	err := svc.DeleteUser("2")
	assert.EqualError(t, err, "del err")
}

func TestLogin_Success(t *testing.T) {
	repoMock := new(UserRepoMock)
	hasherMock := new(HasherMock)
	tokenMock := new(TokenGeneratorMock)
	logger := logrus.New()
	svc := NewUserService(repoMock, hasherMock, tokenMock, logger)

	stored := models.User{Email: "e@e", Password: "hashed"}
	repoMock.On("GetUserByEmail", "e@e").Return(stored, nil)
	hasherMock.On("CheckPasswordHash", "pwd", "hashed").Return(true)
	tokenMock.On("GenerateToken", "e@e").Return("tok", "ref", nil)

	// updated without timestamps
		updated := &models.User{Email: "e@e", Token: "tok", RefreshToken: "ref"}
	repoMock.On("UpdateUserToken", mock.AnythingOfType("*models.User")).Return(updated, nil)

	res, err := svc.Login("e@e", "pwd")
	assert.NoError(t, err)
	// assert key fields for login
	assert.Equal(t, updated.Token, res.Token)
	assert.Equal(t, updated.RefreshToken, res.RefreshToken)
	assert.Equal(t, updated.Email, res.Email)
}

func TestLogin_InvalidPassword(t *testing.T) {
	repoMock := new(UserRepoMock)
	hasherMock := new(HasherMock)
	tokenMock := new(TokenGeneratorMock)
	logger := logrus.New()
	svc := NewUserService(repoMock, hasherMock, tokenMock, logger)

	repoMock.On("GetUserByEmail", "e@e").Return(models.User{Email: "e@e", Password: "hashed"}, nil)
	hasherMock.On("CheckPasswordHash", "pwd", "hashed").Return(false)
	_, err := svc.Login("e@e", "pwd")
	assert.EqualError(t, err, ErrInvalidPassword.Error())
}

func TestLogin_GenerateError(t *testing.T) {
	repoMock := new(UserRepoMock)
	hasherMock := new(HasherMock)
	tokenMock := new(TokenGeneratorMock)
	logger := logrus.New()
	svc := NewUserService(repoMock, hasherMock, tokenMock, logger)

	repoMock.On("GetUserByEmail", "e@e").Return(models.User{Email: "e@e", Password: "hashed"}, nil)
	hasherMock.On("CheckPasswordHash", "pwd", "hashed").Return(true)
	tokenMock.On("GenerateToken", "e@e").Return("", "", errors.New("tok err"))
	_, err := svc.Login("e@e", "pwd")
	assert.EqualError(t, err, "tok err")
}

func TestLogin_UpdateTokenError(t *testing.T) {
	repoMock := new(UserRepoMock)
	hasherMock := new(HasherMock)
	tokenMock := new(TokenGeneratorMock)
	logger := logrus.New()
	svc := NewUserService(repoMock, hasherMock, tokenMock, logger)

	repoMock.On("GetUserByEmail", "e@e").Return(models.User{Email: "e@e", Password: "hashed"}, nil)
	hasherMock.On("CheckPasswordHash", "pwd", "hashed").Return(true)
	tokenMock.On("GenerateToken", "e@e").Return("tok", "ref", nil)
	repoMock.On("UpdateUserToken", mock.Anything).Return(nil, errors.New("upd err"))
	_, err := svc.Login("e@e", "pwd")
	assert.EqualError(t, err, "upd err")
}

func TestGetUserByEmail(t *testing.T) {
	repoMock := new(UserRepoMock)
	logger := logrus.New()
	svc := NewUserService(repoMock, nil, nil, logger)

	repoMock.On("GetUserByEmail", "a@b").Return(models.User{Email: "a@b"}, nil)
	res, err := svc.GetUserByEmail("a@b")
	assert.NoError(t, err)
	assert.Equal(t, &models.User{Email: "a@b"}, res)
}

func TestGetUserByEmail_Error(t *testing.T) {
	repoMock := new(UserRepoMock)
	logger := logrus.New()
	svc := NewUserService(repoMock, nil, nil, logger)

	repoMock.On("GetUserByEmail", "a@b").Return(models.User{}, errors.New("not found"))
	_, err := svc.GetUserByEmail("a@b")
	assert.EqualError(t, err, "not found")
}

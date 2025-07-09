package user

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"pruebaVertice/Api/models"
	"pruebaVertice/Api/dto"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestCreateUser_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	userInput := models.User{Username: "john", Password: "pass", Email: "john@example.com"}
	created := &models.User{Model: gorm.Model{ID: 1}, Username: "john", Email: "john@example.com"}
	serviceMock := &UserServiceMock{}
	serviceMock.On("CreateUser", &userInput).Return(created, nil)
	logger := logrus.New()
	h := NewUserHandler(serviceMock, logger)

	body, _ := json.Marshal(userInput)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request, _ = http.NewRequest(http.MethodPost, "/", bytes.NewReader(body))

	h.CreateUser(c)

	assert.Equal(t, http.StatusCreated, rec.Code)
	var resp models.User
	_ = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.Equal(t, *created, resp)
	serviceMock.AssertExpectations(t)
}

func TestGetUserByID_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	returned := &models.User{Model: gorm.Model{ID: 2}, Username: "alice", Email: "a@b.com"}
	serviceMock := &UserServiceMock{}
	serviceMock.On("GetUserByID", "2").Return(returned, nil)
	logger := logrus.New()
	h := NewUserHandler(serviceMock, logger)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Params = gin.Params{{Key: "id", Value: "2"}}
	c.Request, _ = http.NewRequest(http.MethodGet, "/users/2", nil)

	h.GetUserByID(c)

	assert.Equal(t, http.StatusOK, rec.Code)
	var resp models.User
	_ = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.Equal(t, *returned, resp)
	serviceMock.AssertExpectations(t)
}

func TestGetUserByID_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	serviceMock := &UserServiceMock{}
	serviceMock.On("GetUserByID", "3").Return(nil, assert.AnError)
	logger := logrus.New()
	h := NewUserHandler(serviceMock, logger)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Params = gin.Params{{Key: "id", Value: "3"}}
	c.Request, _ = http.NewRequest(http.MethodGet, "/users/3", nil)

	h.GetUserByID(c)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	serviceMock.AssertExpectations(t)
}

func TestUpdateUser_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	input := models.User{Username: "bob", Email: "b@c.com"}
	updated := &models.User{Model: gorm.Model{ID: 4}, Username: "bob", Email: "b@c.com"}
	serviceMock := &UserServiceMock{}
	serviceMock.On("UpdateUser", &input).Return(updated, nil)
	logger := logrus.New()
	h := NewUserHandler(serviceMock, logger)

	body, _ := json.Marshal(input)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request, _ = http.NewRequest(http.MethodPut, "/", bytes.NewReader(body))

	h.UpdateUser(c)

	assert.Equal(t, http.StatusOK, rec.Code)
	var resp models.User
	_ = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.Equal(t, *updated, resp)
	serviceMock.AssertExpectations(t)
}

func TestDeleteUser_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	serviceMock := &UserServiceMock{}
	serviceMock.On("DeleteUser", "5").Return(nil)
	logger := logrus.New()
	h := NewUserHandler(serviceMock, logger)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Params = gin.Params{{Key: "id", Value: "5"}}
	c.Request, _ = http.NewRequest(http.MethodDelete, "/users/5", nil)

	h.DeleteUser(c)

	assert.Equal(t, http.StatusOK, rec.Code)
	var result map[string]string
	_ = json.Unmarshal(rec.Body.Bytes(), &result)
	assert.Equal(t, "User deleted successfully", result["message"])
	serviceMock.AssertExpectations(t)
}

func TestDeleteUser_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	serviceMock := &UserServiceMock{}
	serviceMock.On("DeleteUser", "6").Return(assert.AnError)
	logger := logrus.New()
	h := NewUserHandler(serviceMock, logger)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Params = gin.Params{{Key: "id", Value: "6"}}
	c.Request, _ = http.NewRequest(http.MethodDelete, "/users/6", nil)

	h.DeleteUser(c)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	serviceMock.AssertExpectations(t)
}

func TestLoginUser_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	input := models.User{Email: "u@e.com", Password: "pwd"}
	logged := &models.User{Token: "t1", RefreshToken: "r1"}
	serviceMock := &UserServiceMock{}
	serviceMock.On("Login", "u@e.com", "pwd").Return(logged, nil)
	logger := logrus.New()
	h := NewUserHandler(serviceMock, logger)

	body, _ := json.Marshal(input)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request, _ = http.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))

	h.LoginUser(c)

	assert.Equal(t, http.StatusOK, rec.Code)
	var resp dto.LoginResponse
	_ = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.Equal(t, dto.LoginResponse{Token: "t1", RefreshToken: "r1"}, resp)
	serviceMock.AssertExpectations(t)
}

func TestGetLoggedInUser_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	userData := &models.User{Model: gorm.Model{ID: 7}, Username: "x", Email: "x@y"}
	serviceMock := &UserServiceMock{}
	serviceMock.On("GetUserByEmail", "x@y").Return(userData, nil)
	logger := logrus.New()
	h := NewUserHandler(serviceMock, logger)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request, _ = http.NewRequest(http.MethodGet, "/me", nil)
	c.Set("userEmail", "x@y")

	h.GetLoggedInUser(c)

	assert.Equal(t, http.StatusOK, rec.Code)
	var resp models.User
	_ = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.Equal(t, *userData, resp)
	serviceMock.AssertExpectations(t)
}

func TestGetLoggedInUser_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	serviceMock := &UserServiceMock{}
	logger := logrus.New()
	h := NewUserHandler(serviceMock, logger)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request, _ = http.NewRequest(http.MethodGet, "/me", nil)

	h.GetLoggedInUser(c)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

package handler

import (
	"pruebaVertice/Api/dto"
	"pruebaVertice/Api/models"
	"pruebaVertice/Api/services"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type UserHandler struct {
	services services.UserService
	logger   *logrus.Logger
}

func NewUserHandler(services services.UserService, logger *logrus.Logger) *UserHandler {
	return &UserHandler{
		services: services,
		logger:   logger,
	}
}
func (h *UserHandler) CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		h.logger.Error("Layer: userHandler, method: CreateUser, error:", err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	createdUser, err := h.services.CreateUser(&user)
	if err != nil {
		h.logger.Error("Layer: userHandler, method: CreateUser, error:", err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, createdUser)
}
func (h *UserHandler) GetUserByID(c *gin.Context) {
	id := c.Param("id")
	user, err := h.services.GetUserByID(id)
	if err != nil {
		h.logger.Error("Layer: userHandler, method: GetUserByID, error:", err)
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, user)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		h.logger.Error("Layer: userHandler, method: UpdateUser, error:", err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	updatedUser, err := h.services.UpdateUser(&user)
	if err != nil {
		h.logger.Error("Layer: userHandler, method: UpdateUser, error:", err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, updatedUser)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if err := h.services.DeleteUser(id); err != nil {
		h.logger.Error("Layer: userHandler, method: DeleteUser, error:", err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "User deleted successfully"})
}

func (h *UserHandler) LoginUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		h.logger.Error("Layer: userHandler, method: LoginUser, error:", err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	loggedInUser, err := h.services.Login(user.Email, user.Password)
	if err != nil {
		h.logger.Error("Layer: userHandler, method: LoginUser, error:", err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, dto.LoginResponse{
		Token:        loggedInUser.Token,
		RefreshToken: loggedInUser.RefreshToken,
	})
}
func (h *UserHandler) GetLoggedInUser(c *gin.Context) {
	email, exists := c.Get("userEmail")
	if !exists {
		h.logger.Error("User ID not found in context")
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	user, err := h.services.GetUserByEmail(email.(string))
	if err != nil {
		h.logger.Error("Layer: userHandler, method: GetLoggedInUser, error:", err)
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	c.JSON(200, user)
}

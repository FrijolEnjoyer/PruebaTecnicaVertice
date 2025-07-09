package user

import (
	"net/http"
	"pruebaVertice/Api/dto"
	"pruebaVertice/Api/models"
	services_user "pruebaVertice/Api/services/user"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type UserHandler struct {
	userService services_user.UserService
	logger      *logrus.Logger
}

func NewUserHandler(userService services_user.UserService, logger *logrus.Logger) *UserHandler {
	return &UserHandler{
		userService: userService,
		logger:      logger,
	}
}

// CreateUser godoc
// @Summary Crear un usuario
// @Description Crea un nuevo usuario en la base de datos
// @Tags Users
// @Accept json
// @Produce json
// @Param user body models.User true "Datos del usuario"
// @Success 201 {object} models.User
// @Failure 400 {object} map[string]string
// @Router /api/auth/register/ [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		h.logger.Error("Layer: userHandler, Method: CreateUser, Error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	created, err := h.userService.CreateUser(&user)
	if err != nil {
		h.logger.Error("Layer: userHandler, Method: CreateUser, Error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, created)
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	id := c.Param("id")
	user, err := h.userService.GetUserByID(id)
	if err != nil {
		h.logger.Error("Layer: userHandler, Method: GetUserByID, Error:", err)
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		h.logger.Error("Layer: userHandler, Method: UpdateUser, Error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updated, err := h.userService.UpdateUser(&user)
	if err != nil {
		h.logger.Error("Layer: userHandler, Method: UpdateUser, Error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updated)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if err := h.userService.DeleteUser(id); err != nil {
		h.logger.Error("Layer: userHandler, Method: DeleteUser, Error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// LoginUser godoc
// @Summary Login de usuario
// @Description Inicia sesión de un usuario y devuelve tokens
// @Tags Users
// @Accept json
// @Produce json
// @Param credentials body models.User true "Credenciales del usuario"
// @Success 200 {object} dto.LoginResponse
// @Failure 400 {object} map[string]string
// @Router /api/auth/login [post]
func (h *UserHandler) LoginUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		h.logger.Error("Layer: userHandler, Method: LoginUser, Error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logged, err := h.userService.Login(user.Email, user.Password)
	if err != nil {
		h.logger.Error("Layer: userHandler, Method: LoginUser, Error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.LoginResponse{Token: logged.Token, RefreshToken: logged.RefreshToken})
}

// GetLoggedInUser godoc
// @Summary Obtener usuario logueado
// @Description Devuelve la información del usuario autenticado
// @Tags Users
// @Produce json
// @Success 200 {object} models.User
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /api/auth/me [get]
func (h *UserHandler) GetLoggedInUser(c *gin.Context) {
	emailVal, exists := c.Get("userEmail")
	if !exists {
		h.logger.Error("User email not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	email := emailVal.(string)
	user, err := h.userService.GetUserByEmail(email)
	if err != nil {
		h.logger.Error("Layer: userHandler, Method: GetLoggedInUser, Error:", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

package server

import (
	"fmt"
	"os"
	"pruebaVertice/Api/handler"
	"pruebaVertice/Api/models"
	"pruebaVertice/Api/repo/orders_repo"
	"pruebaVertice/Api/repo/products_repo"
	user_repo "pruebaVertice/Api/repo/user_repo"
	"pruebaVertice/Api/services"
	"pruebaVertice/Api/utils"
	"time"

	jwtUtils "pruebaVertice/Api/utils/jwt"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Server struct {
	router *gin.Engine
	db     *gorm.DB
	logger *logrus.Logger
}

func NewServer(db *gorm.DB, logger *logrus.Logger) *Server {
	router := gin.Default()
	server := &Server{
		router: router,
		db:     db,
		logger: logger,
	}
	server.setupRoutes()
	return server
}
func (s *Server) setupRoutes() {
	hasher := utils.BcryptHasher{}
	tokenGen := jwtUtils.JWTGenerator{}
	userService := services.NewUserService(
		user_repo.NewUserRepository(s.db, s.logger),
		hasher,
		tokenGen,
		s.logger,
	)

	userHandler := handler.NewUserHandler(userService, s.logger)
	productsService := services.NewProductsService(
		products_repo.NewProductsRepository(s.db, s.logger),
		s.logger,
	)

	productsHandler := handler.NewProductsHandler(productsService, s.logger)
	ordersService := services.NewOrdersService(
		orders_repo.NewOrdersRepository(s.db, s.logger),
		products_repo.NewProductsRepository(s.db, s.logger),
		s.logger,
	)
	ordersHandler := handler.NewOrdersHandler(ordersService, userService, s.logger)

	api := s.router.Group("/api")
	{
		user := api.Group("/auth")
		user.POST("/register", userHandler.CreateUser)
		user.POST("/login", userHandler.LoginUser)

		protected := user.Group("/")
		protected.Use(jwtUtils.GinJWTMiddleware(tokenGen, s.logger))
		{
			protected.GET("/me", userHandler.GetLoggedInUser)

			products := protected.Group("/products")
			{
				products.GET("/", productsHandler.GetAllProducts)
				products.GET("/:id", productsHandler.GetProductByID)
				products.POST("/", productsHandler.CreateProducts)
			}
			orders := protected.Group("/orders")
			{
				orders.POST("/", ordersHandler.CreateOrder)
				orders.GET("/", ordersHandler.GetUserOrders)
			}
		}
	}
}
func InitDB(logger *logrus.Logger) (*gorm.DB, error) {

	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	name := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, pass, host, name)
	var db *gorm.DB
	var err error

	maxRetries := 10
	for i := 0; i < maxRetries; i++ {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		time.Sleep(3 * time.Second)
	}
	if err != nil {
		return nil, err
	}

	if err = db.AutoMigrate(&models.User{}, &models.Product{}, &models.Order{}, &models.OrderProduct{}); err != nil {
		return nil, err
	}

	return db, nil
}

func (s *Server) Run() error {
	return s.router.Run(":8080")
}

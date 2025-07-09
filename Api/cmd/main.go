package main

import (
	"pruebaVertice/Api/server"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	_ "pruebaVertice/Api/docs"
)

// @title API de Prueba Técnica Vértice
// @version 1.0
// @description Esta API gestiona usuarios, productos y órdenes.
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Token JWT en formato Bearer. Ejemplo: "Bearer {token}"
func main() {
	logger := logrus.New()
	err := godotenv.Load(".env")
	if err != nil {
		logger.Fatalf("Error loading .env file: %v", err)
	}

	db, err := server.InitDB(logger)
	if err != nil {
		logrus.Fatalf("Failed to connect to database: %v", err)
	}
	srv := server.NewServer(db, logger)

	if err := srv.Run(); err != nil {
		logrus.Fatalf("Failed to run server: %v", err)
	}
}

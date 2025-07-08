package main

import (
	"pruebaVertice/Api/server"

	"github.com/joho/godotenv"

	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()
	//err := godotenv.Load(".env")
	err := godotenv.Load("../../.env") // Localhost
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

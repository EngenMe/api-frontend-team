package main

import (
	"log"

	"github.com/EngenMe/api-frontend-team/internal/controller"
	"github.com/EngenMe/api-frontend-team/internal/repository"
	"github.com/EngenMe/api-frontend-team/internal/service"
	"github.com/EngenMe/api-frontend-team/pkg/db"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	dbConn := db.InitDB()

	userRepo := repository.NewUserRepository(dbConn)
	userService := service.NewUserService(userRepo)
	userController := controller.NewUserController(userService)

	router := gin.Default()
	userController.SetupRoutes(router)

	//TODO: change port to be got from env
	router.Run(":8080")
}

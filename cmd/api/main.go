package main

import (
	"log"
	"os"

	"github.com/EngenMe/api-frontend-team/internal/controller"
	"github.com/EngenMe/api-frontend-team/internal/middleware"
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
	authService := service.NewAuthService(userRepo)
	authController := controller.NewAuthController(authService)
	userController := controller.NewUserController(userService)

	router := gin.Default()

	apiV1Auth := router.Group("/auth/api/v1")
	apiV1User := router.Group("/user/api/v1")
	apiV1User.Use(middleware.AuthenticationMiddleware())
	userController.SetupUserRoutes(apiV1User)
	authController.SetupAuthRoutes(apiV1Auth)

	router.Run(os.Getenv("PORT"))
}

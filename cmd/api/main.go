package main

import (
	"log"
	"os"

	"github.com/EngenMe/api-frontend-team/cmd/api/docs"
	"github.com/EngenMe/api-frontend-team/internal/controller"
	"github.com/EngenMe/api-frontend-team/internal/middleware"
	"github.com/EngenMe/api-frontend-team/internal/repository"
	"github.com/EngenMe/api-frontend-team/internal/service"
	"github.com/EngenMe/api-frontend-team/pkg/db"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           wiki auth-service API
// @version         1.0
// @description     This is a sample auth service.
// @host      localhost:8080
// @BasePath  /api/v1
// @name Authorization
// @description  Authorization header with Bearer token
// @securityDefinitions.bearer Bearer
func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	dbConn := db.InitDB()

	userRepo := repository.NewUserRepository(dbConn)
	tokenRepo := repository.NewTokenRepo(dbConn)
	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(userRepo, tokenRepo)
	authController := controller.NewAuthController(authService)
	userController := controller.NewUserController(userService)

	router := gin.Default()

	docs.SwaggerInfo.BasePath = "/api/v1"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	apiV1Auth := router.Group("/api/v1/auth")
	apiV1User := router.Group("/api/v1/user")
	apiV1User.Use(middleware.AuthenticationMiddleware())

	apiV1Auth.POST("/register", authController.Register)
	apiV1Auth.POST("/login", authController.Login)
	apiV1Auth.POST("/refresh-token", authController.RefreshToken)
	apiV1User.GET("/me", userController.GetProfile)
	apiV1User.GET("/:email", userController.GetUser)
	apiV1User.PUT("/update", userController.UpdateUser)
	apiV1User.DELETE("/delete", userController.DeleteUser)

	// userController.SetupUserRoutes(apiV1User)
	// authController.SetupAuthRoutes(apiV1Auth)

	router.Run(os.Getenv("PORT"))
}

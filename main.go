package main

import (
	"elearning_api/config"
	_ "elearning_api/docs"
	"elearning_api/handler"
	"elearning_api/middleware"
	"elearning_api/repository"
	"elearning_api/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	filesSwagger "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           My API
// @version         1.0
// @description     Dokumentasi API menggunakan Swagger di Golang dengan Gin.
// @termsOfService  http://swagger.io/terms/
// @contact.name    Support
// @contact.email   support@example.com
// @license.name    MIT
// @license.url     https://opensource.org/licenses/MIT
// @host            localhost:8080
// @BasePath        /api/v1
// @schemes         https
// @securityDefinitions.apikey BearerAuth
// @in              header
// @name            Authorization

func main() {
	r := gin.Default()
	config.LoadEnv()
	db := config.ConnectDatabase()

	userRepository := repository.NewUserRepository(db)
	tokenRepository := repository.NewTokenRepository(db)

	googleService := service.NewGoogleService()
	waService := service.NewWhatsappService()
	authService := service.NewAuthService(userRepository, tokenRepository, googleService, waService)
	userService := service.NewUserService(userRepository, tokenRepository)

	pingHandler := handler.NewPingHandler()
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)

	authMiddleware := middleware.NewAuthMiddleware(tokenRepository)

	r.Use(cors.Default())

	api := r.Group("/api")
	v1 := api.Group("/v1")
	{
		v1.GET("/ping", pingHandler.Ping)

		// auth
		auth := v1.Group("/auth")
		{
			auth.POST("/social", authHandler.RegisterOrLogin)
			auth.POST("/register", authHandler.ManualRegister)
			auth.POST("/login", authHandler.ManualLogin)
			auth.POST("/verification/send", authHandler.SendVerificationLink)
			auth.POST("/refresh-token", authHandler.RefreshToken)
		}

		user := v1.Group("/user")
		{
			user.POST("/confirm", userHandler.Confirm)
			user.GET("/profile", authMiddleware.DecodeToken(), userHandler.GetUserProfile)
		}
	}

	// Endpoint untuk Swagger UI
	r.GET("/docs/*any", ginSwagger.WrapHandler(filesSwagger.Handler))

	r.Run(":8080")
}

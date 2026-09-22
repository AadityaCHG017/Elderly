package main

import (
	"eldercare/internal/config"
	"eldercare/internal/database"
	"eldercare/internal/users"
	"log"
	"net/http"

	"eldercare/internal/auth"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	pool, err := database.Connect(cfg.DatabaseURL)

	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	router := gin.Default()

	userRepository := users.NewRepository(pool)
	userService := users.NewService(userRepository)
	userHandler := users.NewHandler(userService)

	router.POST("/users", userHandler.CreateUser)

	authService := auth.NewService(
		userRepository,
		cfg.JWTSecret,
	)

	authHandler := auth.NewHandler(authService)

	router.POST("/auth/login", authHandler.Login)

	protected := router.Group("/api")
	protected.Use(auth.AuthMiddleware(cfg.JWTSecret))

	protected.GET("/profile", func(c *gin.Context) {

		claims, exists := c.Get("claims")

		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "claims not found",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "You are authenticated",
			"claims":  claims,
		})
	})

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Eldercare Api is Running",
			"databse": "connected",
		})
	})

	router.Run(":" + cfg.Port)
}

package main

import (
	"eldercare/internal/config"
	"eldercare/internal/database"
	"eldercare/internal/users"
	"log"
	"net/http"

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

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Eldercare Api is Running",
			"databse": "connected",
		})
	})

	router.Run(":" + cfg.Port)
}

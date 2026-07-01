package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rayvald29/petfinder-api/database"
	"github.com/rayvald29/petfinder-api/handlers"
	"github.com/rayvald29/petfinder-api/middleware"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	database.ConnectDB()

	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "PetFinder API is running!")
	})

	router.GET("/health", handlers.HealthHandler)
	router.POST("/register", handlers.RegisterHandler)
	router.POST("/login", handlers.LoginHandler)

	protected := router.Group("/")
	protected.Use(middleware.AuthMiddleware())
	protected.GET("/profile", handlers.ProfileHandler)
	router.Run(":8080")
}

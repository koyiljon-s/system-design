// cmd/api/main.go
package main

import (
	"log"
	
	"primejobs/user-service/internal/database"
	"primejobs/user-service/internal/handler"
	"primejobs/user-service/internal/repository"
    "primejobs/user-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	// Connect to database
	if err := database.Connect(); err != nil {
		log.Fatal("Database connection failed:", err)
	}

	userRepo := repository.NewUserRepository()
	userHandler := handler.NewUserHandler(userRepo)

	// Create Gin router
	r := gin.Default() // Includes logger and recovery middleware

	// Public routes
	r.POST("/api/register", userHandler.Register)
	r.POST("/api/login", userHandler.Login)

	protected := r.Group("/api")
    protected.Use(middleware.AuthMiddleware())
	{
	  protected.GET("/me", userHandler.GetMe)
    }

	log.Println("Server running on http://localhost:8080")
	r.Run(":8080")
}
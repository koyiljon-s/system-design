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
	oauthHandler := handler.NewOAuthHandler(userRepo)

    gin.SetMode(gin.ReleaseMode)

	// Create Gin router
	r := gin.Default() 

	// Public routes
	api := r.Group("/api")
	{
		api.POST("/register", userHandler.Register)
		api.POST("/login", userHandler.Login)

		api.GET("/oauth/google", oauthHandler.GoogleLogin)
		api.GET("/oauth/google/callback", oauthHandler.GoogleCallback)
	}
   
    
	protected := api.Group("/")
    protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/me", userHandler.GetMe)
		protected.PUT("/me", userHandler.UpdateMe)
		protected.DELETE("/me", userHandler.DeleteMe)
        // internal usage
		protected.GET("/users/:id", userHandler.GetUserByID)
    }

	log.Println("Server running on http://localhost:8080")
	if err := r.Run(":8000"); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
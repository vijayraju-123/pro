package main

import (
	"log"
	"task-management-system/Backend/config"
	"task-management-system/Backend/handlers"
	"task-management-system/Backend/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize configuration (database, ports, etc.)
	config.Init()

	// Set up Gin router
	router := gin.Default()

	// Apply middleware
	router.Use(middleware.CORSMiddleware())
	router.Use(gin.Logger())

	// Public routes (no auth needed)
	public := router.Group("/api")
	{
		public.POST("/auth/signup", handlers.Signup)
		public.POST("/auth/login", handlers.Login)
	}

	// Protected routes (require JWT)
	protected := router.Group("/api")
	protected.Use(middleware.JWTAuthMiddleware())
	{
		// Task routes
		protected.POST("/tasks", handlers.CreateTask)
		protected.GET("/tasks", handlers.GetTasks)
		protected.PUT("/tasks/:id", handlers.UpdateTask)
		protected.DELETE("/tasks/:id", handlers.DeleteTask)

		// WebSocket for real-time updates
		protected.GET("/ws", handlers.WSHandler)
	}

	// Start the server
	if err := router.Run(":" + config.ConfigInstance.Port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}